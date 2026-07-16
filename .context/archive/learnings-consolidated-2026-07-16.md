# Archived Learnings (consolidated 2026-07-16)

Originals replaced by consolidated entries in LEARNINGS.md.

## Group: OpenCode plugin integration gotchas

## [2026-04-29-050000] BunShell ctx.$ calls echo stdout to OpenCode's process unless .quiet() is set — leaks visible noise

**Context**: After PR #72 wired session.created and session.idle to fire `ctx
system bootstrap`, `ctx agent --budget 4000`, and friends, end users started
seeing chunks of Markdown bleeding into the OpenCode TUI: `## Steering`, `#
Product Context`, `Describe the product...`. These are the contents of
`.context/steering/` template stubs that `ctx agent --budget 4000` includes in
its context packet. The plugin used the shell-level `2>/dev/null || true` to
swallow stderr and force exit 0, but stdout was untouched.

**Lesson**: BunShell's documented behavior: *"By default, the shell will write
to the current process's stdout and stderr, as well as buffering that output."*
So an `await ctx.$\`...\`` call in a plugin echoes its stdout/stderr to
OpenCode's process, which the TUI/agent surfaces. Shell-level `2>/dev/null` only
suppresses stderr; stdout still leaks. The fix is BunShell's `.quiet()` modifier
on the BunShellPromise, which configures the shell to only buffer the output
rather than also writing to the parent process.

**Application**: Always chain `.nothrow().quiet()` on BunShell template literals
in OpenCode plugins, even for fire-and-forget calls where you discard the
result: `await ctx.$\`ctx system bootstrap\`.nothrow().quiet()`. With both
modifiers, you don't need shell-level `2>/dev/null || true` — `.nothrow()`
swallows non-zero exits at the BunShell layer, `.quiet()` keeps every byte of
output buffered. Pattern is the cooperative default for any plugin that spawns
long-output commands during the agent session lifecycle.

---

## [2026-04-29-040000] OpenCode plugin compaction interop is breadcrumb-mediated: own your context preservation explicitly

**Context**: After PR #72 wired `session.created` / `session.idle` /
`tool.execute.after` / `shell.env`, a `/compact` test in OpenCode (with
`oh-my-openagent@3.17.6` also installed) recovered ctx context post-compaction
*only by accident*: oh-my-openagent's `experimental.session.compacting` handler
builds a structured summary template that happens to preserve
`.context/`-prefixed file paths in its "Active Working Context → Files"
section. Combined with our `shell.env` CTX_DIR injection, the agent had enough
breadcrumbs to re-read DECISIONS.md from disk post-compaction. Without that
section, our context would have evaporated silently into the compaction summary.

**Lesson**: Two compaction-aware plugins in the same session can synergize
without either knowing about the other — but the synergy is fragile because it
depends on undocumented serialization choices in the *other* plugin. If the
other plugin's template ever changes (e.g., drops file-path preservation, swaps
the "Active Working Context" section name, condenses paths to basenames), the
breadcrumbs disappear and ctx context is lost without any signal. The `Hooks`
interface in `@opencode-ai/plugin` v1.4.x exposes
`experimental.session.compacting?: (input, output: { context: string[]; prompt?:
string }) => Promise<void>` — pushing to `output.context` is *additive*
(appends to the default prompt), and replacing `output.prompt` is *destructive*
(only one plugin can win that race).

**Application**: Register `experimental.session.compacting` in your own plugin
and push high-signal context strings (e.g., `ctx system bootstrap` output) to
`output.context` so context preservation does not depend on coexisting plugins.
Never set `output.prompt` from a thin shim — that would conflict with primary
compaction harnesses like oh-my-openagent. Composition via `output.context` is
the correct cooperative pattern.

---

## [2026-04-29-030000] @opencode-ai/plugin event hook is a single dispatcher, not an object of named handlers

**Context**: PR #72's first OpenCode plugin shipped with `event: {
"session.created": fn, "session.idle": fn }` — an object keyed by event type.
It compiled clean against `satisfies Plugin` but never fired. End-to-end trace
showed neighboring hooks (`shell.env`, `tool.execute.after`) running while every
event handler silently no-op'd.

**Lesson**: `@opencode-ai/plugin` v1.4.x defines `event?: (input: { event: Event
}) => Promise<void>` — one dispatcher called for every event with
`input.event.type` discriminating. Asymmetric with neighbors because `shell.env`
and `tool.execute.*` *are* top-level named keys; only the dozens of `EventX`
types collapse into the single `event` slot.

**Application**: Use `event: async ({event}) => { if (event.type ===
"session.created") { ... } else if (event.type === "session.idle") { ... } }`.
Type discriminator strings live under each `EventX` type in
`node_modules/@opencode-ai/sdk/dist/gen/types.gen.d.ts`.

---

## [2026-04-29-030100] OpenCode plugin hooks like shell.env take (input, output) and mutate; returned objects are ignored

**Context**: First plugin had `"shell.env": () => ({ CTX_DIR: ".context" })`.
The hook fired but the agent's bash tool never saw `CTX_DIR`; manual export was
required for every ctx call. The returned object was dropped on the floor by the
runtime.

**Lesson**: Multiple hooks in `@opencode-ai/plugin` v1.4.x take two arguments
where the second is an OUT param. Examples: `shell.env: (input, output: {env})
=> void` (mutate `output.env`), `tool.execute.after: (input, output: {title,
output, metadata}) => void`, `chat.params: (input, output: {temperature, ...})
=> void`, `chat.headers: (input, output: {headers}) => void`. Pattern is
consistent across the SDK.

**Application**: Always read the type definition in
`node_modules/@opencode-ai/plugin/dist/index.d.ts` for any hook before wiring.
If a hook signature has two parameters where the second is an object, it's a
mutation hook — return values are discarded.

---

## [2026-04-29-030200] OpenCode shell.env injects env only into agent's shell tool, not into plugin's own ctx.$ calls

**Context**: After fixing `shell.env`'s `(input, output) => mutate output.env`
signature so `CTX_DIR` reached the agent's bash tool, the plugin's own
`ctx.$\`ctx system bootstrap\`` calls still failed silently — they ran without
`CTX_DIR` and ctx fell back to `~/.context`. The hook fired correctly; the
plugin's subprocess side-effects didn't see the env.

**Lesson**: `shell.env` injects env into the agent's shell-tool invocations. The
plugin's own BunShell calls (`ctx.$\`...\``) inherit OpenCode's process env,
which is *separate*. Two shells, two envs.

**Application**: Build an env-aware BunShell once in the plugin factory: `const
$ = ctx.$.env({ ...process.env, CTX_DIR: \`${ctx.directory}/.context\` })`.
Reuse it for every plugin-initiated subprocess call. `ctx.directory` is the
project root from `PluginInput`.

---

## [2026-04-26-180000] OpenCode auto-loads only flat .ts files under .opencode/plugins/; subdirectories are ignored

**Context**: Initial OpenCode integration deployed the plugin as
`.opencode/plugins/ctx/index.ts` (a directory with index.ts inside, mirroring
npm package conventions). End-to-end smoke testing showed the plugin file was
present and the binary was current, yet OpenCode never invoked any of the
plugin's hooks (no `module-load` trace fired even with `--print-logs --log-level
DEBUG`). Copying the same content to a flat `.opencode/plugins/ctx.ts` file made
the plugin load and fire correctly.

**Lesson**: OpenCode's plugin auto-discovery only scans top-level files under
`.opencode/plugins/` and `~/.config/opencode/plugins/`. Subdirectories are
silently skipped — there is no log line indicating a subdirectory was found
and ignored. The official docs at opencode.ai/docs/plugins/ say only "files in
these directories are automatically loaded at startup" without specifying the
rule, so this is easy to miss. The `opencode plugin <module>` CLI registers npm
modules (a different code path) and accepts only npm names, not local paths.

**Application**: Deploy single-file plugins as `.opencode/plugins/<name>.ts`,
not `.opencode/plugins/<name>/index.ts`. No `package.json` is required when the
plugin uses type-only imports (`import type` is erased at compile time) and the
host runtime injects the plugin context. To verify a plugin is actually loaded,
add a top-of-module side effect (e.g. `appendFileSync` to a known path) and
confirm it fires before debugging hook contracts.

---

## [2026-04-26-165500] OpenCode opencode.json MCP shape: command is Array<string>, no separate args field

**Context**: `ctx setup opencode --write` was generating `opencode.json` with
the Copilot CLI MCP shape (`{type: "local", command: "ctx", args: ["mcp",
"serve"]}`). OpenCode rejected the file at startup with `Configuration is
invalid… Expected array, got "ctx" mcp.ctx.command` and `Missing key
mcp.ctx.enabled`.

**Lesson**: OpenCode's `McpLocalConfig` (in `@opencode-ai/sdk`) defines
`command: Array<string>` as a single field that holds the binary AND its
arguments — there is no separate `args` field. It also requires `enabled:
boolean` at runtime even though the TS type marks it optional. The Copilot CLI
MCP shape is similar in spirit but structurally different; do not copy-paste
between them.

**Application**: For OpenCode MCP entries always use `command: ["ctx", "mcp",
"serve"]` and include `enabled: true`. If you add a new editor integration with
its own MCP file format, read the upstream type definitions from
`node_modules/@<vendor>/sdk/dist/gen/types.gen.d.ts` (or equivalent) before
reusing an existing generator.

---

## Group: site/ tracked build output

## [2026-02-27-231228] site/ directory must be committed with docs changes

**Context**: The site/ directory contains generated HTML served directly from
the repo (no CI build step). Multiple sessions have committed docs/ changes
without the corresponding site/ output, or ignored site/ as 'generated noise'.

**Lesson**: site/ is intentionally tracked in git — there is no GitHub Pages
workflow or CI step to build it. When docs change, the regenerated site/ HTML
must be staged and committed alongside the source.

**Application**: Always git add site/ when committing changes under docs/. Never
gitignore site/.

---

## [2026-06-07-162840] site/ is tracked build output — rebuild and bundle it with doc commits

**Context**: A docs change (docs/cli/dream.md, etc.) produced a surprise
189-file site/ drift mid-session; site/ is tracked (zensical build via 'make
site'), not gitignored.

**Lesson**: Any change under docs/ requires regenerating site/ with 'make site'
and committing the rebuilt site/ in the SAME commit (cf. f0f100a0, which bundled
its site rebuild). Otherwise the built output silently drifts and shows up as a
large untracked/modified set later.

**Application**: After editing anything under docs/, run 'make site' and 'git
add site/' and include it in the doc commit. Don't treat site/ as ephemeral —
it's versioned.

---

## Group: Contributor PR handling

## [2026-04-01-074418] Contributor PRs based on older code reintroduce removed features

**Context**: PR #45 brought back prompt templates, PROMPT.md, and
IMPLEMENTATION_PLAN.md that were explicitly removed in March

**Lesson**: When resolving contributor merge conflicts, check decisions history
for intentional removals — do not assume the PR content is additive

**Application**: Cross-reference DECISIONS.md before accepting PR content that
adds files or features

---

## [2026-03-15-101342] Contributor PRs need post-merge follow-up commits for convention alignment

**Context**: PR #42 (MCP v0.2) addressed bulk of review feedback but left ~12
inline strings, no embed_test coverage, and substring matching in
containsOverlap

**Lesson**: Merging with known gaps is fine when the gaps are mechanical, but
the follow-up must be immediate — track in ideas/done/ with a review status
doc

**Application**: For future contributor PRs: create ideas/pr{N}-review-status.md
during review, merge when architecture is sound, fix convention gaps in a
same-day follow-up commit

---

## Group: Drift/detection false positives

## [2026-02-27-230738] Drift detector false positives on illustrative code examples

**Context**: ctx drift flagged 23 warnings for backtick-quoted paths in
CONVENTIONS.md and ARCHITECTURE.md that were prose examples (loader.go,
session/run.go, sync.Once), not real file references.

**Lesson**: Path reference detection should verify the top-level directory
exists on disk before flagging. Bare filenames and paths under non-existent
directories are almost always examples in documentation.

**Application**: The fix checks os.Stat(topDir) on the first path component.
Future drift checks on documentation-heavy files should use the same heuristic.

---

## [2026-03-24-001001] lint-drift false positives from conflating constant namespaces

**Context**: lint-drift.sh checked all string constants in embed/cmd/*.go
against commands.yaml, but Use* constants are cobra syntax strings, not YAML
lookup keys

**Lesson**: Shell grep on constant values cannot distinguish constant types;
only DescKey* constants are YAML keys. AST-based analysis is needed for
type-aware checks

**Application**: Already captured in specs/ast-audit-tests.md; the lint-drift
fix is shipped in v0.8.0

---

## [2026-03-23-165611] Typography detection script needs exclusion lists for intentional uses

**Context**: detect-ai-typography.sh flagged config/token/delim.go (intentional
delimiter constants) and test files (test data containing em-dashes)

**Lesson**: Detection scripts for convention enforcement need exclusion patterns
for files where the flagged patterns are intentional data, not prose

**Application**: Add exclusion patterns proactively when creating detection
scripts; *_test.go and constant-definition files are common false positive
sources

---


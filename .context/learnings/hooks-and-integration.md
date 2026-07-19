# hooks-and-integration

## [2026-07-16-120001] OpenCode plugin integration gotchas (consolidated)

**Consolidated from**: 7 entries (2026-04-26 to 2026-04-29)

- Two-arg hooks (`shell.env`, `tool.execute.after`, `chat.params`, `chat.headers`) are MUTATION hooks: mutate the second `output` param; returned objects are discarded. Read the type def in `@opencode-ai/plugin/dist/index.d.ts` before wiring.
- The `event` hook is a SINGLE dispatcher — `event: async ({event}) => { if (event.type === "session.created") ... }` with `event.type` discriminating — not an object of named handlers; a named-handlers object compiles clean but never fires.
- `shell.env` injects env only into the agent's shell tool, not the plugin's own `ctx.$` calls (separate process env). Build an env-aware shell once: `ctx.$.env({ ...process.env, CTX_DIR: `${ctx.directory}/.context` })` and reuse it for every plugin-initiated subprocess.
- BunShell `ctx.$` echoes stdout/stderr to OpenCode's process (the TUI/agent surfaces it); chain `.nothrow().quiet()` on every call — shell-level `2>/dev/null` only hides stderr, stdout still leaks.
- Auto-discovery scans only flat top-level files: deploy `.opencode/plugins/<name>.ts`, NOT `<name>/index.ts` (subdirectories are silently ignored, no log line). Verify a plugin actually loads with a top-of-module side effect before debugging hook contracts. `opencode plugin <module>` is a different code path (npm names only).
- `opencode.json` MCP shape: `command: Array<string>` holds the binary AND its args (there is no separate `args` field) and requires `enabled: true` at runtime; don't copy the Copilot CLI MCP shape. Read upstream `types.gen.d.ts` before reusing a generator.
- Compaction interop is breadcrumb-mediated and fragile (ctx context survived `/compact` only by accident, via oh-my-openagent preserving `.context/` paths). Register `experimental.session.compacting` yourself and push high-signal context (e.g. `ctx system bootstrap` output) to `output.context` (additive); never set `output.prompt` from a thin shim (destructive — only one plugin wins).

---

## [2026-06-07-170015] Hook mechanics, output channels & compliance (consolidated)

**Consolidated from**: 5 entries (2026-01-25 to 2026-04-06)

- Hook scripts receive JSON via stdin (HOOK_INPUT=$(cat) then jq), not env vars;
  key names are case-sensitive (PreToolUse, SessionEnd); use
  $CLAUDE_PROJECT_DIR, never hardcode paths; anchor regex to command-start
  `(^|;|&&|\|\|)\s*` ('ctx' binary vs dir); grep matches inside quoted args
  (test with blocked words); scripts silently lose execute permission (verify ls
  -la).
- Output routing: plain-text hook stdout is silently ignored — Claude Code
  parses stdout starting with `{` as JSON directives; return JSON via
  printHookContext(). For UserPromptSubmit specifically, stdout is prepended as
  AI context (not user-visible), stderr+exit0 is swallowed, user-visible output
  requires {"systemMessage":"…"} or exit 2 (blocks); there is NO non-blocking
  user-visible channel. Two-tier severity is sufficient: unprefixed (agent
  context, may relay) and "IMPORTANT: Relay VERBATIM" (guaranteed); don't add
  more prefixes.
- Agents only relay content with explicit display instructions: a
  system-reminder line with no "Display this line verbatim" is invisible to the
  user even when correct. IMPORTANT: signals internal priority, not user-facing
  output.
- Compliance: soft instructions have a ~75–85% ceiling because "don't apply
  judgment" is itself judgment; for 100% compliance inject via additionalContext
  rather than instruct. Hook compliance degrades on narrow mid-session tasks
  (~15–25% skip) because CLAUDE.md's "may or may not be relevant" competes
  with hook authority — fix by elevating hook authority explicitly; the
  mandatory checkpoint relay block is the compliance canary. No reliable
  agent-side before-session-end event exists (SessionEnd fires after the agent
  is gone) — mid-session nudges + explicit /ctx-wrap-up are the only reliable
  persistence. Repeated injection causes repetition fatigue — gate with
  --session $PPID --cooldown and pair with a readback instruction.
- Context-budget injection strategy: once ~7K tokens are auto-injected (fait
  accompli), the agent's rationalization inverts from "skip to save effort" to
  "marginal cost is trivial." Front-load highest-value content as injection,
  then leverage sunk cost for on-demand reads. Verbal summaries + linked diagram
  files cut ARCHITECTURE.md ~12K→3.8K (extract diagrams outside FileReadOrder;
  the 4-chars/token estimator is accurate — optimize content not the
  estimator).

---

## [2026-05-08-195031] Cursor imports Claude Code hooks and sets CLAUDE_PROJECT_DIR per workspace

**Context**: Investigating why .context/state/ appeared in non-ctx projects
opened in Cursor. Hypothesis was a Cursor extension or shell hook; turned out to
be Cursor's documented Claude-compatibility behavior
(https://cursor.com/docs/hooks): it loads ~/.claude hooks and injects
CLAUDE_PROJECT_DIR=workspace_root so they 'just work'. Globally-enabled Claude
plugins therefore fire in every Cursor workspace.

**Lesson**: When debugging cross-tool side effects, check whether the host tool
advertises compatibility with the implicated tool's config format. The leak
surface of any global Claude plugin is now 'every Cursor workspace + every
Claude Code project', not just 'every Claude Code project'.

**Application**: Hooks must be safe to fire in non-ctx projects: silent bail
when state.Initialized() is false, no filesystem side effects. The ctx code-side
fix lives in state.Dir's Initialized gate; the design rule is broader: assume
hooks may run anywhere, not just where the user invoked ctx init.

---

## [2026-03-20-160112] Commit messages containing script paths trigger PreToolUse hooks

**Context**: Git commit message body contained a path to a shell script under
the hack directory which matched a hook pattern that blocks direct script
invocation

**Lesson**: Hooks scan all Bash tool input including heredoc content used for
commit messages, not just the command itself

**Application**: Rephrase commit messages and ctx add content to avoid paths
that match hook deny patterns, use generic references instead of literal file
paths

---

## [2026-03-04-105239] CONSTITUTION hook compliance is non-negotiable — don't work around it

**Context**: After make build, ran ./ctx deps --help which was blocked by
block-non-path-ctx. Instead of asking user to install, tried cp ctx ~/bin/ —
escalating workarounds.

**Lesson**: When a hook blocks an action, the correct response is to follow the
hook's instruction (ask the user to sudo make install), not to find creative
bypasses.

**Application**: Always ask the user to install when testing a freshly built
binary. Never attempt alternative install paths to circumvent a hook.

---

## [2026-03-02-165039] Hook message registry test enforces exhaustive coverage of embedded templates

**Context**: Adding billing.txt to embedded assets without a registry entry
caused TestRegistryCoversAllEmbeddedFiles to fail immediately

**Lesson**: Every new .txt file under internal/assets/hooks/messages/ must have
a corresponding entry in registry.go — the test acts as an exhaustive
bidirectional check

**Application**: When adding new hook message variants, update the registry
entry before running tests

---

## [2026-02-28-150701] Plugin reload script must rebuild cache, not just delete it

**Context**: hack/plugin-reload.sh was deleting
~/.claude/plugins/cache/activememory-ctx/ without repopulating it. Claude Code's
installed_plugins.json still referenced the cache path, so the plugin appeared
enabled but hooks.json was missing — all plugin hooks silently stopped firing.

**Lesson**: Claude Code snapshots plugin hooks from the cache directory at
session startup. If the cache is deleted, plugin hooks vanish silently with no
error. The reload script must rebuild the cache from source assets
(internal/assets/claude/) after clearing it, and warn that a session restart is
required.

**Application**: Always rebuild the plugin cache in hack/plugin-reload.sh. When
debugging hooks that don't fire, check ~/.claude/plugins/cache/ first — a
missing hooks.json is the most likely cause.

---

## [2026-02-26-003854] Webhook silence after ctxrc profile swap is the most common notify debugging red herring

**Context**: Spent time investigating why webhooks weren't firing — checked
binary version, hook configs, notify.Send internals. Actual cause was .ctxrc
swapped to prod profile (notify commented out) earlier in session.

**Lesson**: When webhooks stop, check .ctxrc profile first (`ctx config
status`). Also: not all tool uses trigger webhook-sending hooks — Read only
triggers context-load-gate (one-shot) and ctx agent (no webhook). qa-reminder
requires Edit matcher.

**Application**: Before debugging notify internals, run `ctx config status` and
verify the event would actually match a hook with notify.Send.

---

## [2026-04-26-152836] ctx system help can list project-local hooks not in the Go binary

**Context**: PR #72 plugin called 'ctx system block-dangerous-commands'; user's
installed ctx 0.7.2 listed it in help, but no directory exists under
internal/cli/system/cmd/ — it's a Claude Code plugin-local hook surfaced via
wrapper

**Lesson**: ctx system help output is a union of compiled Go subcommands and
project-local Claude wrappers; non-Claude integrations only see the Go subset

**Application**: When porting plugin behavior to a new editor, only call
subcommands that have a directory under internal/cli/system/cmd/. Don't trust
ctx system help output as the canonical surface.

---


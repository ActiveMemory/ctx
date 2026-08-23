# Spec: Pi CLI Integration for ctx

## Context

Pi (`earendil-works/pi`, pi.dev) is a "self-extensible" coding-agent CLI and
agent harness (TypeScript monorepo, Bun-based, MIT, actively developed as of
2026-08-23). Its design principle: keep the core small, and intentionally ship
**no built-in MCP, sub-agents, permission popups, plan mode, or background
bash**. Workflow behavior moves to four surfaces:

1. **Extensions** — TypeScript modules auto-discovered from
   `~/.pi/agent/extensions/` (global) and `.pi/extensions/` (project-local,
   loaded only after project trust; hot-reloadable via `/reload`). Flat
   `.pi/extensions/*.ts` and subdirectory form `.pi/extensions/<name>/index.ts`
   are both supported. Events of interest:
   - `session_start` (reason `startup | reload | new | resume | fork`)
   - `session_before_compact` / `session_compact` / `session_compact_failed`
   - `before_agent_start` — can inject a **persistent message**
     (`{message: {customType, content, display}}`, stored in the session and
     sent to the LLM) and/or modify the system prompt
   - `tool_call` (pre-execution, can block/modify) and `tool_result`
     (post-execution, middleware-style; carries `isError`)
   - `agent_settled` — fired when Pi will not continue running automatically
   - On `/new`, `/resume`, `/fork`, and `/reload` Pi tears down the extension
     instance (`session_shutdown` → reload → `session_start`), so in-memory
     extension state resets — flag loss can only cause an extra injection,
     never a missed one
   - `ExtensionContext` exposes `ctx.cwd`, `ctx.hasUI`,
     `ctx.isProjectTrusted()`, `ctx.sessionManager`, and `ctx.signal`
   - Available imports: `@earendil-works/pi-coding-agent` (types), `typebox`,
     `@earendil-works/pi-ai`, `@earendil-works/pi-tui`, Node built-ins
     (`node:child_process`, `node:path`, …)
   - Built-in tool names: `read`, `bash`, `edit`, `write`, `grep`, `find`,
     `ls`
2. **Skills** — the Agent Skills standard (agentskills.io). Project locations:
   `.pi/skills/` in the project, and `.agents/skills/` in cwd and ancestor
   directories up to the git root (only `.agents` gets the ancestor walk; both
   load post-trust). Global: `~/.pi/agent/skills/`, `~/.agents/skills/`. Also
   package.json `pi.skills` and a `skills` array in `.pi/settings.json`.
   Required frontmatter: `name` (lowercase a-z, 0-9, hyphens; need not match
   the parent directory) and `description`. Skills register as
   `/skill:<name>` commands.
3. **AGENTS.md** — loaded natively as a context file
   (`systemPromptOptions.contextFiles`).
4. **MCP — absent by design.** No MCP registration path exists.

ctx already ships the OpenCode integration on the established blueprint
(`specs/opencode-integration.md`): a Go setup package
(`internal/cli/setup/core/opencode/`), a thin embedded TypeScript shim
(`internal/assets/integrations/opencode/plugin/index.ts`) that shells out to
`ctx system` subcommands, bundled skills
(`internal/assets/integrations/opencode/skills/`), and the shared
`agent.AgentsMd()` project-root template.

## Goal

Add `ctx setup pi [--write]` that deploys, following the OpenCode blueprint:

1. `.pi/extensions/ctx.ts` — thin shim extension (embedded asset, not an npm
   dependency; type-only import of `@earendil-works/pi-coding-agent`,
   `node:child_process` for subprocess). Flat single-file deploy for exact
   OpenCode symmetry; Pi auto-loads flat `.pi/extensions/*.ts`. v1 ships no
   package.json next to the extension (no runtime deps); the subdirectory form
   is left for a future deps-carrying extension.
2. `.pi/skills/ctx-*/SKILL.md` — the same bundled skill set as OpenCode
   (frontmatter already conforms to the Agent Skills standard)
3. AGENTS.md at the project root (shared `coreAgents.Deploy(cmd)`, as
   `internal/cli/setup/core/opencode/opencode.go` does)

**Why this shape:** every ctx integration is a Go package that deploys config +
assets; all real logic stays in Go via `ctx system`. Pi's lack of MCP means the
extension is the only lifecycle channel — the same "thin shim" discipline
applies, just without an MCP leg.

## Extension wiring (`.pi/extensions/ctx.ts`)

### Hook envelope requirement (load-bearing)

The `ctx system` hook subcommands (`post-commit`, `check-task-completion`,
`check-persistence`) parse a **hook-JSON envelope from stdin** via
`coreCheck.FullPreamble`; with no envelope they bail silently. The extension
must therefore spawn these with piped stdin and write
`{"session_id": "<session id>", "tool_input": {"command": "<command>"}}`
before closing stdin (`child.stdin.end(json)`). Closing stdin immediately also
avoids the `hook.StdinReadTimeout` (2s) stall on open pipes. Consequences:

- **`node:child_process` is mandatory** for these calls: Pi's own exec helper
  hardcodes `stdio: ["ignore", "pipe", "pipe"]` (no stdin pipe), so the
  envelope cannot be delivered through it.
- `session_id` comes from `ctx.sessionManager` so per-session counters and
  pause markers are keyed correctly instead of collapsing to `IDUnknown`.
- `ctx agent` is a normal flag-driven command and needs **no** envelope.

### Event map

| Pi event | Action |
|----------|--------|
| `session_start` | Warm-up: run `ctx agent --budget 4000` (cwd-anchored to `ctx.cwd`), cache the packet; keep out of the prompt path |
| `before_agent_start` | Inject the cached packet as a persistent `message` **only when** no ctx-injected `custom_message` (our `customType`) exists in the current branch **after the most recent compaction entry** (branch-scan predicate via `ctx.sessionManager`) |
| `session_compact` (success) | No-op for the flag (the branch-scan predicate subsumes it); optionally drop the cache so the next injection re-warms |
| `tool_result` — tool `bash`, command matches `git commit`, **`!event.isError`** | Envelope call `ctx system post-commit` |
| `tool_result` — tool `edit` or `write`, **`!event.isError`** | Envelope call `ctx system check-task-completion` |
| `agent_settled` | Envelope call `ctx system check-persistence` |

Design notes:

- **`tool_result`, not `tool_call`,** for the post-commit trigger:
  `tool_call` fires before execution (blocking semantics); `tool_result`
  mirrors OpenCode's post-execution `tool.execute.after` trigger.
- **`isError` gating** (Pi-specific; OpenCode has no equivalent signal): a
  failed `git commit` must not fire post-commit (it would re-score the
  *previous* commit via `ScoreCommitViolations()` and emit a spurious nudge).
  Gate both post-commit and edit/write branches on `!event.isError`.
- **Compaction interop is breadcrumb-mediated and stateless across the
  reload:** Pi's own compaction mechanics validate the design — an injected
  `custom_message` is a valid compaction cut point, so once it ages past
  `keepRecentTokens` it is folded into the lossy LLM summary and re-injection
  is genuinely needed; scanning the branch for our `customType` since the last
  compaction entry decides precisely (no duplicate when the packet is still in
  the kept tail, no `/reload` duplicate, no missed injection after
  switch/fork/reload — extension state resets are fail-safe). The rejected
  alternative: a custom `session_before_compact` summary would *replace* Pi's
  LLM summary; cross-extension precedence is undocumented, so we don't take
  ownership of the summary.
  Known edge: after an overflow compaction with `willRetry`, the retried turn
  does not pass through `before_agent_start`, so re-injection waits for the
  next user prompt — accepted.
- **Latency:** the packet is produced once at `session_start` (warm-up) and
  injected from cache on `before_agent_start`, so the prompt path never waits
  on `ctx agent`. All subprocess calls receive `ctx.signal` and a timeout so
  Esc cancels cleanly; swallowed non-zero exits (nothrow equivalent), bounded
  output. If the `ctx` binary is absent, the extension no-ops silently.
- Unrecognized tool names silently no-op. Tool-name sets are pinned to Pi's
  built-ins (`bash`; `edit`, `write`) — do **not** carry OpenCode's
  `shell`/`file_edit` names. Verify against Pi's docs when Pi bumps.
- The agent's own shell tool is not anchored by the extension; users launch
  `pi` from the project root (documented in the quickstart, same caveat as
  OpenCode).

## Files to create

```
internal/assets/integrations/pi/
├── extension/
│   └── ctx.ts            # Thin shim extension (~150 lines with envelope +
│                         # branch-scan + isError gating)
└── skills/               # Same bundled skill set as OpenCode (10 skills)
    ├── ctx-agent/SKILL.md
    ├── ctx-handover/SKILL.md
    ├── ctx-kb-ask/SKILL.md
    ├── ctx-kb-ground/SKILL.md
    ├── ctx-kb-ingest/SKILL.md
    ├── ctx-kb-note/SKILL.md
    ├── ctx-kb-site-review/SKILL.md
    ├── ctx-remember/SKILL.md
    ├── ctx-status/SKILL.md
    └── ctx-wrap-up/SKILL.md

internal/cli/setup/core/pi/
├── doc.go                # package doc (docstring floor: related packages)
├── pi.go                 # Deploy entry: extension (fatal) → coreAgents.Deploy
│                         # (warn) → skills (warn) → InfoPiSummary
├── extension.go          # deploy .pi/extensions/ctx.ts (refresh in place)
├── skill.go              # deploy .pi/skills/ctx-*/SKILL.md (sorted, deterministic)
├── validate.go           # managed-target validation (OpenCode semantics)
└── deploy_test.go        # modeled on opencode's suite

internal/config/asset/asset.go                    # DirIntegrationsPiExtension,
                                                  # DirIntegrationsPiSkill
internal/assets/read/agent/pi.go                  # PiExtension(), PiSkills()
                                                  # accessors (embedded-fs reads;
                                                  # new file is a deliberate
                                                  # deviation from agent.go's
                                                  # single-file pattern — noted,
                                                  # not silent)
internal/config/hook/hook.go                      # ToolPi = "pi" (tool-const block),
                                                  # DirPi = ".pi", DirPiExtensions,
                                                  # FilePiExtension = "ctx.ts",
                                                  # DirPiSkills = "skills"
internal/assets/commands/text/hooks.yaml          # hook.pi instructions text +
                                                  # "pi" line in hook.supported-tools
internal/assets/commands/text/write.yaml          # write.hook-pi-created/-skipped/-summary
internal/config/embed/text/hook.go                # DescKeyHookPi, DescKeyWriteHookPiCreated/
                                                  # Skipped/Summary
internal/write/setup/hook.go                      # InfoPiCreated, InfoPiSkipped,
                                                  # InfoPiSummary
internal/cli/setup/cmd/root/run.go                # case cfgHook.ToolPi branch
internal/cli/setup/cmd/root/doc.go                # docstring update
tools/typecheck/pi/                               # tsconfig.json + package.json
                                                  # (@earendil-works/pi-coding-agent
                                                  # types, tsc --noEmit)
.github/workflows/ci.yml                          # typecheck-pi-extension job
```

`internal/config/setup` needs no Pi entry (Pi has no MCP/global-config leg).

## CLI surface

No new cobra subcommand and no new `Use*` constant — tools are
`case cfgHook.Tool*` branches inside `root.Run()` (existing `--write` flag and
positional tool arg already exist):

- `ctx setup pi` — prints the static `hook.pi` instruction text
  (`writeSetup.InfoPi*` with `desc.Text(text.DescKeyHookPi)`), matching every
  other integration's no-`--write` path. No filesystem-inspecting drift report
  (that machinery does not exist for any integration and is out of scope).
- `ctx setup pi --write` — calls `corePi.Deploy(cmd)`.
- The YAML/DescKey legs (`hooks.yaml`, `write.yaml`,
  `internal/config/embed/text/hook.go`) are enforced by
  `TestDescKeyYAMLLinkage` — they are mandatory, not optional.

## Error cases

- `ctx` binary not found on PATH: extension no-ops (tolerated; matches the
  OpenCode shim's subprocess tolerance).
- Managed-target validation follows **OpenCode semantics verbatim**: symlinks
  and non-regular files at `.pi/extensions/ctx.ts` /
  `.pi/skills/ctx-*/SKILL.md` are refused; a differing regular file at a
  managed path is refreshed in place (managed paths are documented as
  ctx-owned). There is no ownership sentinel — "drifted ctx file" and "foreign
  file" are the same observable state, so the spec deliberately does not claim
  foreign-target refusal.
- Deploy error semantics mirror `opencode.Deploy`: extension deploy failure is
  fatal (returns error); `coreAgents.Deploy` (AGENTS.md) and skills failures
  are `writeErr.WarnFile` warnings that do not halt; then `InfoPiSummary`.
- Project trust: `.pi/extensions` and `.pi/skills` load only after the project
  is trusted; first `pi` launch in a fresh project prompts for trust
  (`defaultProjectTrust` is configurable). Setup itself is filesystem work and
  succeeds regardless; the trust gate gates *execution*.

## Tests

- `deploy_test.go` modeled on `internal/cli/setup/core/opencode/deploy_test.go`:
  fresh deploy, refresh-in-place on drift (seed arbitrary content at the
  managed path, assert overwrite), refuse on symlink/non-regular target,
  deterministic skill ordering.
- Asset linkage: embedded Pi extension + skills present; deploy constants match
  embedded paths (extend the existing assets/embed linkage test pattern).
- `TestDescKeyYAMLLinkage` green (YAML keys ↔ DescKey constants).
- Full validation suite before declaring complete: `make build`, `make lint`,
  `go test ./...`.

## Non-goals

- No MCP registration (Pi has none by design; documented, not a gap).
- No `tool_call` blocking gate (dangerous-command interception) — carried over
  as a *permanent* omission from the OpenCode integration (its `DECISIONS.md`
  entry 2026-04-26-231517: a shim that shelled out to block-dangerous-commands
  would block every shell command on installs without the wrapper).
- No `registerTool` / `registerCommand` in the extension v1 (keeps the shim
  dependency-free, no typebox schemas).
- No global deploy (`~/.pi/agent/...`) in v1 — project-local only.
- No `pi install` package distribution in v1.
- No `.agents/skills/` deploy in v1 (project `.pi/skills/` only).
- No session-file parser for Pi's JSONL session format (context capture stays
  on the `ctx system` nudge path).
- No `resources_discover` skill-path contribution (the deploy-files approach is
  chosen over Pi's other skill-delivery channel).

## Resolved open questions

1. **Injected packet visibility:** `message.display: true` — the user sees the
   ctx packet in the TUI (transparency; decision 2026-08-23).
2. **CI type-check for the extension:** in scope — a `tools/typecheck/pi/`
   job mirroring `tools/typecheck/opencode/` (tsconfig + package.json pulling
   `@earendil-works/pi-coding-agent` types, `tsc --noEmit`), wired into
   `.github/workflows/ci.yml` as a `typecheck-pi-extension` job.

## Verification

1. `make build && make lint && go test ./...` green on the branch.
2. Scratch project: `ctx init`, `ctx setup pi --write`; confirm
   `.pi/extensions/ctx.ts`, `.pi/skills/ctx-*/SKILL.md`, and AGENTS.md
   land with the managed markers.
3. Launch `pi` in the scratch project: accept the project-trust prompt;
   first prompt triggers ctx packet injection; a successful `git commit`
   triggers the post-commit nudge (a *failed* commit does not); `/compact`
   then next prompt re-injects the packet; skills are reachable as
   `/skill:ctx-status` etc.
4. Independent sub-agent review against this spec + TASKS.md before declaring
   complete (playbook requirement).

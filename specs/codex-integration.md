---
title: OpenAI Codex Integration
status: implemented
date: 2026-08-23
owner: parlakisik
scope: integration — assets, setup deployer, hooks, plugin/marketplace, journal parser, steering, docs
related:
  - specs/opencode-integration.md (the write-to-disk blueprint this follows)
  - specs/future-complete/copilot-cli-integration.md
  - specs/future-complete/agents-md.md (Codex reads AGENTS.md natively)
  - specs/steering-sync-drift-respects-configured-tools.md (codex is a non-synced tool)
  - specs/hooks-wiring-guard.md (extended to the Codex hooks manifest)
  - specs/cwd-anchored-context.md (why every hook command `cd`s to the git root)
---

# OpenAI Codex Integration

## Problem

`ctx` already treats `codex` as a first-class tool identifier
(`cfgHook.ToolCodex`, `.ctxrc` `tool:` schema, drift validation,
docs that promise "hook + MCP" parity with Claude Code) — but nothing
backs the promise:

- `ctx setup codex` falls through to `UnsupportedTool`.
- `ctx steering sync` with `tool: codex` in `.ctxrc` errors with
  `unsupported sync tool "codex"` instead of the documented polite skip
  (`.context/journal/2026-08-19-ctx-remember-3564bc24.md:558`).
- No Codex hooks, no Codex plugin/marketplace, no Codex skills, no
  `~/.codex/sessions` journal parser.

Codex CLI 0.148 ships `hooks` and `plugins` as **stable** features with
a lifecycle-hook contract that is a near clone of Claude Code's
(same event names, same stdin fields, same `hookSpecificOutput` /
`decision: block` output shapes). Nothing in the `ctx system` hook
runtime needs to change to serve Codex; what is missing is the
delivery layer (manifests, deployer, parser, docs).

## Approach

Give Codex the **same two delivery routes Claude Code has** plus the
write-to-disk route Copilot CLI / OpenCode have:

1. **Plugin route** — `internal/assets/codex/` is a Codex plugin root
   (mirrors `internal/assets/claude/` being the Claude plugin root):
   `.codex-plugin/plugin.json`, `hooks/hooks.json`, `skills/`,
   `.mcp.json`. A repo marketplace at `.agents/plugins/marketplace.json`
   (mirrors `.claude-plugin/marketplace.json`) points at it, so users
   run `codex plugin marketplace add ActiveMemory/ctx` then
   `codex plugin add ctx@activememory-ctx` and every project gets
   hooks, skills, and the MCP server.
2. **Project-local route** — `ctx setup codex --write` materializes
   the same embedded assets into the project: `.codex/hooks.json`,
   `.codex/config.toml` (`[mcp_servers.ctx]`), `.agents/skills/ctx-*/`,
   and `AGENTS.md` (via the shared `core/agents` deployer). This is
   the route for teams that don't want a user-level plugin, or for
   CI / `codex exec` runs.
3. **Journal route** — a `codex` session parser reads Codex rollout
   transcripts (`$CODEX_HOME/sessions/YYYY/MM/DD/rollout-*.jsonl`)
   so `ctx journal import` and the `SessionEnd` hook capture Codex
   sessions exactly as they capture Claude Code sessions.

The hook commands are identical in both routes (`cd "$(git rev-parse
--show-toplevel)" && ctx …`), so a single embedded `hooks.json` serves
the plugin and the project file. Codex runs hook commands with the
session cwd, and `ctx` is CWD-anchored (`$PWD/.context`), so every
command anchors to the git root exactly as the Claude manifest anchors
to `${CLAUDE_PROJECT_DIR}`.

### Event mapping (Claude Code manifest → Codex manifest)

| Codex event        | matcher                    | ctx command                                   | Why |
|--------------------|----------------------------|-----------------------------------------------|-----|
| `SessionStart`     | *(all sources)*            | `ctx agent --budget 8000`                     | Plain stdout becomes developer context; re-fires on `compact` so context survives compaction. Replaces Claude's PreToolUse `.*` → `ctx agent` (Codex ignores plain text on PreToolUse). |
| `PreToolUse`       | `.*`                       | `ctx system context-load-gate`                | same as Claude |
| `PreToolUse`       | `Bash`                     | `ctx system block-non-path-ctx`               | same; emits legacy `{"decision":"block"}` which Codex accepts |
| `PreToolUse`       | `Bash`                     | `ctx system qa-reminder`                      | same |
| `PreToolUse`       | `update_plan`              | `ctx system specs-nudge`                      | Codex's planning tool is `update_plan` (Claude: `EnterPlanMode`) |
| `PostToolUse`      | `Bash`                     | `ctx system post-commit`                      | same; `tool_input.command` is a string in both |
| `PostToolUse`      | `apply_patch\|Edit\|Write` | `ctx system check-task-completion`            | Codex file edits are `apply_patch`; `Edit`/`Write` are its matcher aliases |
| `UserPromptSubmit` | —                          | the 13 `ctx system check-*` / `heartbeat`     | same list, same order |
| `SessionEnd`       | —                          | `ctx journal import --all -y` (`timeout: 3`)  | Codex caps SessionEnd at 3 s; the import is incremental so this is normally sub-second |

Not mapped: `PermissionRequest`, `PreCompact`, `PostCompact`,
`SubagentStart/Stop`, `Stop` — none carries a ctx behavior today
(`PreCompact`/`PostCompact` cannot return `additionalContext`; `Stop`
requires JSON-only output and ctx has no Stop-shaped nudge).

### Trust model

Codex refuses to run non-managed hooks until the user reviews them in
`/hooks`. Both routes therefore end with the same instruction: start
`codex`, run `/hooks`, trust the `ctx` entries. The deployer and the
setup hint print this; the docs repeat it.

## Behavior

### Happy Path

**Plugin route**

1. `codex plugin marketplace add ActiveMemory/ctx` registers the repo
   marketplace (`.agents/plugins/marketplace.json`, name
   `activememory-ctx`).
2. `codex plugin add ctx@activememory-ctx` installs the plugin from
   `./internal/assets/codex`.
3. User opens `codex`, runs `/hooks`, trusts the ctx hooks.
4. Every new Codex session in a ctx-initialized project receives the
   `ctx agent` packet at start, the `UserPromptSubmit` nudges, the
   tool gates, and a journal import at session end. `$ctx-remember`
   etc. are available as skills; the `ctx` MCP server is registered.

**Project-local route**

1. `ctx setup codex` (no flag) prints the integration overview and the
   detected state (Codex binary? plugin installed? plugin enabled?).
2. `ctx setup codex --write`:
   - writes `.codex/hooks.json` (create, or merge into an existing
     file preserving foreign hook groups and replacing stale
     ctx-managed groups);
   - writes `.codex/config.toml` with `[mcp_servers.ctx]` (create, or
     append the table when the header is absent; skip when present);
   - deploys `AGENTS.md` via `coreAgents.Deploy` (marker merge);
   - writes `.agents/skills/<name>/SKILL.md` for every embedded Codex
     skill (create / refresh-if-stale / reject foreign file), plus the
     skill's `references/` files (deployed only when the SKILL.md at
     that path is ctx-managed);
   - prints a summary and the `/hooks` trust reminder.
3. If the ctx plugin is **enabled** in `~/.codex/config.toml`, the
   deployer skips hooks, MCP, and skills (they would run twice —
   Codex loads every matching hook from every source) and says so;
   only `AGENTS.md` is deployed.

**Journal route**

1. `ctx journal import --all` (or the SessionEnd hook) scans
   `$CODEX_HOME/sessions` (default `~/.codex/sessions`) recursively for
   `rollout-*.jsonl`, matches sessions to the current project by
   `session_meta.cwd` / git origin, and imports them with
   `tool: codex`.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `codex` binary absent | `ctx setup codex` still writes files (`--write`) and hints how to install Codex; no error |
| `.codex/hooks.json` exists with foreign hooks | Merge: foreign matcher groups preserved byte-for-byte in meaning (re-serialized), ctx groups replaced/added; top-level `description` preserved |
| `.codex/hooks.json` is not valid JSON | Refuse to touch it; warn with the path; continue with the other deploy steps |
| `.codex/config.toml` exists without `[mcp_servers.ctx]` | Append a newline-separated `[mcp_servers.ctx]` table at EOF (valid TOML regardless of prior content); never rewrite existing bytes |
| `.codex/config.toml` has `[mcp_servers.ctx]` | Skip (no comparison of body — the user owns it) |
| `.agents/skills/<name>/SKILL.md` exists, identical | Skip |
| `.agents/skills/<name>/SKILL.md` exists, stale ctx content | Refresh in place (frontmatter `name:` matches, body differs) |
| `.agents/skills/<name>/SKILL.md` exists, foreign | Reject with a warning; do not overwrite |
| Plugin enabled in `~/.codex/config.toml` | Skip hooks/MCP/skills with an info line; deploy `AGENTS.md` only; the summary printed is the plugin-mode variant (no project-local paths, no `/hooks` step) |
| User hook group whose commands copy the git-root anchor | Preserved: ctx-managed classification requires the full anchor **plus** a `ctx ` invocation (`HookCommandPrefix`), never the anchor alone |
| Project untrusted in Codex | Codex ignores `.codex/` layers; the setup hint and docs explain `trust_level = "trusted"` |
| `$CODEX_HOME` set | Parser and plugin detection use it instead of `~/.codex` |
| `~/.codex/sessions` absent | Parser contributes no sessions; no error |
| Rollout file with `session_meta` but zero user messages | Session skipped (same rule as the Claude parser) |
| Rollout with injected `<environment_context>` / `<user_instructions>` / `# AGENTS.md instructions for <path>` user items | Filtered out of the message list (not user prose). The AGENTS.md marker is anchored on the full `... instructions for ` form so a user prompt that merely opens with `# AGENTS.md instructions` survives |
| Rollout line over 4 MB (parser buffer) | `ParseFile` errors and the file is skipped by the directory scan; 4 MB matches the schema checker's ceiling and exceeds Codex's observed line sizes |
| `[mcp_servers.ctx]` (or the plugin table header) appearing inside a TOML multi-line string | Accepted limitation of the never-parse-TOML design: detection is a trimmed-line scan, so such a file is treated as already configured (skip). Contrived input; recoverable by adding the table manually |
| `.ctxrc` `tool: codex` + `ctx steering sync` (no `--tool`) | Info line: codex consumes steering via `ctx agent`; exit 0 (same for `claude`) |
| `ctx steering sync --tool codex` (explicit) | Same polite skip, exit 0 |
| SessionEnd import exceeds 3 s | Codex reports a hook failure; the import resumes on the next `UserPromptSubmit` `check-journal` nudge / next session end. Documented. |
| Windows | No `commandWindows` override shipped (parity with the Claude manifest); hooks require a POSIX shell with `git` on PATH. Documented as a known limitation. |

### Validation Rules

- Every command in `internal/assets/codex/hooks/hooks.json` resolves to
  a registered cobra path (`internal/compliance/hooks_wiring_test.go`,
  extended to the Codex manifest).
- Every command in the Codex manifest starts with the git-root anchor
  (`cd "$(git rev-parse --show-toplevel)" &&`).
- Every event key in the Codex manifest is a Codex-supported event
  name.
- `internal/assets/codex/.codex-plugin/plugin.json` `version` ==
  `VERSION` == `.agents/plugins/marketplace.json` version (hack
  sync + test, same as the Claude manifests).
- Every `internal/assets/codex/skills/*/SKILL.md` has frontmatter
  `name` == directory name and a non-empty `description`
  (`frontmatter_test.go` `skillTrees`).
- `internal/assets/codex/skills/` is generated from
  `internal/assets/claude/skills/` by `hack/sync-codex-skills.sh`
  (strip `allowed-tools:`; exclude the Claude-only list);
  `make check-codex-skills` fails when stale.
- Deploy targets are validated to stay inside the project root
  (reuse the OpenCode `validateManagedTarget` shape).

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| `.codex/hooks.json` unparseable | `warning: <path>: <json error> — left untouched` | fix or delete the file, re-run `--write` |
| foreign `SKILL.md` at a ctx skill path | `warning: <path>: not ctx-managed, skipped` | move the file, re-run |
| `AGENTS.md` deploy error | existing `writeErr.WarnFile` | existing behavior |
| rollout line malformed | line skipped; file still imports | none needed |

## Interface

### CLI

| Command | Behavior |
|---------|----------|
| `ctx setup codex` | prints overview + both install routes + state: `configured` when `.codex/hooks.json` exists, else the plugin detection state |
| `ctx setup codex --write` | deploys project-local integration (see Happy Path) |
| `ctx journal import` (all forms) | now also discovers Codex rollouts |
| `ctx steering sync` / `--tool codex` | polite skip for non-synced tools |
| `ctx init` | post-init hint when `codex` is on PATH and neither `.codex/hooks.json` nor the plugin is present |

No new flags.

### Skill

Codex skills are the Claude skills minus `allowed-tools:` and minus
the Claude-only set: `ctx-permission-sanitize` (`.claude/settings`),
`ctx-plan-import` (`~/.claude/plans`), `ctx-dream` (Claude-headless
cron + guard script), `ctx-skill-create` (authors Claude skills).
Codex invokes them as `$ctx-remember` etc.; skill bodies that say
`/ctx-…` remain understandable (same name).

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/config/codex/codex.go`, `doc.go` | constants: binary, `.codex`, `hooks.json`, `config.toml`, `.agents`, `skills`, `plugins`, `marketplace.json`, plugin id, marketplace id, `CODEX_HOME`, `sessions`, hook anchor, TOML headers, event names, `update_plan` |
| `internal/config/asset/asset.go` | `DirCodex*`, `FileHooksJSON`, `FileDotMCPJSON`, `PathCodex*` |
| `internal/config/setup/setup.go` | `DisplayCodex`, `HooksPathCodex`, `MCPConfigPathCodex`, `SkillsPathCodex` |
| `internal/config/session/tool.go` | `ToolCodex` parser id |
| `internal/config/embed/text/hook.go` | `DescKeyHookCodex`, `DescKeyWriteHookCodex*` |
| `internal/assets/embed.go` | embed `codex/.codex-plugin/plugin.json codex/.mcp.json codex/hooks/hooks.json codex/skills/*/SKILL.md` |
| `internal/assets/codex/**` | plugin root (manifest, hooks, mcp, generated skills) |
| `.agents/plugins/marketplace.json` | repo marketplace |
| `internal/assets/read/agent/agent.go` | `CodexHooksJSON()`, `CodexPluginJSON()`, `CodexMCPJSON()`, `CodexSkills()` |
| `internal/cli/setup/core/codex/{codex,hooks,mcp,skill,validate,detect,types,doc}.go` + tests | deployer + detection |
| `internal/cli/setup/cmd/root/run.go` | `case cfgHook.ToolCodex` |
| `internal/cli/initialize/cmd/root/run.go` (+ core) | Codex post-init hint |
| `internal/write/setup/hook.go` | `InfoCodex*` writers |
| `internal/assets/commands/text/hooks.yaml`, `write.yaml`, `ui.yaml`, `commands.yaml` | text keys; supported-tools lists; `ctx setup --help` tool list |
| `internal/journal/parser/codex.go`, `codex_types.go`, `codex_path.go`, tests, `testdata/codex/*.jsonl` | rollout parser |
| `internal/journal/parser/parser.go`, `query.go` | register parser; scan `CodexSessionDirs()` |
| `internal/steering/sync.go` + `internal/cli/steering/cmd/synccmd/run.go` | polite skip for `claude`/`codex` |
| `hack/sync-codex-skills.sh`, `Makefile` (`sync-codex-skills`, `check-codex-skills`, `codex-plugin-install`), `hack/build-all.sh`, `hack/release.sh` | generation + version sync |
| `internal/assets/plugin_test.go`, `internal/assets/codex_test.go`, `internal/assets/read/skill/frontmatter_test.go`, `internal/compliance/hooks_wiring_test.go`, `internal/compliance/ctxctl_isolation_test.go` | guards |
| `docs/home/codex.md` (new), `docs/cli/setup.md`, `docs/cli/journal.md`, `docs/cli/system.md`, `docs/operations/integrations.md`, `docs/recipes/multi-tool-setup.md`, `docs/home/getting-started.md`, `README.md`, `zensical.toml` | docs + nav |
| `.context/TASKS.md`, `.context/DECISIONS.md` | phase + decisions |

### Key Functions

```go
// internal/cli/setup/core/codex
func Deploy(cmd *cobra.Command) error          // orchestrates the 4 steps + plugin short-circuit
func deployHooks(cmd) error                     // read-merge-write .codex/hooks.json
func ensureMCPConfig(cmd) error                 // create-or-append .codex/config.toml
func deploySkills(cmd) error                    // .agents/skills/<name>/SKILL.md
func Detect() State                             // Absent | PluginNotInstalled | PluginInstalledNotEnabled | PluginReady

// internal/journal/parser
func NewCodex() *Codex
func (c *Codex) Matches(path string) bool       // .jsonl + first peeked line is session_meta
func (c *Codex) ParseFile(path string) ([]*entity.Session, error)
func CodexSessionDirs() []string                // $CODEX_HOME/sessions or ~/.codex/sessions
```

### Helpers to Reuse

- `internal/cli/setup/core/agents.Deploy` — `AGENTS.md`
- `internal/cli/setup/core/opencode/{skill,validate}.go` shape — skill deploy + managed-target gate
- `internal/cli/initialize/core/merge/settings.go` raw-JSON round-trip pattern — hooks.json merge
- `internal/cli/initialize/core/claudecheck` shape — detection state machine
- `internal/journal/parser/claude.go`, `parse.go`, `envelope.go` — session assembly
- `internal/cli/system/core/session.FormatContext` — unchanged; Codex consumes the same JSON

## Configuration

- `.ctxrc` `tool: codex` — already in the schema; no new keys.
- `CODEX_HOME` — honored for session discovery and plugin detection.
- Codex side: `[features] hooks = true` (default), project
  `trust_level = "trusted"` for `.codex/` layers.

## Testing

- **Unit**: deployer (create / merge / foreign-reject / plugin
  short-circuit / idempotency), MCP append semantics, detection
  states under a fake `CODEX_HOME`, parser on fixture rollouts
  (messages, tool calls, token totals, filtered injected items,
  zero-user-message skip), `CodexSessionDirs` with/without
  `CODEX_HOME`.
- **Compliance**: Codex manifest wiring guard; no `check-audit`;
  plugin version sync; skill frontmatter; `check-codex-skills`.
- **Live** (manual, recorded in the PR): `ctx setup codex --write` in
  a scratch repo, `codex exec` with `--dangerously-bypass-hook-trust`
  confirming `SessionStart` context injection, `UserPromptSubmit`
  nudges, and a `rollout-*.jsonl` that `ctx journal import` ingests.

## Non-Goals

- **Codex memories bridge** (`~/.codex/memories/`, SQLite-backed,
  off by default, documented as "generated state — don't edit by
  hand"). No stable file contract to mirror the Claude `MEMORY.md`
  bridge onto. Revisit if OpenAI documents the format.
- **Statusline** — Codex has no statusline hook.
- **Steering sync to a Codex-native rules format** — Codex reads
  `AGENTS.md` and gets the packet via `ctx agent`; the deliberate
  exclusion in `specs/steering-sync-drift-respects-configured-tools.md`
  stands.
- **Windows `commandWindows` overrides** — tracked as a follow-up task.
- **Publishing to the universal plugin directory** — the repo
  marketplace is the distribution channel, same as Claude's.

## Resolved Questions

- **SessionStart plain text.** Live test (Codex 0.148, `codex exec`
  with a trusted project): the `ctx agent` plain-text stdout is
  injected verbatim as a developer message ("# Context Packet ...")
  in the session — no JSON wrapping needed. The manifest keeps
  `additionalContextLimit: 10000` as headroom; if a very large packet
  ever spills, Codex's head-and-tail preview plus the saved file path
  is an acceptable degradation.
- **Trust via `-c`.** A `-c 'projects."...".trust_level="trusted"'`
  override does not unlock project `.codex/` layers; the entry must
  be in the real `~/.codex/config.toml`. Documented in
  `docs/home/codex.md`.
- **Unified exec matches `Bash`.** Verified live: `PreToolUse` with
  matcher `Bash` intercepts Codex's code-mode `exec` path, and the
  legacy `{"decision":"block"}` response blocks the command
  (`block-non-path-ctx` produced "Command blocked by PreToolUse
  hook" in the transcript).
- **SessionEnd import fits the 3 s cap.** Verified live: each
  `codex exec` session was imported into `.context/journal/` at
  session end by the hook.

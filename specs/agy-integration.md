# Antigravity (`agy`) Integration — First-Class Peer of Claude Code & Codex

**Branch**: `feat/agy-integration`

## Problem

`ctx` supports Claude Code and OpenAI Codex as full lifecycle peers (plugin + project-local setup + MCP + skills + hooks + journal import), as well as OpenCode, Pi, and GitHub Copilot CLI. Google Antigravity (`agy` CLI and Antigravity IDE/2.0) provides a complete native customization surface:

- **Native Plugins** (`plugin.json` at plugin root, `hooks.json`, `mcp_config.json`, `rules/AGENTS.md`, `skills/ctx-*/SKILL.md`), validated via `agy plugin validate --strict` and installed globally via `agy plugin install <dir>` into `~/.gemini/config/plugins/ctx/`.
- **Lifecycle Hooks** (`SessionStart`, `PreInvocation`, `PreToolUse`, `PostToolUse`, `PostInvocation`, `Stop`) configured in `~/.gemini/config/hooks.json`, `<workspace>/.agents/hooks.json` (or `.agent/hooks.json`), or a plugin's root `hooks.json`.
- **MCP Servers** configured in `~/.gemini/config/mcp_config.json` (or a plugin's root `mcp_config.json`).
- **Agent Rules & Skills** loaded from `AGENTS.md` (or `rules/AGENTS.md` in a plugin) and `.agents/skills/<name>/SKILL.md` (or `skills/<name>/SKILL.md` in a plugin).
- **Session Transcripts** stored as JSONL in `~/.gemini/antigravity-cli/brain/<conversation-id>/.system_generated/logs/transcript.jsonl` with workspace mapping in `~/.gemini/antigravity-cli/history.jsonl`.

Today, `ctx` has no `agy` / `antigravity` integration: `ctx setup agy` is unknown, `ctx system` hooks only speak the Claude Code / Codex snake_case `stdin` and plain-text/hookSpecificOutput `stdout` protocol, and `ctx journal import` cannot parse Antigravity session transcripts.

---

## Scope

### Phase 1 — Native Plugin, Unified Go Hook Protocol, `ctx setup agy`, Skills & Guards

1. **Unified Antigravity Hook Protocol in `ctx system`**:
   - Auto-detect Antigravity hook invocations when `stdin` JSON carries `conversationId` (or `CTX_HOOK_TOOL=agy`).
   - Normalize Antigravity `stdin` fields into `entity.HookInput`:
     - `conversationId` → `SessionID`
     - `workspacePaths[0]` → `CWD` (and `os.Chdir(workspacePaths[0])` + `rc.Reset()` inside `check.ReadInput` **before** `state.Initialized()` / `rc.ContextDir()`, because Antigravity runs hook commands with `CWD` set to the directory containing `hooks.json`).
     - `transcriptPath` → `TranscriptPath`
     - `modelName` → `ModelName`
     - `lastUserInput` → `Prompt`
     - `toolCall.name` / `toolCall.args` → `ToolName` / `ToolInput` (`CommandLine` → `Command`, `Path` → `FilePath`)
     - `error` → `ToolError` (non-empty when a tool call failed; suppresses `post-commit` and `check-task-completion`).
   - Emit Antigravity-compliant JSON on `stdout` for every hook execution in Antigravity mode:
     - `SessionStart`, `PreInvocation`, `PostInvocation`: `{"injectSteps": [{"systemMessage": {"systemMessage": "<content>"}}]}` when context/nudges are present, or `{}` when silent.
     - `PreToolUse`: `{"decision": "deny", "reason": "<reason>"}` when blocking (`block-non-path-ctx`), or `{"decision": "allow"}` when allowing.
     - `PostToolUse` and `Stop`: `{}` on `stdout`.
     - **Pending Nudge Queue Bridge**: Because Antigravity `PreToolUse` and `PostToolUse` only accept `{"decision": ...}` or `{}` (not `injectSteps`), any non-blocking context/nudge emitted during `PreToolUse` (`qa-reminder`, `specs-nudge`) or `PostToolUse` (`post-commit`, `check-task-completion`) in Antigravity mode is appended to `.context/state/agy-pending-<sessionID>.txt` while returning `{"decision": "allow"}` / `{}` on `stdout`, and drained into `injectSteps` by `PostInvocation` (`heartbeat`) and `PreInvocation` (`check-context-size`).
     - **SessionStart Context Packet**: On `SessionStart` in Antigravity mode, `context-load-gate` prepends the assembled `ctx agent` packet to the gate directive so a single JSON stdout response delivers both the context packet and the read gate.
     - **Plain-Text Suppression (`check-reminder` / `session-event`)**: In Antigravity mode (`input.Agy()`), `coreProv.Emit` in `check-reminder` and `wSession.Event` in `session-event` buffer/suppress raw plain-text stdout so every hook response remains a single valid JSON object.
2. **Embedded Native Antigravity Plugin (`internal/assets/agy/`)**:
   - `internal/assets/agy/plugin.json`: valid Antigravity plugin manifest (`name: "ctx"`, `displayName: "ctx"`, `version`, `description`, `logo: "assets/logo.svg"`, `author`, `homepage`, `repository`, `suggestedPrompts`).
   - `internal/assets/agy/hooks.json`: named hook group `"ctx"` wiring `SessionStart`, `PreInvocation`, `PreToolUse` (`run_command`, `edit_file|write_file`), `PostToolUse` (`run_command`, `edit_file|write_file`), `PostInvocation` (`ctx system heartbeat`), and `Stop` (`ctx system session-event --type end --caller agy`).
   - `internal/assets/agy/mcp_config.json`: `mcpServers.ctx` (`ctx mcp serve`).
   - `internal/assets/agy/rules/AGENTS.md`: universal `ctx` agent instructions inside `rules/` (required by `agy plugin validate`).
   - `internal/assets/agy/assets/logo.svg`: minimal SVG icon so `agy plugin validate --strict` passes with zero warnings and zero notes.
   - `internal/assets/agy/skills/ctx-*/SKILL.md` (+ `references/`): synced from `internal/assets/claude/skills/` via `hack/sync-agy-skills.sh` (`allowed-tools:` stripped, Claude-only skills excluded, and any skill whose Markdown body exceeds Antigravity's 24,000-byte limit deterministically splits overflow into `references/extended.md` so `agy plugin validate --strict` passes with 0 warnings).
3. **Project-Local Setup (`ctx setup agy [--write]` / `ctx setup antigravity [--write]`)**:
   - Detection in `internal/agy/` (`Home`, `ConfigDir`, `Detect`, `Unwired`, `MergeHooks`, `EnsureMCPConfig`).
   - If the `ctx` Antigravity plugin is already installed and enabled in `~/.gemini/config/plugins/ctx/`, `ctx setup agy --write` short-circuits hooks/MCP/skills and deploys only `AGENTS.md` (warning if project-local `.agents/hooks.json` also exists).
   - Otherwise, `ctx setup agy --write` deploys:
     - `.agents/hooks.json`: merge-safe update of the `"ctx"` named hook group, preserving any user-defined sibling hook groups.
     - `~/.gemini/config/mcp_config.json` (or `$AGY_CONFIG_DIR/mcp_config.json`): merge-safe registration of `mcpServers.ctx`.
     - `AGENTS.md`: marker-managed `<!-- ctx:agents -->` block via `coreAgents.Deploy`.
     - `.agents/skills/ctx-*/SKILL.md` (+ `references/`): synced skill tree.
   - `ctx init` post-script hint when `agy` is on `PATH` and unwired.
   - `internal/steering/sync.go` registers `agy` and `antigravity` in `directConsumers` so `ctx steering sync` skips politely.
4. **Compliance Tests & Documentation**:
   - `internal/assets/agy_test.go`, `internal/assets/plugin_test.go`, `internal/assets/read/skill/{frontmatter_test,parity_test}.go`, `internal/compliance/hooks_wiring_test.go`.
   - `Makefile` targets: `sync-agy-skills`, `check-agy-skills`, and `sync-version` updates `internal/assets/agy/plugin.json`.
   - `docs/cli/setup.md` and `docs/operations/integrations.md`.

### Phase 2 — Antigravity Session Transcript Parser (`ctx journal import`)

1. **Transcript Parser (`internal/journal/parser/agy.go`)**:
   - Implements `parser.Session` for Antigravity session logs at `~/.gemini/antigravity-cli/brain/<conversation-id>/.system_generated/logs/transcript.jsonl` (overridable via `$AGY_HOME`).
   - Parses JSONL steps (`step_index`, `source`, `type`, `status`, `created_at`, `content`, `thinking`, `tool_calls`):
     - `USER_INPUT` → user `entity.Message`
     - `PLANNER_RESPONSE` → assistant `entity.Message` (`Text` from `content`, `Thinking` from `thinking`, `ToolUses` from `tool_calls`)
   - Resolves session `CWD` (workspace path) and title from `~/.gemini/antigravity-cli/history.jsonl` (`id`, `title`, `workspaces` URI `file:///...`).
   - Skips `transcript_full.jsonl` in directory walks (to avoid duplicate sessions alongside `transcript.jsonl`) and skips malformed/partial trailing JSONL lines (live-transcript tolerance).
2. **Discovery & Registration**:
   - Registers `NewAgy()` in `registeredParsers` (`internal/journal/parser/parser.go`) and scans `AgySessionDirs()` in `findSessionsWithFilter` (`internal/journal/parser/query.go`).
   - Fixture-backed unit tests in `internal/journal/parser/agy_test.go`.

### Non-Goals

- Modifying Antigravity's built-in skills (`~/.gemini/antigravity-cli/builtin/`).
- Bridging Antigravity's internal `.pbtxt` protobuf state files (all needed metadata and messages are in `history.jsonl` and `transcript.jsonl`).

---

## Architecture & Hook Mapping

### Hook Lifecycle Mapping

| Antigravity Event | Matcher | `ctx` Command(s) | Stdin Used | Stdout Emitted (Agy Mode) |
|---|---|---|---|---|
| `SessionStart` | *(flat array)* | `ctx system context-load-gate` | `conversationId`, `workspacePaths` | `{"injectSteps": [{"systemMessage": {"systemMessage": "<agent packet + gate>"}}]}` or `{}` |
| `PreInvocation` | *(flat array)* | `check-context-size`, `check-persistence`, `check-ceremony`, `check-journal`, `check-memory-drift`, `check-knowledge`, `check-map-staleness`, `check-version`, `check-freshness`, `check-skill-discovery`, `check-hub-sync`, `check-reminder`, `check-resource` | `conversationId`, `workspacePaths`, `transcriptPath`, `modelName`, `lastUserInput` | `{"injectSteps": [...]}` (including any drained pending nudges) or `{}` |
| `PreToolUse` | `run_command` | `block-non-path-ctx`, `qa-reminder` | `conversationId`, `workspacePaths`, `toolCall` | `{"decision": "deny", "reason": "..."}` or `{"decision": "allow"}` (queues nudge) |
| `PreToolUse` | `edit_file\|write_file` | `specs-nudge` | `conversationId`, `workspacePaths`, `toolCall` | `{"decision": "allow"}` (queues nudge if triggered) |
| `PostToolUse` | `run_command` | `post-commit` | `conversationId`, `workspacePaths`, `toolCall`, `error` | `{}` (queues commit audit/nudge if `error == ""` and command was `git commit`) |
| `PostToolUse` | `edit_file\|write_file` | `check-task-completion` | `conversationId`, `workspacePaths`, `toolCall`, `error` | `{}` (queues task nudge if `error == ""` and threshold reached) |
| `PostInvocation` | *(flat array)* | `ctx system heartbeat` | `conversationId`, `workspacePaths`, `transcriptPath`, `modelName` | `{"injectSteps": [...]}` (drains pending queue from `PreToolUse`/`PostToolUse`) or `{}` |
| `Stop` | *(flat array)* | `ctx system session-event --type end --caller agy` | `conversationId`, `workspacePaths` | `{}` |

### CWD Resolution Order in `check.ReadInput` / `check.FullPreamble`

Antigravity spawns hook commands with `CWD` set to the directory containing `hooks.json` (`~/.gemini/config/plugins/ctx/` or `<workspace>/.agents/`).
Therefore:
1. `check.FullPreamble(stdin)` (and `contextloadgate.Run`, `checkcontextsize.Run`, `checkreminder.Run`, `heartbeat.Run`, `sessionevent.Run`) reads `stdin` first via `check.ReadInput(stdin)`.
2. When `input.Agy()` is true (`input.ConversationID != ""` or `os.Getenv("CTX_HOOK_TOOL") == "agy"`), `check.ReadInput`:
   - Copies `ConversationID` to `SessionID` (if `SessionID == ""`).
   - Resolves `CWD` from `WorkspacePaths[0]` (if `len(WorkspacePaths) > 0`), calls `os.Chdir(input.CWD)` and `rc.Reset()`, so subsequent calls to `state.Initialized()` and `rc.ContextDir()` resolve against the target workspace's `.context/`.
   - Registers Antigravity output mode in `session` so `writeSetup.Context`, `writeSetup.Block`, and `writeSetup.Nudge` (plus a silent exit flush) emit valid Antigravity JSON.
   - Maps `ToolCall.Name` → `ToolName` and `ToolCall.Args.CommandLine` → `ToolInput.Command`, `ToolCall.Args.Path` → `ToolInput.FilePath`.
3. `state.Initialized()` and `state.Paused()` then run against the workspace's `.context/` directory. If `FullPreamble` bails early in Antigravity mode (uninitialized or paused), it still emits `{}` (or `{"decision": "allow"}` for `PreToolUse`) so Antigravity never sees an empty stdout parse error.

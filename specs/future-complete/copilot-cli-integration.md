# Spec: Copilot CLI Integration — Feature Matrix

## Feature Matrix: Claude Code vs VS Code Extension vs GitHub Copilot CLI

### Legend

- **✅** — Implemented and shipping
- **🔧** — Partially implemented / needs work
- **📋** — Planned / specced
- **—** — Not applicable to this surface

---

### 1. Context Injection (How the agent learns about ctx)

| Feature | Claude Code | VS Code Extension | Copilot CLI |
|---------|-------------|-------------------|-------------|
| Project instructions file | ✅ `CLAUDE.md` | ✅ `.github/copilot-instructions.md` | 📋 `.github/copilot-instructions.md` + `AGENTS.md` |
| Auto-generated on `ctx init` | ✅ Merged into project root | ✅ Via `@ctx /init` (also runs `hook copilot --write`) | 📋 `ctx init` should also generate `AGENTS.md` |
| Marker-based idempotency | ✅ `<!-- ctx:context -->` / `<!-- ctx:end -->` | ✅ `<!-- ctx:copilot -->` / `<!-- ctx:copilot:end -->` | 📋 Same copilot markers |
| Path-specific instructions | — | — | 📋 `.github/instructions/*.instructions.md` |
| Custom agents | — | — | 📋 `.github/agents/ctx.md` |
| Home-dir instructions | — | — | 📋 `~/.copilot/copilot-instructions.md` |
| Reads `CLAUDE.md` natively | ✅ Core feature | — | ✅ Built-in (Copilot CLI reads CLAUDE.md) |
| Reads `AGENTS.md` natively | — | — | ✅ Built-in (primary instructions) |

---

### 2. Hook System (Pre/post tool execution, session lifecycle)

| Feature | Claude Code | VS Code Extension | Copilot CLI |
|---------|-------------|-------------------|-------------|
| Config location | `.claude/settings.local.json` | — (extension handles internally) | 📋 `.github/hooks/ctx-hooks.json` |
| PreToolUse / preToolUse | ✅ Regex matcher + command | — | 📋 `bash` + `powershell` fields |
| PostToolUse / postToolUse | ✅ Command hook | — | 📋 `bash` + `powershell` fields |
| UserPromptSubmit / userPromptSubmitted | ✅ Command hook | — | 📋 `bash` + `powershell` fields |
| SessionEnd / sessionEnd | ✅ Command hook | — | 📋 `bash` + `powershell` fields |
| SessionStart / sessionStart | — | — | 📋 `bash` + `powershell` fields |
| agentStop | — | — | 📋 Available in Copilot CLI |
| subagentStop | — | — | 📋 Available in Copilot CLI |
| errorOccurred | — | — | 📋 Available in Copilot CLI |
| Hook script format | Bash only | N/A | 📋 Dual: bash + PowerShell |
| Block dangerous commands | ✅ `block-hack-scripts.sh` | ✅ Detection ring (deny patterns) | 📋 `ctx-block-commands.sh` + `.ps1` |
| Platform support | Linux/macOS (bash) | All (TypeScript) | 📋 All (bash + powershell) |
| Timeout control | — (Claude manages) | — | 📋 `timeoutSec` per hook |
| Working directory | Implicit (project root) | Implicit | 📋 `cwd` field per hook |
| Environment variables | — | — | 📋 `env` field per hook |

---

### 3. MCP Server (Model Context Protocol)

| Feature | Claude Code | VS Code Extension | Copilot CLI |
|---------|-------------|-------------------|-------------|
| MCP server registration | ✅ Plugin system (`ctx@activememory-ctx`) | ✅ `.vscode/mcp.json` auto-generated | 📋 `~/.copilot/mcp-config.json` |
| Transport | Plugin (in-process) | stdio (`ctx mcp serve`) | 📋 stdio (`ctx mcp serve`) |
| `ctx_status` tool | ✅ | ✅ | 📋 (same server) |
| `ctx_add` tool | ✅ | ✅ | 📋 (same server) |
| `ctx_complete` tool | ✅ | ✅ | 📋 (same server) |
| `ctx_drift` tool | ✅ | ✅ | 📋 (same server) |
| `ctx_recall` tool | ✅ | ✅ | 📋 (same server) |
| `ctx_watch_update` tool | ✅ | ✅ | 📋 (same server) |
| `ctx_compact` tool | ✅ | ✅ | 📋 (same server) |
| `ctx_next` tool | ✅ | ✅ | 📋 (same server) |
| `ctx_check_task_completion` | ✅ | ✅ | 📋 (same server) |
| `ctx_session_event` tool | ✅ | ✅ | 📋 (same server) |
| `ctx_remind` tool | ✅ | ✅ | 📋 (same server) |
| 8 context resources | ✅ | ✅ | 📋 (same server) |
| Resource change notifications | ✅ Poller-based | ✅ Poller-based | 📋 (same server) |
| Prompt templates | ✅ | ✅ | 📋 (same server) |
| Session governance tracking | ✅ | ✅ | 📋 (same server) |

---

### 4. Session Recall (Parsing AI session history)

| Feature | Claude Code | VS Code Extension | Copilot CLI |
|---------|-------------|-------------------|-------------|
| Session parser | ✅ ClaudeCodeParser (JSONL) | ✅ CopilotParser (JSONL) | 📋 CopilotCLIParser (TBD format) |
| Auto-detect session dir | ✅ `~/.claude/projects/` | ✅ Platform-specific `workspaceStorage/` | 📋 `~/.copilot/sessions/` (TBD) |
| Windows path handling | ✅ Drive letter fix | ✅ APPDATA detection | 📋 USERPROFILE / COPILOT_HOME |
| macOS path handling | ✅ ~/Library/... | ✅ ~/Library/Application Support/Code/... | 📋 ~/.copilot/ |
| Linux path handling | ✅ ~/.claude/ | ✅ ~/.config/Code/... | 📋 ~/.copilot/ |
| WSL path handling | — | — | 📋 Must handle WSL ↔ Windows boundary |
| Markdown session export | ✅ | ✅ | 📋 |

---

### 5. Governance & Safety

| Feature | Claude Code | VS Code Extension | Copilot CLI |
|---------|-------------|-------------------|-------------|
| Permission allow-list | ✅ `permissions.allow[]` | — | 📋 `--allow-tool` / `--deny-tool` flags |
| Permission deny-list | ✅ `permissions.deny[]` | — | 📋 `--deny-tool` flags |
| Dangerous command blocking | ✅ PreToolUse hook | ✅ Detection ring (regex) | 📋 preToolUse hook script |
| Sensitive file detection | — | ✅ SENSITIVE_FILE_PATTERNS | 📋 preToolUse hook script |
| Violation recording | — | ✅ `.context/state/violations.json` | 📋 Hook script writes violations |
| Hack script interception | ✅ `block-hack-scripts.sh` | ✅ DENY_COMMAND_SCRIPT_PATTERNS | 📋 preToolUse hook script |
| Tool approval model | Per-session allow | N/A (Copilot manages) | Per-session or `--allow-tool` |
| Trusted directories | Implicit (project root) | Implicit (workspace) | ✅ Built-in trust prompt |

---

### 6. Binary Management

| Feature | Claude Code | VS Code Extension | Copilot CLI |
|---------|-------------|-------------------|-------------|
| Auto-install ctx binary | ✅ Plugin installation | ✅ GitHub releases download | 📋 ctx already on PATH or manual |
| Platform detection | — (Go binary) | ✅ darwin/windows/linux + amd64/arm64 | 📋 Same Go binary |
| Binary verification | — | ✅ Executes `--version` check | — |
| Update mechanism | Plugin update | ✅ GitHub releases (latest) | — (user manages) |

---

### 7. UI & User Experience

| Feature | Claude Code | VS Code Extension | Copilot CLI |
|---------|-------------|-------------------|-------------|
| Chat participant | — (terminal-based) | ✅ `@ctx` with 34 slash commands | — (terminal-based) |
| Status bar | — | 🔧 Reminder status bar (PR pending) | — |
| Diagnostics command | — | ✅ `/diag` with timing | — |
| Progress indicators | — | ✅ `stream.progress()` | — |
| Markdown rendering | Terminal output | ✅ VS Code markdown | Terminal output |
| Interactive mode | ✅ Terminal REPL | ✅ Chat panel | ✅ Terminal REPL |
| Programmatic mode | — | — | ✅ `copilot -p "prompt"` |
| Plan mode | — | — | ✅ Shift+Tab |
| Custom agents | — | — | ✅ `/agent` + `--agent=` flag |
| Skills | ✅ `.claude/skills/` | — | ✅ `.github/skills/` |
| Autopilot mode | — | — | ✅ `--experimental` |

---

### 8. Cross-Platform Support

| Feature | Claude Code | VS Code Extension | Copilot CLI |
|---------|-------------|-------------------|-------------|
| Windows (native) | ✅ | ✅ | ✅ (PowerShell v6+) |
| macOS | ✅ | ✅ | ✅ |
| Linux | ✅ | ✅ | ✅ |
| WSL | — | ✅ (Remote WSL) | ✅ (bash) |
| Hook script: bash | ✅ | N/A | 📋 Required |
| Hook script: PowerShell | — | N/A | 📋 Required |
| Path separator handling | ✅ filepath.Join | ✅ path.join + filepath | 📋 filepath.Join (Go binary) |
| Case-insensitive paths | ✅ (validation pkg) | ✅ (VS Code handles) | 📋 Inherit from ctx binary |
| Home dir detection | `~/.claude/` | Extension globalStorage | 📋 `~/.copilot/` or `$COPILOT_HOME` |

---

### 9. Context System (Shared across all surfaces)

| Feature | Claude Code | VS Code Extension | Copilot CLI |
|---------|-------------|-------------------|-------------|
| `ctx init` | ✅ CLI | ✅ `@ctx /init` | 📋 CLI (same binary) |
| `ctx status` | ✅ CLI | ✅ `@ctx /status` | 📋 CLI (same binary) |
| `ctx agent` | ✅ CLI | ✅ `@ctx /agent` | 📋 CLI (same binary) |
| `ctx drift` | ✅ CLI | ✅ `@ctx /drift` | 📋 CLI (same binary) |
| `ctx recall` | ✅ CLI | ✅ `@ctx /recall` | 📋 CLI (same binary) |
| `ctx add` | ✅ CLI | ✅ `@ctx /add` | 📋 CLI (same binary) |
| `ctx compact` | ✅ CLI | ✅ `@ctx /compact` | 📋 CLI (same binary) |
| `ctx setup <tool>` | ✅ `ctx setup claude` | ✅ `ctx setup copilot` | 📋 `ctx setup copilot-cli` |
| Session persistence | ✅ `.context/sessions/` | ✅ `.context/sessions/` | 📋 `.context/sessions/` |

---

## Implementation Plan: Copilot CLI Integration

### Phase 1 — Hook Generation (cross-platform)

**Goal:** `ctx setup copilot-cli --write` generates:

1. `.github/hooks/ctx-hooks.json` — Hook configuration with dual bash/powershell
2. `.github/hooks/scripts/ctx-preToolUse.sh` — Bash pre-tool gate
3. `.github/hooks/scripts/ctx-preToolUse.ps1` — PowerShell pre-tool gate
4. `.github/hooks/scripts/ctx-sessionStart.sh` — Bash session init
5. `.github/hooks/scripts/ctx-sessionStart.ps1` — PowerShell session init
6. `.github/hooks/scripts/ctx-postToolUse.sh` — Bash post-tool audit
7. `.github/hooks/scripts/ctx-postToolUse.ps1` — PowerShell post-tool audit
8. `.github/hooks/scripts/ctx-sessionEnd.sh` — Bash session teardown
9. `.github/hooks/scripts/ctx-sessionEnd.ps1` — PowerShell session teardown

**Hook JSON structure:**
```json
{
  "version": 1,
  "hooks": {
    "sessionStart": [{
      "type": "command",
      "bash": ".github/hooks/scripts/ctx-sessionStart.sh",
      "powershell": ".github/hooks/scripts/ctx-sessionStart.ps1",
      "cwd": ".",
      "timeoutSec": 10
    }],
    "preToolUse": [{
      "type": "command",
      "bash": ".github/hooks/scripts/ctx-preToolUse.sh",
      "powershell": ".github/hooks/scripts/ctx-preToolUse.ps1",
      "cwd": ".",
      "timeoutSec": 5
    }],
    "postToolUse": [{
      "type": "command",
      "bash": ".github/hooks/scripts/ctx-postToolUse.sh",
      "powershell": ".github/hooks/scripts/ctx-postToolUse.ps1",
      "cwd": ".",
      "timeoutSec": 5
    }],
    "sessionEnd": [{
      "type": "command",
      "bash": ".github/hooks/scripts/ctx-sessionEnd.sh",
      "powershell": ".github/hooks/scripts/ctx-sessionEnd.ps1",
      "cwd": ".",
      "timeoutSec": 15
    }]
  }
}
```

**Script behavior:** All scripts are thin shims that call the ctx binary:
- `ctx-sessionStart` → `ctx system session-event --type start --caller copilot-cli`
- `ctx-preToolUse` → reads JSON stdin, calls `ctx` for dangerous command check
- `ctx-postToolUse` → reads JSON stdin, appends to audit log
- `ctx-sessionEnd` → `ctx system session-event --type end --caller copilot-cli`

### Phase 2 — Agent & Instructions

1. `AGENTS.md` generation in project root (read by Copilot CLI as primary instructions)
2. `.github/agents/ctx.md` — custom agent for context management delegation
3. `.github/instructions/context.instructions.md` — path-specific instructions for `.context/` files

### Phase 3 — MCP & Recall

1. Register ctx MCP server in `~/.copilot/mcp-config.json` (respects `$COPILOT_HOME`)
2. Copilot CLI session parser for `ctx recall`
3. Cross-session continuity (Copilot CLI `--resume` ↔ ctx session files)

### Phase 4 — Deep Integration

1. ACP (Agent Client Protocol) server mode — Copilot CLI can use ctx as an ACP agent
2. Copilot Memory ↔ ctx memory bridge bidirectional sync
3. Skills in `.github/skills/` that wrap ctx operations

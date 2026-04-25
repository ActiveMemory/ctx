```text
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
```

## `ctx`: VS Code Chat Extension

A VS Code Chat Participant that brings [ctx](https://ctx.ist) (persistent
project context for AI coding sessions) directly into GitHub Copilot Chat.

Type `@ctx` in the Chat view to access 45 slash commands, automatic context
hooks, a reminder status bar, and natural language routing, all powered by
the ctx CLI.

## Quick Start

1. Install the extension (or build from source; see [Development](#development))
2. Open a project in VS Code
3. Open Copilot Chat and type `@ctx /init`

The extension auto-downloads the ctx CLI binary if it isn't on your PATH.

## Slash Commands

### Core Context

| Command | Description |
|---------|-------------|
| `/init` | Initialize a `.context/` directory with template files |
| `/status` | Show context summary with token estimate |
| `/agent` | Print AI-ready context packet |
| `/drift` | Detect stale or invalid context |
| `/recall` | Browse and search AI session history |
| `/hook` | Generate AI tool integration configs (copilot, claude) |
| `/add` | Add a task, decision, learning, or convention |
| `/load` | Output assembled context Markdown |
| `/compact` | Archive completed tasks and clean up context |
| `/sync` | Reconcile context with codebase |

### Tasks & Reminders

| Command | Description |
|---------|-------------|
| `/complete` | Mark a task as completed |
| `/remind` | Manage session-scoped reminders (add, list, dismiss) |
| `/tasks` | Archive or snapshot tasks |
| `/next` | Show the next open task from TASKS.md |
| `/implement` | Show the implementation plan with progress |

### Session Lifecycle

| Command | Description |
|---------|-------------|
| `/wrapup` | End-of-session wrap-up with status, drift, and journal audit |
| `/remember` | Recall recent AI sessions for this project |
| `/reflect` | Surface items worth persisting as decisions or learnings |
| `/pause` | Save session state for later |
| `/resume` | Restore a paused session |

### Discovery & Planning

| Command | Description |
|---------|-------------|
| `/brainstorm` | Browse and develop ideas from `ideas/` |
| `/spec` | List or scaffold feature specs from templates |
| `/verify` | Run verification checks (doctor + drift) |
| `/map` | Show dependency map (go.mod, package.json) |
| `/prompt` | Browse and view prompt templates |
| `/blog` | Draft a blog post from recent context |
| `/changelog` | Show recent commits for changelog |

### Maintenance & Audit

| Command | Description |
|---------|-------------|
| `/check-links` | Audit local links in context files |
| `/journal` | View or export journal entries |
| `/consolidate` | Find duplicate entries across context files |
| `/audit` | Alignment audit: drift + convention check |
| `/worktree` | Git worktree management (list, add) |

### Context Metadata

| Command | Description |
|---------|-------------|
| `/memory` | Claude Code memory bridge (sync, status, diff, import, publish) |
| `/decisions` | List or reindex project decisions |
| `/learnings` | List or reindex project learnings |
| `/config` | Manage config profiles (switch, status, schema) |
| `/permissions` | Backup or restore Claude settings |
| `/changes` | Show what changed since last session |
| `/deps` | Show package dependency graph |
| `/guide` | Quick-reference cheat sheet for ctx |
| `/reindex` | Regenerate indices for DECISIONS.md and LEARNINGS.md |
| `/why` | Read the philosophy behind ctx |

### System & Diagnostics

| Command | Description |
|---------|-------------|
| `/system` | System diagnostics and bootstrap |
| `/pad` | Encrypted scratchpad for sensitive notes |
| `/notify` | Send webhook notifications |

Sub-routes for `/system`: `resources`, `doctor`, `bootstrap`, `stats`,
`backup`, `message`.

## Automatic Hooks

The extension registers several VS Code event handlers that mirror
Claude Code's hook system. These run in the background; no user action
needed.

| Trigger | What Happens |
|---------|--------------|
| **File save** | Runs task-completion check on non-`.context/` files |
| **Git commit** | Notification prompting to add a Decision, Learning, run Verify, or Skip |
| **`.context/` file change** | Refreshes reminders and regenerates `.github/copilot-instructions.md` |
| **Dependency file change** | Notification when `go.mod`, `package.json`, etc. change; offers `/map` |
| **Every 5 minutes** | Updates reminder status bar and writes heartbeat timestamp |
| **Extension activate** | Fires `session-event --type start` to ctx CLI |
| **Extension deactivate** | Fires `session-event --type end` to ctx CLI |

## Status Bar

A `$(bell) ctx` indicator appears in the status bar when you have pending
reminders. It updates every 5 minutes. When no reminders are due, it hides
automatically.

## Natural Language

You can also type plain English after `@ctx`: the extension routes
common phrases to the correct handler:

- "What should I work on next?" → `/next`
- "Time to wrap up" → `/wrapup`
- "Show me the status" → `/status`
- "Add a decision" → `/add`
- "Check for drift" → `/drift`

## Auto-Bootstrap

If the ctx CLI isn't found on PATH or at the configured path, the
extension automatically downloads the correct platform binary from
[GitHub Releases](https://github.com/ActiveMemory/ctx/releases):

1. Detects OS and architecture (darwin/linux/windows, amd64/arm64)
2. Fetches the latest release from the GitHub API
3. Downloads and verifies the matching binary
4. Caches it in VS Code's global storage directory

Subsequent sessions reuse the cached binary. To force a specific version,
set `ctx.executablePath` in your settings.

## Follow-Up Suggestions

After each command, Copilot Chat shows context-aware follow-up buttons.
For example:

- After `/init` → "Show status" or "Generate copilot integration"
- After `/drift` → "Sync context" or "Show status"
- After `/reflect` → "Add decision", "Add learning", or "Wrap up"
- After `/spec` → "Show implementation plan" or "Run verification"

## Prerequisites

- VS Code 1.93+
- [GitHub Copilot Chat](https://marketplace.visualstudio.com/items?itemName=GitHub.copilot-chat) extension
- [ctx](https://ctx.ist) CLI on PATH, or let the extension auto-download it

## Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `ctx.executablePath` | `ctx` | Path to the ctx CLI binary. Set this if ctx isn't on PATH and you don't want auto-download. |

## Development

```bash
cd editors/vscode
npm install
npm run watch   # Watch mode
npm run build   # Production build
npm test        # Run tests (53 test cases via vitest)
```

### Architecture

The extension is a single-file implementation
(`src/extension.ts`, ~3 000 lines) that:

- Registers a `ChatParticipant` with `@ctx` as the handle
- Routes slash commands to dedicated `handleXxx()` functions
- Each handler calls the ctx CLI via `execFile` and streams the output
- On Windows, uses `shell: true` so PATH resolution works without `.exe`
- Merges stdout/stderr with deduplication (Cobra prints errors to both)
- A `handleFreeform()` function maps natural language to handlers

### Testing

Tests live in `src/extension.test.ts` and use vitest with a VS Code API
mock. They verify:

- All 45 command handlers exist and are callable
- `runCtx` invokes the correct binary with correct arguments
- Platform detection returns valid GOOS/GOARCH values
- Follow-up suggestions are returned after commands
- Edge cases: missing workspace, cancellation, empty output

## License

Apache-2.0

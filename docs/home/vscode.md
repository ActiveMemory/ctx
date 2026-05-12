---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "ctx for VS Code"
icon: lucide/code-2
---

![ctx](../images/ctx-banner.png)

## The Problem

Every Copilot Chat session in VS Code starts from zero. You re-explain
what you were doing, the AI repeats yesterday's mistakes, and decisions
you spent an hour reasoning through last week get rediscovered instead
of remembered.

**Without `ctx`:**

```
@workspace add the validation middleware we discussed

I don't have context about previous discussions. Could you describe
what validation middleware you're referring to?
```

**With `ctx`:**

```
@ctx Do you remember?

Last session (2026-05-09): you decided on Zod schemas at the route level
(DECISIONS.md #12). Pattern lives in CONVENTIONS.md. Open task: wire
the auth middleware into the new /admin routes (TASKS.md, in-progress).
The reference implementation is src/middleware/auth.ts.
```

That's the whole pitch: **your AI remembers**, right inside the IDE you
already work in.

## Setup

Install the extension and the `ctx` binary, then `ctx init` your project:

1. **Install the extension** from the
   [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=activememory.ctx-context)
   (publisher: `activememory`, display name: *`ctx` — Persistent Context
   for AI*). Or build from source (see
   [editors/vscode/README.md](https://github.com/ActiveMemory/ctx/blob/main/editors/vscode/README.md#development)).
2. **Install the `ctx` CLI** if you haven't already
   ([installation docs](getting-started.md#installation)). If you skip
   this step, the extension will auto-download the right binary for
   your platform on first use (see [Auto-Bootstrap](#auto-bootstrap)
   below).
3. **From your project root**, run:

   ```bash
   ctx init && eval "$(ctx activate)"
   ```

4. **Open Copilot Chat** in VS Code and type `@ctx /init` to verify
   the extension can reach the CLI.

### What Gets Created

| File | Purpose |
|------|---------|
| `.context/` | Project-local context directory (created by `ctx init`) |
| `.github/copilot-instructions.md` | Repository instructions Copilot reads natively; regenerated automatically whenever `.context/` files change |

The extension itself lives in VS Code's extension storage. No project
files are added beyond `.context/` and the Copilot instructions.

## How You Use It

Type `@ctx` in the Copilot Chat view to invoke the chat participant.
Then either:

- **Use a slash command:** `@ctx /status`, `@ctx /wrapup`, etc. There
  are 45 commands; the most common ones live in the [Slash Commands](#slash-commands)
  table below.
- **Use natural language:** `@ctx what should I work on?` routes to
  `/next`; `@ctx time to wrap up` routes to `/wrapup`. See
  [Natural Language](#natural-language).

The extension shows context-aware follow-up suggestions after each
command. For example, after `/init` you'll see buttons for "Show
status" or "Generate copilot integration."

## What Happens Automatically

The extension registers several VS Code event handlers that mirror
Claude Code's hook system. These run in the background; no user action
needed.

| Trigger | What fires |
|---------|------------|
| **File save** | Task-completion check on non-`.context/` files |
| **Git commit** | Notification prompting to add a Decision, Learning, run `/verify`, or Skip |
| **`.context/` file change** | Refreshes pending reminders and regenerates `.github/copilot-instructions.md` |
| **Dependency file change** | When `go.mod`, `package.json`, etc. change, prompts to refresh the dependency map (`/map`) |
| **Every 5 minutes** | Updates the reminder status-bar item and writes a heartbeat timestamp |
| **Extension activate** | Fires `ctx system session-event --type start` |
| **Extension deactivate** | Fires `ctx system session-event --type end` |

### Status Bar

A `$(bell) ctx` indicator appears in the status bar when you have
pending reminders. It refreshes every 5 minutes and hides itself when
nothing is due.

## Slash Commands

The extension surfaces 45 commands across six categories. The most
commonly used:

### Core Context

| Command | When to use |
|---------|-------------|
| `/init` | Initialize a `.context/` directory with template files |
| `/status` | Token estimate, file count, what's recent |
| `/agent` | Print AI-ready context packet |
| `/drift` | Detect stale paths, missing files, dead references |
| `/recall` | Browse and search prior AI session history |
| `/add` | Add a task, decision, learning, or convention |

### Session Lifecycle

| Command | When to use |
|---------|-------------|
| `/wrapup` | End-of-session ceremony: status, drift, journal audit |
| `/remember` | Structured readback (trigger: "Do you remember?") from tasks, decisions, learnings, recent journal |
| `/reflect` | Surface items worth persisting as decisions or learnings |
| `/pause` / `/resume` | Save and restore session state for later |

### Discovery & Planning

| Command | When to use |
|---------|-------------|
| `/brainstorm` | Browse and develop ideas from `ideas/` |
| `/spec` | List or scaffold feature specs from templates |
| `/verify` | Run verification (doctor + drift) |
| `/map` | Show dependency map (go.mod, package.json) |

Full list (with maintenance, audit, metadata, and system commands) is
in [editors/vscode/README.md](https://github.com/ActiveMemory/ctx/blob/main/editors/vscode/README.md#slash-commands).

## Natural Language

Plain English after `@ctx` is routed to the right command:

- "What should I work on next?" → `/next`
- "Time to wrap up" → `/wrapup`
- "Show me the status" → `/status`
- "Add a decision" → `/add`
- "Check for drift" → `/drift`

If the phrase doesn't match a known pattern, the extension surfaces a
short menu of likely matches.

## Auto-Bootstrap

If the `ctx` CLI isn't on PATH (or at a path configured via
`ctx.executablePath`), the extension auto-downloads the right binary:

1. Detects OS and architecture (darwin / linux / windows, amd64 / arm64).
2. Fetches the latest release from
   [GitHub Releases](https://github.com/ActiveMemory/ctx/releases).
3. Downloads and verifies the matching binary.
4. Caches it in VS Code's global storage directory.

Subsequent sessions reuse the cached binary. To pin a specific version,
set `ctx.executablePath` in your VS Code settings.

## Prerequisites

- **VS Code 1.93+**
- **[GitHub Copilot Chat](https://marketplace.visualstudio.com/items?itemName=GitHub.copilot-chat)** extension
- **`ctx` CLI** on PATH, or let the extension auto-download it

## Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `ctx.executablePath` | `ctx` | Path to the `ctx` CLI binary. Set this if `ctx` isn't on PATH and you don't want auto-download. |

## Refreshing the Integration

The extension updates through the VS Code Marketplace like any other
extension; install new versions via the Extensions view. Updates to
the **`ctx` CLI** are independent: bump it via your package manager, or
let the auto-bootstrap fetch the latest release.

Unlike the OpenCode integration, there is **no `ctx setup` step** for
VS Code. The extension carries its own runtime; `ctx`'s role is only to
provide the CLI it shells out to.

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| `@ctx` participant doesn't appear in Copilot Chat | Copilot Chat not installed or not signed in | Install [GitHub Copilot Chat](https://marketplace.visualstudio.com/items?itemName=GitHub.copilot-chat) and ensure you're signed in to a Copilot-eligible account |
| `@ctx /status` says `ctx` not found | CLI not on PATH and auto-download disabled | Either add `ctx` to PATH (`brew install activememory/tap/ctx` or download from [Releases](https://github.com/ActiveMemory/ctx/releases)), or unset `ctx.executablePath` to let the extension auto-download |
| Status-bar reminder never updates | Heartbeat suppressed or `.context/` doesn't exist | Run `ctx init` from your project root; reload VS Code if the indicator still doesn't appear within 5 minutes |
| Commands run but nothing is captured to `.context/` | Workspace folder missing or `.context/` outside the open folder | Make sure your project root (the one with `.context/`) is the workspace root, not a subdirectory of it |

## Verify It Works

Open Copilot Chat and ask:

```
@ctx Do you remember?
```

You should see a structured readback citing specific tasks, decisions,
and recent session topics. If you instead see "I don't have memory" or
"Let me check," something went wrong: confirm the CLI is reachable
(`@ctx /system doctor`) and `.context/` has files in it.

## What's Next

- [Your First Session](first-session.md): step-by-step walkthrough from
  `ctx init` to verified recall.
- [Common Workflows](common-workflows.md): day-to-day commands for
  tracking context, checking health, and browsing history.
- [Context Files](context-files.md): what lives in `.context/` and how
  each file is used.
- [Setup across AI Tools](../recipes/multi-tool-setup.md): wiring `ctx`
  for Claude Code, OpenCode, Cursor, Aider, Copilot, or Windsurf
  alongside VS Code.

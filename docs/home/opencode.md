---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "ctx for OpenCode"
icon: lucide/terminal
---

![ctx](../images/ctx-banner.png)

## The Problem

Every OpenCode session starts from zero. You re-explain your architecture,
the AI repeats mistakes it made yesterday, and decisions get rediscovered
instead of remembered.

**Without ctx:**

```
> "Add the validation middleware we discussed"

I don't have context about previous discussions. Could you describe
what validation middleware you're referring to?
```

**With ctx:**

```
> "Add the validation middleware we discussed"

Yes — from the Jan 15 session. You decided on Zod schemas at the
route level (DECISIONS.md #12), and the pattern is in
CONVENTIONS.md. I'll follow the existing middleware in
src/middleware/auth.ts as a reference.
```

That's the whole pitch: **your AI remembers**.

## Setup (One Command)

Install the `ctx` binary first ([installation docs](getting-started.md#installation)),
then run from your project root:

```bash
ctx setup opencode --write && ctx init && eval "$(ctx activate)"
```

This does three things:

1. **`ctx setup opencode --write`** — generates the project-local OpenCode plugin,
   skills, and `AGENTS.md`, then merges the ctx MCP server into OpenCode's
   global config (`~/.config/opencode/opencode.json` or
   `$OPENCODE_HOME/opencode.json`). This writes outside the project root
   because non-interactive shells (like MCP subprocesses) cannot discover
   project-local config — the same reason the Copilot CLI integration
   writes to `~/.copilot/mcp-config.json`.
2. **`ctx init`** — creates the `.context/` directory with template files
3. **`eval "$(ctx activate)"`** — binds `CTX_DIR` for your shell

### What Gets Created

| File | Purpose |
|------|---------|
| `.opencode/plugins/ctx.ts` | Lifecycle plugin (hooks into `ctx system` commands) |
| `~/.config/opencode/opencode.json` | Global MCP server registration (or `$OPENCODE_HOME/opencode.json`) |
| `AGENTS.md` | Agent instructions (OpenCode reads this natively) |
| `.opencode/skills/ctx-*/SKILL.md` | Slash command skills |

The plugin is a single file with no runtime dependencies — no `bun install`
or `npm install` needed. OpenCode loads it automatically on launch.

## What Happens Automatically

The plugin wires OpenCode lifecycle events to `ctx`. You don't need to
do anything — it just works.

| Event | What fires | What it does |
|-------|-----------|--------------|
| New session | `session.created` | Bootstraps ctx in the background so tools and hooks are ready for on-demand context access |
| Agent idle | `session.idle` | Nudges you to persist context; checks if recent edits completed any pending tasks |
| After `git commit` | `tool.execute.after` | Captures post-commit context state |
| After file edit | `tool.execute.after` | Detects if the edit completed a tracked task |
| Every shell call | `shell.env` | Injects `CTX_DIR` so all `ctx` commands resolve correctly |
| Context compaction | `experimental.session.compacting` | Re-injects context state into the compressed window — your memory survives compaction |

The last one matters most. When OpenCode compresses your context window to
free up tokens, ctx re-injects the full context state. Other tools lose
everything on compaction. ctx doesn't.

### How Compaction Works

When your conversation exceeds the context window, OpenCode runs a
compaction pass (you can trigger one manually with `/compact`). The
compaction agent summarizes older messages and drops the originals. Without
ctx, all accumulated knowledge disappears. With ctx, the plugin intercepts
the `experimental.session.compacting` event and appends `ctx system bootstrap`
output (context directory path and file inventory) into the compaction
context. The result: the compressed summary retains the breadcrumbs the
agent needs to re-read tasks, decisions, learnings, and conventions
on demand, even though the original messages that loaded them are gone.

### What Is *Not* Included

Note: dangerous-command blocking is Claude Code-specific and is not part of
the OpenCode integration. OpenCode's execution model (explicit user
approval for every shell command) makes a pre-execution blocklist
unnecessary.

## Slash Commands

Four skills are available as slash commands:

| Command | When to use |
|---------|-------------|
| `/ctx-agent` | Load full context packet. Use at session start or when context feels stale. |
| `/ctx-remember` | "Do you remember?" — reads tasks, decisions, learnings, and recent journal entries. Returns a structured readback. |
| `/ctx-status` | Context summary at a glance: file count, token estimate, recent activity. |
| `/ctx-wrap-up` | End-of-session ceremony. Captures learnings, decisions, conventions, and outstanding tasks to `.context/` files. |

You don't need to use these often. The plugin handles most context loading
automatically. These are for when you want explicit control.

## MCP Tools

The ctx MCP server exposes tools directly to the agent. These let the AI
read and write your context files without shell commands:

| Tool | Purpose |
|------|---------|
| `ctx_add` | Add a task, decision, learning, or convention |
| `ctx_complete` | Mark a task done by number or text match |
| `ctx_search` | Full-text search across all `.context/` files |
| `ctx_next` | Suggest the next pending task by priority |
| `ctx_drift` | Detect stale context: dead paths, missing files |
| `ctx_compact` | Archive completed tasks, clean empty sections |
| `ctx_remind` | List pending session-scoped reminders |
| `ctx_status` | Context health: file count, token estimate |
| `ctx_steering_get` | Retrieve steering files applicable to the current prompt |
| `ctx_journal_source` | Query recent AI session history |
| `ctx_session_event` | Signal session start/end lifecycle events |
| `ctx_watch_update` | Apply structured updates to `.context/` files |
| `ctx_check_task_completion` | After a write, detect silently completed tasks |

You don't invoke these yourself. The agent uses them as needed.

## Refreshing the Integration

If you re-run `ctx setup opencode --write` (e.g., after updating ctx), the
plugin and skills are rewritten in place. **Restart OpenCode to pick up the
refreshed plugin** — OpenCode only loads plugins at launch, not mid-session.

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| `opencode mcp list` shows `ctx ✗ failed MCP error -32000: Connection closed` | `CTX_DIR` not resolving in the MCP subprocess | Re-run `ctx setup opencode --write` to regenerate the sh-wrapper that sets `CTX_DIR` |
| Plugin installed but no hooks fire | Flat-file vs. subdirectory discovery mismatch (OpenCode requires `.opencode/plugins/<name>.ts`, not a subfolder) | Verify the plugin is at `.opencode/plugins/ctx.ts`. Check with `opencode --print-logs --log-level DEBUG` |
| `ctx agent` markdown leaking into the TUI | BunShell command missing `.nothrow().quiet()` | Update to the latest plugin: `ctx setup opencode --write` and restart |

## Verify It Works

Start a new OpenCode session and ask:

```
Do you remember?
```

The AI should cite specific context: current tasks, recent decisions, or
previous session topics. If it says "I don't have memory" or "Let me
check," something went wrong — check that the plugin installed correctly
and `.context/` has files in it.

## What's Next

- [Your First Session](first-session.md) — step-by-step walkthrough from
  `ctx init` to verified recall
- [Common Workflows](common-workflows.md) — day-to-day commands for
  tracking context, checking health, and browsing history
- [Context Files](context-files.md) — what lives in `.context/` and how
  each file is used

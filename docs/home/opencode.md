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

1. **`ctx setup opencode --write`** — generates the OpenCode plugin, MCP config,
   skills, and `AGENTS.md`
2. **`ctx init`** — creates the `.context/` directory with template files
3. **`eval "$(ctx activate)"`** — binds `CTX_DIR` for your shell

### What Gets Created

| File | Purpose |
|------|---------|
| `.opencode/plugins/ctx.ts` | Lifecycle plugin (hooks into `ctx system` commands) |
| `opencode.json` | MCP server registration (merged into existing config) |
| `AGENTS.md` | Agent instructions (OpenCode reads this natively) |
| `.opencode/skills/ctx-*/SKILL.md` | Slash command skills |

The plugin is a single file with no runtime dependencies — no `bun install`
or `npm install` needed. OpenCode loads it automatically on launch.

## What Happens Automatically

The plugin wires OpenCode lifecycle events to `ctx`. You don't need to
do anything — it just works.

| Event | What fires | What it does |
|-------|-----------|--------------|
| New session | `session.created` | Bootstraps context, loads a 4000-token AI context packet (tasks, decisions, learnings, conventions) |
| Agent idle | `session.idle` | Nudges you to persist context; checks if recent edits completed any pending tasks |
| After `git commit` | `tool.execute.after` | Captures post-commit context state |
| After file edit | `tool.execute.after` | Detects if the edit completed a tracked task |
| Every shell call | `shell.env` | Injects `CTX_DIR` so all `ctx` commands resolve correctly |
| Context compaction | `experimental.session.compacting` | Re-injects context state into the compressed window — your memory survives compaction |

The last one matters most. When OpenCode compresses your context window to
free up tokens, ctx re-injects the full context state. Other tools lose
everything on compaction. ctx doesn't.

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

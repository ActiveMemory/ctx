---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "ctx for Pi"
icon: lucide/terminal
---

![ctx](../images/ctx-banner.png)

## The Problem

Every Pi session starts from zero. You re-explain your architecture,
the agent repeats mistakes it made yesterday, and decisions get
rediscovered instead of remembered.

**Without `ctx`:**

```
> "Add the validation middleware we discussed"

I don't have context about previous discussions. Could you describe
what validation middleware you're referring to?
```

**With `ctx`:**

```
> "Add the validation middleware we discussed"

Yes. From the Jan 15 session. You decided on Zod schemas at the
route level (DECISIONS.md #12), and the pattern is in
CONVENTIONS.md. I'll follow the existing middleware in
src/middleware/auth.ts as a reference.
```

That's the whole pitch: **your AI remembers**.

## Setup (One Command)

Install the `ctx` binary first ([installation docs](getting-started.md#installation)),
then run from your project root:

```bash
ctx setup pi --write && ctx init
```

This does two things:

1. **`ctx setup pi --write`**: generates the project-local Pi
   extension, skills, and `AGENTS.md`.
2. **`ctx init`**: creates the `.context/` directory with template
   files.

### What Gets Created

| File | Purpose |
|------|---------|
| `.pi/extensions/ctx.ts` | Lifecycle extension (hooks into `ctx system` commands) |
| `AGENTS.md` | Agent instructions (Pi reads this natively) |
| `.pi/skills/ctx-*/SKILL.md` | `ctx` skills, available as `/skill:ctx-*` commands |

The extension is a single file with no runtime dependencies; no
`npm install` needed. Pi loads it automatically on launch.

> **Trust note:** Pi loads project-local `.pi/` files only after the
> project is trusted. The first `pi` launch in a new project prompts
> for trust; answer yes for the integration to activate.

> **Launch from the project root:** the extension resolves `.context/`
> relative to Pi's working directory, so start `pi` from the project
> root for the agent-side `ctx` commands to resolve.

## What Happens Automatically

The extension wires Pi lifecycle events to `ctx`. You don't need to
do anything; it just works.

| Event | What fires | What it does |
|-------|-----------|--------------|
| Session start | `session_start` | Warms the ctx agent packet off the prompt path |
| First prompt / after compaction | `before_agent_start` | Injects the packet as a persistent message when no ctx message exists in the live context (summary + kept tail + post-compaction entries) |
| After `git commit` | `tool_result` (bash) | Runs `ctx system post-commit` to capture context state (failed commits are ignored) |
| After file edit | `tool_result` (edit/write) | Runs `ctx system check-task-completion` to detect silent task completions |
| Agent settled | `agent_settled` | Runs `ctx system check-persistence` |

The compaction behavior matters most. When Pi compresses your context
window (`/compact` or auto-compaction) and the injected packet ages
past the kept window, it is folded into the lossy summary. The
extension then re-injects a fresh packet on the next prompt so the
agent keeps breadcrumbs back to your `.context/` directory and its
file inventory.

Pi intentionally has no built-in MCP, so there is no MCP server to
register; the extension is the lifecycle channel.

### What Is *Not* Included

Dangerous-command blocking is Claude Code-specific and is not part of
the Pi integration.

## Skills

Pi ships the same bundled skill set as OpenCode, generated at build time
from the canonical ctx skills (`hack/sync-pi-skills.sh` copies the
generated OpenCode tree). Four skills are Claude Code-only and are not
shipped: `ctx-permission-sanitize` (audits
`.claude/settings.local.json`), `ctx-plan-import` (reads
`~/.claude/plans/`), `ctx-dream` (headless `claude -p` cron), and
`ctx-skill-create` (authors Claude Code skills). Skills that cite
`references/` files ship those files alongside `SKILL.md`.

The ones you'll reach for most, available as `/skill:*` commands:

| Command | When to use |
|---------|-------------|
| `/skill:ctx-remember` | "Do you remember?"; reads tasks, decisions, learnings, and recent journal entries. Returns a structured readback. |
| `/skill:ctx-status` | Context summary at a glance: file count, token estimate, recent activity. |
| `/skill:ctx-wrap-up` | End-of-session ceremony. Captures learnings, decisions, conventions, and outstanding tasks to `.context/` files. |
| `/skill:ctx-agent` | Load full context packet on demand |

The KB editorial pipeline ships too: `/skill:ctx-kb-ingest`,
`/skill:ctx-kb-ask`, `/skill:ctx-kb-note`, `/skill:ctx-kb-ground`,
`/skill:ctx-kb-site-review`.

You don't need to use these often. The extension handles most context
loading automatically. These are for when you want explicit control.

## Refreshing the Integration

If you re-run `ctx setup pi --write` (e.g., after updating `ctx`), the
extension and skills are refreshed in place. Pi auto-discovered
extensions can be hot-reloaded with `/reload`; a restart always works.

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| Extension installed but nothing fires | Project not trusted | Trust the project when prompted; check `defaultProjectTrust` in Pi settings |
| `ctx` commands resolve to the wrong project | Pi launched outside the project root | Launch `pi` from the project root |
| Extension not loading at all | Wrong location | Verify the extension is at `.pi/extensions/ctx.ts` (flat top-level file) |

## Verify It Works

Start a new Pi session and ask:

```
Do you remember?
```

The agent should cite specific context: current tasks, recent
decisions, or previous session topics. If it says "I don't have
memory" or "Let me check," something went wrong; check that the
extension installed correctly and `.context/` has files in it.

## What's Next

- [Your First Session](first-session.md): step-by-step walkthrough
  from `ctx init` to verified recall.
- [Common Workflows](common-workflows.md): day-to-day commands for
  tracking context, checking health, and browsing history.
- [Context Files](context-files.md): what lives in `.context/` and
  how each file is used.

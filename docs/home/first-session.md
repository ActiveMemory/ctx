---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Your First Session
icon: lucide/play
---

![ctx](../images/ctx-banner.png)

Here's what a complete first session looks like, from initialization to
the moment your AI cites your project context back to you.

## Step 1: Initialize Your Project

```bash
cd your-project
ctx init
```

```
Context initialized in .context/

  ✓ CONSTITUTION.md
  ✓ TASKS.md
  ✓ DECISIONS.md
  ✓ LEARNINGS.md
  ✓ CONVENTIONS.md
  ✓ ARCHITECTURE.md
  ✓ GLOSSARY.md
  ✓ AGENT_PLAYBOOK.md

Creating project root files...
  ✓ PROMPT.md
  ✓ IMPLEMENTATION_PLAN.md

Setting up Claude Code permissions...
  ✓ .context.key

Claude Code plugin (hooks + skills):
  Install: claude /plugin marketplace add ActiveMemory/ctx
  Then:    claude /plugin install ctx@activememory-ctx

Next steps:
  1. Edit .context/TASKS.md to add your current tasks
  2. Run 'ctx status' to see context summary
  3. Run 'ctx agent' to get AI-ready context packet
```

This created your `.context/` directory with template files. 

For Claude Code, [install the ctx plugin](getting-started.md#installation) to get automatic hooks and skills.

## Step 2: Populate Your Context

Add a task and a decision — these are the entries your AI will remember:

```bash
ctx add task "Implement user authentication"
```

```
✓ Added to TASKS.md
```

```bash
ctx add decision "Use PostgreSQL for primary database" \
  --context "Need a reliable database for production" \
  --rationale "PostgreSQL offers ACID compliance and JSON support" \
  --consequences "Team needs PostgreSQL training"
```

```
✓ Added to DECISIONS.md
```

These entries are what the AI will recall in future sessions. You don't need
to populate everything now: Context grows **naturally** as you work.

## Step 3: Check Your Context

```bash
ctx status
```

```
Context Status
====================

Context Directory: .context/
Total Files: 8
Token Estimate: 1,247 tokens

Files:
  ✓ CONSTITUTION.md (loaded)
  ✓ TASKS.md (1 items)
  ✓ DECISIONS.md (1 items)
  ○ LEARNINGS.md (empty)
  ✓ CONVENTIONS.md (loaded)
  ✓ ARCHITECTURE.md (loaded)
  ✓ GLOSSARY.md (loaded)
  ✓ AGENT_PLAYBOOK.md (loaded)

Recent Activity:
  - TASKS.md modified 2 minutes ago
  - DECISIONS.md modified 1 minute ago
```

Notice the **token estimate**: This is how much context your AI will load.

The `○` next to `LEARNINGS.md` means it's still empty;
it will fill in as you capture lessons during development.

## Step 4: Start an AI Session

With **Claude Code** (and the ctx plugin), start every session with:

```
/ctx-remember
```

This loads your context and presents a structured 
[**readback**](prompting-guide.md#do-you-remember) so you can
confirm the agent knows what is going on. Context also loads automatically
via hooks, but the explicit ceremony gives you a **readback** to verify.

!!! info "Using VS Code?"
    With **VS Code Copilot Chat** (and the
    [ctx extension](../operations/integrations.md#vs-code-chat-extension-ctx)), 
    type `@ctx /agent` in chat to load your context packet, or `@ctx /status`
    to check your project context. Run `ctx hook copilot --write` once
    to generate `.github/copilot-instructions.md` for automatic context loading.

If you are not using Claude Code, generate a 
**context packet** for your AI tool:

```bash
ctx agent --budget 8000
```

```markdown
# Context Packet
Generated: 2026-02-14T15:30:45Z | Budget: 8000 tokens | Used: ~2450

## Read These Files (in order)
1. .context/CONSTITUTION.md
2. .context/TASKS.md
3. .context/CONVENTIONS.md
...

## Current Tasks
- [ ] Implement user authentication
- [ ] Add rate limiting to API endpoints

## Key Conventions
- Use gofmt for formatting
- Path construction uses filepath.Join

## Recent Decisions
## [2026-02-14-120000] Use PostgreSQL for the primary database

**Context**: Evaluated PostgreSQL, MySQL, and SQLite...
**Rationale**: PostgreSQL offers better JSON support...

## Key Learnings
## [2026-02-14-100000] Connection pool sizing matters

**Context**: Hit connection limits under load...
**Lesson**: Default pool size of 10 is too low for concurrent requests...

## Also Noted
- Use JWT for session management
- Always validate input at API boundary
```

Paste this output into your AI tool's system prompt or conversation start.

## Step 5: Verify It Works

Ask your AI: **"What are our current tasks?"**

A working setup produces a response like:

```text
Based on the project context, you have one active task:

- **Implement user authentication** (pending)

There's also a recent architectural decision to **use PostgreSQL for
the primary database**, chosen for its ACID compliance and JSON support.

Want me to start on the authentication task?
```

That's the success moment:

The AI is citing your exact context entries from Step 2, not hallucinating or
asking you to re-explain.

## What Gets Created

```
.context/
├── CONSTITUTION.md     # Hard rules — NEVER violate these
├── TASKS.md            # Current and planned work
├── CONVENTIONS.md      # Project patterns and standards
├── ARCHITECTURE.md     # System overview
├── DECISIONS.md        # Architectural decisions with rationale
├── LEARNINGS.md        # Lessons learned, gotchas, tips
├── GLOSSARY.md         # Domain terms and abbreviations
└── AGENT_PLAYBOOK.md   # How AI tools should use this
```

Claude Code integration (hooks + skills) is provided by the
**`ctx` plugin**:<br />
See [Integrations/Claude Code](../operations/integrations.md#claude-code-full-integration).

VS Code Copilot Chat integration is provided by the
**`ctx` extension**:<br />
See [Integrations/VS Code](../operations/integrations.md#vs-code-chat-extension-ctx).

See [Context Files](context-files.md) for detailed documentation of each file.

## What to `.gitignore`

!!! tip "Rule of Thumb"
    * If it's knowledge (*decisions, tasks, learnings,
      conventions*), **commit it**.
    * If it's generated output, raw session data, or a secret, `.gitignore` it.

**Commit** your `.context/` knowledge files: **that's the whole point**.

You should **`.gitignore`** the generated and sensitive paths:

```gitignore
# Journal data (large, potentially sensitive)
.context/journal/
.context/journal-site/
.context/journal-obsidian/

# Hook logs (machine-specific)
.context/logs/

# Encryption key (NEVER commit)
.context/.context.key

# Claude Code local settings (machine-specific)
.claude/settings.local.json
```

!!! tip "ctx init Patches Your .gitignore for You"
    `ctx init` automatically adds these entries to your `.gitignore`.

     Review the additions with `cat .gitignore` after init.

*See also*:

* [Security considerations](../security/reporting.md)
* [Scratchpad encryption](../reference/scratchpad.md)
* [Session Journal](../reference/session-journal.md)

----

**Next Up**: [Common Workflows →](common-workflows.md) — day-to-day commands for tracking context, checking health, and browsing history.

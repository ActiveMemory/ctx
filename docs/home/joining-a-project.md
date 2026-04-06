---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Joining a ctx Project
icon: lucide/user-plus
---

![ctx](../images/ctx-banner.png)

You've joined a team or inherited a project, and there's a `.context/`
directory in the repo. Good news: someone already set up persistent
context. This page gets you oriented fast.

## What to Read First

The files in `.context/` have a deliberate priority order. Read them
top-down:

1. **CONSTITUTION.md**: Hard rules. Read this before you touch anything.
   These are inviolable constraints the team has agreed on.
2. **TASKS.md**: Current and planned work. Shows what's in progress,
   what's pending, and what's blocked.
3. **CONVENTIONS.md**: How the team writes code. Naming patterns,
   file organization, preferred idioms.
4. **ARCHITECTURE.md**: System overview. Components, boundaries, data flow.
5. **DECISIONS.md**: Why things are the way they are. Saves you from
   re-proposing something the team already evaluated and rejected.
6. **LEARNINGS.md**: Gotchas, tips, and hard-won lessons. The stuff
   that doesn't fit anywhere else but will save you hours.

See [Context Files](context-files.md) for detailed documentation of
each file's structure and purpose.

## Checking Context Health

Before you start working, check whether the context is current:

```bash
ctx status
```

This shows file counts, token estimates, and recent activity. If files
haven't been touched in weeks, the context may be stale.

```bash
ctx drift
```

This compares context files against recent code changes and flags
potential drift: decisions that no longer match the codebase,
conventions that have shifted, or tasks that look outdated.

If things are stale, mention it to the team. Don't silently fix it
yourself on day one.

## Starting Your First Session

Generate a context packet to prime your AI:

```bash
ctx agent --budget 8000
```

This outputs a token-budgeted summary of the project context, ordered
by priority. With Claude Code and the ctx plugin, context loads
automatically via hooks. You can also use the `/ctx-remember` skill
to get a structured readback of what the AI knows.

The readback is your verification step: if the AI can cite specific
tasks and decisions, the context is working.

## Adding Context

As you work, you'll discover things worth recording. Use the CLI:

```bash
# Record a decision you made or learned about
ctx add decision "Use connection pooling for DB access" \
  --rationale "Reduces connection overhead under load" \
  --session-id abc12345 --branch main --commit 68fbc00a

# Capture a gotcha you hit
ctx add learning "Redis timeout defaults to 5s" \
  --context "Hit timeouts during bulk operations" \
  --application "Set explicit timeout for batch jobs" \
  --session-id abc12345 --branch main --commit 68fbc00a

# Add a convention you noticed the team follows
ctx add convention "All API handlers return structured errors"
```

You can also just tell the AI: "Record this as a learning" or
"Add this decision to context." With the ctx plugin, context-update
commands handle the file writes.

See the [Knowledge Capture recipe](../recipes/knowledge-capture.md) for
the full workflow.

## Session Etiquette

A few norms for working in a ctx-managed project:

- **Respect existing conventions.** If `CONVENTIONS.md` says
  "use `filepath.Join`," use `filepath.Join`. If you disagree, propose
  a change, don't silently diverge.
- **Don't restructure context files without asking.** The file layout
  and section structure are shared state. Reorganizing them affects
  every team member and every AI session.
- **Mark tasks done when complete.** Check the box (`[x]`) in place.
  Don't move tasks between sections or delete them.
- **Add context as you go.** Decisions, learnings, and conventions
  you discover are valuable to the next person (or the next session).

## Common Pitfalls

**Ignoring CONSTITUTION.md.** The constitution exists for a reason.
If a task conflicts with a constitution rule, the task is wrong. Raise
it with the team instead of working around the constraint.

**Deleting tasks.** Never delete a task from TASKS.md. Mark it `[x]`
(done) or `[-]` (skipped with a reason). The history matters for
session replay and audit.

**Bypassing hooks.** If the project uses ctx hooks (pre-commit nudges,
context autoloading), don't disable them. They exist to keep context
fresh. If a hook is noisy or broken, fix it or file a task.

**Over-contributing on day one.** Read first, then contribute. Adding
a dozen learnings before you understand the project's norms creates
noise, not signal.

----

**Related**:

* [Getting Started](getting-started.md): installation and setup from scratch
* [Context Files](context-files.md): detailed file reference
* [Knowledge Capture](../recipes/knowledge-capture.md): recording decisions, 
  learnings, and conventions
* [Session Lifecycle](../recipes/session-lifecycle.md): how a typical AI 
  session flows with ctx

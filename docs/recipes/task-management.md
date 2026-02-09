---
title: "Tracking Work Across Sessions"
icon: lucide/list-checks
---

![ctx](../images/ctx-banner.png)

## Problem

You have work that spans multiple sessions. Tasks get added during one session,
partially finished in another, and completed days later. Without a system,
follow-up items fall through the cracks, priorities drift, and you lose track
of what was done versus what still needs doing. TASKS.md grows cluttered with
completed checkboxes that obscure the remaining work.

How do you manage work items that span multiple sessions without losing context?

## Commands and Skills Used

| Tool                   | Type    | Purpose                                  |
|------------------------|---------|------------------------------------------|
| `ctx add task`         | Command | Add a new task to TASKS.md               |
| `ctx complete`         | Command | Mark a task as done by number or text    |
| `ctx tasks snapshot`   | Command | Create a point-in-time backup of TASKS.md |
| `ctx tasks archive`    | Command | Move completed tasks to archive file     |
| `/ctx-add-task`        | Skill   | AI-assisted task creation with validation |
| `/ctx-archive`         | Skill   | AI-guided archival with safety checks    |
| `/ctx-next`            | Skill   | Pick what to work on based on priorities |

## The Workflow

### Step 1: Add Tasks with Priorities

Every piece of follow-up work gets a task. Use `ctx add task` from the terminal
or `/ctx-add-task` from your AI assistant. Tasks should start with a verb and be
specific enough that someone unfamiliar with the session could act on them.

```bash
# High-priority bug found during code review
ctx add task "Fix race condition in session cooldown when two hooks fire simultaneously" --priority high

# Medium-priority feature work
ctx add task "Add --format json flag to ctx status for CI integration" --priority medium

# Low-priority cleanup
ctx add task "Remove deprecated --raw flag from ctx load" --priority low
```

The `/ctx-add-task` skill validates your task before recording it. It checks
that the description is actionable, not a duplicate, and specific enough for
someone else to pick up. If you say "fix the bug," it will ask you to clarify
which bug and where.

### Step 2: Organize with Phase Sections

Tasks live in phase sections inside TASKS.md. Phases provide logical groupings
that preserve order and enable replay. A task never moves between sections --
it stays in its phase permanently, and status is tracked via checkboxes and
inline tags.

```markdown
## Phase 1: Core CLI

- [x] Implement ctx add command `#done:2026-02-01-143022`
- [x] Implement ctx complete command `#done:2026-02-03-091544`
- [ ] Add --section flag to ctx add task `#priority:medium`

## Phase 2: AI Integration

- [ ] Implement ctx agent cooldown `#priority:high` `#in-progress`
- [ ] Add ctx watch XML parsing `#priority:medium`
  - Blocked by: Need to finalize agent output format

## Backlog

- [ ] Performance optimization for large TASKS.md files `#priority:low`
- [ ] Add metrics dashboard to ctx status `#priority:deferred`
```

Use `--section` when adding a task to a specific phase:

```bash
ctx add task "Add ctx watch XML parsing" --priority medium --section "Phase 2: AI Integration"
```

Without `--section`, the task is inserted before the first unchecked task in
TASKS.md.

### Step 3: Pick What to Work On

At the start of a session, or after finishing a task, use `/ctx-next` to get
prioritized recommendations. The skill reads TASKS.md, checks recent sessions,
and ranks candidates using explicit priority, blocking status, in-progress
state, momentum from recent work, and phase order.

```text
/ctx-next
```

The output looks like this:

> **1. Implement ctx agent cooldown** `#priority:high`
> Still in-progress from yesterday's session. The tombstone file approach is
> half-built. Finishing is cheaper than context-switching.
>
> **2. Add --section flag to ctx add task** `#priority:medium`
> Last Phase 1 item. Quick win that unblocks organized task entry.
>
> ---
>
> *Based on 8 pending tasks across 3 phases. Last session: agent-cooldown (2026-02-06).*

In-progress tasks almost always come first. Finishing existing work takes
priority over starting new work.

### Step 4: Complete Tasks

When a task is done, mark it complete by number or partial text match:

```bash
# By task number (as shown in TASKS.md)
ctx complete 3

# By partial text match
ctx complete "agent cooldown"
```

The task's checkbox changes from `[ ]` to `[x]` and a `#done` timestamp is
added. Tasks are never deleted -- they stay in their phase section so the
history is preserved.

### Step 5: Snapshot Before Risky Changes

Before a major refactor or any change that might break things, snapshot your
current task state. This creates a copy of TASKS.md in `.context/archive/`
without modifying the original.

```bash
# Default snapshot
ctx tasks snapshot

# Named snapshot (recommended before big changes)
ctx tasks snapshot "before-refactor"
```

This creates a file like `.context/archive/tasks-before-refactor-2026-02-08-1430.md`.
If the refactor goes sideways and you need to understand what the task state
looked like before you started, the snapshot is there.

Snapshots are cheap. Take them before any change you might want to undo or
review later.

### Step 6: Archive When TASKS.md Gets Cluttered

After several sessions, TASKS.md accumulates completed tasks that make it hard
to see what is still pending. Use `ctx tasks archive` to move all `[x]` items
to a timestamped archive file.

Start with a dry run to preview what will be moved:

```bash
ctx tasks archive --dry-run
```

Then archive:

```bash
ctx tasks archive
```

Completed tasks move to `.context/archive/tasks-2026-02-08.md`. Phase headers
are preserved in the archive for traceability. Pending tasks (`[ ]`) remain
in TASKS.md.

The `/ctx-archive` skill adds two safety checks before archiving: it verifies
that completed tasks are genuinely cluttering the view and that nothing was
marked `[x]` prematurely.

## Putting It Together

```bash
# Add a task
ctx add task "Implement rate limiting for API endpoints" --priority high

# Add to a specific phase
ctx add task "Write integration tests for rate limiter" --section "Phase 2"

# See what to work on
# (from AI assistant) /ctx-next

# Mark done by text
ctx complete "rate limiting"

# Mark done by number
ctx complete 5

# Snapshot before a risky refactor
ctx tasks snapshot "before-middleware-rewrite"

# Archive completed tasks when the list gets long
ctx tasks archive --dry-run     # preview first
ctx tasks archive               # then archive
```

## Tips

- **Start tasks with a verb.** "Add," "Fix," "Implement," "Investigate" -- not
  just a topic like "Authentication."
- **Include the "why" in the task description.** Future sessions lack the context
  of why you added the task. "Add rate limiting" is worse than "Add rate limiting
  to prevent abuse on the public API after the load test showed 10x traffic spikes."
- **Use `#in-progress` sparingly.** Only one or two tasks should carry this tag
  at a time. If everything is in-progress, nothing is.
- **Snapshot before, not after.** The point of a snapshot is to capture the state
  *before* a change, not to celebrate what you just finished.
- **Archive regularly.** Once completed tasks outnumber pending ones, it is time
  to archive. A clean TASKS.md helps both you and your AI assistant focus.
- **Never delete tasks.** Mark them `[x]` (completed) or `[-]` (skipped with a
  reason). Deletion breaks the audit trail and violates the constitution.

## See Also

- [The Complete Session](session-lifecycle.md) -- full session lifecycle
  including task management in context
- [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md) --
  capturing the "why" behind your work
- [Detecting and Fixing Drift](context-health.md) -- keeping TASKS.md
  accurate over time
- [CLI Reference](../cli-reference.md) -- full documentation for `ctx add`,
  `ctx complete`, `ctx tasks`
- [Context Files: TASKS.md](../context-files.md#tasksmd) -- format and
  conventions for TASKS.md

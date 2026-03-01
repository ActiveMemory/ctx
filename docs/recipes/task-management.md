---
title: "Tracking Work Across Sessions"
icon: lucide/list-checks
---

![ctx](../images/ctx-banner.png)

## The Problem

You have work that spans multiple sessions. Tasks get added during one session,
partially finished in another, and completed days later.

Without a **system**, follow-up items fall through the cracks, priorities drift,
and you lose track of what was done versus what still needs doing. `TASKS.md`
grows cluttered with completed checkboxes that obscure the remaining work.

How do you manage work items that span multiple sessions without losing context?

## TL;DR

**Manage Tasks**:

```bash
ctx add task "Fix race condition" --priority high  # add
ctx add task "Write tests" --section "Phase 2"     # add to phase
ctx complete "race condition"                      # mark done
ctx tasks snapshot "before-refactor"               # backup
ctx tasks archive                                  # clean up
```

**Pick Up the Next Task**:

```text
/ctx-next # pick what's next
```

Read on for the full workflow and conversational patterns.

## Commands and Skills Used

| Tool                 | Type    | Purpose                                     |
|----------------------|---------|---------------------------------------------|
| `ctx add task`       | Command | Add a new task to `TASKS.md`                |
| `ctx complete`       | Command | Mark a task as done by number or text       |
| `ctx tasks snapshot` | Command | Create a point-in-time backup of `TASKS.md` |
| `ctx tasks archive`  | Command | Move completed tasks to archive file        |
| `/ctx-add-task`      | Skill   | AI-assisted task creation with validation   |
| `/ctx-archive`       | Skill   | AI-guided archival with safety checks       |
| `/ctx-next`          | Skill   | Pick what to work on based on priorities    |

## The Workflow

### Step 1: Add Tasks with Priorities

Every piece of follow-up work gets a task. Use `ctx add task` from the terminal
or `/ctx-add-task` from your AI assistant. Tasks should start with a verb and be
specific enough that someone unfamiliar with the session could act on them.

```bash
# High-priority bug found during code review
ctx add task "Fix race condition in session cooldown" --priority high

# Medium-priority feature work
ctx add task "Add --format json flag to ctx status for CI integration" --priority medium

# Low-priority cleanup
ctx add task "Remove deprecated --raw flag from ctx load" --priority low
```

The `/ctx-add-task` skill validates your task before recording it. It checks
that the description is actionable, not a duplicate, and specific enough for
someone else to pick up.

If you say "*fix the bug,*" it will ask you to clarify which bug and where.

!!! tip "Tasks Are Often Created Proactively"
    In practice, many tasks are created **proactively** by the agent rather than by
    explicit CLI commands.

    After completing a feature, the agent will often identify follow-up work:
    tests, docs, edge cases, error handling, and offer to add them as tasks.

    You do not need to dictate `ctx add task` commands; the agent picks up on
    work context and suggests tasks naturally.

### Step 2: Organize with Phase Sections

Tasks live in phase sections inside `TASKS.md`.

Phases provide logical groupings that preserve order and enable replay.

A task does not move between sections. It stays in its phase permanently, and
status is tracked via checkboxes and inline tags.

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
ctx add task "Add ctx watch XML parsing" --priority medium --section \
    "Phase 2: AI Integration"
```

Without `--section`, the task is inserted before the first unchecked task in
`TASKS.md`.

### Step 3: Pick What to Work On

At the start of a session, or after finishing a task, use `/ctx-next` to get
prioritized recommendations. 

The skill reads `TASKS.md`, checks recent sessions, and ranks candidates using 
explicit priority, blocking status, in-progress state, momentum from 
recent work, and phase order.

You can also ask naturally: "*what should we work on?*" or
"*what's the highest priority right now?*"

```text
/ctx-next
```

The output looks like this:

```markdown
**1. Implement ctx agent cooldown** `#priority:high`

    Still in-progress from yesterday's session. The tombstone file approach is
    half-built. Finishing is cheaper than context-switching.

**2. Add --section flag to ctx add task** `#priority:medium`

    Last Phase 1 item. Quick win that unblocks organized task entry.

---

*Based on 8 pending tasks across 3 phases.

Last session: agent-cooldown (2026-02-06).*
```

In-progress tasks almost always come first: 

Finishing existing work takes priority over starting new work.

### Step 4: Complete Tasks

When a task is done, mark it complete by number or partial text match:

```bash
# By task number (as shown in TASKS.md)
ctx complete 3

# By partial text match
ctx complete "agent cooldown"
```

The task's checkbox changes from `[ ]` to `[x]` and a `#done` timestamp is
added. Tasks are never deleted: they stay in their phase section so history is
preserved.

!!! tip "Be Conversational"
    You rarely need to run `ctx complete` yourself during an interactive session.

    When you say something like "*the rate limiter is done*" or "*we finished that*,"
    the agent marks the task complete and moves on to suggesting what is next.

    The CLI commands are most useful for manual housekeeping, scripted workflows,
    or when you want precision.

### Step 5: Snapshot Before Risky Changes

Before a major refactor or any change that might break things, snapshot your
current task state. This creates a copy of `TASKS.md` in `.context/archive/`
without modifying the original.

```bash
# Default snapshot
ctx tasks snapshot

# Named snapshot (recommended before big changes)
ctx tasks snapshot "before-refactor"
```

This creates a file like `.context/archive/tasks-before-refactor-2026-02-08-1430.md`.
If the refactor goes sideways, and you need to confirm what the task state looked
like before you started, the snapshot is there.

**Snapshots are cheap**: Take them before any change you might want to undo or
review later.

### Step 6: Archive When `TASKS.md` Gets Cluttered

After several sessions, `TASKS.md` accumulates completed tasks that make it hard
to see what is still pending.

Use `ctx tasks archive` to move all `[x]` items to a timestamped archive file.

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
in `TASKS.md`.

The `/ctx-archive` skill adds two safety checks before archiving: it verifies
that completed tasks are genuinely cluttering the view and that nothing was
marked `[x]` prematurely.

## Letting the Agent Drive

The steps above show the CLI commands for task management.

In practice, most task management happens conversationally:

An agent that has loaded the context files does not need you to type
`ctx add task`. It tracks work naturally and offers the right operations.

### Use Natural Language Instead of Commands

These conversational prompts replace explicit commands during interactive sessions:

| Instead of typing...                   | Just say...                                        |
|----------------------------------------|----------------------------------------------------|
| `ctx add task "Write tests for X"`     | "We should add tests for this: track that?"        |
| `/ctx-next`                            | "What should we work on?"                          |
| `ctx complete "rate limiting"`         | "The rate limiter is done, what's next?"           |
| `ctx tasks archive`                    | "`TASKS.md` is getting long, can you clean it up?" |
| `ctx add task ... && ctx add task ...` | "Add follow-ups for what we just built."           |

The agent translates these into the right `ctx` operations behind the scenes.

### The Agent Creates Tasks Proactively

After completing a piece of work, the agent does not wait for you to dictate
follow-up items. It identifies them from the work context and offers:

```text
**Agent**: "I've finished implementing the rate limiter middleware. I noticed
three follow-up items:

1. Write integration tests for the rate limiter (the unit tests pass but
   there's no end-to-end coverage)
2. Add rate limit headers to API responses (X-RateLimit-Remaining, etc.)
3. Document the rate limit configuration in the ops runbook

Want me to add these as tasks?"
```

This is how an agent behaves once it has read the Agent Playbook.

The playbook's "**Work, Reflect, Persist**" cycle means the agent asks itself
after every milestone: Did I just complete a task? Did new work emerge?

When the answer is yes, it acts.

### Asking the Agent What to Capture

You don't always know what's worth persisting.

Instead of reviewing the session yourself and typing specific commands, ask the
agent to do the thinking:

| You say                                                        | What happens                                                          |
|----------------------------------------------------------------|-----------------------------------------------------------------------|
| "*What tasks came out of this session?*"                       | Agent reviews the work, identifies follow-ups, and offers to add them |
| "*Are there any tasks, learnings, or decisions worth adding?*" | Agent scans the full session for all three types and proposes a batch |
| "*What should we remember from today?*"                        | Agent surfaces everything persist-worthy: not just tasks              |
| "*Anything we should track before we move on?*"                | Agent checks for loose ends specific to the work just completed       |

This is especially powerful after long or exploratory sessions where multiple
threads were touched. 

The agent has the full conversation in context and can spot follow-ups you 
might miss:

```text
You:   "Are there any tasks, learnings, or decisions worth adding
        from this session?"

Agent: "Yes. Here's what I'd suggest:

  Tasks:
  - Add retry logic to the webhook delivery (we discussed it when
    the timeout test failed but never tracked it)
  - Update the API docs to reflect the new rate limit headers

  Learnings:
  - net/http's default client has no timeout: we burned 15 minutes
    on that before adding an explicit 30s timeout

  Decisions:
  - We chose token-bucket over sliding-window for rate limiting
    because it handles bursts better and the library is simpler

  Want me to add all of these?"

You:   "All of them."
```

The agent fills in structured fields from what it observed during the session.
You get the same quality entries as if you typed verbose commands, without
needing to remember what happened or how to format it.

### The Agent Completes Tasks and Moves On

When you confirm that work is finished, or when the agent can see that it is,
it marks the task done and suggests the next step:

```text
You: 
"The cooldown logic is working now."

Agent: 
"Marked 'Implement ctx agent cooldown' as done. Your highest-
priority remaining task is 'Add --section flag to ctx add task'.

It's the last item in Phase 1 and a quick win. Want to start on that?"
```

The agent chains completion into next-task selection.

One sentence from you triggers two operations.

### End-of-Session Capture

At the end of a session, the agent reviews what happened and catches loose ends:

```
You: 
"I think we're good for today."

Agent: 
"Before we wrap up: during the middleware work you mentioned
wanting to add metrics for rate limit hits. That never made it into TASKS.md.

Want me to add it as a task so it's there next session?"
```

This is the proactive persistence model at work.

The agent monitors the conversation for mentioned-but-not-tracked work and
surfaces it before context is lost.

### A Conversational Session Example

Here is what task management looks like in a real session when the agent is
driving:

```text
You:   What's left to do?
Agent: You have 5 pending tasks across 2 phases. The highest priority is
       "Implement ctx agent cooldown": it's marked in-progress from
       yesterday and half-built. Want to pick that up?

You:   Yeah, let's finish it.
       [... 20 minutes of implementation work ...]

You:   That's working. Tests pass.
Agent: Nice. I'll mark "Implement ctx agent cooldown" as done.
       I noticed two things from the implementation:
       1. The tombstone file cleanup should happen on a timer, not just
          at startup. Want me to add that as a task?
       2. We learned that file locks don't work reliably on NFS. I'll
          add that to LEARNINGS.md.

You:   Yes to both. What's next?
Agent: Added the cleanup task as high priority. Next up is "Add --section
       flag to ctx add task": last item in Phase 1. It should be
       straightforward since the flag parsing pattern already exists
       in the codebase. Ready?
```

!!! tip "It's All Conversational"
    Notice what did not happen:
    The user never typed a `ctx` command.

    The agent handled task completion, follow-up creation, learning capture,
    and next-task selection: all from natural conversation.

## Putting It All Together

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

* Start tasks with a **verb**: "*Add,*" "*Fix,*" "*Implement,*" "*Investigate*": 
  not just a topic like "Authentication."
* Include the **why** in the task description. Future sessions lack the context of
  why you added the task. "Add rate limiting" is worse than "Add rate limiting
  to prevent abuse on the public API after the load test showed 10x traffic spikes."
* Use `#in-progress` **sparingly**. Only one or two tasks should carry this tag at
  a time. If everything is in-progress, nothing is.
* Snapshot **before**, not after. The point of a snapshot is to capture the 
  state before a change, not to celebrate what you just finished.
* **Archive regularly**. Once completed tasks outnumber pending ones, it is time
  to archive. A clean `TASKS.md` helps both you and your AI assistant focus.
* **Never delete tasks**. Mark them `[x]` (completed) or `[-]` (skipped with a
  reason). Deletion breaks the audit trail.
* **Trust the agent's task instincts**. When the agent suggests follow-up items
  after completing work, it is drawing on the full context of what just happened.
* **Conversational prompts beat commands** in interactive sessions. Saying
  "what should we work on?" is faster and more natural than running `/ctx-next`.
  Save explicit commands for scripts, CI, and unattended runs.
* **Let the agent chain operations**. A single statement like "that's done, what's
  next?" can trigger completion, follow-up identification, and next-task
  selection in one flow.
* **Review** proactive task suggestions before moving on. The best follow-ups come
  from items spotted in-context right after the work completes.

## Next Up

**[Using the Scratchpad →](scratchpad-with-claude.md)**:
Store short-lived sensitive notes in an encrypted scratchpad.

## See Also

* [The Complete Session](session-lifecycle.md): full session lifecycle including
  task management in context
* [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md):
  capturing the "why" behind your work
* [Detecting and Fixing Drift](context-health.md):
  keeping `TASKS.md` accurate over time
* [CLI Reference](../cli/context.md):
  full documentation for `ctx add`, `ctx complete`, `ctx tasks`
* [Context Files: `TASKS.md`](../home/context-files.md#tasksmd): 
  format and conventions for `TASKS.md`

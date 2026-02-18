---
name: ctx-next
description: "Suggest what to work on next. Use when starting a session, finishing a task, or when unsure what to prioritize."
allowed-tools: Bash(ctx:*), Read
---

Analyze current tasks and recent session activity, then suggest
1-3 concrete next actions with rationale.

## When to Use

- At session start after loading context ("what should I do?")
- After completing a task ("what's next?")
- When the user asks for priorities or direction
- When multiple tasks exist and it's unclear which to pick

## When NOT to Use

- When the user has already stated what they want to work on
- When actively mid-task (don't interrupt flow with suggestions)
- When no `.context/` directory exists (nothing to analyze)

## Usage Examples

```text
/ctx-next
/ctx-next (just finished the auth refactor)
```

## Process

Do all of this **silently** — do not narrate the steps:

1. **Read TASKS.md** to get the full task list with statuses,
   priorities, and phases
2. **Check recent sessions** to understand what was just worked
   on and avoid suggesting already-completed work:
   ```bash
   ctx recall list --limit 3
   ```
3. **Read the most recent session file** (if any) to understand
   what was accomplished and what follow-up items were noted
4. **Analyze and rank** tasks using the priority logic below
5. **Present 1-3 recommendations** in the output format below

## Priority Logic

Rank candidate tasks using these criteria (in order):

1. **Explicit priority**: `#priority:high` > `#priority:medium`
   > `#priority:low` > untagged
2. **Unblocked**: tasks not tagged `#blocked` or listed under a
   "Blocked" section
3. **In-progress first**: `#in-progress` tasks should be resumed
   before starting new ones (finishing > starting)
4. **Momentum**: prefer tasks related to recent session work
   (continuing a thread is cheaper than context-switching)
5. **Phase order**: earlier phases before later phases (Phase 0
   before Phase 1, etc.) unless priority overrides
6. **Quick wins**: if two tasks have equal priority, prefer the
   one that seems smaller/faster (builds momentum)

### Skip these tasks:

- `[x]` completed tasks
- `[-]` skipped tasks
- Tasks explicitly tagged `#blocked` with no resolution path
- Tasks that were the main focus of the most recent session
  (user likely wants variety or the session ended because it
  was done)

## Output Format

Present your recommendations like this:

### Recommended Next

**1. [Task title or summary]** `#priority:X`
> [1-2 sentence rationale: why this, why now]

**2. [Task title or summary]** `#priority:X`
> [1-2 sentence rationale]

**3. [Task title or summary]** *(optional — only if genuinely
useful)*
> [1-2 sentence rationale]

---

*Based on N pending tasks across M phases. Last session:
[topic] ([date]).*

### Rules for recommendations:

- **1-3 items only** — more than 3 defeats the purpose
- **Be specific** — "Fix `block-non-path-ctx` hook" not
  "work on hooks"
- **Include the priority tag** so the user sees the weight
- **Rationale must reference context** — why *this* task, not
  just what it is. Connect to recent work, priority, or
  dependencies
- If an in-progress task exists, it should almost always be
  recommendation #1 (don't abandon unfinished work)

## Examples

### Good Output

> ### Recommended Next
>
> **1. Fix `block-non-path-ctx` hook** `#priority:high`
> > Still open from yesterday's session. The hook is too
> > aggressive — it blocks `git -C path` commands that don't
> > invoke ctx. Quick fix, clears a blocker.
>
> **2. Add `Context.File(name)` method** `#priority:high`
> > Eliminates 10+ linear scan boilerplate instances across
> > 5 packages. High impact, low effort — good consolidation
> > target.
>
> **3. Topics system (T1.1)** `#priority:medium`
> > Journal site's most impactful remaining feature. Metadata
> > is already in place from the enrichment work.
>
> ---
>
> *Based on 24 pending tasks across 3 phases. Last session:
> doc-drift-cleanup (2026-02-11).*

### Bad Output

> "You have many tasks. Here are some options:
> - Do some stuff with hooks
> - Maybe work on tests
> - There's also some docs to write"

(Too vague, no priorities, no rationale, no connection to
context.)

## Quality Checklist

Before presenting recommendations, verify:
- [ ] TASKS.md was read (not guessed from memory)
- [ ] Recent sessions were checked to avoid re-suggesting
      completed work
- [ ] Each recommendation has a specific task reference
- [ ] Each recommendation has a rationale grounded in context
- [ ] In-progress tasks are prioritized over new starts
- [ ] No more than 3 recommendations
- [ ] Footer shows task count and last session reference

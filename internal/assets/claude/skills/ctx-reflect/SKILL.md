---
name: ctx-reflect
description: "Reflect on session progress. Use at natural breakpoints, after unexpected behavior, or when shifting to a different task."
---

Pause and reflect on this session. Review what has been
accomplished and identify context worth persisting.

## When to Use

- At natural breakpoints (feature complete, bug fixed, task
  done)
- After unexpected behavior or a debugging detour
- When shifting from one task to a different one
- When context is getting full and the session may end soon
- When the user explicitly asks to reflect or wrap up

## When NOT to Use

- At the very start of a session (nothing to reflect on yet)
- After trivial changes (a typo fix does not need reflection)
- When the user is in flow and has not paused; do not interrupt
  with unsolicited reflection

## Usage Examples

```text
/ctx-reflect
/ctx-reflect (after fixing the auth bug)
```

## Reflection Checklist

Before listing items, step back and reason through the session
as a whole: what was the arc, what surprised you, what would
you do differently? This framing surfaces insights that a
mechanical checklist misses.

Work through each category. Skip categories with nothing
to report; do not force empty sections.

### 1. Learnings

- Did we discover any gotchas, bugs, or unexpected behavior?
- Did we learn something about the codebase, tools, or
  patterns?
- Would this help a future session avoid problems?
- Is it specific to this project? (General knowledge does not
  belong in LEARNINGS.md)

### 2. Decisions

- Did we make any architectural or design choices?
- Did we choose between alternatives? What was the trade-off?
- Should the rationale be captured for future sessions?

### 3. Tasks

- Did we complete any tasks? (Mark done in TASKS.md)
- Did we start any tasks that are not yet finished?
- Should new tasks be added for follow-up work discovered
  during this session?

### 4. Session Notes

- Was this a significant session worth a full snapshot?
- Would a future session benefit from the discussion context?
- Are there open threads that a future session needs to pick
  up?

## Output Format

After reflecting, provide:

1. **Summary**: what was accomplished (2-3 sentences)
2. **Suggested persists**: list what should be saved, with
   the specific command or file for each item
3. **Offer**: ask the user which items to persist

### Good Example

> This session implemented the cooldown mechanism for
> `ctx agent` and updated all related docs. We discovered
> that `$PPID` in hook context resolves to the Claude Code
> process PID, which is unique per session.
>
> I'd suggest persisting:
> - **Learning**: `$PPID` in PreToolUse hooks resolves to
>   the Claude Code PID (unique per session)
>   `ctx add learning --context "..." --lesson "..." --application "..."`
> - **Task**: mark "Add cooldown to ctx agent" as done
> - **Decision**: tombstone-based cooldown with 10m default
>   `ctx add decision "..."`
>
> Want me to persist any of these?

### Bad Examples

- "We did some stuff. Want me to save it?" (too vague;
  no specific items or commands)
- Listing 10 trivial learnings that are general knowledge
  (only project-specific insights belong)
- Persisting without asking (always get user confirmation)

## Persistence Commands

| What to persist  | Command                                                               |
|------------------|-----------------------------------------------------------------------|
| Learning         | `ctx add learning --context "..." --lesson "..." --application "..."` |
| Decision         | `ctx add decision "..."`                                              |
| Task completed   | Edit TASKS.md directly                                                |
| New task         | `ctx add task "..."`                                                  |

## Quality Checklist

Before presenting the reflection, verify:
- [ ] Every suggested persist has a concrete command or file
      path (not just "save the learning")
- [ ] Learnings are project-specific, not general knowledge
- [ ] Decisions include the trade-off rationale, not just
      the choice
- [ ] No empty checklist categories (skip what has nothing
      to report)
- [ ] The user is asked before anything is persisted

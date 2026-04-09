---
name: ctx-reflect
description: "Reflect on session progress. Use at natural breakpoints, after unexpected behavior, or when shifting to a different task."
tools: [bash, read]
---

Pause and reflect on this session. Review what has been
accomplished and identify context worth persisting.

## When to Use

- At natural breakpoints (feature complete, bug fixed, task done)
- After unexpected behavior or a debugging detour
- When shifting from one task to a different one
- When the session may end soon
- When the user explicitly asks to reflect

## When NOT to Use

- At the very start of a session (nothing to reflect on)
- After trivial changes (a typo fix does not need reflection)
- When the user is in flow: do not interrupt

## Reflection Checklist

Step back and reason through the session as a whole before
listing items.

### 1. Learnings

- Did we discover any gotchas or unexpected behavior?
- Did we learn something about the codebase or tools?
- Would this help a future session avoid problems?
- Is it specific to this project?

### 2. Decisions

- Did we make any architectural or design choices?
- Did we choose between alternatives? What was the trade-off?
- Should the rationale be captured?

### 3. Tasks

- Did we complete any tasks? (Mark done in TASKS.md)
- Did we start any tasks not yet finished?
- Should new tasks be added for follow-up work?

### 4. Session Notes

- Was this a significant session worth a full snapshot?
- Are there open threads a future session needs to pick up?

## Output Format

1. **Summary**: what was accomplished (2-3 sentences)
2. **Suggested persists**: list what should be saved, with
   the specific command for each item
3. **Offer**: ask the user which items to persist

## Persistence Commands

| What to persist  | Command                                                               |
|------------------|-----------------------------------------------------------------------|
| Learning         | `ctx add learning --context "..." --lesson "..." --application "..."` |
| Decision         | `ctx add decision "..."`                                              |
| Task completed   | Edit TASKS.md directly                                                |
| New task         | `ctx add task "..."`                                                  |

## Quality Checklist

- [ ] Every suggested persist has a concrete command
- [ ] Learnings are project-specific, not general knowledge
- [ ] Decisions include trade-off rationale
- [ ] No empty checklist categories
- [ ] User is asked before anything is persisted

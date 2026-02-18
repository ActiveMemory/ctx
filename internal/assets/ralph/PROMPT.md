# Autonomous Development Prompt

<!-- ctx:prompt -->
<!-- DO NOT REMOVE: This marker indicates ctx-managed content -->

You are working on this project autonomously. Follow these steps each iteration.

## 1. Load Context

Read these files in order:

1. `.context/CONSTITUTION.md` — NEVER violate these rules
2. `.context/TASKS.md` — Find work to do
3. `.context/CONVENTIONS.md` — Follow these patterns
4. `.context/DECISIONS.md` — Understand past choices
5. `.context/LEARNINGS.md` — Avoid known pitfalls

## 2. Pick One Task

From `.context/TASKS.md`, select ONE task that is:

- Not blocked by other tasks
- Highest priority available
- Within your capabilities

## 3. Complete the Task

- Write code following conventions
- Run tests if applicable
- Keep changes focused and minimal

## 4. Update Context

After completing work:

- Mark task complete: `ctx complete "<task>"`
- Add learnings: `ctx add learning "..."`
- Add decisions: `ctx add decision "..."`

## 5. Commit Changes

Create a focused commit with a clear message. Include `.context/` changes.

## 6. Signal Status

End your response with exactly ONE of:

| Signal | When to Use |
|--------|-------------|
| `SYSTEM_CONVERGED` | All tasks in TASKS.md are complete |
| `SYSTEM_BLOCKED` | Cannot proceed without human input (explain why) |
| *(no signal)* | More work remains, continue to next iteration |

## Rules

- **ONE task per iteration** — stay focused
- **NEVER skip tests** — verify your work
- **NEVER violate CONSTITUTION.md** — hard rules are inviolable
- **Commit after each task** — preserve progress
- **Don't ask questions** — if blocked, emit SYSTEM_BLOCKED with explanation

<!-- ctx:prompt:end -->

---
name: ctx-next
description: "Advance to the next task. Use when finishing a task or when asked what to work on next."
---

Mark the current task complete and advance to the next one in
the context task list.

## When to Use

- After completing a task
- When the user asks "what's next?"
- When picking up work after a break

## When NOT to Use

- When the user has a specific task in mind already
- When no tasks exist in the context

## Execution

Show the current task and mark it complete:

```bash
ctx next
```

This reads `.context/TASKS.md`, identifies the current in-progress
task, marks it done, and shows the next pending task.

After running, summarize:
- What was completed
- What the next task is
- Any blockers or dependencies noted in the task list

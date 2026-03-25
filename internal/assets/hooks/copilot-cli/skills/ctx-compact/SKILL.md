---
name: ctx-compact
description: "Archive completed tasks and trim context. Use when context files are growing large."
---

Archive completed tasks and trim stale context entries to keep
the context directory lean and within token budgets.

## When to Use

- When TASKS.md has many completed items
- When context token count is growing large
- When asked to "clean up" or "compact" context
- Before starting a new phase of work

## When NOT to Use

- When all tasks are still active
- When context is already compact

## Execution

Run the compact operation:

```bash
ctx compact
```

This archives completed tasks from TASKS.md into the session
history and trims stale entries from other context files.

After running, confirm:
- How many tasks were archived
- Current token budget usage
- Whether any manual cleanup is recommended

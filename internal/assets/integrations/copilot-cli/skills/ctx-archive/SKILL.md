---
name: ctx-archive
description: "Archive completed tasks. Use when TASKS.md has many completed items cluttering the view."
tools: [bash]
---

Move completed tasks from TASKS.md to the archive.

## Before Archiving

1. **"Are the completed tasks cluttering the view?"** → If TASKS.md is
   still easy to scan, there's no urgency
2. **"Are all `[x]` items truly done?"** → Verify nothing was checked
   off prematurely

## When to Use

- When TASKS.md has many completed `[x]` tasks
- When the task list is hard to navigate
- Periodically to keep context clean

## When NOT to Use

- When there are only a few completed tasks
- When you're unsure if tasks are truly complete (verify first)
- **Never delete tasks**: only archive (CONSTITUTION invariant)

## Execution

```bash
# Preview first (recommended)
ctx tasks archive --dry-run

# Archive after confirming the preview
ctx tasks archive
```

Archived tasks go to `archive/tasks-YYYY-MM-DD.md` in the context
directory, preserving Phase headers for traceability.

## Quality Checklist

- [ ] Previewed with `--dry-run` before archiving
- [ ] All archived items are truly complete
- [ ] No tasks were deleted (only archived)
- [ ] Reported how many tasks were archived

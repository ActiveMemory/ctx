---
name: ctx-archive
description: "Archive completed tasks. Use when TASKS.md has many completed items cluttering the view."
allowed-tools: Bash(ctx:*)
---

Move completed tasks from TASKS.md to the archive.

## Before Archiving

Two questions — if any answer is "no", don't archive:

1. **"Are the completed tasks cluttering the view?"** → If TASKS.md is
   still easy to scan, there's no urgency
2. **"Are all `[x]` items truly done?"** → Verify nothing was checked off
   prematurely

## When to Use

- When TASKS.md has many completed `[x]` tasks
- When the task list is hard to navigate
- Periodically to keep context clean

## When NOT to Use

- When there are only a few completed tasks (not worth the noise)
- When you're unsure if tasks are truly complete (verify first)
- **Never delete tasks** — only archive (CONSTITUTION invariant)

## Constitution Rules

These are inviolable:

- **Archival is allowed, deletion is not** — never delete context history
- **Archive preserves structure** — Phase headers are kept for traceability
- **Never move tasks** — tasks stay in their Phase section; archiving is
  the only sanctioned "move" and it's to `.context/archive/`

## Execution

```bash
ctx tasks archive $ARGUMENTS
```

**Example — preview first (recommended):**
```bash
ctx tasks archive --dry-run
```

**Example — archive after confirming the preview:**
```bash
ctx tasks archive
```

Archived tasks go to `.context/archive/tasks-YYYY-MM-DD.md`, preserving
Phase headers for traceability.

Report how many tasks were archived and where the archive file was written.

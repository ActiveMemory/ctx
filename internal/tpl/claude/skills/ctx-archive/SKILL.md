---
name: ctx-archive
description: "Archive completed tasks. Use when TASKS.md has many completed items cluttering the view."
allowed-tools: Bash(ctx:*)
---

Move completed tasks from TASKS.md to the archive.

## When to Use

- When TASKS.md has many completed `[x]` tasks
- When the task list is hard to navigate
- Periodically to keep context clean

## Usage

```
/ctx-archive
/ctx-archive --dry-run
```

Use `--dry-run` to preview what would be archived.

## Execution

```bash
ctx tasks archive $ARGUMENTS
```

Report how many tasks were archived.

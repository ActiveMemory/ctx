---
description: "Archive completed tasks"
argument-hint: "[--dry-run]"
---

Move completed tasks from TASKS.md to the archive.

**Usage:**
```
/ctx-archive
/ctx-archive --dry-run
```

Moves completed `[x]` tasks to `.context/archive/`. Use `--dry-run` to preview.

```!
ctx tasks archive $ARGUMENTS
```

Report how many tasks were archived.

---
name: ctx-import-plans
description: "Import Claude Code plan files into project specs. Use when plan files in ~/.claude/plans/ should become permanent project specs."
allowed-tools: Bash(ls:*), Bash(stat:*), Read, Write
---

Import Claude Code plan files (`~/.claude/plans/*.md`) into the
project's `specs/` directory so they become part of project memory.

## When to Use

- User says "import plans", "save that plan", "keep the plan"
- User wants to preserve a Claude Code plan as a project spec
- After a planning session produced a plan worth keeping
- User asks to review or archive recent plans

## When NOT to Use

- User wants to create a new spec from scratch (use `/_ctx-spec`)
- User wants to edit an existing spec directly
- No `~/.claude/plans/` directory or it's empty

## Process

### 1. Discover Plans

List plan files with modification dates:

```bash
ls -lt ~/.claude/plans/*.md 2>/dev/null
```

If no files are found, tell the user and stop.

### 2. Filter by Arguments

The user may pass arguments to narrow the selection:

| Argument | Behavior |
|---|---|
| `--today` | Only plans modified today |
| `--since YYYY-MM-DD` | Only plans modified on or after the given date |
| `--all` | Import all plans without prompting |
| *(none)* | Interactive — present the list and ask the user to pick |

**Filtering with `--today`:**
```bash
find ~/.claude/plans/ -name '*.md' -newermt "$(date +%Y-%m-%d)" -type f
```

**Filtering with `--since`:**
```bash
find ~/.claude/plans/ -name '*.md' -newermt "YYYY-MM-DD" -type f
```

### 3. Present for Selection (Interactive Mode)

For each plan file, extract the first H1 heading and show it with
the modification date:

```
1. 2026-02-28  Add authentication middleware
2. 2026-02-27  Refactor database connection pool
3. 2026-02-25  Import plans skill
```

Ask the user which plans to import (comma-separated numbers, or "all").

### 4. Import Each Selected Plan

For each selected plan:

1. **Read the file** to extract the H1 heading (first `# ` line)
2. **Slugify the heading** for the filename:
   - Lowercase
   - Replace spaces and non-alphanumeric characters with hyphens
   - Collapse multiple hyphens
   - Trim leading/trailing hyphens
   - Example: `Add Authentication Middleware` → `add-authentication-middleware`
3. **Check for conflicts** — if `specs/{slug}.md` already exists, ask
   the user whether to overwrite or pick a different name
4. **Copy the file** to `specs/{slug}.md`
5. **Optionally add a task** — ask the user if they want a task in
   TASKS.md referencing the imported spec (use `/ctx-add-task` if yes)

### 5. Report

After importing, summarize what was done:

```
Imported 2 plan(s):
  ~/.claude/plans/abc123.md → specs/add-authentication-middleware.md
  ~/.claude/plans/def456.md → specs/refactor-database-pool.md
```

## Important Notes

- Plan filenames in `~/.claude/plans/` are typically UUIDs or hashes —
  always use the H1 heading for the spec filename, not the original name
- If a plan has no H1 heading, use the original filename (minus extension)
  as the slug
- Do not modify the original plan files — this is a copy, not a move
- The `specs/` directory must exist (it should already be present in
  the project root)

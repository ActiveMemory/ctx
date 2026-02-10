---
name: ctx-journal-enrich-all
description: "Batch-enrich all unenriched journal entries. Use after exporting sessions to process the backlog without manual selection."
allowed-tools: Bash(ctx:*), Read, Glob, Grep, Edit, Write, Task
---

Batch-enrich all unenriched journal entries automatically.

## When to Use

- After `ctx recall export --all` produces a batch of new entries
- When there is a backlog of unenriched sessions
- When the user says "enrich everything" or "process the journal"
- Periodically to catch up on recent sessions

## When NOT to Use

- For a single specific session (use `/ctx-journal-enrich` instead)
- Before exporting (nothing to enrich yet)
- Before normalizing (run `/ctx-journal-normalize` first)

## Process

### Step 1: Ensure Normalization

Check whether `/ctx-journal-normalize` has been run. Look for
entries missing the `<!-- normalized -->` marker:

```bash
grep -rL "normalized" .context/journal/*.md 2>/dev/null | wc -l
```

If there are unnormalized entries, run normalization first.
Clean markdown produces better metadata extraction.

### Step 2: Find Unenriched Entries

List all journal entries that lack YAML frontmatter:

```bash
grep -rL "^---$" .context/journal/*.md 2>/dev/null
```

If all entries already have frontmatter, report that and stop.

### Step 3: Filter Out Noise

Skip entries that are not worth enriching:

- **Suggestion sessions**: files under ~20 lines or containing
  only auto-complete fragments. Check with:
  ```bash
  wc -l <file>
  ```
- **Multi-part continuations**: files ending in `-p2.md`, `-p3.md`
  etc. Enrich only the first part; continuation parts inherit
  the frontmatter topic.

Report how many entries will be processed and how many were
filtered out.

### Step 4: Process Each Entry

For each entry, read the conversation and extract:

1. **Title**: a short descriptive title for the session
2. **Type**: feature, bugfix, refactor, exploration, debugging,
   or documentation
3. **Outcome**: completed, partial, abandoned, or blocked
4. **Topics**: 2-5 topic tags
5. **Technologies**: languages, frameworks, tools used
6. **Summary**: 2-3 sentences describing what was accomplished

Apply YAML frontmatter to each file:

```yaml
---
title: "Session title"
date: 2026-01-27
type: feature
outcome: completed
topics:
  - authentication
  - caching
technologies:
  - go
  - redis
---
```

### Step 5: Report

After processing, report:

- How many entries were enriched
- How many were skipped (already enriched, too short, etc.)
- Remind the user to rebuild: `ctx journal site --build`

## Confirmation Mode

**Interactive** (default when user is present): show a summary
of proposed enrichments before applying. Group by type/outcome
so the user can scan quickly rather than reviewing one by one.

**Unattended** (when running in a loop or explicitly told
"just do it"): apply enrichments directly and report results.

## Large Backlogs

For backlogs of 20+ entries, consider spawning subagents to
process entries in parallel. Each subagent handles a batch of
5-10 entries. The parent agent tracks progress via a task list.

This is optional — sequential processing works fine for smaller
backlogs and avoids coordination overhead.

## Quality Checklist

- [ ] Normalization was run before enrichment
- [ ] Suggestion sessions and multi-part continuations filtered
- [ ] Each enriched entry has all required frontmatter fields
- [ ] Summary is specific to the session, not generic
- [ ] User was shown a summary before applying (unless unattended)

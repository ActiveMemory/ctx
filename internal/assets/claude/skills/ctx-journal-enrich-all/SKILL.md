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

## Process

### Step 1: Find Unenriched Entries

List all journal entries that lack enrichment using the state file:

```bash
# List .md files in journal dir and check state
for f in .context/journal/*.md; do
  name=$(basename "$f")
  ctx journal mark --check "$name" enriched || echo "$f"
done
```

Or read `.context/journal/.state.json` directly and list entries
without an `enriched` date set.

If all entries already have enrichment recorded, report that and stop.

### Step 2: Filter Out Noise

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

### Step 3: Process Each Entry

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

### Step 4: Mark Enriched

After writing frontmatter to each file, update the state file:

```bash
ctx journal mark <filename> enriched
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

- [ ] Suggestion sessions and multi-part continuations filtered
- [ ] Each enriched entry has all required frontmatter fields
- [ ] Summary is specific to the session, not generic
- [ ] User was shown a summary before applying (unless unattended)
- [ ] State file updated for each enriched entry

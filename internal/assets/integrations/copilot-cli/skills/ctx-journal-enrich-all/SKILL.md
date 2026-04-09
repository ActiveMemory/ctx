---
name: ctx-journal-enrich-all
description: "Batch journal pipeline: export unexported sessions then enrich all unenriched entries."
tools: [bash, read, write, edit]
---

Full journal pipeline: import sessions and batch-enrich entries.

## When to Use

- Backlog of unenriched journal entries
- After many sessions without journal maintenance
- When running periodic journal housekeeping

## When NOT to Use

- No journal entries exist
- All entries are already enriched
- Single entry (use `ctx-journal-enrich` instead)

## Process

### 1. Import unexported sessions

```bash
ctx recall export --all
```

### 2. List unenriched entries

```bash
ctx journal list --unenriched
```

### 3. Batch enrich

For each unenriched entry:
1. Read the entry content
2. Generate appropriate frontmatter
3. Write the enriched version
4. Report progress

For large backlogs (20+ entries), use heuristic enrichment:
derive metadata from filename patterns and entry headings
without reading full content.

### 4. Report

```
Enriched: 15/15 entries
Skipped: 3 (already enriched)
```

## Quality Checklist

- [ ] All unexported sessions imported first
- [ ] Each enriched entry has valid frontmatter
- [ ] Progress reported during batch
- [ ] No entries corrupted or lost

---
name: ctx-journal-enrich
description: "Enrich a journal entry with YAML frontmatter metadata. Use to add type, outcome, topics, and technologies to session records."
tools: [bash, read, write, edit]
---

Enrich individual journal entries with structured metadata.

## When to Use

- After exporting a session to the journal
- When journal entries lack metadata for search/filter
- When `ctx journal` shows unenriched entries

## When NOT to Use

- Entry is already fully enriched
- No journal entries exist

## Process

### 1. Identify the entry

If not specified, find unenriched entries:

```bash
ctx journal list --unenriched
```

### 2. Read the entry

Read the full session content to understand what happened.

### 3. Generate frontmatter

Add or update YAML frontmatter with:

```yaml
---
type: feature|bugfix|refactor|research|planning|review
outcome: completed|partial|blocked|abandoned
topics: [topic1, topic2]
technologies: [go, typescript, ...]
summary: "One-line summary of the session"
---
```

### 4. Write enriched entry

Update the file with the new frontmatter while preserving
the body content.

## Quality Checklist

- [ ] Frontmatter is valid YAML
- [ ] Type matches the actual work done
- [ ] Outcome is accurate
- [ ] Topics are specific, not generic
- [ ] Summary is one clear sentence
- [ ] Body content is preserved unchanged

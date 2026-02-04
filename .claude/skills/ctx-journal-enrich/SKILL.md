---
name: ctx-journal-enrich
description: "Enrich journal entry with metadata. Use when journal entries lack frontmatter, tags, or summary for future reference."
---

Enrich a session journal entry with structured metadata.

## When to Use

- When journal entries lack metadata for future reference
- After exporting sessions that need categorization
- When building a searchable session archive

## Input

The user specifies a journal entry by partial match:
- `twinkly-stirring-kettle` (slug)
- `twinkly` (partial slug)
- `2026-01-24` (date)
- `76fe2ab9` (short ID)

Find matching files:
```bash
ls .context/journal/*.md | grep -i "<pattern>"
```

If multiple matches, show them and ask which one.
If no argument given, show recent entries and ask.

## Enrichment Tasks

Read the journal entry and extract:

### 1. Frontmatter (YAML at top of file)
```yaml
---
title: "Session title"
date: 2026-01-27
type: feature|bugfix|refactor|exploration|debugging|documentation
outcome: completed|partial|abandoned|blocked
topics:
  - authentication
  - caching
technologies:
  - go
  - postgresql
libraries:
  - cobra
  - fatih/color
key_files:
  - internal/auth/token.go
  - internal/db/cache.go
---
```

### 2. Summary
If `## Summary` says "[Add your summary...]", replace with 2-3 sentences.

### 3. Extracted Items
Scan the conversation and extract:

**Decisions made** - Link to DECISIONS.md if persisted:
```markdown
## Decisions
- Used Redis for caching ([D12](../DECISIONS.md#d12))
- Chose JWT over sessions (not yet persisted)
```

**Learnings discovered** - Link to LEARNINGS.md if persisted:
```markdown
## Learnings
- Token refresh requires cache invalidation ([L8](../LEARNINGS.md#l8))
- Go's defer runs LIFO (new insight)
```

**Tasks completed/created**:
```markdown
## Tasks
- [x] Implement caching layer
- [ ] Add cache metrics (created this session)
```

## Process

1. Find and read the journal file
2. Analyze the conversation
3. Propose enrichment (type, topics, outcome)
4. Ask user for confirmation/adjustments
5. Show diff and write if approved

---
name: ctx-journal-enrich
description: "Enrich journal entry with metadata. Use when journal entries lack frontmatter, tags, or summary for future reference."
---

Enrich a session journal entry with structured metadata.

## Before Enriching

1. **Check if already enriched**: check the state file via
   `ctx system mark-journal --check <filename> enriched` or read
   `.state.json` in the journal directory; confirm before overwriting

## When to Use

- When journal entries lack metadata for future reference
- After exporting sessions that need categorization
- When building a searchable session archive

## When NOT to Use

- On entries that already have complete frontmatter (unless updating)
- Before normalizing entries with broken formatting
- On suggestion sessions (short auto-complete prompts; not worth enriching)

## Input

The user specifies a journal entry by partial match:
- `twinkly-stirring-kettle` (slug)
- `twinkly` (partial slug)
- `2026-01-24` (date)
- `76fe2ab9` (short ID)

Find matching files in the journal directory:
```bash
ls "$(ctx system bootstrap -q)/journal/"*.md | grep -i "<pattern>"
```

If multiple matches, show them and ask which one.
If no argument given, show recent unenriched entries by reading
`.state.json` in the journal directory and listing entries without
an `enriched` date:

```bash
# List unenriched entries using state file
CTX_DIR=$(ctx system bootstrap -q)
for f in "$CTX_DIR/journal/"*.md; do
  name=$(basename "$f")
  ctx system mark-journal --check "$name" enriched || echo "$f"
done | head -10
```

## Usage Examples

```text
/ctx-journal-enrich twinkly-stirring-kettle
/ctx-journal-enrich twinkly
/ctx-journal-enrich 2026-01-24
/ctx-journal-enrich 76fe2ab9
```

## Enrichment Tasks

Read the journal entry and extract:

### 1. Frontmatter (YAML at top of file)

```yaml
---
title: "Session title"
date: 2026-01-27
model: claude-opus-4-6        # auto-populated at export
tokens_in: 234000             # auto-populated at export
tokens_out: 89000             # auto-populated at export
type: feature
outcome: completed
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

**Auto-populated fields** (set during `ctx journal export`, do NOT overwrite):
`date`, `time`, `project`, `session_id`, `model`, `tokens_in`, `tokens_out`, `branch`

**Type values:**

| Type            | When to use                           |
|-----------------|---------------------------------------|
| `feature`       | Building new functionality            |
| `bugfix`        | Fixing broken behavior                |
| `refactor`      | Restructuring without behavior change |
| `exploration`   | Research, learning, experimentation   |
| `debugging`     | Investigating issues                  |
| `documentation` | Writing docs, comments, README        |

**Outcome values:**

| Outcome     | Meaning                            |
|-------------|------------------------------------|
| `completed` | Goal achieved                      |
| `partial`   | Some progress, work continues      |
| `abandoned` | Stopped pursuing this approach     |
| `blocked`   | Waiting on external dependency     |

### 2. Summary

If `## Summary` says "[Add your summary...]", replace with 2-3 sentences
describing what was accomplished.

### 3. Extracted Items

Scan the conversation and extract:

**Decisions made**: link to DECISIONS.md if persisted:
```markdown
## Decisions
- Used Redis for caching ([D12](../DECISIONS.md#d12))
- Chose JWT over sessions (not yet persisted)
```

**Learnings discovered**: link to LEARNINGS.md if persisted:
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
6. **Mark enriched** in the state file:
   ```bash
   ctx system mark-journal <filename> enriched
   ```
7. Remind user to rebuild: `ctx journal site --build` or `make journal`

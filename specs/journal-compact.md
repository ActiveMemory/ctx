# Spec: Journal Compaction

## Problem

The journal directory (`.context/journal/`) grows monotonically.
At current rates it adds ~8-11 MB/week (~400-500 MB/year). While
not git-controlled, unbounded growth creates three problems:

1. **Disk budget** — manageable today (76 MB) but compounds
2. **Noise ratio** — ~60-70% of each entry is tool call output
   with zero recall value
3. **Re-import explosion** — `ctx journal import --all` regenerates
   from JSONL sources, so compaction gains are lost unless the
   pipeline knows which entries were compacted

The journal markdown files are derived artifacts — the JSONL source
files in `~/.claude/projects/` are the originals. Compaction is a
lossy view optimization, not data destruction.

## Design: Elastic-Style Tiered Storage

Inspired by Elasticsearch's hot/warm/cold/frozen index lifecycle:

| Tier     | Age           | Content                                          | File format         |
|----------|---------------|--------------------------------------------------|---------------------|
| **Hot**  | 0–N days      | Full markdown, all tool output, code blocks      | Individual `.md`    |
| **Warm** | N–M days      | Frontmatter + summary + key outcomes, body stripped | Individual `.md`  |
| **Cold** | Not in v1     | Monthly rollup, one section per session           | `YYYY-MM-rollup.md` |

Default thresholds: N=30, M=90 (configurable via `.ctxrc`).
V1 implements hot → warm only. Cold tier is a future extension.

## Solution: `ctx journal compact`

### CLI Surface

```
ctx journal compact [flags]

Flags:
  --older-than duration   Compact entries older than this (default: 30d)
  --backup                Create timestamped tar.gz before compacting (default: true)
  --backup-dir string     Directory for backup archives (default: .context/archive/journal/)
  --dry-run               Preview what would be compacted and space savings
  --force                 Compact even locked entries (default: false, locked entries skipped)
```

### Behavior

1. **Backup phase** (unless `--no-backup`):
   - Create `journal-YYYYMMDD-HHMMSS.tar.gz` containing all
     `.context/journal/*.md` files at their current state
   - Store in `--backup-dir` (default `.context/archive/journal/`)
   - Print: `Backup: .context/archive/journal/journal-20260331-173000.tar.gz (76 MB)`

2. **Selection phase**:
   - Scan all `.md` files in `.context/journal/`
   - Parse date from filename (`YYYY-MM-DD-*`)
   - Select entries where `today - entry_date > --older-than`
   - Skip entries with `locked: true` in frontmatter (unless `--force`)
   - Skip entries with `compacted: true` in frontmatter (already done)
   - Skip multi-part entries (`-p2.md`, `-p3.md`) — compact the
     parent and leave a redirect stub pointing to it

3. **Compact phase** — for each selected entry:
   - Parse frontmatter (preserve all fields)
   - Add `compacted: true` and `compacted_date: YYYY-MM-DD` to frontmatter
   - Keep only the frontmatter + summary field content as the body
   - If the entry has no `summary` in frontmatter (unenriched):
     skip it and warn — enrichment should run first
   - Write the compacted file in place

4. **State update**:
   - Add `"compacted": "YYYY-MM-DD"` to the entry in `.state.json`

5. **Report**:
   ```
   Compacted 245 of 876 entries (older than 30 days)
   Skipped: 12 locked, 3 unenriched, 8 already compacted
   Space: 52.3 MB → 4.1 MB (saved 48.2 MB, 92% reduction)
   Backup: .context/archive/journal/journal-20260331-173000.tar.gz
   ```

### Compacted Entry Format

Before (hot — 89 KB):
```markdown
---
date: "2026-01-21"
title: "Repository Migration"
summary: "Updated all references from old to new org..."
topics: [repository-migration, go-module-path]
# ... rest of frontmatter
---

# Repository Migration

<details>
<summary>2026-01-21 · 2h5m · claude-opus-4-5-20251101</summary>
... 2000 lines of tool calls, file reads, edits ...
</details>
```

After (warm — ~1 KB):
```markdown
---
date: "2026-01-21"
title: "Repository Migration"
summary: "Updated all references from old to new org..."
topics: [repository-migration, go-module-path]
compacted: true
compacted_date: "2026-03-31"
# ... rest of frontmatter preserved
---

# Repository Migration

Updated all references from old to new org across go.mod, all Go
source files, documentation, specs, and context files. A 2-hour
session involving 22 files updated across the codebase.
```

### Multi-Part Entry Handling

Multi-part entries (`*-p2.md`, `*-p3.md`) are reduced to stubs
pointing to the compacted parent:

```markdown
---
date: "2026-01-21"
title: "Repository Migration (Part 2)"
compacted: true
compacted_date: "2026-03-31"
redirect: "2026-01-21-repository-migration-e3512a48.md"
---

Compacted into parent entry.
```

## Import Interaction

The critical design question: what happens when a user runs
`ctx journal import --all` after compaction?

### Rule: compacted entries are import-immune

Import checks `.state.json` for `"compacted"` and skips those
entries, same as it skips locked entries. The JSONL source is
still available — compaction is a journal-side view decision, not
a source-side deletion.

### Escape hatch: `--force` on import

`ctx journal import --force <slug>` re-generates from JSONL source,
overwriting the compacted version. This clears `compacted` from
`.state.json` and removes `compacted: true` from frontmatter.

Full re-import: `ctx journal import --all --force` regenerates
everything (respecting locks unless `--force-locked` is also set).
This is the "restore from source" path — equivalent to restoring
from the tar.gz backup but using the canonical JSONL sources.

### State precedence

```
locked > compacted > enriched > exported
```

A locked entry is never compacted (unless `--force`).
A compacted entry is never re-imported (unless `--force`).

## Backup Strategy

Since `.context/journal/` is not source-controlled, the tar.gz
backup is the safety net:

- **Created automatically** before each compact run (opt-out
  via `--no-backup`)
- **Timestamped** to allow multiple backups: `journal-YYYYMMDD-HHMMSS.tar.gz`
- **Contains full state** at time of backup: all `.md` files +
  `.state.json`
- **Stored in** `.context/archive/journal/` (already gitignored
  via `.context/archive/`)
- **No automatic pruning** of backups — the user decides when to
  delete old tar.gz files (at ~76 MB each, even 10 backups is
  under 1 GB)

## `.ctxrc` Configuration

```yaml
journal_compact_age: "30d"        # --older-than default
journal_compact_backup: true      # --backup default
journal_compact_backup_dir: ""    # empty = .context/archive/journal/
```

## Non-Goals

- **Cold tier (monthly rollups)**: deferred to a future spec.
  Warm compaction alone delivers ~90% space savings.
- **Automatic compaction on import**: compaction is an explicit
  user action, not a side effect of import
- **JSONL source compaction**: the source files belong to Claude
  Code, not ctx. Don't touch them.
- **Content-aware summarization**: the `summary` field from
  enrichment is the summary. No LLM calls during compaction —
  this is a fast, offline operation.

## Implementation Packages

| Package | What to add |
|---------|-------------|
| `internal/cli/journal/cmd/compact/` | `cmd.go`, `run.go` — CLI wiring |
| `internal/cli/journal/core/` | `compact.go` — core logic (select, transform, backup) |
| `internal/journal/state/` | Add `Compacted` field to entry state, skip logic |
| `internal/cli/journal/cmd/import/` | Respect `compacted` state, add `--force` flag |

## Error Cases

- **No entries qualify**: print "No entries older than 30 days to
  compact" and exit 0
- **All qualifying entries locked**: print count and exit 0
- **Unenriched entries**: skip with warning, suggest running
  enrichment first
- **Backup fails** (disk full, permissions): abort before any
  compaction, exit 1
- **Mid-compact failure**: entries already written are compacted
  (idempotent), remaining are untouched. Re-run is safe.

## Verification

- `ctx journal compact --dry-run` previews changes
- After compaction: `ctx journal site --serve` still renders
  (compacted entries show summary only)
- After compaction: `ctx journal import --all` skips compacted
  entries
- After compaction + `import --force`: full bodies are restored

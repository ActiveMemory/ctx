# Knowledge Scaling: Decisions & Learnings Archival

## Problem

`DECISIONS.md` and `LEARNINGS.md` grow monotonically. Tasks have a full
lifecycle (`ctx tasks archive`, `ctx compact --archive`), but decisions and
learnings accumulate forever with no archival or compaction path.

For a long-lived project, this creates two problems:

1. **Token budget pressure**: `ctx agent --budget` loads files in priority
   order and truncates when the budget is exhausted. Reverse-chronological
   ordering means old entries are cut first, but the index table at the top
   grows linearly and eventually consumes significant budget on its own.

2. **Signal-to-noise decay**: Early decisions and learnings may be obsolete,
   superseded, or internalized into conventions. They dilute attention from
   current, actionable entries.

## Current Mitigations

| Mechanism               | Helps with         | Limitation                         |
|-------------------------|--------------------|------------------------------------|
| Reverse-chronological   | Budget truncation  | Index still grows linearly         |
| Token budget truncation | Loading cost       | Loses old entries silently         |
| Learning deduplication  | Near-duplicates    | Only catches similar content       |
| Reindex commands        | Index freshness    | Doesn't reduce entry count         |

## Design

### 1. Archive Commands

Follow the existing task archive pattern:

```
ctx decisions archive [flags]
ctx learnings archive [flags]
```

**Behavior**:

1. Parse entries from the file (reuse existing entry parser).
2. Select entries older than threshold (default: 90 days, configurable via
   `.contextrc` key `archive_knowledge_after_days`).
3. Write selected entries to `.context/archive/decisions-YYYY-MM-DD.md`
   (or `learnings-YYYY-MM-DD.md`).
4. Remove archived entries from the source file.
5. Rebuild the index (call existing reindex logic).
6. Print summary: `Archived N decisions (M remaining).`

**Flags**:

| Flag         | Short | Default | Description                              |
|--------------|-------|---------|------------------------------------------|
| `--days`     | `-d`  | 90      | Archive entries older than N days         |
| `--dry-run`  |       | false   | Print what would be archived, don't write |
| `--all`      |       | false   | Archive all entries (ignores --days)      |
| `--keep`     | `-k`  | 0       | Always keep the N most recent entries     |

**Safety**:

- Constitution says "archival is allowed, deletion is not" — this moves
  entries to archive, never deletes them.
- `--dry-run` is the default recommendation in docs.
- Archived files are append-only (multiple runs on same day append, not
  overwrite).

### 2. Extend `ctx compact`

Add decisions/learnings archival to the existing compact flow:

```bash
ctx compact --archive  # now also archives old decisions and learnings
```

When `--archive` is set, compact should:
1. Archive old tasks (existing behavior).
2. Archive old decisions (new, same threshold).
3. Archive old learnings (new, same threshold).
4. Deduplicate learnings (existing behavior).
5. Reindex both files.

### 3. Superseded Entries

Add a `superseded_by` marker for decisions:

```markdown
## [2026-01-15-120000] Use JWT for auth
~~Superseded by [2026-02-10-090000] Switch to session cookies~~
```

Entries marked as superseded are archived immediately by `compact --archive`
regardless of age. The marker is a convention (documented in CONVENTIONS.md),
not enforced by code — the archive command checks for `~~Superseded` prefix
on the entry title.

### 4. Configuration

New `.contextrc` keys:

```yaml
archive_knowledge_after_days: 90    # default threshold for decisions/learnings
archive_keep_recent: 5              # always keep N most recent entries
```

## Non-Goals

- **Semantic compaction** (merging related decisions into summaries): too
  complex for v1, would require LLM calls. Revisit after archive proves
  useful.
- **Tiered loading** (index-only mode with on-demand entry fetch): would
  require changes to `ctx agent` output format. Consider as a separate spec
  if archival alone doesn't solve the scaling problem.
- **Automatic archival** (archive on every commit or session end): too
  aggressive. Keep it manual or via `compact`.

## Implementation Order

1. Entry parser extraction (reusable for both decisions and learnings).
2. `ctx decisions archive` command.
3. `ctx learnings archive` command.
4. Extend `ctx compact --archive` to include knowledge files.
5. Superseded marker convention + detection.
6. `.contextrc` configuration keys.
7. Documentation: cli-reference, context-files, recipes/context-health.
8. Update `/consolidate` skill to suggest archival when files are large.

## Testing

- Unit: entry parser, age filtering, archive file writing, reindex after
  archive.
- Integration: full compact cycle with all three file types.
- Edge cases: empty files, no entries old enough, all entries archived,
  superseded entries, multiple runs on same day.

# Knowledge Nudge: Replace Archive with Human Judgment

## Problem

The `ctx decisions archive` and `ctx learnings archive` commands use age as a
proxy for relevance — a fundamentally flawed heuristic. A 90-day-old decision
may be the project's most important architectural choice. Mechanically sweeping
entries into `.context/archive/` creates "out of sight, out of mind" context
loss, which is antithetical to ctx's mission of persistent, accessible context.

## Decision

Replace mechanical knowledge archival with a nudge-based approach:

1. **Detect** — mechanical file/entry checks (entry count thresholds)
2. **Nudge** — VERBATIM relay to the user with actionable suggestions
3. **Decide** — human or agent makes the call

## What's Removed

- `ctx decisions archive` command
- `ctx learnings archive` command
- Knowledge archival from `ctx compact --archive`
- `.ctxrc` keys: `archive_knowledge_after_days`, `archive_keep_recent`
- `RemoveEntryBlocks()` and `OlderThan()` from index package (archive-only)

## What's Kept

- `ctx tasks archive` — task archival is lifecycle-based, not age-based
- `ctx compact` — still archives tasks and removes empty sections
- Entry count drift warnings (`entry_count_learnings`, `entry_count_decisions`)
- Convention line count warning (`convention_line_count`)
- `ParseEntryBlocks()`, `IsSuperseded()`, `BlockContent()` — used by agent
  scoring and drift detection
- `/ctx-consolidate` skill — human-guided merging of overlapping entries

## New Hook: `ctx system check-knowledge`

A UserPromptSubmit hook that nudges when knowledge files exceed entry count
thresholds.

### Design

1. Guard: `isInitialized()`, `isDailyThrottled(markerPath)`
2. Load `.context/DECISIONS.md` and `.context/LEARNINGS.md`
3. Count entries via `index.ParseEntryBlocks()`
4. Load `.context/CONVENTIONS.md` and count lines
5. Compare against thresholds:
   - `rc.EntryCountLearnings()` (default 30)
   - `rc.EntryCountDecisions()` (default 20)
   - `rc.ConventionLineCount()` (default 200)
6. If any exceeds threshold, emit VERBATIM relay with specific counts
7. Touch daily marker, return

### Output Format

```
IMPORTANT: Relay this knowledge health notice to the user VERBATIM before answering their question.

┌─ Knowledge File Growth ──────────────────────────────────
│ LEARNINGS.md has 45 entries (recommended: ≤30).
│
│ Large knowledge files dilute agent context. Consider:
│  • Review and remove outdated entries
│  • Use /ctx-consolidate to merge overlapping entries
│  • Use /ctx-drift for semantic drift (stale patterns)
│  • Move stale entries to .context/archive/ manually
└──────────────────────────────────────────────────────────
```

### Throttling

Daily throttle via marker file in `secureTempDir()`. Same pattern as other
hook commands.

## Migration

- Existing `.context/archive/` files are untouched
- Users who relied on `ctx decisions archive` should use `/ctx-consolidate`
  for human-guided consolidation or manually move entries
- The `archive_knowledge_after_days` and `archive_keep_recent` .ctxrc
  keys are silently ignored (removed from struct)

## Testing

- Below-threshold: silent
- Above-threshold: VERBATIM output with correct counts and units
- Daily throttle: second call same day is silent
- One file over, one under: only the over-threshold file mentioned
- Conventions over line count: mentions CONVENTIONS.md with "lines" unit
- Conventions below threshold: silent
- Uninitialized: silent

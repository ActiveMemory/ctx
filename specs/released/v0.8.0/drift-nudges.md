# Drift Nudges: Entry Count Warnings

Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 2)

## Problem

Context files grow without any feedback signal. A project with 47 learnings
gets the same `ctx drift` output as one with 5. Users don't know when files
are getting unwieldy until the agent packet starts degrading.

Phase 1 (smart retrieval) mitigates the agent-side impact by scoring entries,
but the files themselves still grow unbounded. Users need a nudge — not
enforcement — to consolidate or archive before things get noisy.

## Solution

Add an `entry_count` check to `ctx drift` that warns when knowledge files
exceed configurable thresholds.

### New Check: `checkEntryCount`

Count entries in DECISIONS.md and LEARNINGS.md using the existing
`index.ParseEntryBlocks` parser. Warn when counts exceed thresholds.

**Default thresholds**:

| File            | Default | Rationale                                    |
|-----------------|---------|----------------------------------------------|
| LEARNINGS.md    | 30      | Learnings are situational; many become stale  |
| DECISIONS.md    | 20      | Decisions are more durable but still compound |

**Warning format**:

```
⚠ LEARNINGS.md has 47 entries (recommended: ≤30)
  Run 'ctx learnings archive' or /ctx-consolidate to reduce
```

### Configuration

New `.ctxrc` keys:

```yaml
entry_count_learnings: 30    # warn above this (0 = disable)
entry_count_decisions: 20    # warn above this (0 = disable)
```

These are soft caps — warnings only, never enforcement.

## Changes

### Modified Files

| File | Change |
|------|--------|
| `internal/drift/detector.go` | Add `checkEntryCount` function |
| `internal/drift/types.go` | Add `IssueEntryCount` type, `CheckEntryCount` name |
| `internal/drift/detector_test.go` | Tests for new check |
| `internal/rc/types.go` | Add threshold fields |
| `internal/rc/default.go` | Add default values |

### New Types

```go
const IssueEntryCount IssueType = "entry_count"
const CheckEntryCount CheckName = "entry_count_check"
```

### Implementation

```go
func checkEntryCount(ctx *context.Context, report *Report) {
    checks := []struct {
        file      string
        threshold int
    }{
        {config.FileLearning, rc.EntryCountLearnings()},
        {config.FileDecision, rc.EntryCountDecisions()},
    }
    found := false
    for _, c := range checks {
        if c.threshold <= 0 { continue } // disabled
        f := ctx.File(c.file)
        if f == nil { continue }
        blocks := index.ParseEntryBlocks(string(f.Content))
        if len(blocks) > c.threshold {
            report.Warnings = append(report.Warnings, Issue{...})
            found = true
        }
    }
    if !found { report.Passed = append(report.Passed, CheckEntryCount) }
}
```

## Non-Goals

- **Auto-fix for entry count**: Unlike staleness (which has a clear fix —
  archive), reducing entry count requires judgment. The nudge points users
  to `ctx learnings archive` or `/ctx-consolidate`.
- **Conventions/Tasks count**: Conventions grow slowly and are always
  relevant. Tasks have their own staleness check already.
- **Hard limits**: This is always a warning, never a violation.

## Testing

- Count check with 0 entries (no warning)
- Count check at threshold (no warning)
- Count check above threshold (warning with correct count)
- Disabled threshold (0 = no check)
- Custom threshold from `.ctxrc`
- Both files above threshold (two warnings)
- File missing (no error)

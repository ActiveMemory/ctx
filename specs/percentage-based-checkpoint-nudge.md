# Percentage-based checkpoint nudge

## Problem

The `check-context-size` hook uses two independent trigger mechanisms:

1. **Counter-based** (prompt count): fires every 5 prompts after 15,
   every 3 after 30, gated behind a 20% minimum context usage floor.
2. **Percentage-based** (window usage): fires a single warning at 80%.

The counter-based approach is fundamentally flawed: prompt count is a
poor proxy for session depth. A session with many small prompts burns
through counters quickly while barely touching the context window. A
session with few large prompts (code reviews, long pastes) can fill
the window before the counter triggers.

The 20% gate (HA.2) was a patch — it suppressed the worst noise on
1M windows but didn't fix the underlying problem. The counter logic
still fires at arbitrary points above 20%, and the remaining-tokens
interpretation is disproportional: 100K remaining tokens is 10% of a
1M window but 50% of a 200K window.

Meanwhile, the 80% window warning fires only once, too late for a
useful checkpoint. There's no "you're getting deep, consider
persisting" signal between 20% and 80%.

## Approach

Replace both trigger mechanisms with two percentage-based thresholds:

| Threshold | Name       | Purpose                                  |
|-----------|------------|------------------------------------------|
| **60%**   | Checkpoint | "Past halfway — consider persisting"     |
| **90%**   | Warning    | "Running low — checkpoint and wrap up"   |

### Behavior

- **At 60%**: emit a checkpoint nudge. One-shot per session — once
  the 60% nudge fires, it doesn't repeat. The message encourages
  persisting progress (decisions, learnings, task updates) but is
  not urgent.
- **At 90%**: emit a window warning. Fires on every prompt at or
  above 90%. The message is urgent: checkpoint everything, prepare
  for context compaction or session handoff.
- **Below 60%**: silent. No nudges, no counters.

### What this eliminates

- `CheckpointEarlyThreshold` (15) — deleted
- `CheckpointEarlyInterval` (5) — deleted
- `CheckpointLateThreshold` (30) — deleted
- `CheckpointLateInterval` (3) — deleted
- `ContextCheckpointMinPct` (20) — deleted (subsumed by 60%)
- `ContextWindowThresholdPct` (80) — replaced by 90%
- All counter-based branching in `check_context_size/run.go`
- Counter file (`context-check-{sessionID}`) — no longer needed for
  nudge decisions (still needed for prompt count in stats/logging)

### What stays unchanged

- `EffectiveContextWindow` detection (4-tier fallback)
- Billing threshold (independent, already percentage-agnostic)
- Wrap-up suppression logic
- Pause logic
- Stats recording (still writes prompt count, tokens, pct)
- `check-persistence` hook: continues to read `LatestPct` from stats,
  but its gate changes from `ContextCheckpointMinPct` (20%) to the
  new checkpoint threshold (60%) for consistency

### Edge case: pct = 0

When token data is unavailable (`tokens=0`, `windowSize=0`), `pct`
computes to 0. This means **no nudges fire** — which is correct.
The previous system fired counter-based nudges when pct=0 as a
"fallback." In practice this only happens when the JSONL file hasn't
been written yet (first 1-2 prompts) or the session ID is unknown.
Neither case benefits from checkpoint nudges.

If `EffectiveContextWindow` returns 200K (the floor default) and
token data is available, percentage works correctly. There is no
scenario where we have tokens but no window size — the window always
resolves to at least 200K.

## Storage

### Constants changed (`config/stats/context.go`)

```go
// Delete:
// ContextCheckpointMinPct = 20
// CheckpointLateThreshold = 30
// CheckpointLateInterval = 3
// CheckpointEarlyThreshold = 15
// CheckpointEarlyInterval = 5
// ContextWindowThresholdPct = 80

// Add:
ContextCheckpointPct = 60  // one-shot checkpoint nudge
ContextWindowWarnPct = 90  // recurring urgent warning
```

### State files

- Counter file (`context-check-{sessionID}`) is retained for prompt
  count tracking in stats, but no longer drives nudge decisions.
- Add a one-shot guard file (`checkpoint-nudged-{sessionID}`) in
  state dir to prevent the 60% checkpoint from repeating.

## CLI surface

No changes. All modifications are internal to the hook system.

## Error cases

- **Token data unavailable** (pct=0): no nudges fire. No regression —
  the counter-based fallback was noise, not value.
- **Window detection wrong** (e.g., 200K when actually 1M): nudges
  fire too early but are still proportional. The 4-tier fallback
  makes this unlikely.

## Non-goals

- Changing the billing threshold mechanism
- Making thresholds user-configurable (can be added later via `.ctxrc`
  if needed, but premature now)
- Changing the check-persistence hook's own nudge logic (only its
  minimum-pct gate changes)

## Testing

- Unit test: no nudge fires below 60%
- Unit test: checkpoint fires at 60%, one-shot (not again at 65%)
- Unit test: warning fires at 90%, repeats on subsequent prompts
- Unit test: pct=0 produces no nudges
- Unit test: check-persistence gate updated to 60%
- Update existing tests that reference deleted constants
- Integration: 1M window session — verify no noise below 60%

## Migration

The old constants are internal, not user-facing. Deletion is safe.
The `ContextWindowThresholdPct` constant is referenced by
`nudge/token.go:TokenUsageLine` for the icon swap — update that
reference to `ContextWindowWarnPct`.

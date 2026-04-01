# Spec: ctxrc constant promotion audit

## Problem

125 hardcoded constants control thresholds, intervals, and
limits across `internal/config/`. Only 10 have .ctxrc accessors.
Some of these should be user-configurable; most should not.

## Decision Framework

A constant should be promoted to .ctxrc when:
1. **Users have legitimate reasons to change it** — different
   project sizes, different workflows, different preferences
2. **The value affects user-visible behavior** — nudge
   frequency, display limits, warning thresholds
3. **The value is not a protocol/format constant** — nonce
   sizes, permission bits, column widths are implementation

A constant should stay hardcoded when:
- It's a conversion factor (CharsPerToken, HoursPerDay, GiB)
- It's a security parameter (KeySize, NonceSize, permissions)
- It's a display formatting constant (column widths, separator
  lengths, prefix lengths)
- It's a protocol constant (buffer sizes for known formats)
- Changing it would break correctness (DatePrefixLen, etc.)

## Candidates for Promotion

### Tier 1 — Clear user value (promote)

| Constant | Current | Why configurable |
|----------|---------|-----------------|
| ContextCheckpointPct | 60 | Users with different workflow depth |
| ContextWindowWarnPct | 90 | Urgency preference |
| MapStaleDays | 30 | Architecture refresh cadence |
| AutoPruneStaleDays | 7 | State cleanup frequency |
| WebhookTimeout | 5s | Network conditions |
| MaxMessagesPerPart | 200 | Journal file size preference |
| DefaultRecallListLimit | 20 | Session history depth |
| MaxBlobSize | 64KB | Scratchpad file limits |
| PersistenceEarlyMin | 11 | Nudge sensitivity |
| PersistenceLateInterval | 15 | Nudge frequency |
| DefaultPublishBudget | 80 | Memory publish verbosity |
| PublishRecentDays | 7 | Memory publish lookback |

### Tier 2 — Marginal value (defer)

| Constant | Current | Why defer |
|----------|---------|-----------|
| ThresholdMemory*Pct | 80/90 | Sysinfo is informational |
| ThresholdSwap*Pct | 50/75 | Sysinfo is informational |
| ThresholdDisk*Pct | 85/95 | Sysinfo is informational |
| ThresholdLoad*Ratio | 0.8/1.5 | Sysinfo is informational |
| TaskCompletionWarnPct | 80 | Doctor is advisory |
| ContextSizeWarnPct | 20 | Doctor is advisory |
| AgentRecency* | various | Scoring internals |
| MaxNavTitleLen | 40 | Zensical rendering |
| LineWrapWidth | 80 | Already convention |
| PreviewMaxLen | 100 | Display detail |
| MaxTitleLen | 75 | Journal formatting |

### Tier 3 — Do not promote (implementation)

Everything else: conversion factors, security parameters,
display column widths, buffer sizes, protocol constants,
format-specific lengths, violation scores.

## Implementation

For each Tier 1 constant:
1. Add field to CtxRC struct in `internal/rc/`
2. Add accessor function with default fallback
3. Replace direct constant usage with accessor call
4. Add to .ctxrc schema documentation
5. Add to `ctx config schema` output

## Non-Goals

- Promoting Tier 2 or Tier 3 constants
- Adding validation beyond type checking
- Adding migration for existing .ctxrc files (new fields are
  optional with defaults matching current values)

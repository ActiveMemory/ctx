# Suppress Context Checkpoint Nudges After Wrap-Up

## Problem

After running `/ctx-wrap-up`, the `check-context-size` hook continues firing
context checkpoint nudges. The wrap-up ceremony itself takes several prompts
(gather signal, present candidates, persist each one, offer commit), and those
prompts keep incrementing the counter. A checkpoint fires *during or right
after* wrap-up — pure noise, since everything was already saved.

## Root Cause

`check-context-size` uses a per-session counter file in `secureTempDir()`.
After prompt #30, it fires every 3rd prompt. It has no concept of "wrap-up
already happened."

## Design

Marker file with time-based expiry.

### New plumbing command: `ctx system mark-wrapped-up`

- Writes `ctx-wrapped-up` to `secureTempDir()` (just a touch — mtime is
  the signal)
- Hidden plumbing command, same pattern as `mark-journal`
- No flags, no arguments, no stdin needed

### Change to `check-context-size`

Before emitting a checkpoint, check:

```go
markerPath := filepath.Join(tmpDir, "ctx-wrapped-up")
if info, err := os.Stat(markerPath); err == nil {
    if time.Since(info.ModTime()) < 2*time.Hour {
        logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d suppressed (wrapped up)", count))
        return nil
    }
}
```

If the marker exists and is less than 2 hours old, suppress the nudge.

### Skill change: `/ctx-wrap-up`

After Phase 3 (persist approved candidates), call:

```bash
ctx system mark-wrapped-up
```

This is a single line added to the skill instructions.

## Why 2 Hours?

- Covers even the longest post-wrap-up tail (user asks a few more questions
  after wrapping up)
- A new session starts its counter at 1 — prompts 1-15 are always silent —
  so the marker doesn't bleed into a genuinely new session in practice
- If someone starts a new 3+ hour session within 2 hours of wrapping up,
  nudges resume once the marker expires and the counter crosses the threshold

## Files to Change

| File | Change |
|------|--------|
| `internal/cli/system/mark_wrapped_up.go` | New plumbing command |
| `internal/cli/system/mark_wrapped_up_test.go` | Tests |
| `internal/cli/system/system.go` | Register `markWrappedUpCmd()` |
| `internal/cli/system/check_context_size.go` | Check marker before emitting |
| `internal/cli/system/check_context_size_test.go` | Test suppression |
| `/ctx-wrap-up` skill (plugin) | Add `ctx system mark-wrapped-up` call |

## Scope

- No new dependencies
- No config changes
- No hook registration changes (mark-wrapped-up is a plumbing command, not a hook)
- Backward compatible: without the marker, behavior is unchanged

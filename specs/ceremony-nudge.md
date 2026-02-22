# Ceremony Nudge Hook

## Problem

`/ctx-remember` (session start) and `/ctx-wrap-up` (session end) are the
two bookend ceremonies for productive sessions. But users forget to invoke
them — especially when starting out. The skills exist, the value is clear,
but habit formation requires gentle, repeated reinforcement.

## Approach

A new `ctx system check-ceremonies` hook that:

1. Scans recent session history to detect whether the user has been using
   the ceremony skills.
2. If not, emits a VERBATIM relay nudge at session start — once per day,
   with a teaching tone that explains *why* the ceremonies matter.

This follows the established hook patterns:
- `check-persistence`: nudges when context hasn't been persisted
- `check-journal`: nudges when sessions haven't been exported
- `check-backup-age`: nudges when backups are stale

## Detection

### Signal Source: Journal Files

Scan `.context/journal/` for the N most recent entries (by filename date
prefix). For each, grep for `/ctx-remember` or `/ctx-wrap-up` in the
content.

Journal files are the right signal source because:
- They're already exported and local (no JSONL parsing needed).
- Skill invocations appear in the transcript text.
- The hook can use simple string matching, not JSON parsing.

### No JSONL Fallback

Intentionally **no fallback to JSONL transcripts**. JSONL files are huge
(megabytes each) and scanning them from a hook would eat context budget —
exactly what ctx exists to prevent.

If journals are stale or missing, the hook nudges the user to export
journals instead of silently falling back to an expensive scan. This
doubles as a journal-health nudge (similar to `check-journal`).

### Detection Logic

```
journal_entries = last 3 journal files (by filename date prefix)

if len(journal_entries) == 0:
    # No journals at all — nudge to export, skip ceremony check
    emit journal-export nudge
    exit

remember_used = any entry contains "/ctx-remember" or "ctx-remember"
wrapup_used   = any entry contains "/ctx-wrap-up" or "ctx-wrap-up"

if remember_used AND wrapup_used:
    exit silently
else:
    emit nudge for the missing ceremony(ies)
```

## Hook: `ctx system check-ceremonies`

### Registration

Register under `UserPromptSubmit` in hooks.json, alongside existing hooks.
Should run early in the hook chain (before check-persistence) since it's
a session-start concern.

### Throttling

- **Once per day** (same `isDailyThrottled` pattern as check-journal).
- Uses marker file: `$TMPDIR/ctx-$UID/ceremony-reminded`.

### Output

Uses the VERBATIM relay pattern. Three variants depending on what's missing:

**Both missing:**

```
IMPORTANT: Relay this session tip to the user VERBATIM before answering their question.

┌─ Session Ceremonies ──────────────────────────────
│ Your last 3 sessions didn't use /ctx-remember or
│ /ctx-wrap-up.
│
│ Start sessions with:  /ctx-remember
│   → Loads context, shows active tasks, picks up
│     where you left off. No re-explaining needed.
│
│ End sessions with:    /ctx-wrap-up
│   → Captures learnings and decisions so the next
│     session starts informed, not from scratch.
│
│ These take seconds and save minutes.
└───────────────────────────────────────────────────
```

**Only /ctx-remember missing:**

```
IMPORTANT: Relay this session tip to the user VERBATIM before answering their question.

┌─ Session Start ───────────────────────────────────
│ Try starting this session with /ctx-remember
│
│ It loads your context, shows active tasks, and
│ picks up where you left off — no re-explaining.
└───────────────────────────────────────────────────
```

**Only /ctx-wrap-up missing:**

```
IMPORTANT: Relay this session tip to the user VERBATIM before answering their question.

┌─ Session End ─────────────────────────────────────
│ Your last 3 sessions didn't end with /ctx-wrap-up
│
│ It captures learnings and decisions so the next
│ session starts informed, not from scratch.
└───────────────────────────────────────────────────
```

### Tone

- **Teaching, not nagging**. Explain the benefit, don't just command.
- **Concise**. The box should fit in a terminal without scrolling.
- **Diminishing**. Once the user habitualizes, the nudge stops firing
  because the detection finds ceremony usage in recent sessions.
- **Self-silencing**. The nudge naturally disappears as the user adopts
  the habit — no "dismiss forever" flag needed.

## Implementation

### File: `internal/cli/system/check_ceremonies.go`

```go
func checkCeremoniesCmd() *cobra.Command { ... }

func runCheckCeremonies(cmd *cobra.Command) error {
    // 1. Check initialized
    // 2. Check daily throttle
    // 3. Check journal directory exists and has entries
    //    - If no journals: nudge to export, skip ceremony check
    // 4. Scan last 3 journal entries for ceremony usage
    // 5. Emit appropriate nudge or exit silently
}

// scanJournalsForCeremonies checks the N most recent journal files
// for /ctx-remember and /ctx-wrap-up usage.
func scanJournalsForCeremonies(journalDir string, n int) (remember, wrapup bool) { ... }

// recentJournalFiles returns the N most recent .md files in the journal
// directory, sorted by filename (date prefix gives chronological order).
func recentJournalFiles(journalDir string, n int) []string { ... }
```

### Helpers to Reuse

From `state.go`:
- `secureTempDir()` — temp directory for marker file
- `isDailyThrottled()` — once-per-day check
- `touchFile()` — update marker
- `isInitialized()` — skip if no .context/
- `logMessage()` — diagnostic logging

From `check_journal.go`:
- `newestMtime()` — find recent files
- `countUnenriched()` — pattern for reading journal content
- Journal directory path pattern

### Hook Registration

Add to `system.go` command tree:

```go
systemCmd.AddCommand(checkCeremoniesCmd())
```

Add to `internal/assets/claude/hooks.json`:

```json
{
    "type": "command",
    "command": "ctx system check-ceremonies",
    "event": "UserPromptSubmit"
}
```

## Configuration

No `.ctxrc` keys needed for v1. The behavior is:
- Look back 3 sessions (hardcoded, reasonable default).
- Throttle to once per day.
- Self-silencing when ceremonies are detected.

Future: `ceremony_nudge: false` in `.ctxrc` to disable entirely.

## Testing

- Unit: `scanJournalsForCeremonies` with fixtures containing/missing
  ceremony references.
- Unit: throttling behavior (marker file present/absent/stale).
- Unit: output variants (both missing, one missing, neither missing).
- Integration: full hook run with sample journal directory.
- Edge cases: no journal directory, empty journal, all sessions use
  ceremonies (should be silent), journals stale (should nudge export
  instead of ceremony).

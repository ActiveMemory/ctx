---
title: Consolidation nudge hook (3:1 rule enforcement)
date: 2026-03-23
status: ready
---

# Consolidation Nudge Hook

## Problem

The 3:1 rule (three feature sessions, then one consolidation session)
is documented in the blog and AGENT_PLAYBOOK.md but has no mechanical
enforcement. The v0.8.0 cycle proved the cost: 198 feature commits
before any consolidation, resulting in an 18-day cleanup marathon.

Humans and agents both need a nudge, not a gate.

## Solution

A `UserPromptSubmit` hook that counts sessions since the last
consolidation and nudges after a configurable threshold (default: 6).

The threshold is 6 rather than 3 because not all sessions are full
feature sessions: some are quick fixes, some are research, some get
cut short. 6 sessions approximates "3 full feature sessions" in
practice.

## Design

### Session Counting

Use the existing journal state (`.context/journal/.state.json`) to
count exported sessions. Each session has a date and can be classified
by type (feature, bugfix, refactor, exploration, etc.) via journal
enrichment frontmatter.

Two counting strategies (pick one during implementation):

**Strategy A: Count all sessions.** Simple. Every session increments
the counter. Consolidation resets it. Threshold default: 6.

**Strategy B: Count only feature/bugfix sessions.** Skip sessions
typed as "exploration", "debugging", or under 10 turns. More accurate
to the 3:1 intent but requires enriched journal entries.

Recommendation: start with Strategy A. It's simpler and the threshold
can be tuned via `.ctxrc`.

### Consolidation Detection

A session counts as "consolidation" if any of:

1. A commit message in the session contains keywords: "refactor",
   "consolidate", "cleanup", "convention audit", "lint-drift"
2. The journal entry is enriched with `type: refactor`
3. The user explicitly marks consolidation via
   `ctx system mark-consolidation` (new plumbing command)
4. The session ran `make audit` and it passed (heuristic)

Strategy 3 is the most reliable. The others are heuristics for
automatic detection.

### State Storage

Add a `consolidation` section to `.context/state/session.json`:

```json
{
  "consolidation": {
    "last_session_date": "2026-03-23",
    "last_commit": "692f86cd",
    "sessions_since": 0
  }
}
```

Increment `sessions_since` at each session start. Reset to 0 when
consolidation is detected.

### Hook Behavior

The hook runs during `UserPromptSubmit` (session start). When
`sessions_since >= threshold`:

```
┌─ Consolidation Reminder ─────────────────────────────
│ 7 sessions since last consolidation (threshold: 6).
│ Consider scheduling a consolidation session.
│ Run: make audit
│ Mark done: ctx system mark-consolidation
│ Snooze: ctx system snooze-consolidation
└──────────────────────────────────────────────────
```

The message is a nudge, not a gate. It does not block work.

### Snooze

`ctx system snooze-consolidation` suppresses the nudge for N sessions
(default: 3). This handles "I know, but I need to finish this feature
first." The snooze counter decrements per session; when it hits 0 the
nudge resumes.

### Configuration

In `.ctxrc`:

```yaml
consolidation:
  threshold: 6          # sessions before nudge (default: 6)
  snooze_sessions: 3    # sessions to suppress after snooze
  auto_detect: true     # use commit keyword heuristics
  keywords:             # consolidation detection keywords
    - refactor
    - consolidate
    - cleanup
    - convention audit
```

## CLI Surface

### New plumbing commands

```
ctx system mark-consolidation    # reset counter, record baseline commit
ctx system snooze-consolidation  # suppress nudge for N sessions
```

### Integration with existing commands

- `ctx status` shows sessions-since-consolidation count
- `ctx drift` includes consolidation debt as a warning when over threshold

## Implementation Plan

1. Add `consolidation` state to session state JSON
2. Add `mark-consolidation` and `snooze-consolidation` plumbing commands
3. Add `check-consolidation` hook to `UserPromptSubmit`
4. Add hook message template to `internal/assets/hooks/messages/`
5. Add `.ctxrc` configuration support
6. Add `sessions_since` to `ctx status` output
7. Add consolidation debt warning to `ctx drift`

## Non-Goals

- Blocking feature work (this is a nudge, never a gate)
- Automatically starting consolidation sessions
- Tracking consolidation quality (that's what `make audit` is for)
- Per-file or per-package consolidation tracking

## Open Questions

- Should the counter reset on any refactor-type session, or only on
  sessions where `make audit` passes afterward?
- Should the nudge include a diff of what's drifted since the last
  consolidation (expensive but informative)?

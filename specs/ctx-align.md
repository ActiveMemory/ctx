# ctx-align (course-correction capture)

*Promoted from `ideas/ctx-align.md` by a ctx-dream serendipity round
(2026-06-07). Seed-stage spec: the shape is the author's note; the
contested choices are Open Questions for a later `/ctx-plan` pass.*

## Problem

When a user interrupts an agent with "stop — that's not what I meant,"
the correction is high-signal: it pinpoints exactly where the agent's
heuristic diverged from intent. Today that signal evaporates. The agent
adjusts for the rest of the turn and the *why* (which assumption was
wrong, what the user actually wanted) is never captured — so the same
class of divergence recurs in the next session, and the conventions /
docs that would prevent it are never updated.

## Approach

A `/ctx-align` skill the user invokes at the moment of divergence. The
agent **pauses the task** and:

1. **Records what it was doing and the heuristic behind it** — the
   assumption that led it astray (not just "I did X" but "I did X
   because I inferred Y").
2. **Captures the divergence as a learning** (via `/ctx-learning-add`):
   the gap between inferred intent and actual intent.
3. **Updates the relevant convention/doc** if the divergence reflects a
   durable rule (via `/ctx-convention-add`), or — when it's unclear what
   should change — **asks the user** what to update rather than guessing.
4. **Syncs memory** so the correction is persistent, not turn-local.
5. **Catalogs the frustration** in a running file for future lookup (a
   ledger of "times the agent guessed wrong"), so patterns across
   corrections become visible.
6. **Resumes** the task with the corrected understanding.

It reuses the existing capture skills (`/ctx-learning-add`,
`/ctx-convention-add`, `/ctx-reflect`) rather than reinventing capture;
its distinct value is the **trigger** (a user course-correction) and the
**frustration catalog**.

## Behavior

### Happy Path

1. User: "stop, that's not what I meant — do X instead because Y."
2. `/ctx-align`: agent pauses, states what it was doing + the heuristic.
3. Agent records the divergence as a learning; proposes a convention/doc
   update if one applies; asks the user to confirm what to persist.
4. Agent appends the divergence to the frustration catalog.
5. Agent resumes with the corrected approach.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Correction is one-off (no durable rule) | Record as a learning only; no convention update; still cataloged. |
| Correction implies a convention change | Propose the convention edit; apply only with user confirmation. |
| User invokes mid-multi-step task | Capture without losing task state; resume the in-flight step after. |
| Repeated same divergence | The catalog surfaces the recurrence — escalate to a convention. |

## Interface

### Skill

```
/ctx-align
```

Trigger phrases: "stop, that's not what I meant", "align", "that's
wrong, here's what I meant".

## Non-Goals

- Not a general reflection pass (that's `/ctx-reflect`); the trigger is a
  specific user course-correction.
- Not auto-applying convention changes without confirmation.

## Open Questions

- **Frustration-catalog location/format** — a new `.context/` file, or
  fold into LEARNINGS with a tag? (TBD)
- **Relationship to `/ctx-reflect`** — distinct skill vs. a mode of
  reflect? (TBD)
- **Auto-vs-ask threshold** — when does the agent propose a convention
  edit vs. just record a learning? (TBD — `/ctx-plan` it.)

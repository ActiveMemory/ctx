---
title: "Typical KB Session"
icon: lucide/notebook-pen
---

![ctx](../images/ctx-banner.png)

## The Problem

You set the editorial pipeline up
([Build a Knowledge Base](build-a-knowledge-base.md)). Now you
sit down for a real research session: a transcript to ingest, a
question to answer against existing evidence, a finding to
capture for later. What's the actual flow?

## TL;DR

```text
/ctx-remember                                    # session-start recall
/ctx-kb-ingest ./inputs/transcript.md "topic"    # editorial pass
/ctx-kb-ask "does the kb say X?"                 # grounded Q&A
/ctx-kb-note "follow-up: chase the v1.1 link"    # park a finding
/ctx-wrap-up                                     # ceremony → /ctx-handover
```

## Commands and Skills Used

| Tool              | Type    | Purpose                                           |
|-------------------|---------|---------------------------------------------------|
| `/ctx-remember`   | Skill   | Session-start recall (folds KB state when present) |
| `/ctx-kb-ingest`  | Skill   | Mode-aware editorial pass                         |
| `/ctx-kb-ask`     | Skill   | Q&A grounded in the kb                            |
| `/ctx-kb-note`    | Skill   | Park a finding for the next ingest                |
| `/ctx-wrap-up`    | Skill   | End-of-session ceremony; delegates to the handover step |
| `/ctx-handover`   | Skill   | Writes the per-session handover; called by `/ctx-wrap-up` |

## Step 1: Session Start (Recall)

```text
/ctx-remember
```

`/ctx-remember` reads the latest handover under
`.context/handovers/` (timestamped `<TS>-<slug>.md` so
concurrent agent runs never overwrite); its `## Summary` and
`## Next session` are the authoritative recall surface. The
five canonical files (`TASKS`, `DECISIONS`, etc.) are read as
usual.

When `.context/kb/` exists, `/ctx-remember` additionally folds
editorial state into the readback: any closeouts whose
`generated-at` postdates the handover are read for their
`## What changed` sections (these are unfolded passes the
last handover did not yet consume).

`SESSION_LOG.md` is **not** read at session start; it is
mid-flight working memory, not a recall surface.

## Step 2: Ingest the Sources You Brought

```text
/ctx-kb-ingest ./inputs/2026-05-15-call.md "cursor hooks"
```

The skill declares its mode up front (most often
`topic-page`), resolves sources, scans the
**source-coverage ledger** for adjacent incomplete topics,
and synthesizes prose into the topic page section by section.
Every cited claim mints an `EV-###` row in
`evidence-index.md` with the source short-name + locator +
optional `sha:` pin for in-repo files.

The pass ends with a **circuit-breaker check** (file exists,
cites ≥ 1 `EV-###`, site builds clean, cold-reader rubric at
`pass`) and writes a closeout.

If the skill reports `topic-page: deferred` instead of
`produced`, look at the closeout's `Next pass hint`. It
names the exact resumption invocation.

## Step 3: Ask Grounded Questions

```text
/ctx-kb-ask "does the kb say hooks block until they exit?"
```

`/ctx-kb-ask` reads the kb's prose and answers with `EV-###`
citations. If the kb cannot answer, it opens a `Q-###` row
in `outstanding-questions.md` and reports the gap rather than
inventing.

## Step 4: Park Findings for Later

```text
/ctx-kb-note "check whether SIGTERM behavior changed in v1.2"
```

`/ctx-kb-note` appends one-liners to
`.context/ingest/findings.md`, a lightweight surface for
parking ideas that don't earn a full ingest pass right now.
The next `/ctx-kb-ingest` can choose to absorb them.

## Step 5: Wrap Up

```text
/ctx-wrap-up "Cursor Hooks: lifecycle deep dive"
```

`/ctx-wrap-up` runs the standard capture checklist (learnings,
decisions, conventions, tasks) and delegates to
`/ctx-handover` as its final step. In a KB session it
additionally:

- Surfaces pending closeouts under
  `.context/ingest/closeouts/`.
- Counts `open` rows in `outstanding-questions.md`.

The handover artifact lands at
`.context/handovers/<TS>-<slug>.md` (timestamped so concurrent
agent runs never overwrite). The handover folds postdated
closeouts into a `## Folded closeouts` section and archives
them under `.context/archive/closeouts/`. Editorial work that
was incomplete at wrap-up (open `Q-###` rows, `topic-page:
deferred` passes) is surfaced as recall on the next session
start.

## Common Shapes

### Multiple Topics in One Session

Run `/ctx-kb-ingest` once per topic. Each pass writes its own
closeout; the handover folds all of them at the end.

### Mid-Session Checkpoint

```bash
ctx handover write "Mid-day checkpoint" \
  --summary "..." --next "..." --no-fold
```

`--no-fold` writes the handover without consuming closeouts,
useful when you want a recall anchor mid-session without
ending the editorial chunking.

### Aborted Session

If you close the laptop after an ingest pass but before
`/ctx-wrap-up`, the closeouts stay in place. The next
session's `/ctx-remember` reads them as unfolded postdated
closeouts; the next wrap-up's handover step folds them
normally. See
[Recover an Aborted Session](recover-aborted-session.md) for
the failure-mode detail.

## Reference

- Recipe: [Build a Knowledge Base](build-a-knowledge-base.md)
- Recipe: [Recover an Aborted Session](recover-aborted-session.md)
- Skill: [`/ctx-kb-ingest`](../reference/skills.md#ctx-kb-ingest)
- Skill: [`/ctx-handover`](../reference/skills.md#ctx-handover)
- Editorial constitution: `.context/ingest/KB-RULES.md`

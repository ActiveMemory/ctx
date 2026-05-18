---
title: "Recover an Aborted KB Session"
icon: lucide/life-buoy
---

![ctx](../images/ctx-banner.png)

## The Problem

You ran one or more `/ctx-kb-ingest` passes, then the session
ended before `/ctx-wrap-up`. Maybe you closed the laptop,
the connection dropped, or you just forgot the wrap-up step.

You come back the next day and ask "do you remember?" and
the agent picks up the previous handover, but the editorial
work since the last handover seems to be missing from the
readback.

**It isn't missing. It's unfolded.** Here's how the pipeline
handles it and how to close the loop manually.

## TL;DR

```text
/ctx-remember                                       # picks up the unfolded 
                                                    # closeouts automatically
/ctx-handover "recovery: fold the orphan closeouts" # direct invocation is 
                                                    # appropriate for recovery
```

The recovery path is the one legitimate place to invoke
`/ctx-handover` directly. Normally `/ctx-wrap-up` owns
session-end and delegates to the handover step; the abort
broke that path, so a hand-rolled handover invocation is how
you close the loop without re-running the full wrap-up
ceremony.

## How the Fold Mechanism Survives an Abort

Two artifacts make abort-recovery work without any cleanup:

1. **Closeouts are immutable once written.** Every editorial
   pass writes a closeout under
   `.context/ingest/closeouts/<TS>-<mode>-closeout.md` *before*
   the pass reports `done`. If the session dies, the closeout
   is already on disk.

2. **`/ctx-remember` folds unfolded closeouts into the
   readback.** The skill always reads the latest handover.
   When `.context/kb/` exists, it additionally reads any
   closeouts whose `generated-at` postdates the handover.
   The `## What changed` and `## Source-coverage updates`
   sections from each unfolded closeout are surfaced in
   recall.

So an aborted session never loses editorial work; it just
delays the handover fold by one session.

## Step 1: Confirm the Orphan Closeouts

```bash
ls -la .context/ingest/closeouts/
```

Files there with `generated-at` postdating your latest
handover are the unfolded ones. You can read any closeout
directly to see what it claims about its pass:

```bash
cat .context/ingest/closeouts/<TS>-ingest-closeout.md
```

Look at:

- The **Pass-mode** body block (`Declared / Reason /
  Definition of done / Result`): what the pass committed to
  and whether it claimed success or `deferred`.
- The **Source-coverage updates** section: what state
  transitions hit the ledger.
- The **Next pass hint**: the exact resumption invocation
  the closeout recommends, if the pass deferred.

## Step 2: Run `/ctx-remember`

```text
/ctx-remember
```

The readback will include the editorial-state summary as part
of the standard readback shape. If everything looks
consistent, proceed to Step 3.

If the readback surfaces something surprising (a closeout
claiming `topic-page: produced` for a slug whose file is
missing, a `comprehensive` ledger advance against a source
whose page is `speculative`, etc.), fix the underlying
inconsistency before folding. (Doctor advisories for these
shapes are on the Phase-7 backlog.)

## Step 3: Write the Recovery Handover

This step is the one legitimate direct invocation of
`/ctx-handover`. In normal session-end the call goes through
`/ctx-wrap-up`; here the prior session aborted, so you reach
for the handover step directly to retire the orphan
closeouts:

```text
/ctx-handover "recovery: fold orphan closeouts from yesterday"
```

Or via the CLI:

```bash
ctx handover write "recovery: fold orphan closeouts from yesterday" \
  --summary "Folded N orphan closeouts from the aborted session." \
  --next "Resume <topic> per the closeout's Next pass hint."
```

The handover:

- Reads the latest handover cursor.
- Finds all closeouts whose `generated-at` is after the cursor.
- Folds their summaries into a `## Folded closeouts` section.
- **Archives the source closeout files** under
  `.context/archive/closeouts/` (closeouts are
  append-never-rewrite; archival moves bytes but does not
  modify them).

After the handover lands, the orphan closeouts are now durably
tied to a session boundary; the next `/ctx-remember` reads
*just* the new handover (and any closeouts postdating *it*),
without re-folding the recovered ones.

## Edge Cases

| Case | Behavior |
|---|---|
| Closeout has malformed frontmatter | Handover fold **skips it with a warning** to stderr. Hand-edit the malformed file (typically a missing `generated-at`) and re-run `ctx handover write` to fold it next time. |
| Closeout's `generated-at` is *before* the last handover but was never folded | Treated as already-folded (silently skipped; the cursor is the source of truth). If you genuinely want to re-fold it, hand-edit the closeout's `generated-at` forward. |
| You aborted *during* an ingest pass, before its closeout was written | No closeout exists; the pass left no recall residue. Treat the source(s) as un-ingested and re-run `/ctx-kb-ingest`. The source-coverage ledger row may show stale residue from a prior pass; the next ingest will advance it correctly. |
| Multiple sessions piled up unfolded closeouts | One handover run folds them all in a single shot. The fold is cursor-driven, not session-driven. |
| You want recall without consuming closeouts | `ctx handover write ... --no-fold` writes a handover with frontmatter but leaves the closeouts in place. The next handover (without `--no-fold`) folds everything postdating the latest handover cursor. |

## When This Matters

- After a network drop / laptop close mid-session.
- When you ran `/ctx-kb-ingest` from a sub-agent that
  finished without calling `/ctx-handover`.
- After porting work from another environment (e.g. you
  rsynced `.context/ingest/closeouts/` from a different
  machine) and want to integrate the work into the destination
  project's recall thread.

## Reference

- Recipe: [Build a Knowledge Base](build-a-knowledge-base.md)
- Recipe: [Typical KB Session](typical-kb-session.md)
- Editorial constitution: `.context/ingest/KB-RULES.md`

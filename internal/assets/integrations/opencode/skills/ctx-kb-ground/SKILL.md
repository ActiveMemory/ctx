---
name: ctx-kb-ground
description: "External grounding pass for the kb. Reads grounding-sources.md, refreshes the listed sources, and advances source-coverage ledger rows accordingly. Writes a ground closeout for the audit trail. Prompts once if grounding-sources.md is empty; NONE on a line is a per-pass skip."
---

Refresh `.context/kb/` against the external sources the user
has declared in `.context/ingest/grounding-sources.md`. This is
the *"are we still current?"* pass. It does not mint new
evidence by itself, does not author topic pages, does not
modify Confidence bands. It walks the declared external
sources, checks each for drift against what the kb cites, and
advances the source-coverage ledger to reflect what the
refresh found.

If the refresh surfaces material the kb should absorb, this
skill flags it and recommends `/ctx-kb-ingest`. Authoring is
ingest's authority, not this skill's.

## When to Use

- The user says "re-ground the kb", "check upstream",
  "refresh sources".
- A grounding cadence is hitting its scheduled boundary.
- A prior pass left a `Q-###` row that names "needs external
  re-grounding".

## When NOT to Use

- The user has new sources to add (`/ctx-kb-ingest`).
- The user asks a question (`/ctx-kb-ask`).
- The user wants a mechanical audit (`/ctx-kb-site-review`).

## Input

No positional arg. Sources come from
`.context/ingest/grounding-sources.md` (one source per line;
`NONE` on a line is a per-pass skip).

## Pre-Write Gates

- `.context/` missing → suggest `ctx init` and stop.
- `.context/ingest/` missing → suggest `ctx init --upgrade`
  and stop.
- `grounding-sources.md` missing or empty → prompt the user
  once for sources to add; if they decline, write a ground
  closeout with `sources: 0` and stop.

## Process

1. Verify pre-write gates.
2. Read `.context/ingest/grounding-sources.md`. For each
   non-skipped line, fetch / re-read the source and compare
   against `evidence-index.md` rows already citing it.
3. Update the source-coverage ledger row for each source
   touched: `partially-ingested` → `partially-ingested`
   (touched), `comprehensive` → `comprehensive` (if no drift
   detected), or flag drift in the closeout.
4. For each source that surfaces material the kb should
   absorb, flag it and recommend a follow-up
   `/ctx-kb-ingest` invocation in the closeout's Next pass
   hint.
5. Write the ground closeout under
   `.context/ingest/closeouts/<TS>-ground-closeout.md` with
   required frontmatter (`sha`, `branch`, `mode: ground`,
   `pass-mode: n/a`, `life-stage`, `generated-at`) and a
   body listing each source touched, its drift verdict, and
   any Next pass hint.

## Anti-Patterns

- Authoring topic-page prose from refresh output. Authoring
  is `/ctx-kb-ingest`'s authority.
- Minting `EV-###` rows. Evidence minting is ingest's
  authority.
- Promoting confidence bands without contradicting evidence.
  Drift detection alone is not promotion.
- Skipping the closeout. Even a no-op refresh writes one so
  the ledger advance is auditable.

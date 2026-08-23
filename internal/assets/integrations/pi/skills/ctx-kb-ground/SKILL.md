---
name: ctx-kb-ground
description: "Read-only freshness audit over the kb's tracked sources (URLs, in-tree paths, MCP resources) declared in grounding-sources.md. Classifies each source's drift state, annotates the source-coverage ledger, and writes a ground closeout; flags drifted or new-to-kb sources for /ctx-kb-ingest. Never mints evidence, authors prose, or transitions ledger states."
---

Walk the sources declared in `.context/ingest/grounding-sources.md`
and report whether the kb's claims are still current. This is the
*"are we still current?"* pass — a read-only freshness audit, not
a re-ingest.

Each tracked source — URLs, in-tree paths, or MCP resources —
gets resolved and classified as `unchanged`, `drifted`, `gone`,
`freshness opaque`, or `new to kb`. The skill annotates the
source-coverage ledger's `Residue` / `Next action` cells and
writes a ground closeout. It does NOT mint `EV-###` rows, author
topic-page prose, transition ledger states, or modify Confidence
bands; those are `/ctx-kb-ingest`'s authority.

If a tracked source drifted or is new to the kb, flag it and
recommend a follow-up `/ctx-kb-ingest`. The declarative watch
list in `grounding-sources.md` persists across sessions and
tracks sources from anywhere the kb cites — the web, this repo's
tree, an MCP server. Distance from the repo is irrelevant; what
matters is that the kb depends on them for evidence.

## When to Use

- The user says "re-ground the kb", "check upstream",
  "refresh sources".
- A grounding cadence is hitting its scheduled boundary.
- A prior pass left a `Q-###` row that names "needs
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

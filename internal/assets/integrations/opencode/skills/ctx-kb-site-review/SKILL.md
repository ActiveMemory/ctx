---
name: ctx-kb-site-review
description: "Mechanical structural audit of the kb. Coerces malformed capitalization, flags malformed closeout frontmatter, and refuses to make judgment calls that require evidence. Writes a site-review closeout for the audit trail."
---

Walk `.context/kb/` and `.context/ingest/closeouts/`
mechanically. Fix what is unambiguous (capitalization drift,
missing frontmatter fields the CLI knows how to coerce). Flag
what is not (claims that read as broken but require evidence
to fix). Never invent prose. Never mint `EV-###` rows. Never
modify a claim's Confidence band.

This is a janitor pass, not an editorial pass. Editorial
judgment lives in `/ctx-kb-ingest`.

## When to Use

- The user says "audit the kb", "check kb for rot",
  "run a site-review".
- Before a release or at the end of a long editorial run,
  to catch drift.

## When NOT to Use

- The user has new sources to ingest (`/ctx-kb-ingest`).
- The user wants Q&A (`/ctx-kb-ask`).
- The user wants external re-grounding (`/ctx-kb-ground`).

## Authority Boundary

- Coerces formatting drift the CLI can fix unambiguously.
- Flags structural problems that require evidence to resolve
  (open `Q-###` rows; never invent answers).
- Never writes prose, mints `EV-###`, or changes confidence
  bands.

## Pre-Write Gates

- `.context/` missing → suggest `ctx init` and stop.
- `.context/kb/` missing → suggest `ctx init --upgrade` and
  stop.
- Kb scope undeclared → refuse with the scope message and
  stop.

## Process

1. Verify pre-write gates.
2. Walk `.context/kb/` and `.context/ingest/closeouts/`:
   coerce capitalization drift in Confidence bands; fix
   missing frontmatter fields the CLI can supply
   deterministically; flag malformed closeouts (missing
   `generated-at`, malformed YAML).
3. For each structural problem that requires evidence to
   resolve, open a `Q-###` row in
   `outstanding-questions.md` naming the file and the
   shape of evidence that would close it.
4. Write the site-review closeout under
   `.context/ingest/closeouts/<TS>-site-review-closeout.md`
   with required frontmatter (`sha`, `branch`,
   `mode: site-review`, `pass-mode: n/a`, `life-stage`,
   `generated-at`) and a body listing what was coerced,
   what was flagged, and any `Q-###` rows opened.

## Anti-Patterns

- Inventing prose or claims.
- Minting `EV-###` rows.
- Modifying Confidence bands.
- Hand-editing closeouts post-write (closeouts are
  append-never-rewrite).
- Skipping the `Q-###` row when a structural flag would
  otherwise vanish without trace.

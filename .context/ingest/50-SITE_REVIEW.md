# Site-Review Mode

Mechanical structural audit. Walks `.context/kb/` and (if it
exists) the rendered site under `.context/site/kb/` looking
for structural issues that can be fixed **without making
editorial judgment calls**. Defers anything requiring evidence
weighing to a follow-up `ctx kb ingest` or `ctx kb ground`
pass.

Read `KB-RULES.md` first; it carries the authority boundary,
the closeout shape, the validation rules for confidence
bands, and the demotion policy this mode refuses to apply.
This file describes the per-mode procedure only.

---

## Inputs

No arguments required. The skill walks `.context/kb/`
automatically. The operator may pass a scope hint inline
(e.g. *"focus on glossary.md"*) which the skill records in
the closeout's `Inputs` section.

---

## Pre-Write Gates

1. `.context/` and `.context/kb/` exist.
2. `.context/kb/index.md` has a non-placeholder `## Scope`
   section.

---

## What Site-Review Fixes (No Truth-Content Needed)

Mechanical issues. The skill resolves them in place and
records the fix in the closeout's `What changed`:

- **Broken internal links** — relative paths under
  `.context/kb/` that no longer resolve. Retarget to the
  renamed file when the target is unambiguous; otherwise
  remove the link and open an `outstanding-questions.md`
  row.
- **Orphaned source short names** — referenced in
  `evidence-index.md` but missing from `source-map.md`.
  Add the source-map entry when the source identifier is
  obvious from context; flag otherwise.
- **Confidence bands missing or malformed** — must be one of
  `high|medium|low|speculative`. Coerce malformed
  capitalization in place (`High` → `high`, `MEDIUM` →
  `medium`). Flag any other malformation (typos, unknown
  bands) without changing the band.
- **Missing closeout frontmatter fields** in
  `.context/ingest/closeouts/`. Required fields per
  `KB-RULES.md` "Closeout shape": `sha`, `branch`, `mode`,
  `pass-mode`, `life-stage`, `generated-at`. Generate the
  missing field where derivable from filename + git context;
  flag where not derivable.
- **Dangling `outstanding-questions.md` entries** —
  questions whose linked contradictions or claims no longer
  exist after a rename. Update the link or close the row
  with a `(superseded)` annotation citing the rename.
- **Source-coverage ledger inconsistencies** that are purely
  mechanical (a row's `Updated` predates the file's mtime, a
  row references a source not in `source-map.md`). Refresh
  the `Updated` cell; flag the source-map mismatch.

---

## What Site-Review Defers (Requires Evidence Judgment)

These are recorded in `outstanding-questions.md` and as
`Next pass hint`s in the closeout, but **not** resolved in
this pass:

- Confidence-band changes that would require new evidence.
- Promotion or demotion of any claim based on
  stale-vs-fresh judgment.
- Contradictions that surface during the walk but lack a
  clear resolution path.
- Glossary entries with wording the reviewer "doesn't like"
  but that aren't factually wrong.
- Illegal source-coverage ledger state transitions (e.g.
  `comprehensive → highlights-extracted` without an explicit
  `superseded` step). Flag for `ctx kb ingest`; do not
  rewrite the state.
- Sub-page splits suggested by oversized `index.md` files.
  Site-review surfaces the candidate but never auto-splits.

---

## Procedure

1. **Plan** the walk: list the files the skill will touch in
   the closeout's `Inputs` section.
2. **Walk** each file in order. Per file: detect, fix what's
   mechanical, record what's deferred.
3. **Aggregate** into the closeout's `What changed`
   (mechanical fixes) and `New questions` (deferred items).
4. **Write the closeout** under
   `.context/ingest/closeouts/<TS>-site-review-closeout.md`
   with `mode: site-review` in the frontmatter. Set
   `pass-mode: n/a` — site-review does not carry a pass-mode
   contract.
5. **Append a SESSION_LOG line** at the closeout phase
   boundary.

---

## Constraints

- Mechanical fixes are batchable; do **not** write one
  closeout per file. One pass = one closeout.
- Do not touch the rendered site under `.context/site/kb/`.
  The renderer (`zensical`) owns that directory; site-review
  only reads it for cross-referencing.
- If a "fix" turns out to need judgment mid-walk (e.g. a
  broken link could be retargeted to one of two files),
  defer it rather than guessing. The bias in this mode is
  toward refusing to act, not toward acting hastily.
- Site-review never demotes or promotes claims. It never
  edits `evidence-index.md` rows except to fix malformed
  confidence-band capitalization in place.

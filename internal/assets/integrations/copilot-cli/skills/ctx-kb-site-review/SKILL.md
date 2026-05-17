---
name: ctx-kb-site-review
description: "Mechanical structural audit of the kb. Coerces malformed capitalization, flags malformed closeout frontmatter, and refuses to make judgment calls that require evidence. Writes a site-review closeout for the audit trail."
tools: [bash]
---

# Site-Review Pass

Walk `.context/kb/` and `.context/ingest/closeouts/` mechanically.
Fix what is unambiguous (capitalization drift, missing frontmatter
fields the CLI knows how to coerce). Flag what is not (claims
that read as broken but require evidence to fix). Never invent
prose. Never mint `EV-###` rows. Never modify a claim's
Confidence band.

This is a janitor pass, not an editorial pass. Editorial judgment
lives in `/ctx-kb-ingest`.

Authoritative background reading:
`.context/ingest/KB-RULES.md` §Authority boundary.

## When to Use

- The user says "audit the kb", "check kb for rot", "run a
  site-review", or invokes the explicit slash form.
- Before a release / handover where structural cleanliness
  matters.
- After bulk ingest where drift may have accumulated.
- When the doctor advisory has surfaced structural warnings the
  user wants triaged.

## When NOT to Use

- The user wants new material extracted (use `/ctx-kb-ingest`).
- The user wants kb claims re-grounded against external sources
  (use `/ctx-kb-ground`).
- The user is asking a content question (use `/ctx-kb-ask`).
- The user wants to capture a quick finding (use
  `/ctx-kb-note`).

## Authority Boundary (vs Other Skills)

- **`/ctx-kb-site-review`**: mechanical structural audit. May
  coerce capitalization that the spec deems lossless (e.g.
  `Confidence: High` → `high`). May flag any other malformation
  in the closeout's `What changed` block. **May not** modify a
  claim, an `EV-###` row's content, a Confidence band, a topic
  page's prose, or a ledger state. Those require evidence
  judgment.
- **`/ctx-kb-ingest`**: handles anything this skill flags as
  evidence-dependent.
- **`/ctx-kb-ground`**: handles anything this skill flags as
  source-staleness.

## Usage Examples

```text
/ctx-kb-site-review
```

No arguments. The pass walks the kb in full.

## Pre-Write Gates

Three distinct refusals, each leaves zero residue:

- `.context/` missing → suggest `ctx init` and stop.
- `.context/kb/` missing → suggest `ctx init --upgrade` and
  stop.
- Kb scope undeclared (placeholder in `.context/kb/index.md`)
  → refuse with the scope message and stop.

## Process

1. **Verify pre-write gates.** Refuse cleanly if any gate fails.
   Zero residue on refusal.

2. **Walk topic pages.** For every
   `.context/kb/topics/<slug>/index.md` and every sibling
   sub-page:
   - **Status block check**: does the page have the
     four-field Status block (`Subject`, `Last verified`,
     `Author`, `Confidence`)? Missing fields → flag in
     `What changed`. Do not synthesize.
   - **Author field check**: `Author: hand-authored` is
     prohibited per `KB-RULES.md`. Flag (do not auto-coerce;
     human intent matters).
   - **Confidence band coercion**: `high|medium|low|speculative`
     are the only valid values. Coerce capitalization
     (`High` → `high`, `MEDIUM` → `medium`) silently and record
     in `What changed`. Any other malformation (e.g.
     `Confidence: probable`) is flagged for the user.
   - **`TBD-cite` markers**: count them per page. The
     Confidence floor for any page with `TBD-cite` is
     `speculative`. If the page's Confidence is above
     `speculative` while `TBD-cite` is present, flag (do not
     auto-demote; demotion is evidence work).
   - **`EV-###` citation resolution**: every `EV-###` cited on
     the page must resolve to a row in
     `.context/kb/evidence-index.md`. Unresolved IDs → flag.
   - **`## Related concepts in this kb` presence**: if the
     page is more than the lede + Status block AND the kb has
     plausibly adjacent topics, the section should be present.
     Absence is a soft flag (not auto-fixable).

3. **Walk `evidence-index.md`.**
   - **Duplicate `EV-###` IDs**: flag every duplicate; name
     both files / line numbers. The LLM cleanup pass (per
     spec's P1) handles renumbering, not this skill.
   - **Three-digit padding**: `EV-12` should be `EV-012`. Flag
     (do not auto-coerce; renumbering cascades to citations on
     topic pages, which is ingest work).
   - **Confidence band coercion**: same rule as topic pages.
   - **`occurred:` field on dated sources**: if the source-map
     row for the cited source has a `dated:` field but the
     evidence row lacks `occurred:`, flag. The temporal-
     precedence rule needs it.

4. **Walk `source-coverage.md`.**
   - **Ledger row mtime check**: for every row, compare the
     row's `Updated` cell against the actual file mtime of the
     source it points to (when the source is in-tree). Mismatch
     → flag (lying-to-the-ledger advisory). Do not auto-edit.
   - **Illegal state transitions**: flag any row whose state
     does not match an allowed transition from the prior state.
     Examples: `comprehensive → highlights-extracted` without
     an explicit `superseded` step. Do not auto-correct.
   - **Schema integrity**: every row must have the seven
     columns (`Source`, `Topic`, `State`, `EV coverage`,
     `Residue`, `Next action`, `Updated`). Missing columns →
     flag.

5. **Walk closeouts in `.context/ingest/closeouts/`.**
   - **Frontmatter integrity**: every closeout must have
     `sha`, `branch`, `mode`, `pass-mode`, `life-stage`,
     `generated-at`. Missing fields → flag (the handover-fold
     skips malformed closeouts; surface them so the user can
     fix or delete).
   - **Pass-mode body block**: every ingest closeout must
     have a `Pass-mode` body block whose `Declared:` value
     matches the frontmatter's `pass-mode:` field. Drift
     between the two is exactly the false-finish signal the
     redundancy exists to surface. Flag any drift.
   - **Adjacency pre-flight block**: every ingest closeout in
     `topic-page` mode must have an `Adjacency pre-flight`
     block whose value is either `none surfaced` or a
     structured slug-list. Free-prose values fail validation;
     flag.
   - **Cold-reader rubric**: every ingest closeout in
     `topic-page` mode must include the four-item rubric in
     `What changed`. Missing → flag.

6. **Walk `.context/kb/index.md`.**
   - **`CTX:KB:TOPICS` managed block**: should list every
     `.context/kb/topics/<slug>/index.md` currently on disk.
     Drift (slug on disk not in the block, or block entry with
     no matching folder) → recommend `ctx kb reindex` in the
     closeout's `Next pass hint`. Do not run the CLI from this
     skill.

7. **Write the site-review closeout.** Create
   `.context/ingest/closeouts/<TIMESTAMP>-site-review-closeout.md`
   with required frontmatter:

   ```yaml
   ---
   sha: <short>
   branch: <name>
   mode: site-review
   pass-mode: mechanical
   life-stage: <bootstrap|maintenance>
   generated-at: <RFC-3339>
   ---
   ```

   Body sections:
   - **Inputs**: count of topic pages, evidence rows, ledger
     rows, closeouts walked.
   - **What changed**: every coercion this pass actually
     applied (capitalization fixes); cite the file and the
     before/after. Empty if zero coercions.
   - **Flags**: every issue this pass detected but did not
     fix. Group by category: malformed Status blocks,
     unresolved `EV-###`, ledger mismatches, malformed
     closeouts, etc. Each flag names file + line + nature.
   - **Next pass hint**: explicit invocations to address each
     flag category (e.g. *"`/ctx-kb-ingest <slug>` to restore
     missing `EV-###` citation on `<page>`"*).

## Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `.context/` missing | Refuse; suggest `ctx init`. No residue. |
| `.context/kb/` missing | Refuse; suggest `ctx init --upgrade`. No residue. |
| Kb scope undeclared | Refuse with the scope message. No residue. |
| Zero topic pages on disk | Walk closeouts and ledger anyway. Note the bootstrap state in the closeout. Not a failure mode. |
| Zero closeouts on disk | Walk topic pages and ledger anyway. Note in the closeout body. Not a failure mode. |
| `Confidence: High` (capitalization drift) | Coerce to `high` silently; record in `What changed`. |
| `Confidence: probable` (unknown band) | Flag for the user; do not coerce. |
| `Author: hand-authored` | Flag for the user; do not coerce (human intent matters). |
| Duplicate `EV-###` ID across files | Flag both files; defer renumbering to the LLM cleanup pass per spec's P1. |
| `EV-12` (missing zero-pad) | Flag for the user; do not auto-pad (cascades to citations). |
| Unresolved `EV-###` on a topic page | Flag; recommend `/ctx-kb-ingest <slug>` in `Next pass hint`. |
| Ledger row `Updated` predates source file mtime | Flag (lying-to-the-ledger advisory). Do not auto-edit. |
| Illegal ledger transition (e.g. `comprehensive → highlights-extracted` without `superseded`) | Flag; recommend the corrective ingest invocation. Do not auto-correct. |
| Closeout missing `pass-mode` frontmatter field | Flag; the handover-fold skips malformed closeouts so the user can fix or delete. |
| Closeout body's `Pass-mode` `Declared:` disagrees with frontmatter `pass-mode:` | Flag (false-finish signal); recommend hand-edit. |
| Closeout's `Adjacency pre-flight` is free prose instead of `none surfaced` or a slug-list | Flag; recommend hand-edit to structured form. |
| `CTX:KB:TOPICS` managed block drift | Recommend `ctx kb reindex` in `Next pass hint`; do not run the CLI from this skill. |
| `TBD-cite` on a page with Confidence above `speculative` | Flag; do not auto-demote (demotion is evidence work for `/ctx-kb-ingest`). |
| Sibling sub-page exists with no link from `index.md` | Flag; recommend hand-edit or `/ctx-kb-ingest <slug>` to extend. |

## Anti-Patterns

- Auto-fixing anything that requires evidence judgment
  (Confidence promotion/demotion, claim text edits, `EV-###`
  renumbering, ledger state changes, prose synthesis).
- Skipping the closeout once pre-write gates pass.
- Hand-editing `INBOX.md` or `SESSION_LOG.md` (other skills'
  surfaces; never this one's).
- Coercing `Author: hand-authored` to anything else. The user's
  intent matters; flag and wait.
- Auto-renumbering duplicate `EV-###` IDs. The cascade to
  citations is ingest work; this skill flags only.

## Output Contract

For pre-write refusals, return only the specified refusal text
and stop. No residue.

For passes that clear pre-write gates, end with this structured
summary:

- **Walked**: counts (topic pages, evidence rows, ledger rows,
  closeouts).
- **Coercions applied**: count + one-line categories (e.g.
  *"3 capitalization fixes on Confidence bands"*).
- **Flags raised**: count + categories (e.g. *"2 unresolved
  EV-### citations; 1 ledger mtime mismatch"*).
- **Closeout**: filename on its own line.
- **Next-recommended-action**: explicit invocations to address
  each flag category (or `none` if the kb is clean).

## Quality Checklist

Before reporting completion, verify:

- [ ] Pre-write gates passed (or the matching refusal was
      returned with zero residue).
- [ ] Every coercion applied is recorded in `What changed` with
      file + before/after.
- [ ] Every flag is recorded in `Flags` with file + line +
      nature.
- [ ] No topic-page prose was edited, no `EV-###` row was
      modified, no Confidence band was promoted/demoted, no
      ledger state was changed.
- [ ] Closeout written with all required frontmatter fields.

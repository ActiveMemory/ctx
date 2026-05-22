# KB Rules

Stable workflow contract for the editorial knowledge-ingestion
pipeline shipped by `ctx`. These rules govern the
**`.context/ingest/`** working directory and the mode skills
(`/ctx-kb-ingest`, `/ctx-kb-ask`, `/ctx-kb-site-review`,
`/ctx-kb-ground`, `/ctx-kb-note`).

This file is project memory. Hand-edit it the same way you
hand-edit the five canonical `.context/` files. Round-trips
through git like any Markdown file.

The filename is `KB-RULES.md`, not `CONSTITUTION.md`, to avoid
collision with `.context/CONSTITUTION.md`. Per DECISIONS.md
"Editorial constitution at `.context/ingest/KB-RULES.md`, not
`CONSTITUTION.md`" (2026-05-10).

---

## Four Inviolable Rules

These four rules are load-bearing for the entire workflow.
Breaking any one of them silently corrupts the knowledge base.

1. **Skills are the sole writers of `INBOX.md`.** The inbox is
   an audit trail of the last skill run, never a hand-edit
   form. Args or inline natural-language context drive what the
   skill writes; cold invocation with no input causes the skill
   to refuse cleanly. The pipeline takes no hand-edited config
   file; sources are supplied at invocation time. To configure
   external grounding inputs, edit `grounding-sources.md`.

2. **The handover is the sole authoritative recall artifact.**
   `/ctx-wrap-up` writes a handover at the end of every session
   (always delegating to `/ctx-handover` as its final step,
   regardless of whether `.context/kb/` exists). Handovers live
   at `.context/handovers/<TS>-<slug>.md` (timestamped so
   concurrent agent runs never overwrite). `Do you remember?`
   reads the latest handover unconditionally; when
   `.context/kb/` exists, it additionally folds any closeouts
   whose `generated-at` postdates the handover.
   `SESSION_LOG.md` is mid-flight working memory; it is **not**
   read on session start.

3. **Provenance lives on operational artifacts, not knowledge
   artifacts.** `SESSION_LOG.md` entries, closeouts, and
   handovers carry `sha=<short>` and `branch=<name>`. KB prose,
   glossary, domain-decisions, contradictions,
   outstanding-questions, timeline, schemas, and the rendered
   site stay SHA-free. Single exception: `evidence-index.md`
   entries pointing at in-repo files pin to a SHA at extraction
   time so cited bytes are recoverable via `git show`.

4. **Schemas ship with shape, not content.** Each
   `schemas/*.md` lists fields and one worked example; no
   domain entries. The drift-prevention property the schemas
   folder exists for is preserved without prescribing values
   per project.

---

## Pass-Mode Contract

Every `ctx kb ingest` invocation declares its mode **before any
source extraction begins**:

| Mode             | Mints prose? | Mints EV-### ? | Touches topic page? | Default? |
|------------------|--------------|----------------|---------------------|----------|
| `topic-page`     | yes          | yes            | yes (create/extend) | yes      |
| `triage`         | no           | **no**         | no                  | no       |
| `evidence-only`  | no           | yes (tagged)   | no                  | no       |

The declaration is a contract, not a label. The skill emits a
three-line block in the response stream:

> **Pass-mode:** `<mode>`
> **Reason:** `<one sentence; required when non-default>`
> **Definition of done:** `<mode-specific completion criterion>`

**Mid-pass mode-switching is forbidden.** If the work in flight
no longer fits, abort with a partial closeout citing what was
done, and recommend re-invocation under the correct mode.

**Inferring `evidence-only` to avoid topic-page validation is a
hard anti-pattern.** Mode is explicit-request-only; size,
ambiguity, time pressure, and operator convenience are NOT
valid triggers.

---

## Topic-Page Circuit Breaker

A pass in `topic-page` mode MAY NOT report
`topic-page: produced` or `topic-page: extended` unless ALL of
the following are true at completion:

1. `.context/kb/topics/<slug>/index.md` (or a sibling sub-page)
   exists and was created or extended in this pass.
2. The page cites at least one `EV-###` row that resolves to
   `evidence-index.md`.
3. `ctx kb site build` ran clean (or its failure is named in
   the closeout's `Next pass hint` AND the pass reports
   `topic-page: deferred`).
4. The cold-reader orientation rubric records `Result: pass`
   in the closeout (all four items at `yes`).

Any failure → `topic-page: deferred` and the source-coverage
ledger advances to `topic-page-drafted` (not `comprehensive`).

---

## Source-Coverage Ledger (State Machine)

`.context/kb/source-coverage.md` is a state machine over every
source the kb has touched. Allowed transitions:

| state                   | next states |
|-------------------------|-------------|
| `discovered`            | `admitted`, `skipped` |
| `admitted`              | `highlights-extracted`, `partially-ingested`, `topic-page-drafted`, `comprehensive` |
| `highlights-extracted`  | `partially-ingested`, `topic-page-drafted`, `comprehensive` |
| `partially-ingested`    | `topic-page-drafted`, `comprehensive` |
| `topic-page-drafted`    | `comprehensive` |
| `comprehensive`         | terminal until source updates |
| `superseded`            | terminal |
| `skipped`               | terminal until scope changes |

Every pass updates the ledger before writing the closeout.
**Lying to the ledger is a hard anti-pattern.** Set the state
honestly even when it means recording incomplete work.

---

## Topic-Adjacency Pre-Flight

Before resolving the topic in `topic-page` mode, scan the
ledger for rows whose state is **not** in
`{comprehensive, skipped, superseded}` AND whose `Topic` is
plausibly adjacent (shared first segment of a slug,
shared vendor / surface, explicit cross-references).

For each adjacent incomplete topic surfaced, the pass MUST:

1. Acknowledge it in `## Related concepts in this kb` on the
   topic page being authored.
2. Surface it in the closeout's `Adjacency pre-flight` block.
3. Surface it in the response contract's
   `Adjacent topics noted` field.

**Do NOT enumerate `EV-###` IDs by name in the adjacency
block.** Use count + location (*"seventeen rows in
`evidence-index.md`"*) instead; naming an EV row from a
lower-confidence sibling demotes the floor of cited bands.

Silence is not the same as a clean pre-flight; when zero
matches, explicitly record *"no incomplete adjacent topics
surfaced"*.

---

## Cold-Reader Orientation Rubric

Four yes/no items recorded in the closeout's `What changed`
section, in `topic-page` mode:

```
Cold-reader orientation:
- Concept clear?                yes|no: <short note>
- Why this kb cares clear?      yes|no: <short note>
- Canonical evidence reachable? yes|no: <short note>
- Boundaries clear?             yes|no: <short note>
Result: pass | fail
```

`Result: pass` requires all four `yes`. Any `no` →
`Result: fail` → circuit-breaker fails → `topic-page:
deferred`.

---

## Life-Stage Discipline

Count `.context/kb/topics/*/index.md` before the pass begins
synthesizing:

- `< 5` topic pages → **bootstrap** mode. Skip reconciliation
  ceremony; synthesize topic pages aggressively. Exception:
  surface a contradiction even in bootstrap if the new
  material plainly contradicts existing kb claims.
- `>= 5` topic pages → **maintenance** mode. Apply full
  reconciliation: claim laddering, demotion, contradiction
  detection.

Document the life-stage call in the closeout's frontmatter
(`life-stage:`).

---

## Evidence Discipline

- Every claim in `.context/kb/` carries at least one citation
  in `evidence-index.md`. Claims without citations stay in
  `outstanding-questions.md` until grounded.
- Citations name the source by **short name** (defined in
  `source-map.md`) plus a locator: line range for files,
  timestamp for transcripts, anchor for URLs, symbol for code.
- In-repo citations include a `sha:` field set to the short
  SHA at extraction time. Out-of-repo citations omit `sha:`.
- A claim with **only** verbal-source backing (transcript,
  meeting note) is `confidence: low` until reinforced by a
  written source.

---

## Confidence Bands

Every claim and definition in `.context/kb/` carries an
explicit confidence band:

- **`high`**: corroborated by ≥ 2 independent sources, or by
  one source plus a working code-level check (compiles, runs,
  matches).
- **`medium`**: single authoritative source, no independent
  corroboration; or two sources that agree on the conclusion
  but not the mechanism.
- **`low`** — single source only; or sources disagree but the
  claim is the best current synthesis.
- **`speculative`** — no source backing; an inference the
  human or agent flagged for follow-up.

`speculative` content does **not** ship in the rendered site.
`low`-confidence content ships only when paired with the
matching `outstanding-questions.md` entry.

**Floor-of-cited-bands rule:** a topic page's Confidence is
the *lowest* band cited on the page. Refuse to set Confidence
above the floor. Refuse to set above `speculative` while any
`TBD-cite` remains.

---

## Demotion Policy

When new evidence contradicts an existing claim:

1. Add a row to `contradictions.md` naming both EV-### IDs and
   one-line summary of the disagreement.
2. Demote the older claim one band
   (`high → medium → low → speculative`), citing the new EV
   row as the cause.
3. Open an `outstanding-questions.md` entry naming both sides
   and what evidence would resolve.

**Never renumber or delete `EV-###` rows.** Demote in place.

---

## Authority Boundary

- `/ctx-kb-ingest` writes prose AND evidence rows AND topic
  scaffold AND cross-links AND ledger updates in the same
  pass. That combination is unique to ingest.
- `/ctx-kb-ask` is Q&A grounded in the kb; read-only on prose;
  refuses to web-jump; flags gaps the kb cannot answer.
- `/ctx-kb-site-review` is a structural audit; mechanical
  fixes only. Defers anything that requires evidence judgment.
- `/ctx-kb-ground` re-grounds against external sources listed
  in `grounding-sources.md`.
- `/ctx-kb-note` is a lightweight capture into
  `ingest/findings.md`; never writes to a topic page or to
  evidence-index.

**Topic-page file creation is performed only by `ctx kb topic
new`.** Skills invoke that CLI; they do not synthesize a
scaffold by hand.

---

## Closeout Shape

Every pass that clears pre-write gates writes a closeout under
`.context/ingest/closeouts/<TIMESTAMP>-<mode>-closeout.md`
with required frontmatter:

```yaml
---
sha: <short>
branch: <name>
mode: <ingest|ask|site-review|ground|note>
pass-mode: <topic-page|triage|evidence-only>
life-stage: <bootstrap|maintenance>
generated-at: <RFC-3339>
---
```

Body sections (mode-aware): Inputs, Pass-mode block, Topic(s)
touched, What changed (including the Cold-reader rubric in
topic-page mode), New questions, New contradictions,
Confidence drift, Source-coverage updates, Overflow, Adjacency
pre-flight, Next pass hint.

Closeouts are append-never-rewrite. Archived closeouts are
immutable.

---

## SESSION_LOG Line Shape

Each phase-boundary line:

```
[YYYY-MM-DD HH:MM:SS sha=<short> branch=<name>] phase=<name> status=<done|partial|blocked> note=<≤120 chars>
```

`<short>` is the 7-char git short SHA. `<name>` is the current
symbolic ref or `detached` (HEAD is not on a branch).
`<phase>` is one of `resolve`, `synthesise`, `reconcile`,
`closeout`.

---

## Hard Anti-Patterns

- Treating closeout existence as topic-page validation.
- Skipping the topic-page circuit breaker in topic-page mode.
- Inferring `evidence-only` mode from source size, ambiguity,
  time pressure, or operator convenience.
- Mid-pass mode-switching (abort and re-invoke instead).
- Hiding incomplete coverage under a comprehensive-looking
  closeout (lying to the ledger).
- Skipping the topic-adjacency pre-flight or running it but
  failing to acknowledge surfaced incomplete adjacent topics.
- Inventing claims beyond what the source backs.
- Inventing `EV-###` citations to make a page look complete.
- Promoting claims above `speculative` without an
  `evidence-index.md` row.
- Promoting a topic page above its weakest cited band.
- Setting `Confidence` above `speculative` while `TBD-cite`
  remains.
- Renumbering or deleting `EV-###` rows when reconciling.
- Skipping the closeout once the pass clears pre-write gates.
- Bypassing `ctx kb topic new` when scaffolding a page.
- Running maintenance discipline against a bootstrap-stage kb.
- Hand-editing `INBOX.md`.

---

## See Also

- `specs/kb-editorial-pipeline.md` — full spec, including
  failure analysis and open questions.
- `OPERATOR.md` — human-facing framing.
- `PROMPT.md` — fallback auto-router when no skill is
  installed.

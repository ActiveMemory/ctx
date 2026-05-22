# Ingest Mode

The primary editorial pass. Turns semi-structured inputs
(notes, transcripts, source code, diagrams, upstream docs)
into evidence-backed Markdown kb content with explicit
confidence bands.

Read `KB-RULES.md` first; it carries the inviolable rules,
the pass-mode contract, the topic-page circuit breaker, the
source-coverage state machine, the topic-adjacency pre-flight,
the cold-reader rubric, evidence discipline, confidence
bands, the demotion policy, and the closeout shape. This
file describes **how a single ingest pass executes**. It does
not duplicate `KB-RULES.md`; it cites it.

The `/ctx-kb-ingest` skill reads this file before doing
anything else.

---

## Inputs

The skill takes either a folder (recurse), a list of paths,
a URL, an MCP resource, or inline natural-language context
naming the materials. **No input ⇒ refuse cleanly:**

> no sources provided; pass a folder, a URL, an MCP resource,
> or describe the materials inline.

Refuse-on-empty is the contract; do not prompt the user for
sources mid-skill. Exit non-zero so the operator re-invokes
with arguments.

---

## Pre-Write Gates (All Must Pass Before Any Extraction)

1. `.context/` and `.context/kb/` exist.
2. `.context/kb/index.md` has a non-placeholder `## Scope`
   section (no `<!-- TODO: ... -->` marker present).
3. `.context/ingest/KB-RULES.md` exists and was read this
   pass.
4. Source paths supplied on the command line all resolve. A
   missing path is a refuse condition, not a partial-extract
   condition.

Any gate failure ⇒ refuse with the recovery hint named in
`KB-RULES.md`'s error table. Do not partially execute.

---

## Up-Front Pass-Mode Declaration (MANDATORY)

After pre-write gates pass and **before** topic resolution,
emit a visible three-line block:

> **Pass-mode:** `<topic-page|triage|evidence-only>`
> **Reason:** `<one sentence; required for non-default modes>`
> **Definition of done:** `<mode-specific completion criterion>`

The declaration is a contract, not a label. The closeout's
body block restates these three lines verbatim, and the
closeout's frontmatter records `pass-mode:` to match. Doctor
advisory flags mismatches between the body and frontmatter.

**Mode selection (see `KB-RULES.md` "Pass-mode contract"):**

- Default is `topic-page`.
- `triage` fires when the user supplied multiple disparate
  sources with no clear single topic, OR explicitly invoked
  triage language ("triage these"; "just classify").
- `evidence-only` fires **only** on explicit user request
  matching valid-trigger criteria ("just mint EV rows";
  "backfill evidence"). Size, ambiguity, time pressure, and
  operator convenience are NOT valid triggers. Inferring
  `evidence-only` to dodge topic-page validation is a hard
  anti-pattern.

**Mid-pass mode-switching is forbidden.** If the work in
flight no longer fits, abort with a partial closeout citing
what was done so far and recommend re-invocation under the
correct mode. Never silent-switch.

---

## Procedure

### 1. Resolve Topic and Run Adjacency Pre-Flight

Identify the slug. In `topic-page` mode, before scaffolding
anything:

1. Read `.context/kb/source-coverage.md`.
2. Scan for rows whose state is **not** in
   `{comprehensive, skipped, superseded}` and whose `Topic`
   is plausibly adjacent (shared first slug segment, shared
   vendor / surface, explicit cross-reference).
3. Record the result in the closeout's `Adjacency pre-flight`
   block. Use count + location (*"seventeen rows in
   `evidence-index.md`"*); do **not** name individual
   `EV-###` IDs (naming a lower-confidence sibling's row
   demotes the floor of cited bands on this page).
4. If zero matches, record explicit
   *"no incomplete adjacent topics surfaced"*. Silence is
   not allowed.

Each surfaced incomplete adjacent topic MUST be acknowledged
in `## Related concepts in this kb` on the page being
authored.

### 2. Life-Stage Check

Count `.context/kb/topics/*/index.md` before synthesizing:

- `< 5` topic pages → **bootstrap**. Skip reconciliation
  ceremony; synthesize aggressively. Exception: surface a
  contradiction even in bootstrap if the new material
  plainly contradicts existing kb claims.
- `>= 5` topic pages → **maintenance**. Apply full
  reconciliation discipline (laddering, demotion,
  contradictions).

Record the life-stage call in the closeout frontmatter
(`life-stage:`) and `What changed` section.

### 3. Scaffold (Topic-Page Mode Only)

If `.context/kb/topics/<slug>/index.md` does not exist,
shell out to `ctx kb topic new "<name>"`. **Topic-page file
creation is performed only by `ctx kb topic new`.** Skills
do not synthesize the scaffold by hand.

If the slug already exists, append/extend the existing
`index.md` rather than creating a new one.

### 4. Extract and Synthesise

For each input, extract atomic claims. One claim = one row
in `evidence-index.md`. Each row carries: claim text (single
sentence, declarative, present tense), source short name +
locator, `sha:` for in-repo citations only, confidence band
per `KB-RULES.md`.

Mode-specific output:

- `topic-page`: mint prose AND `EV-###` rows AND topic
  scaffold AND cross-links AND ledger updates.
- `triage`: classify sources into the inbox, advance ledger
  rows to `discovered` or `admitted`. **No EV rows minted.**
  No topic page touched. No prose written.
- `evidence-only`: mint `EV-###` rows tagged with the source.
  No topic page. No prose. No glossary / domain-decisions
  edits.

### 5. Reconcile (Maintenance Life-Stage Only)

For each new claim, decide:

- **Net new** ⇒ append to the appropriate kb file
  (`glossary.md` for terms; `domain-decisions.md` for
  editorial calls; `timeline.md` for events; the relevant
  topic-page prose for narrative context).
- **Reinforces existing** ⇒ promote the existing claim's
  confidence band per `KB-RULES.md`; cross-link the new
  evidence row.
- **Contradicts existing** ⇒ add a `contradictions.md` row;
  demote per the demotion policy; open an
  `outstanding-questions.md` entry.

In bootstrap, the reconciliation ceremony is skipped except
for the contradiction exception.

### 6. Advance the Source-Coverage Ledger

Update `.context/kb/source-coverage.md` for every source
touched in this pass, per the state-machine transitions in
`KB-RULES.md`. Illegal transitions are refused at write time
by `AdvanceLedger`. **Lying to the ledger is a hard
anti-pattern**: record `topic-page-drafted` honestly when
the page is incomplete; do not claim `comprehensive` to make
the closeout look clean.

### 7. Topic-Page Circuit Breaker (Topic-Page Mode Only)

Before claiming `topic-page: produced` or
`topic-page: extended`, verify ALL FOUR invariants from
`KB-RULES.md` "Topic-page circuit breaker":

1. The page exists and was created/extended this pass.
2. The page cites at least one `EV-###` row that resolves to
   `evidence-index.md`.
3. `ctx kb site build` ran clean (or its failure is named in
   `Next pass hint` AND the pass reports `topic-page:
   deferred`).
4. The cold-reader orientation rubric records
   `Result: pass`.

Any failure ⇒ `topic-page: deferred`; ledger advances to
`topic-page-drafted` (not `comprehensive`).

### 8. Cold-Reader Orientation Rubric (Topic-Page Mode Only)

Record in the closeout's `What changed` section:

```
Cold-reader orientation:
- Concept clear?                yes|no: <short note>
- Why this kb cares clear?      yes|no: <short note>
- Canonical evidence reachable? yes|no: <short note>
- Boundaries clear?             yes|no: <short note>
Result: pass | fail
```

All four `yes` ⇒ `Result: pass`. Any `no` ⇒ `Result: fail`
⇒ circuit-breaker fails ⇒ `topic-page: deferred`. Name
which items returned `no` in the closeout body.

If `Boundaries clear?` is `no` because `index.md` has
outgrown a single page, **propose** a sub-page split (e.g.
`.context/kb/topics/<slug>/security.md`) and wait for
operator confirmation. Never auto-split.

### 9. Write the Closeout

Path:
`.context/ingest/closeouts/<TS>-ingest-closeout.md`.

Required frontmatter (see `KB-RULES.md` "Closeout shape"):
`sha`, `branch`, `mode: ingest`, `pass-mode`, `life-stage`,
`generated-at`.

Body sections (mode-aware): Inputs, Pass-mode block, Topic(s)
touched, What changed (including the Cold-reader rubric in
topic-page mode), New questions, New contradictions,
Confidence drift, Source-coverage updates, Overflow,
Adjacency pre-flight, Next pass hint.

Closeouts are append-never-rewrite. Archived closeouts are
immutable.

### 10. Append the SESSION_LOG Line

At the closeout phase boundary, append a line to
`.context/ingest/SESSION_LOG.md` in the exact shape defined
by `KB-RULES.md` "SESSION_LOG line shape".

---

## What This Mode Does NOT Do

- Does not web-fetch new external sources; that is
  `ctx kb ground`'s job. Inputs are project-internal or
  user-supplied URLs.
- Does not mechanically audit existing kb structure; that
  is `ctx kb site-review`.
- Does not answer questions from the kb; that is
  `ctx kb ask`.
- Does not write to canonical files (`TASKS.md`,
  `DECISIONS.md`, `LEARNINGS.md`, `CONVENTIONS.md`). The
  authority boundary in `KB-RULES.md` is strict.

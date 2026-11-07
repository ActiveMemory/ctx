---
name: ctx-kb-ingest
description: Editorial knowledge-ingestion pass. Reads sources the user supplies, declares its pass-mode (topic-page / triage / evidence-only) before extraction, and is held to mode-specific completion semantics. The topic page is the deliverable; the closeout is the audit trail.
---

# Editorial Ingestion Pass

This skill is the **single editorial pass** for adding knowledge
to `.context/kb/`. It reads materials the user supplies, decides
which topic page(s) they belong to, finds-or-creates those pages,
writes synthesized prose section by section, mints `EV-###` rows
in the structured layer as it cites them, cross-links neighboring
topics, updates the source-coverage ledger, and writes a closeout
file under `.context/ingest/closeouts/`.

The split between "extract claims" and "write the topic page" is
mechanical, not editorial. A student reading a book does not
extract a glossary first and synthesize later, they read and write
at the same time. This skill matches that model: the user supplies
*intent and material*; the skill does *judgment and typing*.

**The topic page is the deliverable. The closeout is the audit
trail. The closeout never substitutes for the page.** Intermediate
artifacts (EV rows, glossary entries, candidate-source registries,
closeouts) are valuable, but they do not validate topic-page work
by themselves; only the topic page does.

Authoritative background reading lives at
`.context/ingest/KB-RULES.md` and `specs/kb-editorial-pipeline.md`.
This skill encodes the workflow contract; the rules file is the
constitution. Hand-edit `KB-RULES.md` to evolve the contract; do
not paraphrase it here.

## When to Use

- The user supplies one or more sources (paths, URLs, MCP
  resources, inline natural-language descriptions) and wants them
  read into the kb.
- The user says "ingest the transcripts", "pull this into the
  kb", "add evidence from <source>", "extract claims from this
  call", or invokes the explicit slash form with paths.
- A prior pass left residue (a `topic-page-drafted` ledger row,
  a `Next pass hint` in a closeout) and the user is resuming.

## When NOT to Use

- The user asked a question about the kb (use `/ctx-kb-ask`).
- The user wants a structural audit of the kb (use
  `/ctx-kb-site-review`).
- The user wants to re-ground existing kb claims against
  external sources (use `/ctx-kb-ground`).
- The user wants to park a quick finding for the next ingest
  (use `/ctx-kb-note`).
- No sources were supplied (refuse-on-empty; see §Refuse-on-empty
  below).
- `.context/kb/` does not exist (refuse with the no-pipeline
  message in §Pre-write gates).

## Authority Boundary (vs Other Skills)

- **`/ctx-kb-ingest`**: primary editorial pass. Reads materials
  (in-tree paths, out-of-tree paths, URLs, MCP resources, inline
  references); writes topic pages
  (`.context/kb/topics/<slug>/index.md`, plus optional sibling
  sub-pages); mints evidence, glossary, source-map, timeline,
  contradictions, outstanding questions; cross-links into existing
  kb topology; updates the source-coverage ledger; writes
  closeout. **Topic-page file creation is performed only by
  `ctx kb topic new`**: this skill MAY invoke that CLI as part
  of a topic-page pass, but it MUST NOT synthesize or write a
  scaffold directly. This preserves the public editorial workflow
  (`/ctx-kb-ingest`) and the actual scaffold authority
  (`ctx kb topic new`) as two separate facts.
- **`/ctx-kb-ask`**: Q&A grounded in the kb. Read-only on prose;
  refuses to web-jump; flags gaps the kb cannot answer.
- **`/ctx-kb-site-review`**: structural audit; mechanical fixes
  only. Defers anything that requires evidence judgment.
- **`/ctx-kb-ground`**: external grounding against
  `grounding-sources.md`; advances ledger rows for sources it
  refreshes.
- **`/ctx-kb-note`**: lightweight capture into
  `.context/ingest/findings.md`; never writes to a topic page or
  to `evidence-index.md`.

This skill writes prose AND evidence rows AND scaffold (via CLI)
AND cross-links AND ledger updates in the same pass; that
combination is unique to ingest.

## Usage Examples

```text
/ctx-kb-ingest ./inputs/2026-04-12-call.md "cursor hooks"
/ctx-kb-ingest ./inputs/your-domain/
/ctx-kb-ingest https://cursor.com/docs/hooks
/ctx-kb-ingest ./a.md ./b.md "incident retros"
/ctx-kb-ingest --inline "the four transcripts under inputs/ \
  and the pool.go file" "connection pooling"
```

## Input Contract

**Sources**, supplied as one or more of:

- **Paths**: folder to recurse, single file, list of files.
- **URLs**: primary-source web pages.
- **MCP resources**: named resources from connected MCP servers.
- **Inline gestures**: natural-language naming the materials.
- **Open invitation**: *"feel free to search for more"*. The
  skill gets web-search and MCP-discovery authority for this
  pass; hard cap of 50 total sources.

**Optional second argument, topic name**, e.g. *"cursor hooks"*.
When omitted, the skill proposes one at §3 of Process and
confirms with the user before any extraction work. Naming the
topic up front skips that round-trip.

## Refuse-on-Empty

The skill writes to the kb; refuse-on-empty is the default. If
the invocation supplied no sources and no inline gesture, return
exactly:

> no sources provided; pass a folder, a URL, an MCP resource, or
> describe the materials inline.

Stop. Do not prompt for sources interactively, do not invent a
topic, do not propose a triage pass on imagined material. The CLI
enforces this independently via `cmd/ingest`.

## Pre-Write Gates

Three distinct refusals, each leaves zero residue (no
`INBOX.md` rewrite, no `SESSION_LOG.md` entry, no claim
extraction, no ledger update, no closeout, no topic-page edits):

- `.context/` missing entirely → suggest `ctx init` and stop.
- `.context/ingest/` missing (project initialised before this
  spec shipped) → suggest `ctx init --upgrade` and stop.
- Kb scope undeclared (`.context/kb/index.md` missing, contains
  the `TODO: declare what this kb covers` placeholder, has no
  `## Scope` H2, or `## Scope` lacks substantive
  non-placeholder prose):

  > kb scope is undeclared. Open `.context/kb/index.md` and
  > replace the TODO placeholder with a one-paragraph scope
  > statement that names what is in scope and what is out.
  > `/ctx-kb-ingest` refuses to ingest until scope is declared.

## Pass-Mode Contract

Every invocation MUST classify itself as exactly one of three
modes **before any source extraction begins**. The mode commits
the pass to a specific definition of done; the skill is held to
that definition and may not narrate success on residue belonging
to a different mode. Full mode semantics live in
`.context/ingest/KB-RULES.md` §Pass-mode contract.

| Mode             | Mints prose? | Mints `EV-###`? | Touches topic page?   | Default? |
|------------------|--------------|------------------|------------------------|----------|
| `topic-page`     | yes          | yes              | yes (create/extend)    | yes      |
| `triage`         | no           | **no**           | no                     | no       |
| `evidence-only`  | no           | yes (tagged)     | no                     | no       |

**Mode selection rules.** Default is `topic-page`. `triage` fires
only when the user supplied multiple disparate sources with no
clear single topic, OR explicitly invoked triage language.
`evidence-only` fires only on explicit user request matching the
valid-trigger criteria (*"just mint EV rows"*, *"backfill
evidence"*); the skill MAY NOT infer it from source size,
ambiguity, time pressure, or operator convenience.

**Up-front declaration (mandatory).** Before extraction begins
(after pre-write gates pass and **before** topic resolution), emit
a visible pre-work declaration in the response stream:

> **Pass-mode:** `<mode>`
> **Reason:** `<one sentence; required when non-default>`
> **Definition of done:** `<mode-specific completion criterion>`

The declaration is a contract, not a label. The skill is bound to
it for the rest of the pass.

**Mid-pass mode-switching is forbidden.** If the work in flight
no longer fits, abort with a partial closeout citing what was
done, and recommend re-invocation under the correct mode. Silent
mode-drift is a hard anti-pattern.

## Topic-Page Circuit Breaker

A pass operating in `topic-page` mode MAY NOT report
`topic-page: produced` or `topic-page: extended` unless ALL of
the following are true at completion:

1. `.context/kb/topics/<slug>/index.md` (or a sibling sub-page
   like `.context/kb/topics/<slug>/<sub>.md`) exists and was
   created or extended in this pass.
2. The page cites at least one `EV-###` row that resolves to
   `evidence-index.md`.
3. `ctx kb site build` ran clean (or its failure is named in the
   closeout's `Next pass hint` AND the pass reports
   `topic-page: deferred`).
4. The cold-reader orientation rubric records **`Result: pass`**
   in the closeout's `What changed` section. All four rubric
   items must be `yes`.

Any failure → `topic-page: deferred` and the source-coverage
ledger advances to `topic-page-drafted` (not `comprehensive`).
This invariant prevents intermediate residue from being treated
as topic-page success. **Topic-page validation requires the
topic page.**

## Source-Coverage Ledger

`.context/kb/source-coverage.md` is a state machine over every
source the kb has touched. Allowed transitions live in
`KB-RULES.md` §Source-coverage ledger; do not paraphrase them
here. Every pass updates the ledger before writing the closeout.
**Lying to the ledger is a hard anti-pattern.** Set the state
honestly even when it means recording incomplete work.

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
`Result: fail` → circuit-breaker fails → `topic-page: deferred`.

## Life-Stage Check

Count `.context/kb/topics/*/index.md` pages **before** this pass
begins synthesizing:

- `< 5` topic pages → **bootstrap** mode. Skip reconciliation
  ceremony; synthesize topic pages aggressively. Exception:
  surface a contradiction even in bootstrap if the new material
  plainly contradicts existing kb claims.
- `>= 5` topic pages → **maintenance** mode. Apply full
  reconciliation discipline (laddering, demotion, contradiction
  detection).

Document the life-stage call in the closeout's frontmatter
(`life-stage:`) and `What changed` section.


## Extended Reference

Read `references/extended.md` in this skill directory for the remaining specification and detailed steps.

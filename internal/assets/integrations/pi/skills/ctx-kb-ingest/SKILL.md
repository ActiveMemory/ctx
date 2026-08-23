---
name: ctx-kb-ingest
description: "Editorial knowledge-ingestion pass. Reads sources the user supplies, declares its pass-mode (topic-page / triage / evidence-only) before extraction, and is held to mode-specific completion semantics. The topic page is the deliverable; the closeout is the audit trail."
---

Single editorial pass that adds knowledge to `.context/kb/`.
Reads materials the user supplies, decides which topic page(s)
they belong to, finds-or-creates those pages, writes synthesized
prose section by section, mints `EV-###` rows as it cites them,
cross-links neighbouring topics, updates the source-coverage
ledger, and writes a closeout under `.context/ingest/closeouts/`.

Authoritative background reading lives at
`.context/ingest/KB-RULES.md`. This skill encodes the workflow
contract; the rules file is the constitution. Hand-edit
`KB-RULES.md` to evolve the contract; do not paraphrase it
here.

## When to Use

- The user supplies one or more sources (paths, URLs, MCP
  resources, inline natural-language descriptions) and wants
  them read into the kb.
- The user says "ingest the transcripts", "pull this into the
  kb", "add evidence from <source>", or invokes the slash form
  with paths.

## When NOT to Use

- The user asked a question about the kb (use `/ctx-kb-ask`).
- The user wants a structural audit (use `/ctx-kb-site-review`).
- The user wants external re-grounding (use `/ctx-kb-ground`).
- The user wants to park a quick finding (use `/ctx-kb-note`).
- No sources were supplied (refuse-on-empty).

## Input

Sources, supplied as one or more of: paths (file or folder),
URLs, MCP resources, or inline natural-language gestures.
Optional second argument is the topic name; when omitted the
skill proposes one and confirms before extraction.

## Refuse-on-Empty

If the invocation supplied no sources, return exactly:

> no sources provided; pass a folder, a URL, an MCP resource, or
> describe the materials inline.

Stop. The CLI enforces this independently.

## Pre-Write Gates

- `.context/` missing → suggest `ctx init` and stop.
- `.context/ingest/` missing → suggest `ctx init --upgrade`
  and stop.
- Kb scope undeclared (`.context/kb/index.md` missing, or its
  `## Scope` H2 holds the `TODO` placeholder, or lacks
  substantive non-placeholder prose):

  > kb scope is undeclared. Open `.context/kb/index.md` and
  > replace the TODO placeholder with a one-paragraph scope
  > statement that names what is in scope and what is out.

## Pass-Mode Contract

Every invocation classifies itself as exactly one mode before
extraction begins. Full semantics in
`.context/ingest/KB-RULES.md` §Pass-mode contract.

| Mode             | Mints prose? | Mints `EV-###`? | Touches topic page?   | Default? |
|------------------|--------------|------------------|------------------------|----------|
| `topic-page`     | yes          | yes              | yes (create/extend)    | yes      |
| `triage`         | no           | no               | no                     | no       |
| `evidence-only`  | no           | yes (tagged)     | no                     | no       |

Default is `topic-page`. `triage` fires when sources are
disparate with no clear single topic, or the user explicitly
asks for triage. `evidence-only` fires only on explicit user
request ("just mint EV rows", "backfill evidence"); never
inferred from source size or operator convenience.

Before extraction, emit the declaration in the response stream:

> **Pass-mode:** `<mode>`
> **Reason:** `<one sentence; required when non-default>`
> **Definition of done:** `<mode-specific completion criterion>`

The declaration is a contract. Mid-pass mode-switching is
forbidden: abort with a partial closeout and recommend
re-invocation under the correct mode.

## Topic-Page Circuit Breaker

A pass in `topic-page` mode may not report `topic-page:
produced` or `topic-page: extended` unless:

1. `.context/kb/topics/<slug>/index.md` exists and was
   created or extended in this pass (topic-page file
   creation is performed only by `ctx kb topic new`; this
   skill MAY invoke it but MUST NOT write the scaffold
   directly).
2. The page cites at least one `EV-###` row that resolves
   to `evidence-index.md`.
3. `ctx kb site build` ran clean, or its failure is named
   in the closeout's `Next pass hint` and the pass reports
   `topic-page: deferred`.
4. The cold-reader orientation rubric records `Result: pass`.

Any failure → `topic-page: deferred`; ledger advances to
`topic-page-drafted` (not `comprehensive`).

## Process

1. Verify pre-write gates. Refuse cleanly on failure.
2. Emit the pass-mode declaration.
3. Resolve sources and the topic name; confirm with the user
   when topic is not supplied.
4. Run the topic-adjacency pre-flight against
   `.context/kb/source-coverage.md`; record the result in
   the closeout's Adjacency pre-flight block.
5. Life-stage check (< 5 topic pages = bootstrap;
   reconciliation ceremony skipped except for contradictions;
   >= 5 = maintenance, full discipline).
6. Scaffold (topic-page mode only): if the topic folder does
   not exist, shell out to `ctx kb topic new "<name>"`.
7. Extract atomic claims; mint `EV-###` rows in
   `evidence-index.md`. Confidence band per `KB-RULES.md`;
   topic page never claims more certainty than its weakest
   cited band.
8. Reconcile (maintenance only): net-new claims append;
   reinforcing claims promote; contradicting claims demote
   per the demotion policy and open a paired
   `outstanding-questions.md` row.
9. Advance the source-coverage ledger per the state-machine
   transitions in `KB-RULES.md`. Illegal transitions are
   refused at write time.
10. Run the topic-page circuit breaker (topic-page mode only).
11. Record the cold-reader orientation rubric in the closeout.
12. Write the closeout under
    `.context/ingest/closeouts/<TS>-ingest-closeout.md` with
    required frontmatter (`sha`, `branch`, `mode: ingest`,
    `pass-mode`, `life-stage`, `generated-at`).
13. Append one line to `.context/ingest/SESSION_LOG.md` at
    the closeout phase boundary, in the exact shape from
    `KB-RULES.md` §SESSION_LOG line shape.

## Anti-Patterns

- Synthesizing the topic-page scaffold by hand. Only `ctx kb
  topic new` writes scaffolds.
- Mid-pass mode-switching without aborting.
- Claiming `comprehensive` ledger advance when the topic page
  is incomplete.
- Inventing `EV-###` IDs to make a topic page look complete.
- Demoting in `evidence-index.md` (rows are append-only;
  retire by demoting the confidence band, not by deletion).

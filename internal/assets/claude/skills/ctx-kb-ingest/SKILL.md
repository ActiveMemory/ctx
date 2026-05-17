---
name: ctx-kb-ingest
description: Editorial knowledge-ingestion pass. Reads sources the user supplies, declares its pass-mode (topic-page / triage / evidence-only) before extraction, and is held to mode-specific completion semantics. The topic page is the deliverable; the closeout is the audit trail.
allowed-tools: Bash(ctx:*), Bash(ls:*), Bash(git status), Read, Edit, Write
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

## Process

1. **Verify pre-write gates.** Refuse cleanly with the matching
   message from §Pre-write gates if `.context/`,
   `.context/ingest/`, or kb scope is missing. No residue on
   refusal.

2. **Declare pass-mode and surface the up-front declaration.**
   Determine the mode per §Pass-mode contract. Emit the
   three-line declaration block in the response stream **before
   any further work**. Mid-pass mode-switching is forbidden;
   abort and re-invoke if the work no longer fits.

3. **Resolve the topic.** *(Topic-page mode only; skipped in
   `triage` and `evidence-only`.)*

   - **Read `.context/kb/source-coverage.md` in full first.** It
     answers *"what does this kb already know about which
     sources, and at what completeness?"*: a precondition for
     honest topic resolution, not an afterthought.
   - **Topic-adjacency pre-flight (mandatory).** Scan the ledger
     for rows whose state is **not** in
     `{comprehensive, skipped, superseded}` AND whose `Topic` is
     plausibly *adjacent*. Heuristics:
     - **Shared first segment of a slash- or hyphen-separated
       slug**: `cursor/skills` is adjacent to `cursor/hooks`.
     - **Shared product / vendor / surface** in the source URL or
       description.
     - **Explicit cross-references** in the named topic's
       existing sub-pages or this pass's source set.

     For each adjacent incomplete topic surfaced, this pass MUST:
     1. Acknowledge it in `## Related concepts in this kb` on the
        topic page being authored.
     2. Surface it in the closeout's `Adjacency pre-flight`
        block.
     3. Surface it in the response contract's `Adjacent topics
        noted` field.

     **Do NOT enumerate `EV-###` IDs by name in the adjacency
     block.** Use *count + location* (*"seventeen rows in
     `evidence-index.md`"*). Naming an EV row from a
     lower-confidence sibling demotes the floor of cited bands.

     Silence is not a clean pre-flight; if zero matches, record
     *"no incomplete adjacent topics surfaced"* explicitly.

   - **Named vs unnamed branches.** If the user named a topic,
     accept it and map to slug (lowercase + kebab-case). If not,
     scan the inputs *just enough* to propose one and confirm:

     > you haven't named a topic; based on the inputs this looks
     > like **"<proposed name>"**. Confirm or correct.

     One question. Wait for confirmation. If material spans
     multiple topics, ask once for the splits. Do not auto-split.

4. **Resolve sources (and discover, if invited).** *(All modes.)*

   - Resolve every supplied source: fetch URLs, recurse folders,
     enumerate MCP resources.
   - If the user invited discovery, do bounded web/MCP search.
   - **Hard cap: 50 total sources** (supplied + discovered) per
     pass. Quality of synthesis collapses past it.
   - **If discovery exceeds 50**, keep the 50 highest-judged
     sources for this pass; append the overflow to
     `.context/ingest/candidate-sources.md` under a "Pending
     (overflow from <date> ingest of `<topic-slug>`)" heading.
   - **Update the source-coverage ledger**: every supplied source
     moves from absent → `discovered` (if newly seen) →
     `admitted` (if scope-conformant) or → `skipped` (if not).
     Discovered sources kept for this pass also land at
     `admitted`; overflow stays at `discovered` with a pointer.

   Append a `SESSION_LOG.md` line:

   ```
   [YYYY-MM-DD HH:MM:SS sha=<short> branch=<name>] phase=resolve status=<done|partial|blocked> note=<<=120 chars>
   ```

5. **Survey kb topology and determine life-stage.** *(All
   modes.)*

   - List `.context/kb/topics/*/index.md`; glance for sibling
     sub-pages so the cross-link palette includes them.
   - Read `.context/kb/index.md` for the canonical scope.
   - Skim recent sections of `evidence-index.md`, `glossary.md`,
     `outstanding-questions.md`, `contradictions.md`,
     `timeline.md` for prior claims relevant to this pass.
   - **Life-stage check**: count `kb/topics/*/index.md`. `< 5`
     is bootstrap; `>= 5` is maintenance. Document the call in
     the closeout's frontmatter (`life-stage:`).

6. **Find or create the topic page.** *(Topic-page mode only.)*

   Topic pages are folder-shaped from day one:
   `.context/kb/topics/<slug>/index.md`, with optional sibling
   sub-pages.

   - **If `.context/kb/topics/<slug>/index.md` exists**, read it
     AND enumerate any sibling sub-pages. The pass **extends**
     the topic: append/extend prose; reuse existing `EV-###`
     rows where possible; preserve human edits; do not reformat
     to match a newer template. Choose the right file:
     - Lede / "What it is" overview → edit `index.md`.
     - Existing sibling sub-page material → edit that sub-page.
     - **Sub-page split is lazy.** Do NOT pre-emptively split.
       Only split when `index.md` has grown to fail the
       cold-reader "boundaries clear?" check; at which point,
       propose the split (one question, wait for confirmation;
       sub-page topology affects long-term shape).
   - **If `.context/kb/topics/<slug>/` does not exist**, scaffold
     by invoking `ctx kb topic new "<concept name>"`. The CLI is
     the sole writer of the scaffold; do not synthesize it by
     hand. The CLI creates the folder, writes `index.md`, AND
     registers the new slug in `.context/kb/index.md`'s
     `CTX:KB:TOPICS` managed block.

   After revising the page's H1 or Confidence band in §10, run
   `ctx kb reindex` so the managed block refreshes.

7. **Synthesise.** Body depends on declared mode.

   ### `topic-page` mode

   For each template section (Status block, lede, "What it is",
   "Why this kb cares", "Sources and further reading", optional
   sections):

   - **Read the source(s) carefully**: full pass, not skim.
   - **Write paraphrased prose that captures the understanding**,
     not a transcription.
   - **For each claim needing citation, mint or reuse `EV-###`:**
     - Re-read `evidence-index.md` immediately before writing to
       find the highest existing `EV-NNN`; append the next
       integer. Pad to three digits (`EV-012`, not `EV-12`).
       Duplicate IDs are a hard refusal: abort and re-read.
     - **If the claim is already pinned** by an existing row,
       reuse the ID verbatim. If the existing claim no longer
       matches, treat as a contradiction (§8).
     - **If the existing row carries the `evidence-only` tag**,
       treat as review-required: re-read the source, confirm the
       claim, then promote onto the page. Leave the tag in
       place; it is audit trail.
     - Append the row to `evidence-index.md` per its schema
       (claim, source short name + locator, optional `sha:` for
       in-repo citations, confidence band, tags, extracted
       date).
     - If the source is new, append a row to `source-map.md`.
     - Cite `EV-###` inline in the prose.
   - **Cross-link** to existing kb topics, DECISIONS.md,
     LEARNINGS.md, and `docs/` entries when applicable.
   - **Mandatory `## Related concepts in this kb` entries** for
     adjacent incomplete topics surfaced by §3's pre-flight. The
     acknowledgement must read as a forward pointer (state +
     count + location), not as trivia.
   - **Mark unbacked claims with `TBD-cite`** and open
     `outstanding-questions.md` entries for each.
   - **Update `glossary.md`** for net-new terms.
   - **Update `timeline.md`** if the pass surfaces a dateable
     event.

   **Never invent citations.** **Never** promote a claim above
   `speculative` without an `evidence-index.md` row backing it.

   ### `triage` mode

   For each admitted source, judge admission/skip against the
   scope paragraph and propose topic routing in the closeout. Do
   NOT write to any topic page. **Do NOT mint `EV-###` rows.** Do
   NOT touch `evidence-index.md`, `glossary.md`, or
   `timeline.md`.

   Triage is routing and admission, not extraction. If the user
   asks to *"triage and grab obvious facts as you go,"* abort
   with a partial closeout and recommend re-invocation under
   either `topic-page` or `evidence-only` mode. Triage MAY update
   `source-coverage.md` and `candidate-sources.md`. That is the
   full write surface for triage.

   ### `evidence-only` mode

   For each admitted source, mint `EV-###` rows + `source-map.md`
   rows + `glossary.md` entries for terms encountered. **Do not
   touch any topic page.** Do not write prose synthesis.

   Every minted `EV-###` row MUST include the literal tag
   `evidence-only` in its tags column. The tag is **additive**;
   it does not replace topical tags.

   Append a `SESSION_LOG.md` line:

   ```
   [YYYY-MM-DD HH:MM:SS sha=<short> branch=<name>] phase=synthesise status=<done|partial|blocked> note=<topic slug + <=80 chars>
   ```

8. **Apply life-stage reconciliation discipline.** *(All modes;
   behavior depends on life-stage.)*

   **Bootstrap (`< 5` topic pages)**: skip except for the
   contradiction exception in §5. Append a `SESSION_LOG.md` line
   with `status=skipped-bootstrap`.

   **Maintenance (`>= 5` topic pages)**: for each EV row minted
   in §7:

   - **Reinforces an existing claim** → promote per the
     laddering rules in `KB-RULES.md` §Confidence bands
     (`speculative → low → medium → high`); cross-link the new
     row to the prior one.
   - **Contradicts an existing claim** → add a row to
     `contradictions.md`; demote the older claim per the
     demotion policy in `KB-RULES.md` §Demotion policy; open an
     `outstanding-questions.md` entry naming both sides and what
     evidence would resolve.

9. **Set the topic page's Confidence floor.** *(Topic-page mode
   only.)* Inspect every `EV-###` cited on the page; the page's
   Status-block `Confidence` is the **lowest** of those cited
   bands. Refuse to set Confidence above the floor. Refuse to
   set above `speculative` while any `TBD-cite` remains.

10. **Update the topic page's Status block.** *(Topic-page mode
    only.)* Substitute `Subject:`, `Last verified:`, `Author:`
    (`agent-ingested` if untouched by a human in this pass;
    `mixed` if a human revised prose; **never**
    `hand-authored`), and `Confidence:` per §9.

11. **Update the source-coverage ledger.** *(All modes.)* For
    every source touched, advance its row in
    `.context/kb/source-coverage.md` per the state machine.
    Update `EV coverage`, `Residue`, `Next action`, `Updated`
    columns honestly. Lying to the ledger is a hard
    anti-pattern.

12. **Topic-page circuit breaker check.** *(Topic-page mode
    only.)* Verify all four invariants from §Topic-page circuit
    breaker. Any failure → `topic-page: deferred` and ledger to
    `topic-page-drafted` (NOT `comprehensive`).

13. **Write the closeout.** *(All modes; mode-aware body.)*
    Create
    `.context/ingest/closeouts/<TIMESTAMP>-ingest-closeout.md`
    with required frontmatter:

    ```yaml
    ---
    sha: <short>
    branch: <name>
    mode: ingest
    pass-mode: <topic-page|triage|evidence-only>
    life-stage: <bootstrap|maintenance>
    generated-at: <RFC-3339>
    ---
    ```

    Body sections (mode-aware): **Inputs**, **Pass-mode** (block
    repeated from §2 declaration so reviewers can compare promise
    vs. result), **Topic(s) touched**, **What changed**
    (including the Cold-reader rubric in topic-page mode),
    **New questions**, **New contradictions**, **Confidence
    drift**, **Source-coverage updates**, **Overflow**,
    **Adjacency pre-flight**, **Next pass hint**.

    Append a final `SESSION_LOG.md` line:

    ```
    [YYYY-MM-DD HH:MM:SS sha=<short> branch=<name>] phase=closeout status=done note=<topic slug + <=80 chars>
    ```

## Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Empty input | Refuse with the standard no-sources text. No residue. |
| `.context/` missing | Refuse; suggest `ctx init`. No residue. |
| `.context/ingest/` missing | Refuse; suggest `ctx init --upgrade`. No residue. |
| Kb scope undeclared | Refuse with the scope message; point at `.context/kb/index.md`. No residue. |
| Source returns nothing usable (404, binary, paywall) | Record in closeout's `Next pass hint` AND the topic page's "Open questions"; advance the ledger row to `skipped` with the failure reason. Do not invent claims. |
| All sources skipped during admission | Write a short closeout with empty `Topic(s) touched` and `What changed`; `Next pass hint` lists every skipped source with scope-citation. |
| Material spans 3+ topics and user can't decide | Ask once in §3; if still unresolved, abort with a partial closeout recommending re-invocation under `triage`. |
| Discovery turns up zero additional sources | Note in closeout's `Inputs` section. Not a failure. |
| Stale Status block on existing page | Flag in closeout's `Next pass hint`; do not silently overwrite the verification cursor unless the source was actually re-verified. |
| Multiple sessions filling the same page | Read existing prose first; do not overwrite human edits; append/extend rather than replace. |
| Page scaffolded long ago with older template | Fill what's there; do not reformat to match a newer template. Open a task if drift is significant. |
| `ctx kb topic new` fails or refuses (slug exists, kb missing) | Resolve the underlying condition and retry; do not hand-write a scaffold. |
| `ctx kb site build` fails during §12 | Report `topic-page: deferred`; name the build failure in `Next pass hint`; ledger to `topic-page-drafted` (NOT `comprehensive`). |
| Cold-reader rubric returns `Result: fail` | Report `topic-page: deferred` AND `validation: deferred (cold-reader orientation failed)`; name failed items in `Next pass hint`; ledger to `topic-page-drafted`. |
| Adjacency pre-flight surfaces zero matches | Record *"no incomplete adjacent topics surfaced"* explicitly in closeout's `Adjacency pre-flight`; response contract reads `none surfaced`. Silence is not allowed. |
| Mid-pass mode-switching tempted | Forbidden. Abort, write a partial closeout citing the mismatch, recommend re-invocation under the correct mode. Never silent-switch. |
| `evidence-only` pass discovers a contradiction | Still mint the contradiction row (truth surface always wins); flag in `Next pass hint` that a topic-page pass is needed to resolve. |
| Inferring `evidence-only` from source size / time pressure | Hard anti-pattern. Refuse to set `evidence-only` without explicit user trigger. |

## Hard Anti-Patterns

- Treating closeout existence as topic-page validation.
- Skipping the topic-page circuit breaker in `topic-page` mode.
- Inferring `evidence-only` from source size, complexity,
  ambiguity, time pressure, or operator convenience.
- Mid-pass mode-switching (abort and re-invoke instead).
- Hiding incomplete coverage under a comprehensive-looking
  closeout (lying to the ledger).
- Skipping the topic-adjacency pre-flight, or running it but
  failing to acknowledge surfaced incomplete adjacent topics.
- Claiming `topic-page: produced` when the cold-reader
  orientation result is missing.
- Asking the human mid-pass beyond the §3 naming gate, unless
  continuing would change durable kb topology, evidence
  confidence, source admission, or scope.
- Inventing claims beyond what the source backs.
- Inventing `EV-###` citations to make a page look complete.
- Promoting claims above `speculative` without an
  `evidence-index.md` row.
- Promoting a topic page above its weakest cited band.
- Setting `Confidence` above `speculative` while any
  `TBD-cite` remains.
- Setting `Author: hand-authored` on agent-ingested prose.
- Re-extracting from a source that already has `EV-###` rows
  instead of reusing the IDs.
- Citing an `evidence-only`-tagged row in a topic page without
  re-reading the source first.
- Renumbering or deleting `EV-###` rows when reconciling.
- Skipping the closeout once the pass clears pre-write gates.
- Bypassing `ctx kb topic new` when scaffolding a page.
- Running maintenance discipline against a bootstrap-stage kb.
- Hand-editing `INBOX.md`.

## Output Contract

For pre-write refusals, return only the specified refusal text
and stop. No closeout, no residue.

For passes that clear pre-write gates, **emit the up-front
declaration first** (between §1 and §3):

> **Pass-mode:** `<mode>`
> **Reason:** `<one sentence; required when non-default>`
> **Definition of done:** `<mode-specific criterion>`

Then proceed. At completion, end with this structured summary:

- **Pass-mode**: as declared, with reason if non-default.
- **Topic-page**: `produced [<slug>]`, `extended [<slug>]`,
  `deferred (<reason>)`, or `not-applicable` (triage /
  evidence-only).
- **Validation**: `passed`, `not-attempted`, or
  `deferred (<reason>)`.
- **Coverage**: current state(s) from `source-coverage.md` for
  sources touched.
- **EV range minted**: e.g. `EV-035..EV-051`, or `none`.
- **Counts**: glossary entries added, source-map rows added,
  cross-links written, contradictions surfaced, questions
  opened.
- **Life-stage**: `bootstrap` or `maintenance`, with the
  topic-page count it was based on.
- **Closeout**: filename on its own line.
- **Adjacent topics noted** *(topic-page mode only; mandatory)*:
  either `none surfaced` or a slug-list with states. Free prose
  fails validation; the doctor advisory parses this field.
- **Next-recommended-action**: explicit invocation that would
  resume incomplete work. Adjacent topics surfaced by the
  pre-flight MUST appear here too (deliberate redundancy).
- **Review-required**: `true` for `evidence-only` passes;
  otherwise omit.

The structured summary, the closeout's body, and the
source-coverage ledger MUST agree. Discrepancies between the
three are a hard anti-pattern; the doctor advisory detects them
and surfaces a non-fatal warning on next `ctx doctor` run.

## Quality Checklist

Before reporting completion, verify:

- [ ] Pre-write gates passed (or the matching refusal was
      returned with zero residue).
- [ ] Pass-mode declaration was emitted in the response stream
      before any extraction.
- [ ] Source-coverage ledger advanced honestly for every source
      touched.
- [ ] Topic-adjacency pre-flight ran in `topic-page` mode and its
      result is in the closeout AND the page AND the response
      contract.
- [ ] Cold-reader rubric is recorded in `topic-page` mode.
- [ ] Circuit-breaker check ran in `topic-page` mode; failure
      → `topic-page: deferred`, NOT `produced`.
- [ ] Closeout written with all required frontmatter fields
      (`sha`, `branch`, `mode`, `pass-mode`, `life-stage`,
      `generated-at`).
- [ ] Structured response summary matches the closeout body and
      the ledger.

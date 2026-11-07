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

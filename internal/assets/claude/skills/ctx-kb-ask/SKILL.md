---
name: ctx-kb-ask
description: Q&A grounded in the existing kb. Read-only on prose; refuses to web-jump; if the kb cannot answer, opens a Q-### row in outstanding-questions.md and reports the gap. Writes an ask closeout for the audit trail.
allowed-tools: Bash(ctx:*), Bash(ls:*), Read, Edit
---

# Ask the KB

Answer a question using only what `.context/kb/` already contains.
Cite by `EV-###`. Do not web-jump, do not invent prose, do not
modify topic pages. If the kb cannot answer, open a `Q-###` row
in `.context/kb/outstanding-questions.md` and report the gap.

This is the read side of the editorial pipeline. The write side
is `/ctx-kb-ingest`. Authority for prose synthesis lives there;
this skill is read-only on prose.

Authoritative background reading:
`.context/ingest/KB-RULES.md` §Authority boundary and
§Evidence discipline; `specs/kb-editorial-pipeline.md` §Interface.

## When to Use

- The user asks "does the kb say...", "according to evidence...",
  "what do we know about <topic>", or invokes the explicit slash
  form with the question.
- The user wants a citation-backed answer before deciding whether
  to ingest more material.
- The user is auditing what is already known versus what is
  asserted elsewhere (DECISIONS.md, LEARNINGS.md, conversation).

## When NOT to Use

- The user wants new material extracted (use `/ctx-kb-ingest`).
- The user wants the kb structurally audited (use
  `/ctx-kb-site-review`).
- The user wants kb claims re-grounded against external sources
  (use `/ctx-kb-ground`).
- The question is about `ctx` itself or the editorial pipeline
  contract (answer from `KB-RULES.md` / spec directly).

## Authority Boundary (vs Other Skills)

- **`/ctx-kb-ask`**: read-only Q&A over `.context/kb/` prose,
  `evidence-index.md`, `glossary.md`, `contradictions.md`,
  `outstanding-questions.md`, `timeline.md`, `source-map.md`,
  `domain-decisions.md`. Writes are limited to opening a
  `Q-###` row in `outstanding-questions.md` when the kb cannot
  answer, plus the ask closeout.
- **`/ctx-kb-ingest`**: writes prose, evidence, scaffold. Only
  ingest may add citations or extend topic pages.
- **`/ctx-kb-ground`**: refreshes external sources via
  `grounding-sources.md`; this skill never web-jumps to fill a
  gap. If the gap matters, recommend `/ctx-kb-ground` or
  `/ctx-kb-ingest`.

## Usage Examples

```text
/ctx-kb-ask "what does the kb say about cursor hooks failure modes?"
/ctx-kb-ask "how do we cite a transcript locator?"
/ctx-kb-ask "are there contradictions on backup retention windows?"
```

## Input Contract

A single question, supplied as the slash argument or inline. No
flags. No sources. No URLs.

## Refuse-on-Empty

If the invocation supplied no question (empty slash arg, empty
inline body), return exactly:

> no question provided; pass a question or describe it inline.

Stop. Do not prompt interactively. The CLI enforces this
independently via `cmd/ask`.

## Pre-Write Gates

Three distinct refusals, each leaves zero residue (no
`Q-###` row opened, no closeout):

- `.context/` missing → suggest `ctx init` and stop.
- `.context/kb/` missing → suggest `ctx init --upgrade` and
  stop.
- `.context/kb/index.md` exists but `## Scope` is undeclared →
  refuse with the scope message (same wording as `/ctx-kb-ingest`
  uses) and stop.

## Process

1. **Verify pre-write gates.** Refuse cleanly if any gate fails.
   Zero residue on refusal.

2. **Read the question.** Parse for the concept(s) it names.

3. **Survey the kb.** Read in this order, stopping early when an
   answer surfaces with adequate citation coverage:
   - `.context/kb/index.md` for scope.
   - `.context/kb/topics/<slug>/index.md` and any sibling
     sub-pages for any slug that plausibly matches the question.
   - `.context/kb/evidence-index.md` for `EV-###` rows whose
     claim text matches.
   - `.context/kb/glossary.md` for term definitions.
   - `.context/kb/contradictions.md` for known disagreements
     relevant to the question.
   - `.context/kb/outstanding-questions.md` for prior
     unanswered questions on the topic.

4. **Decide answer vs gap.** One of three outcomes:

   - **Answerable with citations.** The kb's prose plus
     `EV-###` rows cover the question. Compose a concise answer.
     Cite every load-bearing claim by `EV-###`. Name the topic
     page(s) where the prose lives. Note the Confidence floor
     of the cited rows.
   - **Partial answer.** Some of the question is covered; the
     rest is not. Answer the covered part with citations. Name
     the gap explicitly. Open a `Q-###` row for the gap (see §6).
   - **Not answerable.** The kb has no prose and no `EV-###`
     coverage. Do not invent. Do not web-jump. Open a `Q-###`
     row (see §6) and report the gap.

5. **Do not jump.** This skill is read-only on prose AND
   web-quiet. If the kb cannot answer:

   - Do **not** fetch a URL.
   - Do **not** propose synthesized prose without citations.
   - Do **not** call MCP search tools.
   - Do **not** quote LLM training-data recall as if it were
     kb evidence.

   The correct response to a gap is to name the gap, open a
   `Q-###` row, and recommend `/ctx-kb-ground` (if external
   refresh is the right path) or `/ctx-kb-ingest <sources>` (if
   the user has materials to feed in).

6. **Open a `Q-###` row if there is a gap.** Append a row to
   `.context/kb/outstanding-questions.md` per its schema. The
   row's question text is the user's question (or a faithful
   paraphrase). The row notes what the kb does cover (if
   partial) and what evidence would resolve. Do NOT mint
   `EV-###` rows from this skill; that is ingest's authority.

7. **Write the ask closeout.** Create
   `.context/ingest/closeouts/<TIMESTAMP>-ask-closeout.md` with
   required frontmatter:

   ```yaml
   ---
   sha: <short>
   branch: <name>
   mode: ask
   pass-mode: read-only
   life-stage: <bootstrap|maintenance>
   generated-at: <RFC-3339>
   ---
   ```

   Body sections:
   - **Question**: what the user asked, verbatim.
   - **Answer**: the answer given, or `none (gap)` if not
     answerable.
   - **Citations**: `EV-###` IDs cited, with topic-page paths.
   - **Gaps**: `Q-###` opened in `outstanding-questions.md`,
     with a one-line rationale.
   - **Next pass hint**: explicit invocation for the next
     pipeline step (e.g. `/ctx-kb-ground` to refresh,
     `/ctx-kb-ingest <sources>` to extend).

## Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Empty question | Refuse with the standard no-question text. No `Q-###` opened, no closeout. |
| `.context/` missing | Refuse; suggest `ctx init`. No residue. |
| `.context/kb/` missing | Refuse; suggest `ctx init --upgrade`. No residue. |
| Kb scope undeclared | Refuse with the scope message; point at `.context/kb/index.md`. No residue. |
| Multiple topics relevant | Cite each topic page; do not synthesize a new cross-topic claim (that would be ingest work). Surface the seam as a `Q-###` if it merits one. |
| Contradiction surfaces during answer | Answer with the lower-confidence side noted; cite both `EV-###` rows; point at `contradictions.md`. |
| Cited rows are all `speculative` or `low` | Surface the confidence band in the answer. Recommend `/ctx-kb-ground` to corroborate. Do not promote in this pass. |
| Question matches an existing `Q-###` row | Cite the existing row's ID; report status (`open`, `partially-answered`); do not open a duplicate. |
| Question requires external evidence the kb does not have | Open a `Q-###` row; recommend `/ctx-kb-ground` with the gap named; do not fetch the source. |
| Question is meta (about the pipeline itself) | Answer from `KB-RULES.md` / spec directly; this skill is for kb content, not pipeline contract. State that explicitly. |

## Anti-Patterns

- Web-jumping when the kb cannot answer. The contract is
  read-only on prose AND web-quiet.
- Inventing citations or claims to make the answer look fuller.
- Modifying a topic page to extend an answer mid-pass. Topic-page
  authoring is `/ctx-kb-ingest`'s authority.
- Minting `EV-###` rows from this skill. Evidence authoring is
  `/ctx-kb-ingest`'s authority.
- Skipping the `Q-###` row when the kb cannot answer. The gap
  is the audit trail; silence on a gap is invisible.
- Skipping the closeout once the pre-write gates pass. The
  closeout is the residue wrap-up's handover step folds into
  the next session's recall.

## Output Contract

For pre-write refusals, return only the specified refusal text
and stop. No residue.

For passes that clear pre-write gates, end with this structured
summary:

- **Question**: verbatim or faithful paraphrase.
- **Answer**: concise; cites every load-bearing claim by
  `EV-###`.
- **Confidence floor**: lowest band among cited rows
  (`high|medium|low|speculative`), or `n/a` if no rows cited.
- **Gaps**: `Q-### opened` (one bullet per opened row), or
  `none`.
- **Closeout**: filename on its own line.
- **Next-recommended-action**: explicit invocation if a gap
  was opened (e.g. `/ctx-kb-ground` or `/ctx-kb-ingest
  <sources>`), or `none` if the answer is complete.

## Quality Checklist

Before reporting completion, verify:

- [ ] Pre-write gates passed (or the matching refusal was
      returned with zero residue).
- [ ] Every load-bearing claim in the answer cites at least one
      `EV-###` row from `evidence-index.md`.
- [ ] If a gap exists, a `Q-###` row was opened (or an existing
      row was cited).
- [ ] No URL was fetched, no MCP search was called, no LLM
      training-data recall was quoted as kb evidence.
- [ ] No topic page was modified, no `EV-###` row was minted.
- [ ] Closeout written with all required frontmatter fields.

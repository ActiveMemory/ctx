# Ask Mode

Q&A pass grounded in the kb. The operator has a question;
the skill answers from `.context/kb/` content with explicit
confidence bands and citation chains. **Read-only on kb
prose.** The only kb-side write allowed is a new row in
`outstanding-questions.md` when the answer surfaces a gap.

Read `KB-RULES.md` first; it carries the authority boundary,
evidence discipline, confidence bands, and the closeout
shape. This file describes the per-mode procedure only.

---

## Inputs

The skill takes a question as its single argument, or inline
natural-language context with the question embedded. **No
question ⇒ refuse cleanly:**

> no question provided; pass a question or describe it inline.

Refuse-on-empty is the contract; the skill does not prompt
mid-flight.

---

## Pre-Write Gates

1. `.context/` and `.context/kb/` exist.
2. `.context/kb/index.md` has a non-placeholder `## Scope`
   section.

If either gate fails, refuse with the recovery hint from
`KB-RULES.md`.

---

## Procedure

1. **Restate** the question in one sentence. If the operator
   gave a multi-part question, split it into sub-questions and
   answer each in turn. Record the decomposition in the
   closeout's `Inputs` section.
2. **Search** `.context/kb/` and `evidence-index.md` for
   relevant claims. Use the source short names from
   `source-map.md` to track citation chains. Walk topic
   `index.md` files first; descend into sub-pages only when
   the lead claim is in one.
3. **Compose** the answer in this shape:

   - **Direct answer** in 1-3 sentences.
   - **Confidence** band (`high` / `medium` / `low` /
     `speculative`), determined by the floor of the cited
     `EV-###` bands per `KB-RULES.md` "Confidence bands".
   - **Sources**: bullet list of evidence rows that ground
     the answer (cite by `EV-###` plus the source short
     name).
   - **Caveats**: any contradictions or open questions that
     bear on the answer.

4. **Detect gaps.** If the answer is `speculative` or relies
   on `low`-confidence claims, OR if no evidence row backs
   the answer at all:

   - Append a `Q-###` row to `outstanding-questions.md` with
     the question text and a one-line summary of why current
     evidence is insufficient.
   - Suggest a concrete next pass in the closeout's
     `Next pass hint` (usually `ctx kb ground` for external
     evidence or `ctx kb ingest <source>` for internal).

5. **Write the closeout** under
   `.context/ingest/closeouts/<TS>-ask-closeout.md` with
   `mode: ask` in the frontmatter and the body shape from
   `KB-RULES.md`.

6. **Append a SESSION_LOG line** at the closeout phase
   boundary.

---

## Constraints

- **Read-only on kb prose.** This mode does not append
  paragraphs to `glossary.md`, `domain-decisions.md`,
  `timeline.md`, or any topic-page prose. The single
  kb-side write allowed is a new row in
  `outstanding-questions.md`.
- **No web-jumping.** If the answer is not in the kb,
  surface that fact and recommend `ctx kb ground` as the
  next pass. Do not fetch external sources to fabricate an
  answer.
- **Never invent citations.** If a claim is in working
  memory but not in `evidence-index.md`, mark the answer
  `speculative` and open a `Q-###` row rather than inflating
  the confidence band with a fictional row.
- Multi-part questions get a single closeout, not one per
  sub-question. The `Inputs` section captures the
  decomposition.
- The `pass-mode:` frontmatter field is required by the
  closeout shape, but ask mode does not carry a pass-mode
  contract; set `pass-mode: n/a`. Doctor advisory tolerates
  `n/a` for non-ingest modes.

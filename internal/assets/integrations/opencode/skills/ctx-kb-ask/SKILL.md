---
name: ctx-kb-ask
description: "Q&A grounded in the existing kb. Read-only on prose; refuses to web-jump; if the kb cannot answer, opens a Q-### row in outstanding-questions.md and reports the gap. Writes an ask closeout for the audit trail."
---

Answer a question using only what `.context/kb/` already contains.
Cite by `EV-###`. Do not web-jump, do not invent prose, do not
modify topic pages. If the kb cannot answer, open a `Q-###` row
in `.context/kb/outstanding-questions.md` and report the gap.

This is the read side of the editorial pipeline. The write side
is `/ctx-kb-ingest`. Authority for prose synthesis lives there;
this skill is read-only on prose.

## When to Use

- The user asks "does the kb say...", "according to evidence...",
  "what do we know about <topic>".
- The user wants a citation-backed answer before deciding
  whether to ingest more material.

## When NOT to Use

- The user wants new material extracted (`/ctx-kb-ingest`).
- The user wants the kb structurally audited
  (`/ctx-kb-site-review`).
- The user wants external re-grounding (`/ctx-kb-ground`).
- The question is about `ctx` itself (answer from
  `KB-RULES.md` directly).

## Input

A single question, supplied as the slash argument or inline.
No flags, no sources, no URLs.

## Refuse-on-Empty

If no question was supplied, return exactly:

> no question provided; pass a question or describe it inline.

Stop. The CLI enforces this independently.

## Pre-Write Gates

- `.context/` missing → suggest `ctx init` and stop.
- `.context/kb/` missing → suggest `ctx init --upgrade` and stop.
- Kb scope undeclared → refuse with the scope message and stop.

## Process

1. Verify pre-write gates.
2. Survey the kb in this order, stopping early when an answer
   surfaces with adequate citation coverage: `index.md` for
   scope, topic-page indexes and sub-pages for matching slugs,
   `evidence-index.md`, `glossary.md`, `contradictions.md`,
   `outstanding-questions.md`.
3. Decide answer vs gap:
   - **Answerable with citations**: cite every load-bearing
     claim by `EV-###`. Name the topic page(s). Note the
     confidence floor of cited rows.
   - **Partial answer**: answer the covered part; open a
     `Q-###` row for the gap.
   - **Not answerable**: open a `Q-###` row; report the gap.
     Do not invent. Do not web-jump.
4. If a gap exists, append a `Q-###` row to
   `outstanding-questions.md`. Do NOT mint `EV-###` rows;
   evidence authoring is `/ctx-kb-ingest`'s authority.
5. Write the ask closeout under
   `.context/ingest/closeouts/<TS>-ask-closeout.md` with
   required frontmatter (`sha`, `branch`, `mode: ask`,
   `pass-mode: read-only`, `life-stage`, `generated-at`) and
   body sections: Question, Answer (or `none (gap)`),
   Citations (`EV-###` + topic-page paths), Gaps (`Q-###`
   opened with one-line rationale), Next pass hint.

## Anti-Patterns

- Web-jumping when the kb cannot answer. The contract is
  read-only on prose AND web-quiet.
- Inventing citations or claims.
- Modifying a topic page to extend an answer mid-pass.
- Minting `EV-###` rows from this skill.
- Skipping the `Q-###` row when the kb cannot answer.
- Skipping the closeout once pre-write gates pass.

## Output Contract

For pre-write refusals, return only the refusal text and stop.

For passes that clear pre-write gates, end with:

- **Question**: verbatim or faithful paraphrase.
- **Answer**: concise; cites every load-bearing claim by
  `EV-###`.
- **Confidence floor**: lowest band among cited rows.
- **Gaps**: `Q-### opened`, or `none`.
- **Closeout**: filename.
- **Next-recommended-action**: explicit invocation if a gap
  was opened (`/ctx-kb-ground` or
  `/ctx-kb-ingest <sources>`), or `none`.

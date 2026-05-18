---
name: ctx-kb-note
description: "Lightweight capture into .context/ingest/findings.md. Single argument is the note text. Never writes to a topic page or to evidence-index.md. The pipeline's ad-hoc escape hatch for 'park this for the next ingest'."
---

Append a short note to `.context/ingest/findings.md` so a later
`/ctx-kb-ingest` pass can pick it up. The pipeline's escape
hatch for *"I want to remember this, but I'm not running a full
ingest right now."* No closeout, no ledger update, no
topic-page edit, no `EV-###` minting. Just typed memory landing
in one well-known file.

## When to Use

- The user says "drop a note", "capture this for the next
  ingest", "park this finding".
- A conversation surfaces a fact, link, or observation that
  should land in the kb later but does not justify a full
  ingest pass right now.

## When NOT to Use

- The user has sources to ingest (`/ctx-kb-ingest`).
- The user is asking a content question (`/ctx-kb-ask`).
- The note is actually a task / decision / learning /
  convention for the code-dev side; use the matching
  canonical-capture path instead.
- The note is empty (refuse-on-empty).

## Input

A single argument: the note text. Free-form prose. No flags.

## Refuse-on-Empty

If the invocation supplied no note text, return exactly:

> no note text provided; pass the note inline.

Stop. The CLI enforces this independently.

## Pre-Write Gates

- `.context/` missing → suggest `ctx init` and stop.
- `.context/ingest/` missing → suggest `ctx init --upgrade`
  and stop.

Kb scope declaration is not required; notes land pre-scope.

## Process

1. Verify pre-write gates.
2. Append the note to `.context/ingest/findings.md` as a
   single bulleted line, prefixed with the UTC timestamp
   (RFC-3339), short SHA, and branch:

   ```
   - 2026-05-16T14:32:11Z sha=88d52870 branch=main
     | <note text>
   ```

   If `findings.md` does not exist, create it with a brief
   header explaining its purpose.
3. No closeout. Notes are intentionally lightweight; the
   audit trail is the file itself.

## Anti-Patterns

- Writing a topic page or minting `EV-###` from this skill.
- Skipping the timestamp + provenance prefix; the audit
  trail depends on it.
- Hand-editing prior entries to "consolidate" them. Append-only.

---
name: ctx-kb-note
description: Lightweight capture into .context/ingest/findings.md. Single argument is the note text. Never writes to a topic page or to evidence-index.md. The pipeline's ad-hoc escape hatch for "park this for the next ingest".
---

# Park a Finding for the Next Ingest

Append a short note to `.context/ingest/findings.md` so a later
`/ctx-kb-ingest` pass can pick it up. This is the pipeline's
escape hatch for *"I want to remember this, but I'm not running
a full ingest right now."* No closeout, no ledger update, no
topic-page edit, no `EV-###` minting. Just typed memory landing
in one well-known file.

Authoritative background reading:
`.context/ingest/KB-RULES.md` §Authority boundary;
`specs/kb-editorial-pipeline.md` §Interface.

## When to Use

- The user says "drop a note", "capture this for the next
  ingest", "park this finding", or invokes the explicit slash
  form with note text.
- A conversation surfaces a fact, link, or observation that
  should land in the kb later but does not justify running
  `/ctx-kb-ingest` right now.
- Mid-session, a sibling skill (architecture, brainstorm, etc.)
  surfaces something kb-shaped and the user wants it parked
  cheaply.

## When NOT to Use

- The user has sources in hand and wants them ingested (use
  `/ctx-kb-ingest`).
- The user is asking a content question (use `/ctx-kb-ask`).
- The note is actually a task / decision / learning / convention
  for the code-dev side (use `/ctx-task-add` /
  `/ctx-decision-add` / `/ctx-learning-add` /
  `/ctx-convention-add`; those write to canonical files, this
  one does not).
- The note is empty (refuse-on-empty; see below).

## Authority Boundary (vs Other Skills)

- **`/ctx-kb-note`** appends to
  `.context/ingest/findings.md` only. Never writes anywhere
  else. No closeout. No ledger update.
- **`/ctx-kb-ingest`** reads `findings.md` opportunistically
  when scoping its source set; the user controls when notes get
  promoted into evidence.
- **Canonical capture skills** (`/ctx-task-add`,
  `/ctx-decision-add`, `/ctx-learning-add`,
  `/ctx-convention-add`) write to the five canonical
  `.context/` files. Strict authority boundary: this skill
  never touches them.

## Usage Examples

```text
/ctx-kb-note "cursor.com/changelog mentions hook lifecycle bump in v1.2"
/ctx-kb-note "check whether your-domain RTO claim still cites the 2024 audit"
/ctx-kb-note "Volkan said in chat: the 50-source cap was lifted from the upstream design"
```

## Input Contract

A single argument: the note text. Free-form prose. No flags.

## Refuse-on-Empty

If the invocation supplied no note text (empty slash arg, empty
inline body, whitespace-only), return exactly:

> no note text provided; pass the note inline.

Stop. Do not prompt interactively. The CLI enforces this
independently via `cmd/note`.

## Pre-Write Gates

Two distinct refusals, each leaves zero residue:

- `.context/` missing → suggest `ctx init` and stop.
- `.context/ingest/` missing → refuse:

  > kb not initialized; run `ctx init` first

  Stop.

Kb scope declaration is **not** required for this skill. Notes
land in `.context/ingest/findings.md`, which is pre-kb-scope
territory; the user may be parking notes precisely because they
have not yet decided the kb's scope.

## Process

1. **Verify pre-write gates.** Refuse cleanly if any gate fails.
   Zero residue on refusal.

2. **Append the note** to `.context/ingest/findings.md` as a
   single bulleted line. Prefix with the current UTC timestamp
   (RFC-3339, date-time precision) and a short SHA + branch
   from `gitmeta.ResolveHead` so the note carries minimal
   provenance:

   ```
   - 2026-05-16T14:32:11Z sha=88d52870 branch=main
     | <note text>
   ```

   If `findings.md` does not yet exist, create it with a brief
   header explaining its purpose (one paragraph; the embedded
   template ships at `internal/assets/kb/templates/ingest/`
   handles this for fresh inits, so this fallback applies only
   when the file was deleted by hand).

3. **No closeout.** Notes are intentionally lightweight; the
   audit trail is the file itself. The next `/ctx-kb-ingest`
   pass reads `findings.md` opportunistically.

## Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Empty note text | Refuse with the standard no-note text. No residue. |
| `.context/` missing | Refuse; suggest `ctx init`. No residue. |
| `.context/ingest/` missing | Refuse with the not-initialized message. No residue. |
| `findings.md` missing but `.context/ingest/` exists | Create the file with a brief header; append the note. |
| Multi-line note text | Append as a single bullet with embedded line breaks; preserve the user's formatting. |
| Note text contains a URL | Preserve verbatim; do not auto-fetch (this skill does not web-jump). |
| Note text is structurally a claim that should be evidence | Append as a note anyway; mention in the response that `/ctx-kb-ingest` is the next step if the user wants it minted as `EV-###`. |
| User invokes twice in a row with similar text | Append both; deduplication is the user's call, not this skill's. |

## Output Contract

For refusals, return only the specified refusal text and stop.

For successful appends, return:

- One line confirming the append, with the line number of the
  new entry in `findings.md`.
- A pointer to `/ctx-kb-ingest` as the path for promoting the
  note into evidence when the user is ready.

Example:

```
appended to .context/ingest/findings.md line 42.
run /ctx-kb-ingest with the source materials when ready to mint EV.
```

## Quality Checklist

Before reporting completion, verify:

- [ ] Pre-write gates passed (or the matching refusal was
      returned with zero residue).
- [ ] The note landed in `.context/ingest/findings.md` and
      nowhere else.
- [ ] No `EV-###` row was minted, no topic page was touched, no
      ledger row was advanced, no closeout was written.
- [ ] The appended line carries the timestamp + sha + branch
      provenance prefix.

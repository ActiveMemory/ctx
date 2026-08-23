---
name: ctx-handover
description: "Per-session handover artifact writer. Wraps `ctx handover write` with required `--summary` and `--next`. Always invoked as the final step of `/ctx-wrap-up`; not the user-facing trigger. When `.context/kb/` exists, also folds postdated closeouts into the handover and archives them."
---

Write the per-session handover under
`.context/handovers/<TS>-<slug>.md`. The handover is the
former agent's note to the next agent (or human): what
happened, and what should come next. `/ctx-remember` reads
it at the start of the next session.

## When to Use

`/ctx-wrap-up` owns the user-facing session-end trigger and
delegates to this skill as its final step. Direct invocation
is reserved for:

- `--no-fold` mid-session checkpoint when the user wants to
  pause without consuming closeouts.
- Recovery, when a prior session aborted before wrap-up.

Otherwise, the user invokes `/ctx-wrap-up`, not this skill.

## Input Contract

Wraps `ctx handover write`. Empty `TBD`, `see chat`,
whitespace-only values for required flags are rejected by
the CLI.

| Flag | Required | Description |
|------|----------|-------------|
| `--summary` | yes | Past tense; what happened this session. |
| `--next` | yes | Future tense; the specific first action for the next agent. |
| `--highlights` | no | Notable artifacts produced this session. |
| `--open-questions` | no | Things that remain undecided. |
| `--no-fold` | no | Skip closeout consumption (mid-session checkpoint). |
| `--commit` | no | Override resolved git HEAD for the Provenance line (CI replay). |

Positional argument: handover title (becomes filename slug).

## Pre-Write Gates

- `.context/` missing → suggest `ctx init` and stop.
- `.context/handovers/` missing → suggest `ctx init --upgrade`
  and stop.

`.context/kb/` is not required for handover; KB-state folding
is conditional on the directory's existence.

## Process

1. Verify pre-write gates. Refuse cleanly on failure.
2. Gather signal: `git status --short`, `git diff --stat`,
   `git log --oneline @{upstream}..HEAD || git log --oneline -5`,
   and scan the conversation for the session's arc, concrete
   artifacts, open questions, and the specific next action.
3. Draft `--summary` (past tense, concrete) and `--next`
   (future tense, specific). Surface to the user for
   confirmation before running the CLI.
4. Run:

   ```bash
   ctx handover write "<title>" \
     --summary "<...>" --next "<...>" \
     [--highlights "<...>"] [--open-questions "<...>"] \
     [--no-fold] [--commit <sha>]
   ```

   The CLI validates flags, resolves git HEAD, reads the
   latest-handover cursor, folds postdated closeouts into
   the new handover's `## Folded Closeouts` section, and
   archives them under `.context/archive/closeouts/`.

5. Report the handover filename, count of closeouts folded,
   count of malformed closeouts skipped (with filenames),
   and the `--next` value verbatim so the operator sees what
   the next agent will read first.

## Anti-Patterns

- Hand-writing a handover file. The CLI is the sole writer.
- Skipping the fold to "keep closeouts available." Use
  `--no-fold` explicitly when the user wants the checkpoint
  behavior; do not infer it.
- Inventing `--highlights` or `--open-questions` content the
  session did not actually produce.

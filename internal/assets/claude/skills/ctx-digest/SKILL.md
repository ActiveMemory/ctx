---
name: ctx-digest
description: "Dry-run the progressive-disclosure pass: inspect a knowledge file's staging zone, propose a theme per staged entry, author gists, and show the digest plan — WITHOUT moving anything. Use to preview how LEARNINGS.md / DECISIONS.md would fold into themes when a file grows large."
allowed-tools: Bash(ctx:*), Read
---

Preview how a bounded knowledge root would digest its staging zone into
themes — and **move nothing**. This is the dry-run half of progressive
disclosure (see specs/progressive-disclosure.md): it proposes and shows
a plan; a later milestone's apply step does the moving.

**This skill never edits a knowledge file, never creates a theme file,
and never runs an apply/move.** Its entire output is a proposed plan for
a human to review.

## When to Use

- A knowledge file (LEARNINGS.md, DECISIONS.md) has grown large and you
  want to see how it *would* be folded into themes
- The knowledge-growth nudge fired and you want a concrete digest plan
  before committing to restructuring
- To sanity-check theme groupings before the apply step exists

## When NOT to Use

- To actually move entries — the mover is a later milestone; this skill
  is preview-only
- On CONSTITUTION.md or TASKS.md — out of scope (small by design /
  auto-archived)
- On CONVENTIONS.md — its entries are `##` sections, not `## [` entries;
  staged-entry digestion there is a separate milestone

## Procedure

### 1. Inspect the root

```bash
ctx disclosure inspect .context/LEARNINGS.md --json
```

This reports, as JSON: the `kind`, the `staging` entries (timestamp +
title, un-digested), and the current `themes` (name, gist, link). Read
it; do not parse the file by hand.

If `staging` is empty, there is nothing to digest — say so and stop.

### 2. Propose a theme per staged entry

For each staged entry, assign it to a **theme** — an existing theme (by
name) or a new one. This is a semantic judgement: group entries that
share a subject (e.g. "hook mechanics", "error handling", "OpenCode
integration"), not by date.

- Keep themes **coarse**: a handful covering many entries beats one
  theme per entry.
- Prefer an **existing** theme when an entry fits it.
- Surface your proposed grouping to the human and let them **rename,
  merge, split, or reassign** before anything is final. Themes are the
  human's call; you propose.

### 3. Author a gist per touched theme

For each theme in the plan, write the **gist** — one line, soft ceiling
~140 chars, saying what the theme *covers* (the shape of its knowledge),
not listing its entries. "hook mechanics: output channels, key names,
compliance wiring" — not "entry A; entry B". The gist tells a future
reader *whether to drill in*, nothing more. (Spec: `### Gist format`.)

### 4. Present the plan — then stop

Show the human, per theme:

```
Theme: <name>   →  .context/<noun>/<slug>.md   (would create/append)
  gist: <proposed one-line gist>
  entries (N):
    - [<ts>] <title>
    - …
```

End with an explicit note: **"Dry run — nothing was moved or written.
The apply step (a later milestone) performs the move."** Do not edit any
file. Do not run an apply command (there is none yet).

## Quality Checklist

- [ ] Ran `ctx disclosure inspect --json`; did not hand-parse the file
- [ ] Every staged entry is assigned to exactly one theme in the plan
- [ ] Each theme has a one-line gist describing coverage (not a list)
- [ ] The plan was surfaced for the human to adjust themes
- [ ] **Nothing was edited, created, or moved** — output is a plan only

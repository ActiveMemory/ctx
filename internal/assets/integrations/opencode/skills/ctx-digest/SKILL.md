---
name: ctx-digest
description: "Run the progressive-disclosure pass: inspect a knowledge file's staging zone, propose a theme per staged entry, author gists, present the plan for human approval, then apply it — moving entries into theme files and folding gists into the root. Use when LEARNINGS.md / DECISIONS.md / CONVENTIONS.md grow large and should fold into themes."
---

Digest a bounded knowledge root: fold its staging zone into themes, moving
entry bodies out to per-theme files and leaving a compact gist + link in
the root. This is the full progressive-disclosure pass (see
specs/progressive-disclosure.md) — propose, get the human's sign-off, then
apply.

**The apply is human-gated.** You propose themes and gists; the human
approves before anything moves. The move itself is guarded by the CLI
(`ctx disclosure apply`): entry bodies are appended to their theme files
and verified byte-present before the root is touched, and the root is
rewritten once. Any failure leaves the root byte-identical. **You never
hand-edit a knowledge file or a theme file — the CLI does every write.**

## When to Use

- A knowledge file (LEARNINGS.md, DECISIONS.md, CONVENTIONS.md) has grown
  large and should be folded into themes
- The knowledge-growth nudge fired and it is time to digest the staging
  zone
- You want to preview the grouping first: run through step 4 and stop —
  presenting the plan without approval is a valid dry run

## When NOT to Use

- On CONSTITUTION.md or TASKS.md — out of scope (small by design /
  auto-archived)

## What counts as a staged entry

The layout is the same for every kind — `preamble | staging | ## Themes` —
but what the pass moves differs, and so does how an entry is named:

| Kind | Staged entry | Identity |
|---|---|---|
| learning, decision | a `## [<ts>] Title` entry | timestamp **+** title |
| convention | a `## Title` section, bullets and all | title alone |

A convention has no timestamp: its title *is* its identity. That has one
consequence you must handle — see the duplicate-title guard in step 1.

Folding a convention moves the **whole section**, every bullet under it,
into the theme file. You are grouping sections into themes, not splitting
them.

## Procedure

### 1. Inspect the root

```bash
ctx disclosure inspect .context/LEARNINGS.md --json
```

This reports, as JSON: the `kind`, the `staging` entries (un-digested), and
the current `themes` (name, gist, link). Read it; do not parse the file by
hand. For a convention root each staged entry carries an empty
`timestamp`.

If `staging` is empty, there is nothing to digest — say so and stop.

**Duplicate-title guard (conventions).** If two sections share a title,
the CLI refuses the whole pass with `ErrDuplicateStagedTitle` — with
title-only identity it cannot tell which section a plan means. Do not
work around it and do not rename anything yourself: report the duplicate
and ask the human to rename one section first.

### 2. Propose a theme per staged entry

For each staged entry, assign it to a **theme** — an existing theme (by
name) or a new one. This is a semantic judgement: group entries that
share a subject (e.g. "hook mechanics", "error handling", "OpenCode
integration"), not by date.

- Keep themes **coarse**: a handful covering many entries beats one
  theme per entry.
- Prefer an **existing** theme when an entry fits it.
- Give each new theme a **slug**: a short kebab-case form of the name
  (`hook mechanics` → `hook-mechanics`) — it becomes the theme file's
  basename (`.context/<noun>/<slug>.md`). Reuse the existing slug for an
  existing theme (read it from that theme's link).

The `<noun>` follows the kind: `learnings/`, `decisions/`, `conventions/`.

### 3. Author a gist per touched theme

For each theme in the plan, write the **gist** — one line, soft ceiling
~140 chars, saying what the theme *covers* (the shape of its knowledge),
not listing its entries. "hook mechanics: output channels, key names,
compliance wiring" — not "entry A; entry B". The gist tells a future
reader *whether to drill in*, nothing more. (Spec: `### Gist format`.)

### 4. Present the plan and ask for approval

Show the human, per theme:

```
Theme: <name>   →  .context/<noun>/<slug>.md   (create | append)
  gist: <proposed one-line gist>
  entries (N):
    - [<ts>] <title>        # conventions: just <title>
    - …
```

Let them **rename, merge, split, reassign, or reword gists** — themes are
the human's call. Then ask plainly: **"Apply this plan?"** Do not proceed
to step 5 without an explicit yes. If they only wanted a preview, stop
here — nothing has moved.

### 5. Apply the approved plan

Write the approved plan to a JSON file, then apply it. Each entry is an
**object** with `timestamp` and `title` — the same shape `inspect` reports
under `staging`, so lift them verbatim rather than re-typing them:

```json
{
  "kind": "learning",
  "assignments": [
    {
      "theme": "hook mechanics",
      "slug": "hook-mechanics",
      "gist": "hook mechanics: output channels, key names, compliance wiring",
      "entries": [
        {"timestamp": "2026-07-15-120000", "title": "a staged entry"}
      ]
    }
  ]
}
```

For a convention root, `kind` is `"convention"` and each entry carries the
title alone:

```json
{
  "kind": "convention",
  "assignments": [
    {
      "theme": "error handling",
      "slug": "error-handling",
      "gist": "error handling: wrapping, sentinels, user-facing message shape",
      "entries": [{"title": "Error Handling"}]
    }
  ]
}
```

```bash
ctx disclosure apply .context/LEARNINGS.md --plan /tmp/digest-plan.json
```

`apply` moves every listed entry into its theme file, folds the gists into
`## Themes`, and rewrites the root once. It prints how many entries moved
into how many themes. If it returns an error, **stop and relay it
verbatim** — the root is untouched; do not retry by hand.

### 6. Confirm

Re-inspect to confirm the staging zone shrank and the themes grew:

```bash
ctx disclosure inspect .context/LEARNINGS.md --json
```

The moved entries should be gone from `staging` and their themes present.
Report what moved.

## Related

Declaring a theme *without* moving anything into it is a different
operation, and it is not this skill:

```bash
ctx convention add --section Themes "error handling — how failures surface"
```

That writes the gist bullet and creates the theme file in one step. Use it
to name a theme up front; use this skill to fold entries into themes.

## Quality Checklist

- [ ] Ran `ctx disclosure inspect --json`; did not hand-parse the file
- [ ] Every staged entry is assigned to exactly one theme, each with a slug
- [ ] Entries in the plan JSON are **objects**, lifted verbatim from
      `inspect` (conventions: `title` only, no `timestamp`)
- [ ] Each theme has a one-line gist describing coverage (not a list)
- [ ] Presented the plan and got an **explicit approval** before applying
- [ ] Applied via `ctx disclosure apply` — never hand-edited a file
- [ ] A duplicate convention title was reported to the human, not renamed
- [ ] Re-inspected to confirm the move; relayed any apply error verbatim

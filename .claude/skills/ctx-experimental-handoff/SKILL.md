---
name: ctx-experimental-handoff
description: "EXPERIMENTAL (discardable). Hand a loose intent spec (.context/specs/intent-<slug>.md) off to spec-kit's /speckit-specify with a prose synopsis. Optional and graceful — warns and continues if spec-kit is not installed; the intent spec stands either way. Third step of the experimental chain."
allowed-tools: Bash(specify:*), Read, Glob
---

# Hand off an intent spec to spec-kit — experimental

> **Experimental / discardable.** ctx-native port of an external
> spec-handoff skill. It is the delegation seam canonical ctx does
> not have: ctx's chain terminates at `/ctx-implement`, whereas this one
> hands the intent spec to spec-kit. Trial it; if the seam is worth
> keeping, promote it to a real `/ctx-handoff`.

This skill delegates a loose **intent spec**
(`.context/specs/intent-<slug>.md`, written by `/ctx-experimental-spec`)
to a downstream spec-driven pipeline (spec-kit) by invoking its
**`/speckit-specify`** skill with a prose synopsis that points at the
intent spec. ctx owns the pre-spec debate; spec-kit owns
`specify → plan → tasks → implement` and owns REQ-ID numbering.

**This delegation is OPTIONAL and GRACEFUL:** if spec-kit is not
installed, the skill **warns and continues** — the intent spec still
stands and the user can run `/speckit-specify` themselves later. It
never hard-depends on spec-kit and never calls it over a network; the
handoff is on-disk + slash-command only.

## When to use this skill

- After `/ctx-experimental-spec` has written an
  `intent-<slug>.md` and you want to start spec-kit's build-out from it.
- User says: "hand this off to spec-kit", "kick off speckit", "pipe the
  spec onward".

## When NOT to use this skill

- The intent spec isn't written yet — run `/ctx-experimental-plan` →
  `/ctx-experimental-spec` first.
- The project doesn't use spec-kit — there's nothing to hand off to; the
  intent spec is already the deliverable.

## Procedure

1. **Resolve the intent spec.** Determine the slug (from the spec just
   written, or ask). Check whether `.context/specs/intent-<slug>.md`
   exists:
   - exists → continue with that path.
   - missing → it isn't written yet. Tell the user to run
     `/ctx-experimental-spec` first, and **stop**.

2. **Detect spec-kit (graceful gate).** Check whether spec-kit is
   present: the `/speckit-specify` skill is available, or the `specify`
   binary is on PATH, or a `.specify/` directory exists. If **absent**,
   say plainly:

   > Spec-kit not detected. The intent spec stands at `<path>` — run
   > `/speckit-specify` with it when spec-kit is set up.

   …and **stop cleanly** (this is success, not failure — graceful
   degradation).

3. **Compose the seed.** Read the intent spec and lift a **2–4 sentence
   synopsis** from its one-line promise + Goals — the spec's own words,
   not invented. `/speckit-specify` re-derives its own structured
   `spec.md` and REQ-IDs from this prose; do **not** pre-shape the
   synopsis into spec-kit's template.

4. **Delegate.** Invoke `/speckit-specify` with:

   > Implement `<slug>`. Per the intent spec at `<path>`:
   > `<synopsis>`

   Optionally export `SPECIFY_FEATURE_DIRECTORY` if the user wants to pin
   a specific `specs/<dir>`; otherwise let spec-kit auto-number.

5. **Close the loop.** After spec-kit mints `REQ-<PFX>-NNN`, remind the
   user to thread them back into ctx memory so the chain
   `intent spec → task → commit` stays traceable: capture the REQ-IDs in
   a task via `/ctx-task-add` (referencing this intent spec), and use
   `/ctx-decision-add` for any architectural call. ctx stores REQ-IDs; it
   never generates them.

## Ownership boundary

ctx owns everything up to and including the intent spec; spec-kit owns
`specify → plan → tasks → implement` and REQ-ID numbering. The two keep
independent on-disk trees (`.context/specs/intent-*.md` vs repo-root
`specs/<NNN-slug>/`); the boundary is crossed only by the prose argument
to `/speckit-specify`, never by spec-kit reading `.context/`.

## Anti-patterns

- **Hard-failing when spec-kit is absent** — it's optional; warn and
  continue, leaving the intent spec as the standing deliverable.
- **Passing the intent-spec FILE as the seed** — `/speckit-specify` reads
  a prose `$ARGUMENTS` description, not a file path.
- **Pre-minting REQ-IDs** in ctx — spec-kit owns numbering; ctx only
  stores them after the fact.
- **Pre-shaping the synopsis** into spec-kit's spec-template — it
  duplicates spec-kit's own derivation and splits REQ-ID authority.
- **Calling spec-kit over a network** — composition is on-disk +
  slash-command only.

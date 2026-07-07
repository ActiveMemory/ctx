# Spec: remove the experimental spec-kit delegation chain

**Status:** accepted (impl 2026-07-06)

## Problem

Three project-level skills under `.claude/skills/` —
`ctx-experimental-plan`, `ctx-experimental-spec`, and
`ctx-experimental-handoff` — were carried as an explicitly
**EXPERIMENTAL (discardable)** trio. Their own frontmatter and
lede state the contract: they are ctx-native ports of external
skills, kept as local project skills so the chain
`/ctx-experimental-plan → /ctx-experimental-spec →
/ctx-experimental-handoff` could be trialed, and *"if it earns its
keep, fold the good parts into the real skills; otherwise `rm -rf`
this directory."*

The chain's sole reason for existing is to **bridge ctx to
spec-kit** (an external spec-driven-development framework):
`ctx-experimental-handoff` invokes spec-kit's `/speckit-specify`,
and `ctx-experimental-spec` deliberately writes a *loose* intent
spec to `.context/specs/intent-<slug>.md` as that bridge's input.
This project does not use spec-kit — no spec, skill, or Makefile
target references `speckit`/`spec-kit` (grep: zero hits outside the
experimental trio itself). With the bridge unused:

- `ctx-experimental-plan` is a pure duplicate of the canonical
  `/ctx-plan` (adversarial interview → debated brief under
  `.context/briefs/<TS>-<slug>.md`).
- `ctx-experimental-spec` is a lesser variant of `/ctx-spec` whose
  only distinguishing purpose (seeding spec-kit) does not apply.
- `ctx-experimental-handoff` has nothing to hand off to.

The trio lives **only** in this repo's `.claude/skills/`; it is
**not** shipped in `internal/assets/`, so removal affects only this
repo's local dev surface and nothing downstream for ctx users.

## Design

`rm -rf` the three skill directories. No code change; no shipped
asset change. The skills' content is preserved in git history, so
removal loses no institutional memory (consistent with the Context
Preservation Invariant, which forbids deleting *history*, not
pruning a superseded live artifact).

If the spec-kit seam is ever wanted, the canonical path is to
promote a real `/ctx-handoff` rather than resurrect the trial trio;
the debate and spec halves already exist as `/ctx-plan` and
`/ctx-spec`.

## Scope

- In: remove the three `.claude/skills/ctx-experimental-*`
  directories; add this spec as the rationale.
- Out: no change to canonical `/ctx-plan`, `/ctx-spec`,
  `/ctx-task-out`, or `/ctx-implement`; no change to shipped assets
  under `internal/assets/`.

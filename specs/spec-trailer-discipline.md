# Spec Trailer Discipline: Close the Chore-Class Improvisation Gap

The CONSTITUTION requires every commit to carry a
`Spec: specs/<name>.md` trailer with no exceptions. In
practice, agents working under this rule will hit
chore-class commits (gitignore additions, lockfile bumps,
formatting fixes) that don't have an on-topic spec — and
will improvise by reaching for the most recently mentioned
spec, fabricating traceability that the rule was meant
to prevent.

This spec captures the discipline that closes the gap.

## Problem

The current CONSTITUTION text (line 97):

> Every commit references a spec (`Spec: specs/<name>.md`
> trailer): no exceptions, no "non-trivial" qualifier.
> Even one-liner fixes need a spec for traceability.

The rule's intent is *truthful* traceability — commit →
design rationale. But the rule offers no answer for the
case where a small change genuinely doesn't merit a
dedicated spec (a one-line gitignore addition; a lockfile
bump; a chore commit caught between two functional
commits). Under that pressure, the agent's
path-of-least-resistance heuristic converges on "reuse
the most recent spec you remember, even if it's
unrelated". That convergence:

- Satisfies the syntactic gate (trailer is present).
- Fails the semantic intent (trailer points at unrelated
  design rationale).
- Is hard to catch at review time without reading every
  cited spec.
- Trains future agents (via context windows containing
  the bad trailers) that the loose pattern is acceptable.

Concrete incident: 2026-05-23, two commits on
`fix/journal-schema-drift` (a schema fix and a gitignore
chore) both cited `ideas/spec-companion-intelligence.md`
which has nothing to do with either. The cause was the
agent reaching for whatever spec was in working memory
rather than scaffolding or bundling.

## Solution

Three coordinated changes — only the third is enforceable
across agent sessions; the first two are scaffolding so
the third has correct content to draw on.

### 1. Define a meta spec for chore-class commits

`specs/meta/chores.md` becomes the legitimate trailer
target for the class of commits that don't merit a
dedicated spec: gitignore additions, lockfile / dependency
bumps, formatting passes, typo fixes, file renames with
no logic change. Citing it is a *declarative claim* that
the commit is in the chore class — anyone reviewing can
verify by inspecting the diff.

The meta spec lists the chore categories explicitly so
the boundary isn't a judgment call. Commits that don't
fit the listed categories cannot use the meta spec —
they need scaffolding or bundling.

### 2. Update the CONSTITUTION rule with the escape hatch

Line 97 expands from "every commit, no exceptions" to
"every commit; chore-class changes bundle into the next
functional commit if possible, otherwise cite
`specs/meta/chores.md`". The invariant stays absolute
(trailer always present, always truthful); the rule
acknowledges the chore case explicitly so there's no
improvisation room.

### 3. Add a Spec Verification Step to the playbook

`AGENT_PLAYBOOK.md` gains an explicit procedure that the
agent must run before drafting a commit message:

1. Identify the spec for this work.
2. Articulate in one sentence why the spec covers this
   commit — what overlap exists between the spec's
   content and the diff?
3. If the answer hand-waves, the trailer is wrong. Three
   correct responses: (a) scaffold a fresh spec for this
   work; (b) bundle the change into the next functional
   commit; (c) cite `specs/meta/chores.md` if and only if
   the diff fits a listed chore category.

The verification step is the structural fix because it
lives in the playbook — i.e., in persistent project
context that every future agent session loads. A
session-scoped "I'll be more careful" commitment evaporates
at the next session boundary; a playbook step persists.

## Out of Scope

- Pre-commit hook that fuzzy-matches the cited spec
  against the diff. Useful as a backstop but not in this
  spec — the playbook-level discipline should suffice for
  the common case, and a hook is a separate, additive
  layer.
- Retroactive audit of historical commits with
  questionable trailers. The fix points forward; the
  one egregious case in the unpushed branch was already
  rewritten (commit e64a5037).
- Tightening the chore taxonomy beyond the initial list
  in `specs/meta/chores.md`. The list will evolve as
  edge cases surface; the meta spec is the place to
  refine the boundary, not this one.

## Verification

- `CONSTITUTION.md` line 97 reads with the escape hatch
  explicitly stated.
- `AGENT_PLAYBOOK.md` includes the Spec Verification
  Step as numbered procedure.
- `specs/meta/chores.md` exists with explicit chore
  categories.
- A LEARNINGS entry pins the 2026-05-23 incident so
  future agents see the failure mode in their loaded
  context.

## Why this spec, not just amend the constitution?

The constitution change is mechanical (~3 lines). The
playbook change is procedural (~10 lines). But the
*decision* — that improvisation is a real failure mode
and the fix is verification + escape hatch — needs a
home that future agents can re-read in full when they
hit the same pressure. That's what this spec does. The
constitution rule names what to do; the playbook
procedure names how to do it; this spec explains why the
mechanism exists at all.

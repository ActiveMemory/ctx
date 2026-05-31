# Spec: remove the deprecated session-handoff scratchpad

**Status:** accepted (impl 2026-05-30)

## Problem

`.context/scratchpad-handoff.md` is a one-off, hand-written session
handoff artifact (created in `dcfd3772`, "Add session handoff for
context window continuation"; content dated 2026-03-28). Session
handovers are now produced by `/ctx-wrap-up` → `/ctx-handover` and
live under `.context/handovers/<TS>-<slug>.md` — timestamped so
concurrent agent runs never overwrite (per the root `CLAUDE.md`).
The root-level scratchpad is the superseded predecessor of that
mechanism.

It is not merely redundant, it is actively misleading. Its body
still claims "NOT YET COMMITTED … Very large diff" for work
(EH.x error-handling, line-width audit) that has long since
landed; the tree is clean. `ctx status` and `ctx system bootstrap`
enumerate `.context/*.md`, so this stale file is loaded into every
session's start-up context, where a "do you remember?" recall can
resurrect a phantom uncommitted diff.

No code reads it: a repo-wide search of `*.go`, `*.md`, `*.yaml`,
and `*.json` finds zero references to `scratchpad-handoff` outside
the file itself. Removing it cannot break a code path; it only
stops the context pollution.

## Design

`git rm .context/scratchpad-handoff.md`. No code change.

The file's content is preserved in git history (`dcfd3772`,
`0df91654`), so removal loses no institutional memory — consistent
with the Context Preservation Invariant, which forbids *deleting
history*, not pruning a superseded live artifact whose history
remains in the reflog. Going forward, `.context/handovers/` is the
sole handoff channel.

## Scope

- In: remove the single stale file; add this spec as its rationale.
- Out: no change to the `/ctx-wrap-up` → `/ctx-handover` mechanism
  or to `.context/handovers/`; the scratchpad's obsolete content is
  not migrated anywhere (it is already false).

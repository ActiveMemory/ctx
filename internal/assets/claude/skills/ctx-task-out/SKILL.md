---
name: ctx-task-out
description: "Decompose a committed spec into a per-milestone implementation plan at specs/plans/<milestone>.md — data model, contracts, invariant-test matrix, and tasks with falsifiable acceptance criteria — that /ctx-implement consumes. Use after /ctx-spec when a spec is too large to implement in one session."
---

## Canonical Chain

The project's design-to-implementation pipeline is:

```text
/ctx-brainstorm → /ctx-plan → /ctx-spec → /ctx-task-out → /ctx-implement
    (vague)     (contested)  (committed)   (decomposed)     (execution)
```

`/ctx-task-out` is the fourth step. It consumes a committed spec
(`--spec <path>`) and produces the *plan document* that
`/ctx-implement` executes. It closes a gap the chain otherwise
leaves unowned: `/ctx-plan` explicitly disclaims implementation
planning, `/ctx-spec` commits the what/why at spec altitude, and
`/ctx-implement` opens with "use when you have a plan document" —
this skill is what produces that document.

Small specs skip this step. If the whole spec is implementable in
roughly one session, go straight to `/ctx-implement` with the spec
itself.

## Role

You decompose; you do not redesign the bet. The spec is
committed; do not relitigate scope, behavior, or the bet here —
disagreements with the spec go back through `/ctx-plan`.
Implementation structure is different: choosing schema shapes,
signatures, and index strategies is exactly the job, because
resolving those decisions *before* execution is the point of
this skill. Make every task falsifiable and every design detail
the implementer needs explicit before execution starts, so no
large decision is made mid-flight.

Authority boundary: invariants, validation rules, and behavior
come *from the spec*. If decomposition surfaces an invariant the
spec never states, that is a spec gap — surface it and mark it
`TBD`; do not mint it here.

## When to Use

- After `/ctx-spec`, when the spec spans milestones/phases or
  exceeds ~one session of implementation
- When a TASKS.md phase references a spec but its tasks are coarse
  and carry no acceptance criteria
- When `/ctx-implement` is invoked without a plan document
  (redirect here first)

## When NOT to Use

- Single-session features — the spec *is* the plan
- The spec is not committed yet (`/ctx-spec` first)
- The bet is still contested (`/ctx-plan` first)
- Decomposing milestone N+1 while milestone N's DoD is unmet
  (see rolling-wave gate)

## Usage

```text
/ctx-task-out --spec specs/v1-substrate.md --milestone m0a
/ctx-task-out --spec specs/rss-feed.md        # single-milestone: whole spec
```

Without `--milestone`, the plan file takes the spec's basename:
`specs/plans/rss-feed.md`.

## Preconditions (hard gates — refuse, do not degrade)

1. **Spec exists.** If `--spec` is missing or the file is absent,
   stop and report; no interactive fallback.
2. **Blocking-TBD gate.** Enumerate the spec's Open Questions /
   `TBD` entries and classify each as *blocking* or *deferrable*
   for the target milestone. A TBD is blocking if any task in the
   milestone would embed an assumption about its answer (language
   choice, storage engine, schema format…). Refuse to decompose
   past a blocking TBD: list the blockers, name who can resolve
   them, and stop. Deferrable TBDs do not vanish: carry each into
   the plan (Out of scope or Risks), annotated with the milestone
   at which it becomes blocking. Resolving a blocking TBD is a
   spec edit or a DECISIONS.md entry *first*; the plan only
   points at that record. A resolution that exists nowhere but
   the plan is minting.
3. **Rolling-wave gate.** If a prior milestone's plan exists and
   its DoD is not checked off, refuse to decompose the next
   milestone. The user may override explicitly; log the
   override in the plan's Amendments section. Tasking distant milestones
   produces fiction — the current milestone's measurements are
   allowed to reshape everything downstream.

Milestone boundaries belong to the spec. If decomposition shows
the cut is wrong — one "milestone" hiding several, or a boundary
in the wrong place — stop and route the resize through the spec;
do not mint sub-milestones here.

## Process

0. **Detect mode.** If the target plan file already exists, this
   is an amendment run, not a fresh decomposition: read the
   existing plan, classify the change as obsolete/append, re-run
   only the blocking-TBD gate against the delta, and log the
   change in the plan's Amendments section. The rolling-wave
   gate does not fire when amending the current milestone.
   Steps 1–8 below describe a fresh run.

1. Read the spec in full. Read TASKS.md, DECISIONS.md, and
   CONVENTIONS.md from the context directory.
2. Run the blocking-TBD gate; surface the classification to the
   user before proceeding.
3. Draft the plan sections (structure below). Lift from the spec
   verbatim where it speaks; where it is silent, prefer asking or
   marking `TBD` over inventing — same authority discipline as
   `/ctx-spec --brief`.
4. Break down tasks: typically 15–40 per milestone. Each task
   carries: id, title, dependencies (by id), the files/paths it
   is expected to touch, a `[P]` marker when
   parallelizable with its siblings, a **falsifiable acceptance
   criterion** (a command to run, a test that must pass, an
   observable behavior), and a reference to the spec section it
   implements. Size each task to roughly one commit — small
   enough that a failed acceptance check localizes the fault to
   that task. `[P]` is mechanical, not aspirational: no
   dependency edge, no file touched by a sibling `[P]` task, no
   shared sequence (e.g. migration numbers). File disjointness
   is checkable from the files column — which is also how an
   amendment run detects a new task colliding with one in
   flight.
5. Build the test matrix: every invariant, validation rule, and
   edge case the milestone touches × the attempted violation × the
   expected failure mode × the task id whose acceptance criterion
   exercises it. A matrix row no task exercises is documentation,
   not execution.
6. Write `specs/plans/<milestone>.md` (create `specs/plans/` if
   absent).
7. Sync anchors to TASKS.md: **epic-level anchors only** — one per
   task cluster, each annotated `Plan: specs/plans/<milestone>.md`.
   One-way sync, plan → TASKS.md. Never duplicate the full task
   list into TASKS.md; never move or delete existing entries
   (CONSTITUTION).
8. Hand off: report blockers resolved/remaining and suggest
   `/ctx-implement` against the plan.

## Plan Document Structure

```markdown
# <Milestone> Plan — <short name>

**Spec:** <path> · **Status:** Ready | Blocked
**Blocking TBDs resolved:** <list, with where each was decided>

## Scope & DoD          (lifted from the spec's milestone entry)
## Data model & storage (DDL, migrations, indexes)
## Contracts            (API signatures, schemas, CLI surface)
## Test matrix          (invariant × violation attempt × expected failure × task ref)
## Task breakdown       (table: id · task · deps · files · [P] · acceptance criterion · spec ref)
## Risks & measurement gates  (results that may reshape later tasks)
## Out of scope         (deferred to later milestones, with pointers)
## Amendments           (date · what · why — appended by amendment runs)
```

The plan is the **execution ledger**: the task table carries
per-task checkbox state, Scope & DoD carries the DoD checkboxes,
and `/ctx-implement` updates both as it executes. DoD is
confirmed by measurement or by the user — never derived from
task completion — and the rolling-wave gate reads the DoD
checkboxes only. No other record of milestone progress exists.
`Status: Blocked` is reachable only by amendment: a fresh run
refuses instead of writing a Blocked plan; the status marks a
deferrable TBD that graduated to blocking mid-milestone.

## Amendments (Mid-Milestone Changes)

Plans meet reality; the plan document owns that contact. When a
measurement gate fires or the implementer hits a wall:

- Tasks may be marked `obsolete` with a one-line reason; never
  deleted.
- New tasks are appended with fresh ids; ids are never reused.
- An acceptance criterion is **never edited in place** once its
  task has started — weakening the test until it passes is the
  failure mode this rule exists to prevent. A criterion change
  is a re-invocation of `/ctx-task-out` against the same
  milestone: with the plan already present, the skill operates
  in amendment mode — read the existing plan, apply the change
  as obsolete-and-append, and log date · what · why in the
  plan's Amendments section.
- Disagreements with the *spec* discovered mid-flight still
  route through `/ctx-plan`; amendments cover implementation
  reality, not the bet.

## Quality Checklist

Before writing the file, verify:

- [ ] Every task has a falsifiable acceptance criterion — no
      "implement X" without a way to check it happened
- [ ] No task depends on an unresolved blocking TBD
- [ ] Every invariant the milestone touches appears in the test
      matrix, and every matrix row is exercised by a task's
      acceptance criterion (by id)
- [ ] Task ids admit a topological order — verify by listing
      execution waves; `[P]` siblings share no files, edges, or
      sequences
- [ ] TASKS.md gained anchors only — nothing moved, nothing deleted
- [ ] The plan is implementable-alone: a fresh agent holding only
      the plan and the spec can state the acceptance check for
      any task without asking a question

---
title: "Spec-Driven Development"
icon: lucide/milestone
---

![ctx](../images/ctx-banner.png)

## The Problem

A feature big enough to span several milestones doesn't fail at the
keyboard. It fails at the *seams*: the bet gets re-argued halfway
through implementation, a "plan" for milestone three turns out to be
fiction by the time you reach it, a decision the code silently assumed
was never written down, and two different files each claim to be the
authoritative list of what's done.

The five skills that make up the design-to-implementation pipeline each
solve one seam. But the pipeline only holds together if you understand
*which skill owns which decision, and at what altitude*. Read the skill
texts in isolation and the chain looks like five ways to write a
Markdown file. Run them without the mental model and you end up
reverse-engineering the whole thing from error messages.

**This recipe is that mental model, from the operator's seat.** It
walks one invented-but-realistic feature — a weekly context digest —
through all five stages, and calls out the five load-bearing rules that
aren't obvious from any single skill.

!!! note "Relationship to *Design Before Coding*"
    [Design Before Coding](design-before-coding.md) is the gentle
    on-ramp: brainstorm → spec → task-out → implement, four skills, one
    small feature. This recipe is the full chain **including the
    debated-brief step** (`/ctx-plan`), aimed at multi-milestone work
    where the seams actually bite. If you only ever ship
    single-session features, the on-ramp is enough.

## TL;DR

```text
/ctx-brainstorm                                       # shape the vague idea
/ctx-plan                                             # debate the bet → a brief
/ctx-spec --brief .context/briefs/<TS>-<slug>.md      # commit the whole spec
/ctx-task-out --spec specs/<feature>.md --milestone m0   # decompose ONE milestone
/ctx-implement specs/plans/m0.md                      # execute, verify, checkpoint
```

Five skills, one direction. The canonical chain, with the altitude each
step works at:

```text
/ctx-brainstorm → /ctx-plan → /ctx-spec → /ctx-task-out → /ctx-implement
    (vague)     (contested)  (committed)   (decomposed)     (execution)
```

`/ctx-plan` is not optional decoration. It is where the bet is
*attacked* and written down as a debated brief, before the spec commits
to it. Skip it and the spec inherits an unexamined bet; the argument you
avoided resurfaces mid-implementation, where it is most expensive.

## Commands and Skills Used

| Tool                | Type  | Purpose                                                        |
|---------------------|-------|----------------------------------------------------------------|
| `/ctx-brainstorm`   | Skill | Turn a vague idea into a validated design (conversation only)  |
| `/ctx-plan`         | Skill | Attack the bet; write a *debated brief* to `.context/briefs/`  |
| `/ctx-spec`         | Skill | Absorb the brief into a committed spec covering all milestones |
| `/ctx-task-out`     | Skill | Decompose **one** milestone into `specs/plans/<milestone>.md`  |
| `/ctx-implement`    | Skill | Execute the plan step-by-step, updating the execution ledger   |
| `/ctx-decision-add` | Skill | Record a blocking decision the milestone forces into the open  |

## The Running Example

Every stage below moves the same feature forward. Follow it end to end
rather than reading the stages as isolated commands.

> **Weekly Context Digest** (`ctx digest`). Over dozens of sessions,
> `DECISIONS.md` and `LEARNINGS.md` pile up and nobody re-reads them. A
> periodic rollup resurfaces *what changed this week* so it doesn't rot.

The feature has four milestones. Only the first is planned in detail
right now; the rest exist in the spec as boundaries, not task lists.

| Milestone | Scope                                                  |
|-----------|--------------------------------------------------------|
| **m0**    | Render: `ctx digest --since <date>` prints a Markdown rollup to stdout |
| **m1**    | Schedule: surface the digest at session start, at most weekly |
| **m2**    | Deliver: push the digest through the `ctx hook notify` webhook |
| **m3**    | *(deferred)* Roll up across projects via `ctx hub`     |

## Stage 1 — Debate the bet (`/ctx-plan`)

You have a rough shape (`/ctx-brainstorm` got you here) but you have not
committed. `/ctx-plan` is a skeptical collaborator: it attacks the bet
one question at a time — scope, failure modes, rejected alternatives,
sequencing, reversibility — and stops when *you* can state the bet, what
you rejected, the top failure modes, the cheapest validation, and the
unwind cost, without help.

```text
You: /ctx-plan
     I want a weekly digest that summarizes what changed in the
     context files.

Agent: One question at a time. Generated from the structured files, or
       summarized by an LLM from free text?

You:   ...generated. An LLM summary would drift and I couldn't cite it.

Agent: Then the cheapest thing that proves the bet is m0 render-only —
       no schedule, no delivery. What would tell you the digest is
       worth building delivery for?

You:   Whether I actually read it. Ship m0, live with it a week, then
       decide on m1.
```

The deliverable is **not** a task list. It is a *debated brief* at
`.context/briefs/<TS>-context-digest.md`: the bet (render from
structured files, not LLM summary), the rejections (LLM summary; a
separate digest datastore), the failure modes (empty window, giant
window, digest nobody reads), the validation route (ship m0, measure
readership), and the unwind cost.

The bet is debated **once, here, for the whole feature.** Not once per
milestone. That is the altitude rule — see below.

## Stage 2 — Commit the spec (`/ctx-spec --brief`)

`/ctx-spec --brief <path>` treats the brief as authoritative and skips
the interactive Q&A. It lifts the brief into the spec template
verbatim where the brief speaks, and writes `TBD` where the brief is
silent — it does not invent. The result is `specs/context-digest.md`,
covering **all four milestones**, plus the open questions the bet left
unresolved:

```text
Open Questions
  TBD-1  Default digest window — 7 days, or "since last run"?
  TBD-2  Delivery transport — webhook only, or email too?
  TBD-3  Cross-project identity/dedup for the hub rollup
```

Note what the spec is *not*: it is not a task breakdown, and it does not
resolve the TBDs. It commits the what/why at spec altitude and hands the
TBDs downstream, each to be forced open at the milestone that actually
embeds it. A small, single-session spec would stop here and go straight
to `/ctx-implement`. This one spans milestones, so it goes to
`/ctx-task-out`.

## Stage 3 — Task out one milestone (`/ctx-task-out`)

This is the step the other four skills leave unowned: `/ctx-plan`
disclaims implementation planning, `/ctx-spec` stops at spec altitude,
and `/ctx-implement` opens with "use when you have a plan document."
`/ctx-task-out` produces that plan — for **one** milestone at a time.

```text
You: /ctx-task-out --spec specs/context-digest.md --milestone m0

Agent: Blocking-TBD gate for m0:
        TBD-1 (default window)  — DEFERRABLE. m0 takes an explicit
              --since; no task assumes a default. Blocking at m1.
        TBD-2 (transport)       — DEFERRABLE. m0 prints to stdout.
              Blocking at m2.
        TBD-3 (hub identity)    — DEFERRABLE. Blocking at m3.
       No blocker for m0. Decomposing.

       Wrote specs/plans/m0.md: data model, CLI contract, test matrix,
       10 tasks (T01–T10) with falsifiable acceptance criteria and an
       `st` column. TASKS.md gained one epic anchor:
       "Digest m0: render — Plan: specs/plans/m0.md (T01–T10)".
```

Two gates run here, and both are refusals, not warnings:

- **Blocking-TBD gate.** Every open question is classified *blocking or
  deferrable for this milestone*. A TBD is blocking if any task in the
  milestone would embed an assumption about its answer. The skill
  refuses to decompose past a blocker — it lists them, names who can
  resolve them, and stops. Resolution lands as a spec edit or a
  `DECISIONS.md` entry **first**; the plan only points at that record. A
  decision that exists nowhere but the plan is *minting*, and the gate
  exists to stop it.
- **Rolling-wave gate.** If a prior milestone's plan exists and its
  definition of done (DoD) is not checked off, the skill refuses to
  decompose the next one. Tasking a distant milestone produces fiction:
  the current milestone's measurements are allowed to reshape everything
  downstream.

## Stage 4 — Implement the milestone (`/ctx-implement`)

`/ctx-implement specs/plans/m0.md` breaks the plan into atomic steps,
verifies each (build, test, or the task's own acceptance criterion,
run verbatim), and checkpoints every 3–5 steps. Its bookkeeping is what
keeps the two surfaces honest:

- It flips a task's `st` cell to `[x]` **only** when the acceptance
  criterion demonstrably passed — the command ran, the test is green,
  the behavior was observed. `st` never moves backward silently.
- When every id in an epic's range is `[x]` (or `[o]`, obsoleted by
  amendment), it marks that TASKS.md epic `[x]`. Sync is one-way, plan →
  TASKS.md.
- It leaves the **DoD boxes alone.** DoD is confirmed by measurement or
  by you — never derived from the tasks that "cover" it being done.

For the digest, m0's DoD is not "T01–T10 are all `[x]`." It is: *you
ran `ctx digest --since <date>` against a real project and confirmed the
output is readable and correct.* You check that box after reading a real
digest — because the whole bet was "is this worth reading?" and only
measurement answers it.

## Stage 5 — Roll the wave (the next milestone)

A week later you have read a few digests and decided m1 is worth
building. Now — and only now — you check m0's DoD box. That single act
opens the rolling-wave gate for m1:

```text
You: /ctx-task-out --spec specs/context-digest.md --milestone m1

Agent: Rolling-wave gate: m0 DoD is checked. Proceeding.
       Blocking-TBD gate for m1:
        TBD-1 (default window) — NOW BLOCKING. m1 schedules, so a task
              must assume a default window. Resolve before I decompose.
       Stopping. Resolve TBD-1 (spec edit or DECISIONS.md), then re-run.

You: /ctx-decision-add  "Digest default window is 'since last run',
     falling back to 7 days on first run"

Agent: Recorded D-014. Re-run /ctx-task-out --milestone m1 when ready.
```

The debate that would have been a *per-milestone re-litigation of the
bet* is instead a single, scoped decision — exactly the one m1 embeds —
forced into `DECISIONS.md` before any task can silently assume an
answer. That is the blocking-TBD gate doing the job per-milestone
debates used to do, without reopening the bet.

## The Rules the Diagram Doesn't Show

The arrows tell you the *order*. These five rules tell you *why the
order holds* — and they are what a newcomer has to reverse-engineer.

### 1. Altitude: the bet is debated once

The brief is **per-bet, never per-milestone.** `/ctx-plan` debates the
bet one time; `/ctx-spec` commits it across every milestone;
`/ctx-task-out` decomposes — it *does not redesign the bet*. If
decomposition makes you want to re-argue scope or behavior, that is a
signal to route back up to `/ctx-plan`, not to quietly change course in
the plan. Milestones are altitudes of *execution*, not fresh betting
opportunities.

### 2. Plans are just-in-time, behind the rolling-wave gate

You plan the milestone you are about to build, and no further. A plan
for a milestone three steps out is written against measurements you
have not taken yet — it is fiction with a task table. The rolling-wave
gate enforces this mechanically: milestone N+1 cannot be tasked out
while milestone N's DoD is unmet. (You *can* override explicitly; the
override is logged in the plan's Amendments section, so the fiction is
at least on the record.)

### 3. Blocking-TBD gates replace per-milestone debates

Because the bet is debated once, milestones don't get their own
debates. What they get is the blocking-TBD gate: each `/ctx-task-out`
run forces open *exactly* the decisions that milestone's tasks would
otherwise embed as silent assumptions — no more, no fewer. A deferrable
TBD doesn't vanish; it is carried into the plan (Out of scope or Risks),
annotated with the milestone at which it graduates to blocking. This is
how a big, half-decided spec becomes buildable without a design
committee at every step.

### 4. Two surfaces, one truth

There are two places milestone progress appears, and only one is
authoritative:

- **The plan (`specs/plans/<milestone>.md`) is the execution ledger.**
  Its task table has an `st` column — `[ ]` pending, `[x]` done, `[o]`
  obsoleted — and its Scope & DoD section carries the DoD checkboxes.
  This is the single source of truth for what's done.
- **TASKS.md epics are one-way projections.** Each epic anchor carries a
  *disjoint* task-id range (`Plan: specs/plans/m0.md (T01–T10)`); the
  ranges partition the plan's ids with none double-counted. An epic is
  checked `[x]` only when its whole range is `[x]`/`[o]` in the plan.
  Sync flows plan → TASKS.md, never back.

And the load-bearing exception: **DoD is confirmed by measurement or by
you, never derived from task completion.** All ten tasks green does not
check the DoD box. The rolling-wave gate reads *only* the DoD box — so
if you let task completion auto-derive it, you have quietly disabled the
gate that stops you from planning fiction.

### 5. When a new brief is legitimate

Going back to `/ctx-plan` mid-feature is not failure — but only for the
right reasons. A new brief is warranted when:

- **A deferred bet returns.** m3 (the hub rollup) was parked as Out of
  scope. Months later you want it. That is a *new bet* — deferred
  machinery coming back — so it earns a fresh `/ctx-plan` pass and its
  own brief. It is not an amendment to m0.
- **Evidence falsifies the committed bet.** If mid-m2 the measurements
  show webhook delivery is the wrong transport entirely, that
  disagreement is with the *spec*, and it routes *up* through
  `/ctx-plan`.

What is **never** legitimate is relitigating the bet *from below* — at
the implement seat, by weakening a task's acceptance criterion until it
passes, or by inventing a decision in the plan that the spec never made.
Amendments cover implementation reality (a task obsoleted, a new task
appended, a measurement gate that fired); the bet is contested only at
plan altitude, in the open.

## Tips

* **Don't skip `/ctx-plan` to "save time."** The argument you skip
  doesn't disappear; it moves to implementation, where it costs the
  most. Ten minutes of adversarial interview is cheap insurance.
* **Let the DoD box be earned.** The temptation to tick it when the
  tasks are all green is exactly the failure the rolling-wave gate
  guards against. Leave it for measurement or your own confirmation.
* **A blocking TBD is a feature, not a blocker.** When `/ctx-task-out`
  refuses, it just told you the one decision this milestone can't fake.
  Record it (`/ctx-decision-add`) and re-run — that is the workflow
  working.
* **Never edit an acceptance criterion in place once its task has
  started.** Weakening the test until it passes is the exact failure the
  amendment rule exists to prevent. A criterion change is an
  `/ctx-task-out` amendment run, logged with date · what · why.
* **One milestone in flight at a time.** If you find yourself wanting to
  task out two milestones before finishing the first, that is the
  rolling-wave gate telling you the first isn't actually done.

## See Also

* [Design Before Coding](design-before-coding.md): the four-skill
  on-ramp; start there if the feature fits one session.
* [Scrutinizing a Plan](scrutinizing-a-plan.md): a deeper look at the
  `/ctx-plan` adversarial interview and the debated brief.
* [Tracking Work Across Sessions](task-management.md): the TASKS.md
  epic anchors the plan projects into.
* [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md):
  where a blocking TBD gets recorded when the gate forces it open.
* [Skills Reference: /ctx-plan](../reference/skills.md#ctx-plan):
  the debated-brief contract.
* [Skills Reference: /ctx-task-out](../reference/skills.md#ctx-task-out):
  blocking-TBD and rolling-wave gates, the execution ledger.
* [Skills Reference: /ctx-implement](../reference/skills.md#ctx-implement):
  ledger duties and step verification.

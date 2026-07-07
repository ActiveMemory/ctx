---
title: "Design Before Coding"
icon: lucide/drafting-compass
---

![ctx](../images/ctx-banner.png)

## The Problem

You start coding a feature. Halfway through, you realize the approach
doesn't handle a key edge case. You refactor. Then you discover the
CLI interface doesn't fit the existing patterns. More refactoring.

The design work happened *during* implementation, mixed in with
debugging and trial-and-error. The result works, but the spec
was never written down, the trade-offs were never recorded, and
the next session has no idea why things are shaped this way.

**How do you front-load design so the implementation is straightforward?**

## TL;DR

```text
/ctx-brainstorm          # explore the design space
/ctx-spec                # write the spec document
/ctx-task-out            # decompose into a milestone plan
/ctx-implement           # execute step-by-step
```

Four skills, used in sequence. Each produces an artifact that feeds
the next. For specs small enough to implement in one session, the
spec doubles as the plan: skip `/ctx-task-out` and break the work
into TASKS.md entries with `/ctx-task-add` instead.

## Commands and Skills Used

| Tool                | Type  | Purpose                                            |
|---------------------|-------|----------------------------------------------------|
| `/ctx-brainstorm`   | Skill | Structured design dialogue: explore approaches     |
| `/ctx-spec`         | Skill | Scaffold and fill out a spec from the template     |
| `/ctx-task-out`     | Skill | Decompose a spec into `specs/plans/<milestone>.md` |
| `/ctx-task-add`     | Skill | Add implementation tasks to TASKS.md               |
| `/ctx-implement`    | Skill | Execute a plan step-by-step with verification      |
| `/ctx-decision-add` | Skill | Record design choices made during brainstorming    |

## The Workflow

### Step 1: Brainstorm the Design

Start with `/ctx-brainstorm` when you have a rough idea but haven't
committed to an approach. The skill follows a gated process:

1. **Understand context**: what exists today, what's the gap
2. **Clarify the idea**: one question at a time, not a wall of questions
3. **Surface non-functional requirements**: performance, security, compatibility
4. **Lock understanding**: confirm before exploring solutions
5. **Explore 2-3 approaches**: with trade-offs for each
6. **Stress-test**: poke holes in the chosen approach
7. **Present the design**: detailed enough to implement from

```text
You: "I want to add RSS feed generation to the blog."

Agent: "Before we design, let me understand: is this for the journal
  site, the docs site, or both?"

You: "Journal site only."

Agent: [asks about update frequency, content scope, existing infra...]
       [presents 3 approaches with trade-offs]

You: "Option 2 - generate at build time."

Agent: "Let me stress-test that. What happens when..."
       [probes edge cases]
       "Design looks solid. Want me to record the approach as a
       decision and write a spec?"
```

The brainstorm produces a *validated design* in conversation. It doesn't
produce a file - that's what the spec is for.

### Step 2: Write the Spec

Once the design is clear, use `/ctx-spec` to produce the written artifact.
The skill reads the project's spec template (`specs/tpl/spec-template.md`)
and walks through each section:

```text
You: /ctx-spec (rss feed generation)

Agent: "Creating specs/rss-feed.md from template."

       "Problem: What user-visible problem does this solve? Why now?"

You: "Users can't subscribe to blog updates. They have to check
      the site manually."

Agent: [writes the Problem section, moves to Approach]

       "Approach: High-level, how does this work?"

You: "Generate feed.xml during ctx journal site. Atom format."

Agent: [writes Approach, moves to Happy Path, then Edge Cases...]
```

The skill spends extra time on **Edge Cases** and **Error Handling**: These
are where specs earn their value. The skill pushes for at least three edge cases
and won't accept "none" without a challenge.

Sections that don't apply can be skipped. The result is a complete spec
at `specs/{feature-name}.md`.

### Step 3: Task Out the Work

After the spec is written, the skill routes by size. A small,
single-session spec gets tasks in TASKS.md directly:

```text
Agent: "Spec written to specs/rss-feed.md. Want me to break this
  into tasks in TASKS.md?"

You: "Yes."

Agent: [creates a Phase section referencing the spec, adds tasks]
       "Added Phase 0.8: RSS/Atom Feed Generation with 4 tasks:
         - [ ] Implement feed.xml generation in ctx journal site
         - [ ] Add --feed flag and .ctxrc feed config
         - [ ] Write tests for feed generation
         - [ ] Add recipe and update docs"
```

Each task is independently completable and references the spec via
`Spec: specs/rss-feed.md` in the Phase header.

You can also invoke `/ctx-task-add` directly at any point to add
individual tasks.

A spec that spans multiple milestones (or more than about one
session of implementation) gets `/ctx-task-out` instead:

```text
You: /ctx-task-out --spec specs/v1-substrate.md --milestone m0a

Agent: "Open questions check: D-001 (language choice) is blocking
  for m0a — three tasks would embed an assumption about its
  answer. Resolve it before I decompose."

You: [resolves D-001 in the spec]

Agent: [writes specs/plans/m0a.md: data model, contracts, test
  matrix, 24 tasks with falsifiable acceptance criteria]
       "Plan written. TASKS.md gained 4 epic anchors, each
       annotated Plan: specs/plans/m0a.md."
```

The plan document owns the fine-grained tasks; TASKS.md carries
epic-level anchors only. Two hard gates apply: decomposition
refuses to run past a blocking open question in the spec, and
milestone N+1 is not tasked out while milestone N's definition of
done is unmet.

### Step 4: Implement Step-by-Step

Pick up the work with `/ctx-implement`. The skill loads the plan —
`specs/plans/<milestone>.md` when `/ctx-task-out` produced one, the
spec itself for small features — breaks it into atomic steps, and
checkpoints after every 3-5 steps with a build and test
verification:

```text
You: /ctx-implement (specs/rss-feed.md)

Agent: "Loaded spec. Breaking into steps:
         1. Add FeedEntry type to journal package
         2. Implement feed.xml template
         3. Wire into site generation pipeline
         [builds and tests after step 3]
         4. Add --feed flag to ctx journal site
         5. Add .ctxrc feed configuration
         [builds and tests after step 5]
         ..."
```

If a build or test fails, the agent stops, diagnoses, and fixes
before continuing.

## When to Skip Steps

Not every feature needs all four steps. Use your judgment:

| Situation                              | Start at           |
|----------------------------------------|--------------------|
| Vague idea, multiple valid approaches  | Step 1: Brainstorm |
| Clear approach, need to document it    | Step 2: Spec       |
| Spec already exists, need to plan work | Step 3: Task out   |
| Tasks exist, ready to code             | Step 4: Implement  |

A brainstorm without a spec is fine for small decisions. A spec without
a brainstorm is fine when the design is obvious. The full chain is for
features complex enough to warrant front-loaded design.

## Conversational Approach

You don't need skill names. Natural language works:

| You say                            | What happens         |
|------------------------------------|----------------------|
| "Let's think through this feature" | `/ctx-brainstorm`    |
| "Spec this out"                    | `/ctx-spec`          |
| "Write a design doc for..."        | `/ctx-spec`          |
| "Task this out"                    | `/ctx-task-out`      |
| "Break this into tasks"            | `/ctx-task-add`      |
| "Implement the spec"               | `/ctx-implement`     |
| "Let's design before we build"     | Starts at brainstorm |

## Tips

* **Brainstorm first when uncertain**. If you can articulate the approach in
  two sentences, skip to spec. If you can't, brainstorm.
* **Specs prevent scope creep**. The Non-Goals section is as important as the
  approach. Writing down what you *won't* do keeps implementation focused.
* **Edge cases are the point**. A spec that only describes the happy path
  isn't a spec - it's a wish. The `/ctx-spec` skill pushes for at least 3
  edge cases because that's where designs break.
* **Record decisions during brainstorming**. When you choose between
  approaches, the agent offers to persist the trade-off via
  `/ctx-decision-add`. Accept - future sessions need to know *why*, not
  just *what*.
* **Specs are living documents**. Update them when implementation reveals
  new constraints. A spec that diverges from reality is worse than no spec.
* **The spec template is customizable**. Edit `specs/tpl/spec-template.md`
  to match your project's needs. The `/ctx-spec` skill reads whatever
  template it finds there.

## See Also

* [Skills Reference: /ctx-brainstorm](../reference/skills.md#ctx-brainstorm):
  structured design dialogue
* [Skills Reference: /ctx-spec](../reference/skills.md#ctx-spec):
  spec scaffolding from template
* [Skills Reference: /ctx-task-out](../reference/skills.md#ctx-task-out):
  spec decomposition into a per-milestone plan
* [Skills Reference: /ctx-implement](../reference/skills.md#ctx-implement):
  step-by-step execution with verification
* [Scrutinizing a Plan](scrutinizing-a-plan.md): the adversarial
  interview that belongs between brainstorm and spec
* [Spec-Driven Development](spec-driven-development.md): the full
  operator's manual for the chain — the debated brief, per-milestone
  tasking, and the gates — when a feature spans several milestones
* [Tracking Work Across Sessions](task-management.md): task lifecycle
  and archival
* [Importing Claude Code Plans](import-plans.md): turning ephemeral plans
  into permanent specs
* [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md):
  capturing design trade-offs

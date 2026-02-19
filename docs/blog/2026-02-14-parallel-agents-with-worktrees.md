---
title: "Parallel Agents with Git Worktrees"
date: 2026-02-14
author: Jose Alekhinne
topics:
  - agent teams
  - parallelism
  - git worktrees
  - context engineering
  - task management
---

# Parallel Agents with Git Worktrees

![ctx](../images/ctx-banner.png)

## The Backlog Problem

*Jose Alekhinne / 2026-02-14*

!!! question "What do you do with 30 open tasks?"
    You could work through them one at a time.

    One agent, one branch, one commit stream.

    Or you could ask: **which of these don't touch each other?**

I had 30 open tasks in `TASKS.md`. Some were docs. Some were a new
encryption package. Some were test coverage for a stable module. Some
were blog posts.

They had almost zero file overlap.

Running one agent at a time meant serial execution on work that was
fundamentally parallel:

I was bottlenecking on **me**, not on the machine.

## The Insight: File Overlap Is the Constraint

This is not a scheduling problem. It is a **conflict avoidance** problem.

Two agents can work simultaneously on the same codebase if and only if
they don't touch the same files. The moment they do, you get merge
conflicts: And merge conflicts on AI-generated code are expensive
because the human has to arbitrate choices they didn't make.

So the question becomes: **can you partition your backlog into
non-overlapping tracks?**

For `ctx`, the answer was obvious:

| Track        | Touches                    | Tasks                              |
|--------------|----------------------------|------------------------------------|
| `work/docs`  | `docs/`, `hack/`           | Blog posts, recipes, runbooks      |
| `work/pad`   | `internal/cli/pad/`, specs | Scratchpad encryption, CLI, tests  |
| `work/tests` | `internal/cli/recall/`     | Recall test coverage               |

Three tracks. Near-zero overlap. Three agents.

## Git Worktrees: The Mechanism

Git has a feature that most people don't use: **worktrees**.

A **worktree** is a second (*or third, or fourth*) working directory that
shares the same `.git` object database as your main checkout. Each
worktree has its own branch, its own index, its own working tree. But
they all share history, refs, and objects.

```bash
git worktree add ../ctx-docs -b work/docs
git worktree add ../ctx-pad -b work/pad
git worktree add ../ctx-tests -b work/tests
```

Three directories. Three branches. One repository.

This is **cheaper** than three clones. And because they share objects,
`git merge` afterwards is fast: It's a local operation on shared data.

## The Setup

The workflow I landed on:

**1. Group tasks by blast radius.**

Read `TASKS.md`. For each pending task, estimate which files and
directories it touches. Group tasks that share files into the same
track. Tasks with no overlap go into separate tracks.

This is the part that requires **human judgment**: 

An agent can propose groupings, but you need to verify that the boundaries are 
real. A task that says "*update docs*" but actually touches Go code will poison a
docs track.

**2. Create worktrees as sibling directories.**

Not subdirectories: Siblings. 

If your main checkout is at `~/WORKSPACE/ctx`, worktrees go 
at `~/WORKSPACE/ctx-docs`, `~/WORKSPACE/ctx-pad`, etc.

Why siblings? Because some tools (*and some agents*) walk up the directory
tree looking for `.git`. A worktree inside the main checkout confuses
them.

**3. Launch one agent per worktree.**

```bash
# Terminal 1
cd ../ctx-docs && claude

# Terminal 2
cd ../ctx-pad && claude

# Terminal 3
cd ../ctx-tests && claude
```

Each agent gets a full working copy with `.context/` intact. It reads
the same `TASKS.md`, the same `DECISIONS.md`, the same `CONVENTIONS.md`.
It knows the full project state. It just works on a different slice.

**4. Do NOT run `ctx init` in worktrees.**

This is the gotcha. The `.context/` directory is tracked in git. Running
`ctx init` in a worktree would overwrite shared context files — wiping
decisions, learnings, and tasks that belong to the whole project.

The worktree already has everything it needs. Leave it alone.

## What Actually Happened

I ran three agents for about 40 minutes. Here is roughly what each
track produced:

**`work/docs`**: Parallel worktrees recipe, blog post edits, recipe
index reorganization, IRC recipe moved from `docs/` to `hack/`.

**`work/pad`**: `ctx pad show` subcommand, `--append` and `--prepend`
flags on `ctx pad edit`, spec updates, 28 new test functions.

**`work/tests`**: Recall test coverage, edge case tests.

Merging took about **five** minutes. Two of the three merges were clean.

The third had a conflict in `TASKS.md` — both the docs track and the
pad track had marked different tasks as `[x]`.

## The `TASKS.md` Conflict

This deserves its own section because it will happen **every time**.

When two agents work in parallel, they both read `TASKS.md` at the
start and mark tasks complete as they go. When you merge, git sees two
branches that modified the same file differently.

The resolution is always the same: **accept all completions from both
sides**. No task should go from `[x]` back to `[ ]`. The merge is
additive.

This is one of those conflicts that sounds scary but is trivially
mechanical: You are not arbitrating design decisions; you are combining
two checklists.

## Limits

**3-4 worktrees, maximum.** 

I tried four once. By the time I merged the third track, the fourth had drifted
far enough that its changes needed rebasing. 

The merge complexity grows faster than the parallelism benefit.

Three is the sweet spot:
 
* Two is conservative but safe;
* Four is possible if the tracks are truly independent;
* Anything more than four, you are in the danger zone.

**Group by directory, not by priority.**

It is tempting to put all the high-priority tasks in one track: **Don't**. 

Two high-priority tasks that touch the same files **must** be in the same track, 
regardless of urgency. The constraint is file overlap, not importance.

**Commit frequently.** 

Smaller commits make merge conflicts easier to resolve. An agent that writes 
500 lines in a single commit is harder to merge than one that commits every 
logical step.

**Name tracks by concern.** 

* `work/docs` and `work/pad` tell you what's happening;
* `work/track-1` and `work/track-2` tell you nothing.

## The Pattern

This is the same pattern that shows up everywhere in `ctx`:

[**The attention budget**][attention-post] taught me that you can't dump
everything into one context window. You have to partition, prioritize,
and load selectively.

Worktrees are the same principle applied to **execution**: You can't
dump every task into one agent's workstream. You have to **partition by
blast radius**, assign **selectively**, and merge **deliberately**.

[attention-post]: 2026-02-03-the-attention-budget.md

The [codebase audit][audit-post] that generated these 30 tasks used
eight parallel agents for *analysis*. Worktrees let me use parallel
agents for *implementation*. Same coordination pattern, different
artifact.

[audit-post]: 2026-02-08-not-everything-is-a-skill.md

And the [IRC bouncer post][irc-post] from earlier today argued that
stateless protocols need stateful wrappers. Worktrees are the same:
git branches are stateless forks; `.context/` is the stateful wrapper
that gives each agent the project's full memory.

[irc-post]: 2026-02-14-irc-as-context.md

## Should This Be a Skill?

I asked myself the [same question I asked about the codebase
audit][audit-post]: should this be a `/ctx-worktree` skill?

This time the answer is yes. Unlike the audit prompt — which I tweak
every time and run quarterly — the worktree workflow is:

| Criterion | Worktree workflow     | Codebase audit          |
|-----------|-----------------------|-------------------------|
| Frequency | Weekly                | Quarterly               |
| Stability | Same steps every time | Tweaked every time      |
| Scope     | Mechanical, bounded   | Bespoke, 8 agents       |
| Trigger   | Large backlog         | "I feel like auditing"  |

The commands are mechanical: `git worktree add`, `git worktree remove`,
branch naming, safety checks. This is exactly what skills are for:
**stable contracts** for repetitive operations.

Ergo, `/ctx-worktree` exists. 

It enforces the 4-worktree limit, creates sibling directories, 
uses `work/` branch prefixes, and reminds you not to run `ctx init`
in worktrees.

## The Takeaway

Serial execution is the default. But serial is not always necessary.

If your backlog partitions cleanly by file overlap, you can multiply
your throughput with nothing more exotic than `git worktree` and
a second terminal window.

The hard part is not the `git` commands; it is the **discipline**:

* Grouping by blast radius instead of priority; 
* Accepting that `TASKS.md` will conflict; 
* And knowing when three tracks is enough.

---

!!! quote "If you remember one thing from this post..."
    **Partition by blast radius, not by priority.**

    Two tasks that touch the same files belong in the same track,
    no matter how important the other one is.

    **The constraint is file overlap. Everything else is scheduling.**

---

*The practical setup — skill invocation, worktree creation, merge
workflow, and cleanup — lives in the recipe:
[Parallel Agent Development with Git Worktrees](../recipes/parallel-worktrees.md).*

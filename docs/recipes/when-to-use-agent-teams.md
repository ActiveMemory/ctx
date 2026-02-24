---
title: "When to Use a Team of Agents"
icon: lucide/users
---

![ctx](../images/ctx-banner.png)

## The Problem

You have a task, and you are wondering: "*should I throw more agents at it?*"

More agents *can* mean faster results, but they also mean coordination
overhead, merge conflicts, divergent mental models, and wasted tokens
re-reading context. 

The wrong setup costs more than it saves.

This recipe is a **decision framework**: It helps you choose between a single
agent, parallel worktrees, and a full agent team, and explains what `ctx`
provides at each level.

## TL;DR

**Single agent** for most work. **Parallel worktrees** when tasks touch
disjoint file sets. **Agent teams** only when tasks need real-time
coordination. When in doubt, start with one agent.

## The Spectrum

There are three modes, ordered by complexity:

### 1. Single Agent (*Default*)

One agent, one session, one branch. This is correct for most work.

**Use when**:

* The task has linear dependencies (*step 2 needs step 1's output*)
* Changes touch overlapping files
* You need tight feedback loops (*review each change before the next*)
* The task requires deep understanding of a single area
* Total effort is less than a few hours of agent time

**What `ctx` provides**: Full `.context/`: tasks, decisions, learnings,
conventions, all in one session. 

The agent builds a coherent mental model and persists it as it goes.

**Example tasks**: Bug fixes, feature implementation, refactoring a
module, writing documentation for one area, debugging.

### 2. Parallel Worktrees (*Independent Tracks*)

2-4 agents, each in a separate git worktree on its own branch, working
on non-overlapping parts of the codebase.

**Use when**:

* You have 5+ independent tasks in the backlog
* Tasks group cleanly by directory or package
* File overlap between groups is zero or near-zero
* Each track can be completed and merged independently
* You want parallelism without coordination complexity

**What `ctx` provides**: Shared `.context/` via `git` (*each worktree sees
the same tasks, decisions, conventions*). `/ctx-worktree` skill for
setup and teardown. `TASKS.md` as a lightweight work queue.

**Example tasks**: Docs + new package + test coverage (*three tracks
that don't touch the same files*). Parallel recipe writing. Independent
module development.

**See:** [Parallel Agent Development with Git Worktrees](parallel-worktrees.md)

### 3. Agent Team (*Coordinated Swarm*)

Multiple agents communicating via messages, sharing a task list, with
a lead agent coordinating. Claude Code's team/swarm feature.

**Use when:**

* Tasks have dependencies but can still partially overlap
* You need research and implementation happening simultaneously
* The work requires different roles (*researcher, implementer, tester*)
* A lead agent needs to review and integrate others' work
* The task is large enough that coordination cost is justified

**What `ctx` provides**: `.context/` as shared state that all agents
can read. Task tracking for work assignment. Decisions and learnings
as team memory that survives individual agent turnover.

**Example tasks:** Large refactor across modules where a lead reviews
merges. Research and implementation where one agent explores options
while another builds. Multi-file feature that needs integration testing
after parallel implementation.

## The Decision Framework

Ask these questions in order:

```
Can one agent do this in a reasonable time?
  YES → Single agent. Stop here.
  NO  ↓

Can the work be split into non-overlapping file sets?
  YES → Parallel worktrees (2-4 tracks)
  NO  ↓

Do the subtasks need to communicate during execution?
  YES → Agent team with lead coordination
  NO  → Parallel worktrees with a merge step
```

### The File Overlap Test

This is the critical decision point. Before choosing multi-agent, list the
files each subtask would touch. If two subtasks modify the same file, they
belong in the same track (or the same single-agent session).

```text
You: "I want to parallelize these tasks. Which files would each one touch?"

Agent: [reads TASKS.md, analyzes codebase]
       "Task A touches internal/config/ and internal/cli/initialize/
        Task B touches docs/ and site/
        Task C touches internal/config/ and internal/cli/status/

        Tasks A and C overlap on internal/config/ — they should be
        in the same track. Task B is independent."
```

When in doubt, keep things in one track. A merge conflict in a critical
file costs more time than the parallelism saves.

## When Teams Make Things Worse

More agents is not always better. Watch for these patterns:

**Merge hell**. If you are spending more time resolving conflicts than
the parallel work saved, you split wrong: Re-group by file overlap.

**Context divergence**: Each agent builds its own mental model. After
30 minutes of independent work, agent A might make assumptions that
contradict agent B's approach. Shorter tracks with frequent merges
reduce this.

**Coordination theater**: A lead agent spending most of its time
assigning tasks, checking status, and sending messages instead of
doing work. If the task list is clear enough, worktrees with no
communication are cheaper.

**Re-reading overhead**: Every agent reads `.context/` on startup.
A team of 4 agents each reading 4000 tokens of context = 16000 tokens
before anyone does any work. For small tasks, that overhead dominates.

## What ctx Gives You at Each Level

| ctx Feature         | Single Agent         | Worktrees            | Team                   |
|---------------------|----------------------|----------------------|------------------------|
| `.context/` files   | Full access          | Shared via git       | Shared via filesystem  |
| `TASKS.md`          | Work queue           | Split by track       | Assigned by lead       |
| Decisions/Learnings | Persisted in session | Persisted per branch | Persisted by any agent |
| `/ctx-next`         | Picks next task      | Picks within track   | Lead assigns           |
| `/ctx-worktree`     | N/A                  | Setup + teardown     | Optional               |
| `/ctx-commit`       | Normal commits       | Per-branch commits   | Per-agent commits      |

## Tips

* **Start with one agent**. Only add parallelism when you have identified
  the bottleneck. "This would go faster with more agents" is usually
  wrong for tasks under 2 hours.
* **The 3-4 agent ceiling is real**. Coordination overhead grows
  quadratically. 2 agents = 1 communication pair. 4 agents = 6 pairs.
  Beyond 4, you are managing agents more than doing work.
* **Worktrees > teams for most parallelism needs.** If agents don't
  need to talk to each other during execution, worktrees give you
  parallelism with zero coordination overhead.
* **Use `ctx` as the shared brain**. Whether it's one agent or four, the
  `.context/` directory is the single source of truth. Decisions go in
  `DECISIONS.md`, **not** in chat messages between agents.
* **Merge early, merge often**. Long-lived parallel branches diverge.
  Merge a track as soon as it's done rather than waiting for all tracks
  to finish.
* **`TASKS.md` conflicts are normal**. Multiple agents completing different
  tasks will conflict on merge. The resolution is always additive: accept
  all `[x]` completions from both sides.

## Next Up

**[Parallel Agent Development with Git Worktrees →](parallel-worktrees.md)**:
Run multiple agents on independent task tracks using git worktrees.

## Go Deeper

* [CLI Reference](../reference/cli-reference.md) — all commands and flags
* [Integrations](../operations/integrations.md) — setup for Claude Code, Cursor, Aider
* [Session Journal](../reference/session-journal.md) — browse and search session history

## See Also

* [Parallel Agent Development with Git Worktrees](parallel-worktrees.md):
  the mechanical "how" for worktree-based parallelism
* [Running an Unattended AI Agent](autonomous-loops.md): serial autonomous
  loops — a different scaling strategy
* [Tracking Work Across Sessions](task-management.md): managing the task
  backlog that feeds into any multi-agent setup

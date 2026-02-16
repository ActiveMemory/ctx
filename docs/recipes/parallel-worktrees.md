---
title: Parallel Agent Development with Git Worktrees
icon: lucide/git-branch
---

![ctx](../images/ctx-banner.png)

## Problem

You have a large backlog — 10, 20, 30 open tasks — and many of them are
independent: docs work that doesn't touch Go code, a new package that
doesn't overlap with existing ones, test coverage for a stable module.

Running one agent at a time means serial execution. You want 3-4 agents
working in parallel, each on its own track, without stepping on each
other's files.

Git worktrees solve this. Each worktree is a separate working directory
with its own branch, but they share the same `.git` object database.
Combined with ctx's persistent context, each agent session picks up the
full project state and works independently.

!!! tip "TL;DR"
    ```text
    /ctx-worktree                                     # 1. group tasks by file overlap
    ```
    ```bash
    git worktree add ../myproject-docs -b work/docs   # 2. create worktrees
    cd ../myproject-docs && claude                     # 3. launch agents (one per track)
    ```
    ```text
    /ctx-worktree teardown docs                        # 4. merge back and clean up
    ```

    TASKS.md will conflict on merge — accept all `[x]` completions from both sides.

## Commands and Skills Used

| Tool             | Type    | Purpose                                     |
|------------------|---------|---------------------------------------------|
| `/ctx-worktree`  | Skill   | Create, list, and tear down worktrees       |
| `/ctx-next`      | Skill   | Pick tasks from the backlog for each track  |
| `git worktree`   | Command | Underlying git worktree management          |
| `git merge`      | Command | Merge completed tracks back to main         |

## The Workflow

### Step 1: Assess the Backlog

Start in your main checkout. Ask the agent to analyze your tasks and
group them by blast radius — which files and directories each task
touches.

```text
/ctx-worktree
Look at TASKS.md and group the pending tasks into 2-3 independent
tracks based on which files they'd touch. Show me the grouping
before creating anything.
```

The agent reads TASKS.md, estimates file overlap, and proposes groups:

```text
Proposed worktree groups:

  work/docs    — recipe updates, blog post (touches: docs/)
  work/crypto  — scratchpad encryption infra (touches: internal/crypto/)
  work/tests   — recall test coverage (touches: internal/cli/recall/)
```

### Step 2: Create the Worktrees

Once you approve the grouping, the agent creates worktrees as sibling
directories:

```text
Create the worktrees for those three groups.
```

Behind the scenes:

```bash
git worktree add ../myproject-docs -b work/docs
git worktree add ../myproject-crypto -b work/crypto
git worktree add ../myproject-tests -b work/tests
```

Each worktree is a full working copy on its own branch.

### Step 3: Launch Agents

Open a separate terminal (or editor window) for each worktree and
start a Claude Code session:

```bash
# Terminal 1
cd ../myproject-docs
claude

# Terminal 2
cd ../myproject-crypto
claude

# Terminal 3
cd ../myproject-tests
claude
```

Each agent sees the full project, including `.context/`, and can work
independently. Do **not** run `ctx init` in worktrees — the context
directory is already tracked in git.

### Step 4: Work

Each agent works through its assigned tasks. They can read TASKS.md to
know what's assigned to their track, use `/ctx-next` to pick the next
item, and commit normally on their `work/*` branch.

### Step 5: Merge Back

As each track finishes, return to the main checkout and merge:

```text
/ctx-worktree teardown docs
```

The agent checks for uncommitted changes, merges `work/docs` into your
current branch, removes the worktree, and deletes the branch.

### Step 6: Handle TASKS.md Conflicts

TASKS.md will almost always conflict when merging — multiple agents
marked different tasks as `[x]`. This is expected and easy to resolve:

**Accept all completions from both sides.** No task should go from
`[x]` back to `[ ]`. The merge resolution is always additive.

```text
The merge has a conflict in TASKS.md. Both branches completed
different tasks. Accept all [x] completions from both sides.
```

### Step 7: Cleanup

After all tracks are merged, verify everything is clean:

```text
/ctx-worktree list
```

Should show only the main working tree. All `work/*` branches should
be gone.

## Conversational Approach

You don't have to use the skill directly for every step. These natural
prompts work:

- *"I have a big backlog. Can we split it across worktrees?"*
- *"Which of these tasks can run in parallel without conflicts?"*
- *"Merge the docs track back in."*
- *"Clean up all the worktrees, we're done."*

## Tips

- **3-4 worktrees max.** Beyond that, merge complexity outweighs the
  parallelism benefit. The skill enforces this limit.
- **Group by package or directory**, not by priority. Two high-priority
  tasks that touch the same files must be in the same track.
- **TASKS.md will conflict** on merge. This is normal. Accept all `[x]`
  completions — the resolution is always additive.
- **Don't run `ctx init` in worktrees.** The `.context/` directory is
  tracked in git. Running init overwrites shared context files.
- **Name worktrees by concern**, not by number. `work/docs` and
  `work/crypto` are more useful than `work/track-1` and `work/track-2`.
- **Commit frequently** in each worktree. Smaller commits make merge
  conflicts easier to resolve.

## Next Up

**[When to Use a Team of Agents →](when-to-use-agent-teams.md)**: Decide
when multi-agent coordination is worth the overhead.

## See Also

- [Running an Unattended AI Agent](autonomous-loops.md) — for serial
  autonomous loops instead of parallel tracks
- [Tracking Work Across Sessions](task-management.md) — managing the
  task backlog that feeds into parallelization
- [The Complete Session](session-lifecycle.md) — the session workflow
  each agent follows inside its worktree

---
name: ctx-worktree
description: "Manage git worktrees for parallel agent development. Use when splitting work across independent task tracks."
allowed-tools: Bash(git:*), Bash(ls:*), Read
---

Manage git worktrees to parallelize agent work across independent
task tracks. Supports creating, listing, and tearing down worktrees
with ctx-aware guardrails.

## When to Use

- User wants to parallelize a backlog across multiple agents
- Multiple independent task tracks with non-overlapping files
- User says "create worktree", "let's parallelize", "split the work"
- 3+ independent tasks that can be worked concurrently

## When NOT to Use

- Single task or tightly coupled tasks
- Tasks that touch overlapping files (high merge conflict risk)
- Fewer than 3 independent tasks (overhead exceeds benefit)
- Already inside a worktree (manage from the main checkout only)
- User just wants concurrent Claude Code sessions in the same tree

## Operations

### `create <name>`

Create a new worktree as a sibling directory with a `work/` branch.

**Process:**

1. **Check count** — refuse if 4 worktrees already exist:
   ```bash
   git worktree list
   ```
   Count lines. If >= 5 (1 main + 4 worktrees), stop and explain
   the limit.

2. **Determine project name** from the current directory basename:
   ```bash
   basename "$(git rev-parse --show-toplevel)"
   ```

3. **Create the worktree** as a sibling directory:
   ```bash
   git worktree add "../<project>-<name>" -b "work/<name>"
   ```

4. **Verify** the worktree was created:
   ```bash
   ls "../<project>-<name>"
   ```

5. **Remind the user**:
   > Do NOT run `ctx init` in the worktree. The context
   > directory is already tracked in git and will be present.
   > Launch a separate Claude Code session there and work
   > normally.

### `list`

Show all active worktrees:

```bash
git worktree list
```

### `teardown <name>`

Merge a completed worktree back and clean up.

**Process:**

1. **Check for uncommitted changes** in the worktree:
   ```bash
   git -C "../<project>-<name>" status --porcelain
   ```
   If output is non-empty, warn and stop. The user must commit or
   discard changes first.

2. **Merge the work branch** into the current branch:
   ```bash
   git merge "work/<name>"
   ```
   If there are conflicts, stop and help the user resolve them.
   TASKS.md conflicts are common — see guidance below.

3. **Remove the worktree**:
   ```bash
   git worktree remove "../<project>-<name>"
   ```

4. **Delete the branch**:
   ```bash
   git branch -d "work/<name>"
   ```

5. **Verify** cleanup:
   ```bash
   git worktree list
   git branch | grep "work/<name>"
   ```

## Guardrails

- **Max 4 worktrees** — more than 4 parallel tracks makes merge
  complexity outweigh productivity gains
- **Sibling directories only** — worktrees go in `../<project>-<name>`,
  never inside the project tree
- **`work/` branch prefix** — all worktree branches use `work/<name>`
  for easy identification and cleanup
- **No `ctx init` in worktrees** — the context directory is tracked
  in git; running init would overwrite shared context files
- **Manage from main checkout only** — create and teardown worktrees
  from the main working tree, not from inside a worktree
- **TASKS.md conflict resolution** — when merging, TASKS.md will
  often conflict because multiple agents marked different tasks as
  complete. Resolution: accept all `[x]` completions from both sides.
  No task should go from `[x]` back to `[ ]`.

## What Does NOT Work in Worktrees

The scratchpad key (`.scratchpad.key` in the context directory) is
gitignored. It exists only in the main checkout. This means
key-dependent features silently fail in worktrees:

- **`ctx pad`** — cannot decrypt the scratchpad. Commands fail
  gracefully (no key found), but the pad is inaccessible. Use the
  pad from the main checkout only.
- **`ctx notify`** — cannot decrypt the webhook URL. Notifications
  silently do nothing. Agents in worktrees are opaque: no webhook
  alerts for loop completions, nudges, or custom events.
- **Journal enrichment** — `ctx recall export` and `ctx journal enrich`
  resolve paths relative to the current working directory. Files
  created in a worktree stay in that worktree and are discarded on
  teardown. Enrich journals on the main branch after merging — the
  JSONL session logs are intact regardless.

## Task Grouping Guidance

Before creating worktrees, analyze the backlog to group tasks into
non-overlapping tracks:

1. **Read TASKS.md** and identify all pending tasks
2. **Estimate blast radius** — which files/directories does each
   task touch?
3. **Group by non-overlapping directories** — tasks that touch the
   same package or file must go in the same track
4. **Present the grouping** to the user before creating worktrees:

```text
Proposed worktree groups:

  work/docs     — recipe updates, blog post, getting started guide
                  (touches: docs/)
  work/crypto   — P3.1-P3.3 encrypted scratchpad infra
                  (touches: internal/crypto/, internal/config/)
  work/pad-cli  — P3.4-P3.9 pad CLI commands
                  (touches: internal/cli/pad/)
```

Let the user approve or adjust before proceeding.

## Quality Checklist

Before any operation, verify:
- [ ] Worktree count checked (max 4)
- [ ] Branch uses `work/` prefix
- [ ] Worktree is a sibling directory (`../`)
- [ ] User reminded not to run `ctx init` in worktree
- [ ] Uncommitted changes checked before teardown
- [ ] Merge completed before worktree removal
- [ ] Branch deleted after worktree removal

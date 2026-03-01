---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "Detecting and Fixing Drift"
icon: lucide/stethoscope
---

![ctx](../images/ctx-banner.png)

## The Problem

Context files drift: you rename a package, delete a module, or finish a sprint,
and suddenly `ARCHITECTURE.md` references paths that no longer exist,
`TASKS.md` is 80 percent completed checkboxes, and `CONVENTIONS.md` describes
patterns you stopped using two months ago.

**Stale context is worse than no context**: 

An AI tool that trusts outdated references will **hallucinate confidently**.

This recipe shows how to detect drift, fix it, and keep your `.context/`
directory lean and accurate.

## TL;DR

```bash
ctx drift                      # detect problems
ctx drift --fix                # auto-fix the easy ones
ctx sync --dry-run && ctx sync # reconcile after refactors
ctx compact --archive          # archive old completed tasks
ctx status                     # verify
```

Or just ask your agent: *"Is our context clean?"*

## Commands and Skills Used

| Tool                   | Type    | Purpose                                        |
|------------------------|---------|------------------------------------------------|
| `ctx drift`            | Command | Detect stale paths, missing files, violations  |
| `ctx drift --fix`      | Command | Auto-fix simple issues                         |
| `ctx sync`             | Command | Reconcile context with codebase structure      |
| `ctx compact`          | Command | Archive completed tasks, clean up empty sections |
| `ctx status`           | Command | Quick health overview                          |
| `/ctx-drift`           | Skill   | Structural plus semantic drift detection       |
| `/ctx-map`             | Skill   | Refresh `ARCHITECTURE.md` from actual codebase |
| `/ctx-alignment-audit` | Skill   | Audit doc claims against agent instructions    |
| `/ctx-status`          | Skill   | In-session context summary                     |
| `/ctx-prompt-audit`    | Skill   | Audit prompt quality and token efficiency      |

## The Workflow

The best way to maintain context health is **conversational**: Ask your agent,
guide it, and let it detect problems, explain them, and fix them with your
approval. CLI commands exist for CI pipelines, scripting, and fine-grained
control. 

For day-to-day maintenance, **talk to your agent**.

!!! tip "Your Questions Reinforce the Pattern"
    Asking "is our context clean?" does two things:

    * It triggers a drift check right now
    * It reinforces the habit

    This is **reinforcement**, *not* enforcement.

    **Do not wait** for the agent to be proactive on its own: 

    **Guide your agent**, especially in early sessions.

    Over time, you will ask less and the agent will start offering more.

### Step 1: Ask Your Agent

The simplest way to check context health:

```text
Is our context clean?
Anything stale?
How healthy are our context files?
```

Or invoke the skill directly:

```text
/ctx-drift
```

The agent performs two layers of analysis:

**Layer 1**, structural checks (*via `ctx drift`*): Dead paths, missing files,
completed task counts, constitution violations. Fast and programmatic.

**Layer 2**, semantic analysis (*agent-driven*): Does `CONVENTIONS.md` describe
patterns the code no longer follows? Does `DECISIONS.md` contain entries whose
rationale no longer applies? Are there learnings about bugs that are now fixed?
This is where the agent adds value the CLI cannot: It reads both context files
and source code and compares them.

The agent reports both layers together, explains each finding in plain language,
and offers to fix what it can.

### Step 2: Maintenance at Session Start

You do not need to ask explicitly. 

!!! tip "Using Claude Code"
    `ctx` ships with Claude Code hooks that remind the agent 
    *at the right time* to take initiative. 

Checking context health at the session start, offering to persist learnings 
before you quit, and flagging drift when it matters. The agent stays 
proactive without you having to prompt it:

```text
Agent: Good morning. I've loaded the context files. A few things
       before we start:

       - ARCHITECTURE.md references `pkg/auth/` which is now empty
       - DECISIONS.md hasn't been updated in 40 days
       - There are 18 completed tasks ready for archival

       Want me to run a quick maintenance pass, or should we jump
       straight into today's work?
```

☝️️ this is what persistent, initiative-driven sessions feel like when context is 
treated as a system instead of a prompt.

If the agent does not offer this on its own, a gentle nudge is enough:

```text
Anything stale before we start?
How's the context looking?
```

This turns maintenance from a scheduled chore into a conversation that happens
**when** it matters.

### Step 3: Real-Time Detection During Work

Agents *can* notice drift while working: When a mismatch is directly in the
path of their current task. If an agent reads `ARCHITECTURE.md` to find where
to add a handler and `internal/handlers/` doesn't exist, it will notice because
the stale reference blocks its work:

```text
Agent: ARCHITECTURE.md references `internal/handlers/` but that directory
       doesn't exist. I'll look at the actual source tree to find where
       handlers live now.
```

This happens reliably when the drift intersects the task. What is *less*
reliable is the agent generalizing from one mismatch to "*there might be more
stale references; let me run drift detection*" That leap requires the agent
to know `/ctx-drift` exists and to decide the current task should pause for
maintenance.

If you want that behavior, **reinforce** it:

```text
Good catch. Yes, run /ctx-drift and clean up any other stale references.
```

Over time, agents that have seen this pattern will start offering proactively.
But do not expect it from a cold start.

### Step 4: Archival and Cleanup

`ctx drift` detects when `TASKS.md` has more than 10 completed items and
flags it as a staleness warning. Running `ctx drift --fix` archives completed
tasks automatically. 

You can also run `/ctx-archive` to compact on demand.

### Step 5: Alignment Audits

A related problem is **alignment drift**: Documentation that makes claims about
agent behavior not backed by actual playbook or skill instructions. 

Over time, docs accumulate aspirational statements that no instruction teaches 
the agent to do.

Use `/ctx-alignment-audit` to trace behavioral claims in documentation against
the playbook and skill files. The skill identifies gaps, proposes fixes, and
checks instruction file health (*token budgets, bloat signals*).

To avoid confusion with `/ctx-prompt-audit`:

* `/ctx-alignment-audit` checks whether documentation claims are supported by
  actual instructions (*playbook, skills, `CLAUDE.md`*).
* `/ctx-prompt-audit` checks whether your context files are clear, compact, and
  token-efficient for the model.

---

## `ctx doctor`: The Superset Check

`ctx doctor` combines drift detection with hook auditing, configuration checks,
event logging status, and token size reporting in a single command. If you want
one command that covers structural health, hooks, and state:

```bash
ctx doctor          # everything in one pass
ctx doctor --json   # machine-readable for scripting
```

For agent-driven diagnosis that adds semantic analysis on top of the structural
checks, use `/ctx-doctor`. See the [Troubleshooting](troubleshooting.md) recipe
for the full workflow.

---

## CLI Reference

The conversational approach above uses CLI commands under the hood. When you
need direct control, use the commands directly.

### `ctx drift`

Scan context files for structural problems:

```bash
ctx drift
```

Sample output:

```text
Drift Report
============

Warnings (3):
  ARCHITECTURE.md:14  path "internal/api/router.go" does not exist
  ARCHITECTURE.md:28  path "pkg/auth/" directory is empty
  CONVENTIONS.md:9    path "internal/handlers/" not found

Violations (1):
  TASKS.md            31 completed tasks (recommend archival)

Staleness:
  DECISIONS.md        last modified 45 days ago
  LEARNINGS.md        last modified 32 days ago

Exit code: 1 (warnings found)
```

| Level     | Meaning                                             | Action         |
|-----------|-----------------------------------------------------|----------------|
| Warning   | Stale path references, missing files                | Fix or remove  |
| Violation | Constitution rule heuristic failures, heavy clutter | Fix soon       |
| Staleness | Files not updated recently                          | Review content |

Exit codes: `0` equals clean, `1` equals warnings, `3` equals violations.

For CI integration:

```bash
ctx drift --json | jq '.warnings | length'
```

### `ctx drift --fix`

Auto-fix mechanical issues:

```bash
ctx drift --fix
```

This handles removing dead path references, updating unambiguous renames, 
clearing empty sections. Issues requiring judgment are flagged but left for you.

Run `ctx drift` again afterward to confirm what remains.

### `ctx sync`

After a refactor, reconcile context with the actual codebase structure:

```bash
ctx sync --dry-run   # preview first
ctx sync             # apply
```

`ctx sync` scans for **structural changes**, compares with `ARCHITECTURE.md`, 
checks for new dependencies worth documenting, and identifies context referring 
to code that no longer exists.

### `ctx compact`

Consolidate completed tasks and clean up empty sections:

```bash
ctx compact            # move completed tasks to Completed section,
                       # remove empty sections
ctx compact --archive  # also archive old tasks to .context/archive/
```

* Tasks: moves completed items (with all subtasks done) into the Completed
  section of `TASKS.md`
* All files: removes empty sections left behind
* With `--archive`: writes tasks older than 7 days to
  `.context/archive/tasks-YYYY-MM-DD.md`

Without `--archive`, nothing is deleted: Tasks are reorganized in place.

### `ctx status`

Quick health overview:

```bash
ctx status --verbose
```

Shows file counts, token estimates, modification times, and drift warnings in a
single glance.

### `/ctx-alignment-audit` and `/ctx-prompt-audit`

These are both audits, but they answer different questions:

* `/ctx-alignment-audit`: are our behavioral claims backed by actual
  instructions?
* `/ctx-prompt-audit`: are our context files readable, compact, and efficient?

Run them inside your AI assistant:

```text
/ctx-alignment-audit
/ctx-prompt-audit
```

## Putting It All Together

Conversational approach (recommended):

```text
Is our context clean?  -> agent runs structural plus semantic checks
Fix what you can       -> agent auto-fixes and proposes edits
Archive the done tasks -> agent runs ctx compact --archive
How's token usage?     -> agent checks ctx status
```

CLI approach (*for CI, scripts, or direct control*):

```bash
ctx drift                      # 1. Detect problems
ctx drift --fix                # 2. Auto-fix the easy ones
ctx sync --dry-run && ctx sync # 3. Reconcile after refactors
ctx compact --archive          # 4. Archive old completed tasks
ctx status                     # 5. Verify
```

## Tips

Agents cross-reference context files with source code during normal work. When
drift intersects their current task, they will notice: a renamed package, a
deleted directory, a path that doesn't resolve. But they rarely generalize from
one mismatch to a full audit on their own. Reinforce the pattern: when an agent
mentions a stale reference, ask it to run `/ctx-drift`. Over time, it starts
offering.

When an agent says "*this reference looks stale,*" it is **usually right**.

**Semantic drift is more damaging than structural drift**: 
`ctx drift` catches dead paths. But `CONVENTIONS.md` describing a pattern your 
code stopped following three weeks ago is worse. When you ask 
"*is our context clean?*", the agent can do both checks.

**Use `ctx status` as a quick check**: It shows file counts, token estimates, and
drift warnings in a single glance. Good for a fast "is everything ok?" before
diving into work.

**Drift detection in CI**: add `ctx drift --json` to your CI pipeline and fail on
exit code 3 (*violations*). This catches constitution-level problems before they
reach upstream.

**Do not over-compact**: Completed tasks have historical value. The `--archive`
flag preserves them in `.context/archive/` so you can search past work without
cluttering active context.

**Sync is cautious by default**: Use `--dry-run` after large refactors, then apply.

## Next Up

**[Claude Code Permission Hygiene →](claude-code-permissions.md)**:
Recommended permission defaults and maintenance workflow for Claude Code.

## See Also

* [Troubleshooting](troubleshooting.md): full diagnostic workflow using
  `ctx doctor`, event logs, and `/ctx-doctor`
* [Tracking Work Across Sessions](task-management.md): task lifecycle and archival
* [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md):
  keeping knowledge files current
* [The Complete Session](session-lifecycle.md): where maintenance fits in the daily workflow
* [CLI Reference](../cli/index.md): full flag documentation for all commands
* [Context Files](../home/context-files.md): structure and purpose of each `.context/` file

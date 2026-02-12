---
#   /    Context:                     https://ctx.ist
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

**Stale context is worse than no context**: an AI tool that trusts outdated
references will **hallucinate confidently**.

This recipe shows how to detect drift, fix it, and keep your `.context/`
directory lean and accurate.

## Commands and Skills Used

| Tool                   | Type    | Purpose                                        |
|------------------------|---------|------------------------------------------------|
| `ctx drift`            | Command | Detect stale paths, missing files, violations  |
| `ctx drift --fix`      | Command | Auto-fix simple issues                         |
| `ctx sync`             | Command | Reconcile context with codebase structure      |
| `ctx compact`          | Command | Archive completed tasks, deduplicate learnings |
| `ctx status`           | Command | Quick health overview                          |
| `/ctx-drift`           | Skill   | Structural plus semantic drift detection       |
| `/ctx-alignment-audit` | Skill   | Audit doc claims against agent instructions    |
| `/ctx-status`          | Skill   | In-session context summary                     |
| `/ctx-prompt-audit`    | Skill   | Audit prompt quality and token efficiency      |

## The Workflow

The best way to maintain context health is **conversational**: ask your agent,
guide it, and let it detect problems, explain them, and fix them with your
approval. CLI commands exist for CI pipelines, scripting, and fine-grained
control. 

For day-to-day maintenance, **talk to your agent**.

!!! tip "Your Questions Reinforce the Pattern"
    Asking "is our context clean?" does two things:

    * it triggers a drift check right now
    * it reinforces the habit

    This is reinforcement, not enforcement.

    Do not wait for the agent to be proactive on its own. Guide it, especially
    in early sessions. Over time, you will ask less and the agent will start
    offering more.

### Step 1: Ask Your Agent

The simplest way to check context health:

```text
Is our context clean?
Anything stale?
How healthy are our context files?
````

Or invoke the skill directly:

```text
/ctx-drift
```

The agent performs two layers of analysis:

Layer 1, structural checks (via `ctx drift`): dead paths, missing files,
completed task counts, constitution violations. Fast and programmatic.

Layer 2, semantic analysis (agent-driven): does `CONVENTIONS.md` describe
patterns the code no longer follows? Does `DECISIONS.md` contain entries whose
rationale no longer applies? Are there learnings about bugs that are now fixed?
This is where the agent adds value the CLI cannot: it reads both context files
and source code and compares them.

The agent reports both layers together, explains each finding in plain language,
and offers to fix what it can.

### Step 2: Maintenance at Session Start

You do not need to ask explicitly. A proactive agent can check context health
as part of its startup routine and mention anything worth addressing:

```text
Agent: Good morning. I've loaded the context files. A few things
       before we start:

       - ARCHITECTURE.md references `pkg/auth/` which is now empty
       - DECISIONS.md hasn't been updated in 40 days
       - There are 18 completed tasks ready for archival

       Want me to run a quick maintenance pass, or should we jump
       straight into today's work?
```

If the agent does not offer this on its own, a gentle nudge is enough:

```text
Anything stale before we start?
How's the context looking?
```

This turns maintenance from a scheduled chore into a conversation that happens
**when** it matters.

### Step 3: Real-Time Detection During Work

Agents notice drift while working, not just during explicit checks. When an
agent reads `ARCHITECTURE.md` to understand a module and then opens the actual
source tree, it can see mismatches immediately:

```text
Agent: I notice ARCHITECTURE.md still references `internal/handlers/`
       which we deleted in the last refactor. Want me to run drift
       detection so we can clean up any other stale references?
```

The agent is already reading context files to do its job. Noticing staleness is
a natural byproduct of that.

### Step 4: Archival and Cleanup

Agents spot clutter too. When `TASKS.md` is loaded and most of the content is
completed checkboxes, the agent can offer to tidy up:

```text
Agent: TASKS.md has 25 completed tasks obscuring the 4 pending ones.
       Want me to archive the finished work so we can focus on what's
       active?
```

The same applies to `LEARNINGS.md` with near-duplicate entries, or `DECISIONS.md`
with entries that were superseded.

### Step 5: Alignment Audits

A related problem is **alignment drift**: documentation that makes claims about
agent behavior not backed by actual playbook or skill instructions. 

Over time, docs accumulate aspirational statements that no instruction teaches 
the agent to do.

Use `/ctx-alignment-audit` to trace behavioral claims in documentation against
the playbook and skill files. The skill identifies gaps, proposes fixes, and
checks instruction file health (token budgets, bloat signals).

To avoid confusion with `/ctx-prompt-audit`:

* `/ctx-alignment-audit` checks whether documentation claims are supported by
  actual instructions (*playbook, skills, `CLAUDE.md`*).
* `/ctx-prompt-audit` checks whether your context files are clear, compact, and
  token-efficient for the model.

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

This handles removing dead path references, updating unambiguous renames, clearing
empty sections. Issues requiring judgment are flagged but left for you.

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

Archive completed tasks and deduplicate learnings:

```bash
ctx compact --archive
```

* Tasks: moves completed tasks older than 7 days to
  `.context/archive/tasks-YYYY-MM-DD.md`
* Learnings: deduplicates entries with similar content
* All files: removes empty sections left behind

The `--archive` flag preserves old content.

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
Is our context clean?   -> agent runs structural plus semantic checks
Fix what you can        -> agent auto-fixes and proposes edits
Archive the done tasks  -> agent runs ctx compact --archive
How's token usage?      -> agent checks ctx status
```

CLI approach (for CI, scripts, or direct control):

```bash
ctx drift                      # 1. Detect problems
ctx drift --fix                # 2. Auto-fix the easy ones
ctx sync --dry-run && ctx sync # 3. Reconcile after refactors
ctx compact --archive          # 4. Archive old completed tasks
ctx status                     # 5. Verify
```

## Tips

Your agent is your first line of defense. It cross-references context files with
source code during normal work. It will often notice a renamed package, a
deleted directory, or an outdated convention before `ctx drift` runs. 

When an agent says "*this reference looks stale,*" it is **usually right**.

Semantic drift is more damaging than structural drift. `ctx drift` catches dead
paths. But `CONVENTIONS.md` describing a pattern your code stopped following
three weeks ago is worse. When you ask "*is our context clean?*", the agent can 
do both checks.

Use `ctx status` as a quick check. It shows file counts, token estimates, and
drift warnings in a single glance. Good for a fast "is everything ok?" before
diving into work.

Drift detection in CI: add `ctx drift --json` to your CI pipeline and fail on
exit code 3 (violations). This catches constitution-level problems before they
reach upstream.

Do not over-compact. Completed tasks have historical value. The `--archive`
flag preserves them in `.context/archive/` so you can search past work without
cluttering active context.

Sync is cautious by default. Use `--dry-run` after large refactors, then apply.

## Next Up

**[Browsing and Enriching Past Sessions](session-archaeology.md)**:
Export session history to a browsable journal and enrich entries with metadata.

## See Also

* [Tracking Work Across Sessions](task-management.md): task lifecycle and archival
* [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md): 
  keeping knowledge files current
* [The Complete Session](session-lifecycle.md): where maintenance fits in the daily workflow
* [CLI Reference](../cli-reference.md): full flag documentation for all commands
* [Context Files](../context-files.md): structure and purpose of each `.context/` file

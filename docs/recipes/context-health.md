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

## Problem

Context files drift. You rename a package, delete a module, or finish a sprint
-- and suddenly ARCHITECTURE.md references paths that no longer exist, TASKS.md
is 80% completed checkboxes, and CONVENTIONS.md describes patterns you stopped
using two months ago.

Stale context is worse than no context: an AI tool that trusts outdated
references will hallucinate confidently. This recipe shows how to detect drift,
fix it, and keep your `.context/` directory lean and accurate.

## Commands and Skills Used

| Tool                | Type    | Purpose                                        |
|---------------------|---------|------------------------------------------------|
| `ctx drift`         | Command | Detect stale paths, missing files, violations  |
| `ctx drift --fix`   | Command | Auto-fix simple issues                         |
| `ctx sync`          | Command | Reconcile context with codebase structure      |
| `ctx compact`       | Command | Archive completed tasks, deduplicate learnings |
| `ctx status`        | Command | Quick health overview                          |
| `/ctx-status`       | Skill   | In-session context summary                     |
| `/ctx-prompt-audit` | Skill   | Audit prompt quality and token efficiency      |

## The Workflow

The maintenance cycle follows a consistent progression: detect problems, fix
what can be automated, reconcile with the real codebase, then slim down
accumulated clutter.

```
drift detection --> auto-fix --> sync with codebase --> compact --> audit
```

### Step 1: Run Drift Detection

Start with `ctx drift` to scan every context file for problems.

```bash
ctx drift
```

Sample output:

```
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

The report has three severity levels:

| Level         | Meaning                                             | Action         |
|---------------|-----------------------------------------------------|----------------|
| **Warning**   | Stale path references, missing files                | Fix or remove  |
| **Violation** | Constitution rule heuristic failures, heavy clutter | Fix soon       |
| **Staleness** | Files not updated recently                          | Review content |

For CI integration or scripting, use `--json`:

```bash
ctx drift --json | jq '.warnings | length'
```

Exit codes tell you the severity at a glance: `0` means all checks passed,
`1` means warnings were found, and `3` means violations were detected.

### Step 2: Auto-Fix Simple Issues

Many drift warnings have mechanical fixes. Passing `--fix` handles the
straightforward ones:

```bash
ctx drift --fix
```

Auto-fixable issues include:

- Removing references to paths that no longer exist
- Updating directory paths after renames (when the new path is unambiguous)
- Clearing empty sections left behind after manual edits

Issues that require judgment -- like deciding whether a referenced module was
deleted intentionally or moved -- are flagged but left for you to resolve.

After the auto-fix, run `ctx drift` again to confirm the remaining issues
need manual attention.

### Step 3: Sync with the Codebase

After a refactor (renamed packages, moved files, restructured directories),
the context files may describe an architecture that no longer matches reality.
`ctx sync` compares your context with the actual codebase structure.

Always preview first:

```bash
ctx sync --dry-run
```

Sample output:

```
Sync Preview
============

ARCHITECTURE.md:
  + internal/cache/       (new directory, not documented)
  ~ internal/api/         (structure changed)
  - internal/handlers/    (referenced but missing)

Suggestions:
  * Document internal/cache/ in ARCHITECTURE.md
  * Update internal/api/ component description
  * Remove internal/handlers/ references

No changes applied (dry-run mode).
```

If the preview looks correct, apply:

```bash
ctx sync
```

The sync scans for structural changes, compares them with ARCHITECTURE.md,
checks whether package files (go.mod, package.json, etc.) have new dependencies
worth documenting, and identifies context that refers to code that no longer
exists.

### Step 4: Compact Bloated Files

Over time, TASKS.md accumulates completed checkboxes and LEARNINGS.md
collects near-duplicate entries. `ctx compact` cleans this up.

```bash
ctx compact --archive
```

What compact does:

- **Tasks**: Moves completed tasks older than 7 days to
  `.context/archive/tasks-YYYY-MM-DD.md`
- **Learnings**: Deduplicates entries with similar content
- **All files**: Removes empty sections left behind

The `--archive` flag creates the `.context/archive/` directory and preserves
old content there rather than deleting it. You can review archived tasks
later or search them with grep.

If you want to skip the automatic session save that compact performs before
modifying files:

```bash
ctx compact --archive --no-auto-save
```

### Step 5: Audit Prompt Quality

The final step is qualitative. Inside your AI coding assistant, run:

```
/ctx-prompt-audit
```

This skill analyzes your context files for:

- Token efficiency (are you spending tokens on low-value content?)
- Clarity (can an AI tool parse the structure unambiguously?)
- Completeness (are there gaps between what TASKS.md says and what the
  codebase shows?)
- Staleness (decisions marked "Accepted" that should be "Superseded")

The audit produces specific recommendations. Apply them manually, then run
`ctx status` to confirm the overall health:

```bash
ctx status --verbose
```

## Putting It Together

Here is a complete maintenance session combining all five steps:

```bash
# 1. Detect problems
ctx drift

# 2. Auto-fix the easy ones
ctx drift --fix

# 3. Preview sync changes after a refactor
ctx sync --dry-run

# 4. Apply sync
ctx sync

# 5. Archive old completed tasks
ctx compact --archive

# 6. Verify
ctx status
```

Then inside your AI assistant:

```
/ctx-prompt-audit
```

Review the recommendations, make any final edits, and run `ctx drift` one
more time to confirm a clean bill of health.

## Tips

**When to run each command:**

| Command             | When                                                   |
|---------------------|--------------------------------------------------------|
| `ctx drift`         | Regularly -- weekly, or before starting a new feature  |
| `ctx sync`          | After refactors, renames, or major structural changes  |
| `ctx compact`       | When TASKS.md feels cluttered or token budget is tight |
| `/ctx-prompt-audit` | Monthly, or when AI responses seem confused            |

**Use `ctx status` as a quick check.** It shows file counts, token estimates,
and drift warnings in a single glance. Good for a fast "is everything okay?"
before diving into a session.

**Drift detection in CI.** Add `ctx drift --json` to your CI pipeline and
fail on exit code 3 (violations). This catches constitution-level problems
before they reach main.

**Don't over-compact.** Completed tasks have historical value. The `--archive`
flag preserves them in `.context/archive/` so you can search past work without
cluttering active context. Only use compact without `--archive` if you truly
want to discard old items.

**Sync is non-destructive by default.** It suggests changes but never rewrites
files without your confirmation. The `--dry-run` flag exists for extra caution
after large refactors.

## See Also

- [Tracking Work Across Sessions](task-management.md) -- task lifecycle and archival
- [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md) -- keeping knowledge files current
- [The Complete Session](session-lifecycle.md) -- where maintenance fits in the daily workflow
- [CLI Reference](../cli-reference.md) -- full flag documentation for all commands
- [Context Files](../context-files.md) -- structure and purpose of each `.context/` file

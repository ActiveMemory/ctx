---
name: ctx-status
description: "Show context summary. Use at session start or when unclear about current project state."
allowed-tools: Bash(ctx:*)
---

Show the current context status: files, token budget, tasks,
and recent activity.

## When to Use

- At session start to orient before doing work
- When confused about what is being worked on or what context
  exists
- To check token usage and context health
- When the user asks "what's the state of the project?"

## When NOT to Use

- When you already loaded context via `/ctx-agent` in this
  session (status is a subset of what agent provides)
- Repeatedly within the same session without changes in between

## Usage Examples

```text
/ctx-status
/ctx-status --verbose
/ctx-status --json
```

## Flags

| Flag        | Short | Default | Purpose                          |
|-------------|-------|---------|----------------------------------|
| `--json`    |       | false   | Output as JSON (for scripting)   |
| `--verbose` | `-v`  | false   | Include file content previews    |

## What It Shows

The output has three sections:

### 1. Overview

- Context directory path
- Total file count
- Token estimate (sum across all `.context/*.md` files)

### 2. Files

Each `.md` file in `.context/` with:

| Indicator | Meaning                                 |
|-----------|-----------------------------------------|
| check     | File has content (loaded)               |
| circle    | File exists but is empty                |

File-specific summaries:
- `CONSTITUTION.md`: number of invariants
- `TASKS.md`: active and completed task counts
- `DECISIONS.md`: number of decisions
- `GLOSSARY.md`: number of terms
- Others: "loaded" or "empty"

With `--verbose`: adds token count, byte size, and a 3-line
content preview per file.

### 3. Recent Activity

The 3 most recently modified files with relative timestamps
(e.g., "5 minutes ago", "2 hours ago").

## Execution

```bash
ctx status
```

After running, summarize the key points for the user:
- How many active tasks remain
- Whether any context files are empty (might need populating)
- Token budget usage (is context lean or bloated?)
- What was recently modified (gives a sense of momentum)

## Interpreting Results

| Observation             | Suggestion                                                  |
|-------------------------|-------------------------------------------------------------|
| Many empty files        | Context is sparse; populate core files (TASKS, CONVENTIONS) |
| High token count (>30k) | Consider `ctx compact` or archiving completed tasks         |
| No recent activity      | Context may be stale; check if files need updating          |
| TASKS.md has 0 active   | All work done, or tasks need to be added                    |

## Quality Checklist

After running status, verify:
- [ ] Summarized the output for the user (do not just dump
      raw output without commentary)
- [ ] Flagged any empty core files that should be populated
- [ ] Noted token budget if it seems high or low

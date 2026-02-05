---
name: ctx-add-learning
description: "Record a learning. Use when discovering gotchas, bugs, or unexpected behavior that future sessions should know about."
allowed-tools: Bash(ctx:*)
---

Record a learning in LEARNINGS.md.

## Before Recording

Three questions — if any answer is "no", don't record:

1. **"Could someone Google this in 5 minutes?"** → If yes, skip it
2. **"Is this specific to this codebase?"** → If no, skip it
3. **"Did it take real effort to discover?"** → If no, skip it

Learnings should capture **principles and heuristics**, not code snippets.

## When to Use

- After discovering a gotcha or unexpected behavior
- When a debugging session reveals root cause
- When finding a pattern that will help future work

## Execution

```bash
ctx add learning "Learning text" --context "..." --lesson "..." --application "..."
```

Or with just the learning text (will prompt for details):

```bash
ctx add learning $ARGUMENTS
```

Confirm the learning was added.

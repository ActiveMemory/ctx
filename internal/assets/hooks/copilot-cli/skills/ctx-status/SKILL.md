---
name: ctx-status
description: "Show context summary and health. Use at session start or when unclear about project state."
---

Show the current context status: files, token budget, tasks, and
recent activity.

## When to Use

- At session start to orient before doing work
- When confused about what is being worked on
- To check token usage and context health
- When the user asks "what's the state of the project?"

## When NOT to Use

- When you already loaded context via `ctx agent` in this session
- Repeatedly within the same session without changes in between

## Execution

```bash
ctx status
```

For verbose output with file previews:

```bash
ctx status --verbose
```

For machine-readable output:

```bash
ctx status --json
```

After running, summarize the key points:
- How many active tasks remain
- Whether any context files are empty or stale
- What was most recently modified

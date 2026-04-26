---
name: ctx-status
description: "Show context summary. Use at session start or when unclear about current project state."
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

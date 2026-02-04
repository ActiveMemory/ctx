---
name: ctx-recall
description: "Browse session history. Use when referencing past discussions or finding context from previous work."
allowed-tools: Bash(ctx:*)
---

List recent AI sessions from conversation history.

## When to Use

- When needing to reference a previous discussion
- When looking for context from past work
- When the user asks "what did we do last time?"

## Execution

```bash
ctx recall list --limit 10
```

Show the user's recent sessions with project, time, and turn count.

If the user asks about a specific session, use:

```bash
ctx recall show <id>
```

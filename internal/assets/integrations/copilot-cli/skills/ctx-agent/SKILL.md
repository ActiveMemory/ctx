---
name: ctx-agent
description: "Load full context packet. Use at session start or when context seems stale or incomplete."
tools: [bash]
---

Load the full context packet for AI consumption.

## When to Use

- At the start of a session to load all context
- When context seems stale or incomplete
- When switching between different areas of work

## When NOT to Use

- The session start hook already runs `ctx agent` automatically:
  you rarely need to invoke this manually
- Don't run it just to "refresh" if you already have context loaded

## After Loading

**Read the files listed in "Read These Files (in order)"**: the
packet is a summary, not a substitute. In particular, read
CONVENTIONS.md before writing any code.

Confirm to the user: "I have read the required context files and
I'm following project conventions."

## Flags

| Flag         | Default | Description                                   |
|--------------|---------|-----------------------------------------------|
| `--budget`   | 8000    | Token budget for context packet               |
| `--format`   | md      | Output format: `md` or `json`                 |
| `--cooldown` | 10m     | Suppress repeated output within this duration |
| `--session`  | (none)  | Session ID for cooldown isolation             |

## Execution

```bash
ctx agent --budget 4000
```

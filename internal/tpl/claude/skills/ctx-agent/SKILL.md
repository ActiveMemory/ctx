---
name: ctx-agent
description: "Load full context packet. Use at session start or when context seems stale or incomplete."
allowed-tools: Bash(ctx:*)
---

Load the full context packet for AI consumption.

## When to Use

- At the start of a session to load all context
- When context seems stale or incomplete
- When switching between different areas of work

## Usage

```
/ctx-agent
/ctx-agent --budget 4000
/ctx-agent --budget 8000
```

Use `--budget` to limit token count.

## Execution

```bash
ctx agent $ARGUMENTS
```

This provides the complete context state optimized for AI assistants.

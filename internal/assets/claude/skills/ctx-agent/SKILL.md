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

## When NOT to Use

- The PreToolUse hook already runs `ctx agent` automatically with a cooldown
  — you rarely need to invoke this manually
- Don't run it just to "refresh" if you already have the context loaded in
  this session

## After Loading

**Read the files listed in "Read These Files (in order)"** — the packet is a
summary, not a substitute. In particular, read CONVENTIONS.md before writing
any code.

Confirm to the user: "I have read the required context files and I'm
following project conventions." Do not begin implementation until you have
done so.

## Flags

| Flag         | Default | Description                                       |
|--------------|---------|---------------------------------------------------|
| `--budget`   | 8000    | Token budget for context packet                   |
| `--format`   | md      | Output format: `md` or `json`                     |
| `--cooldown` | 10m     | Suppress repeated output within this duration     |
| `--session`  | (none)  | Session ID for cooldown isolation (e.g., `$PPID`) |

## Execution

```bash
ctx agent $ARGUMENTS
```

**Example — default load:**
```bash
ctx agent
```

**Example — smaller packet for limited contexts:**
```bash
ctx agent --budget 4000
```

**Example — with cooldown (how the PreToolUse hook invokes it):**
```bash
ctx agent --budget 4000 --session $PPID
```

**Example — JSON for programmatic use:**
```bash
ctx agent --format json --budget 8000
```

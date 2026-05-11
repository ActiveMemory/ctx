---
name: ctx-agent
description: "Load full context packet. Use at session start or when context seems stale or incomplete."
---

Load the full context packet for AI consumption.

## When to Use

- At the start of a session to load all context
- When context seems stale or incomplete
- When switching between different areas of work

## When NOT to Use

- The plugin hook already runs `ctx agent` on session start:
  you rarely need to invoke this manually
- Don't run it just to "refresh" if you already have the context loaded in
  this session

## After Loading

**Read the files listed in "Read These Files (in order)"**: the packet is a
summary, not a substitute. In particular, read CONVENTIONS.md before writing
any code.

Confirm to the user: "I have read the required context files and I'm
following project conventions." Read and confirm before beginning
implementation.

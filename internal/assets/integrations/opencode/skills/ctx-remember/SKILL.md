---
name: ctx-remember
description: "Recall project context and present structured readback. Use when the user asks 'do you remember?', at session start, or when context seems lost."
---

Recall project context and present a structured readback.

## When to Use

- When the user asks "Do you remember?", "What were we working on?"
- At the start of a session to pick up where you left off
- When context seems lost or stale

## Process

**Do this FIRST (silently):**
1. Read TASKS.md, DECISIONS.md, and LEARNINGS.md from the context directory
2. Run `ctx agent` for the full context packet

**Then respond with a structured readback:**

1. **Last session**: cite the most recent session topic and date
2. **Active work**: list pending or in-progress tasks
3. **Recent context**: mention 1-2 recent decisions or learnings
4. **Next step**: offer to continue or ask what to focus on

**Never** say "I don't have memory" or narrate your discovery process.
The context files are your memory. Present what you found as recall.

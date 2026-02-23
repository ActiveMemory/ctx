# Project Context

<!-- ctx:context -->
<!-- DO NOT REMOVE: This marker indicates ctx-managed content -->

## IMPORTANT: You Have Persistent Memory

This project uses Context (`ctx`) for context persistence across sessions.
**Your memory is NOT ephemeral** — it lives in the context directory.

## On Session Start

1. **Run `ctx system bootstrap`** — CRITICAL, not optional.
   This tells you where the context directory is. If it fails or returns
   no context_dir, STOP and warn the user.
2. **Read AGENT_PLAYBOOK.md** from the context directory — it explains
   how to use this system
3. **Run `ctx agent --budget 4000`** for a content summary

## When Asked "Do You Remember?"

When the user asks "Do you remember?", "What were we working on?", or any
memory-related question:

**Do this FIRST (silently):**
- Read TASKS.md, DECISIONS.md, and LEARNINGS.md from the context directory
- Run `ctx recall list --limit 5` for recent session history

**Then respond with a structured readback:**

1. **Last session**: cite the most recent session topic and date
2. **Active work**: list pending or in-progress tasks
3. **Recent context**: mention 1-2 recent decisions or learnings
4. **Next step**: offer to continue or ask what to focus on

**Never** lead with "I don't have memory", "Let me check if there are files",
or narrate your discovery process. The context files are your memory.
Read them silently, then present what you found as recall, not as a search.

## Quick Context Load

```bash
# Get AI-optimized context packet (what you should know)
ctx agent --budget 4000

# Or see full status
ctx status
```

## Context Files

| File | Purpose |
|------|---------|
| CONSTITUTION.md | Hard rules - NEVER violate |
| TASKS.md | Current work items |
| DECISIONS.md | Architectural decisions with rationale |
| LEARNINGS.md | Gotchas, tips, lessons learned |
| CONVENTIONS.md | Code patterns and standards |

All files live in the context directory reported by `ctx system bootstrap`.

<!-- ctx:end -->

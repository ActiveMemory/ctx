# Project Context

<!-- ctx:context -->
<!-- DO NOT REMOVE: This marker indicates ctx-managed content -->

## IMPORTANT: You Have Persistent Memory

This project uses Context (`ctx`) for context persistence across sessions.
**Your memory is NOT ephemeral** - it lives in `.context/` files.

## On Session Start

1. **Read `.context/AGENT_PLAYBOOK.md`** first - it explains how to use this system
2. **Check `.context/sessions/`** for session dumps from previous sessions
3. **Run `ctx status`** to see current context summary

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
| `.context/CONSTITUTION.md` | Hard rules - NEVER violate |
| `.context/TASKS.md` | Current work items |
| `.context/DECISIONS.md` | Architectural decisions with rationale |
| `.context/LEARNINGS.md` | Gotchas, tips, lessons learned |
| `.context/CONVENTIONS.md` | Code patterns and standards |
| `.context/sessions/` | **Session dumps** - check here for deep context |

## Before Session Ends

**ALWAYS offer to persist context before the user quits:**

1. Add learnings: `ctx add learning "..."`
2. Add decisions: `ctx add decision "..."`
3. Save full session: Write to `.context/sessions/YYYY-MM-DD-<topic>.md`

<!-- ctx:end -->

# Project Session Prompt

<!-- ctx:prompt -->
<!-- DO NOT REMOVE: This marker indicates ctx-managed content -->

## On Session Start

1. Read this file (you're doing it now)
2. Run `ctx status` to see current context summary
3. Check `.context/TASKS.md` for active work items

## Context Files

| File | Purpose |
|------|---------|
| `.context/CONSTITUTION.md` | Hard rules — NEVER violate |
| `.context/TASKS.md` | Current work items |
| `.context/DECISIONS.md` | Architectural decisions with rationale |
| `.context/LEARNINGS.md` | Gotchas and lessons learned |
| `.context/CONVENTIONS.md` | Code patterns and standards |
| `.context/AGENT_PLAYBOOK.md` | How to persist context, session patterns |

## Working Style

- **Ask questions** when requirements are unclear
- **Persist context** as you work (don't wait for session end)
- **Use `ctx add`** for learnings, decisions, tasks
- **Check existing patterns** before writing new code

## Persist as You Go

After completing meaningful work, capture what matters:

| Trigger | Action |
|---------|--------|
| Completed a task | Mark done in TASKS.md, add learnings if any |
| Made a decision | `ctx add decision "..."` |
| Discovered a gotcha | `ctx add learning "..."` |
| Significant code changes | Consider what's worth capturing |

Don't wait for the session to end — it may never come cleanly.

<!-- ctx:prompt:end -->

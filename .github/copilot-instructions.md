# Project Context

<!-- ctx:copilot -->
<!-- DO NOT REMOVE: This marker indicates ctx-managed content -->

## Context System

This project uses Context (`ctx`) for persistent AI context
management. Your memory is NOT ephemeral — it lives in `.context/` files.

## On Session Start

Read these files **in order** before starting any work:

1. `.context/CONSTITUTION.md` — Hard rules, NEVER violate
2. `.context/TASKS.md` — Current work items
3. `.context/CONVENTIONS.md` — Code patterns and standards
4. `.context/ARCHITECTURE.md` — System structure
5. `.context/DECISIONS.md` — Architectural decisions with rationale
6. `.context/LEARNINGS.md` — Gotchas, tips, lessons learned
7. `.context/GLOSSARY.md` — Domain terms and abbreviations
8. `.context/AGENT_PLAYBOOK.md` — How to use this context system

After reading, confirm: "I have read the required context files and I'm
following project conventions."

## When Asked "Do You Remember?"

**Do this FIRST (silently):**
- Read `.context/TASKS.md`
- Read `.context/DECISIONS.md` and `.context/LEARNINGS.md`
- Check `.context/sessions/` for recent session files

**Then respond with a structured readback:**
1. **Last session**: cite the most recent session topic and date
2. **Active work**: list pending or in-progress tasks
3. **Recent context**: mention 1-2 recent decisions or learnings
4. **Next step**: offer to continue or ask what to focus on

**Never** lead with "I don't have memory." The `.context/` files are
your memory. Read them silently, then present what you found as recall.

## Session Persistence

After completing meaningful work, save a session summary to
`.context/sessions/`.

### Session File Format

Create a file named `YYYY-MM-DD-topic.md`:

```markdown
# Session: YYYY-MM-DD — Brief Topic Description

## What Was Done
- Describe completed work items

## Decisions
- Key decisions made and their rationale

## Learnings
- Gotchas, tips, or insights discovered

## Next Steps
- Follow-up work or remaining items
```

### When to Save

- After completing a task or feature
- After making architectural decisions
- After a debugging session
- Before ending the session
- At natural breakpoints in long sessions

## Context Updates During Work

Proactively update context files as you work:

| Event                       | Action                              |
|-----------------------------|-------------------------------------|
| Made architectural decision | Add to `.context/DECISIONS.md`  |
| Discovered gotcha/bug       | Add to `.context/LEARNINGS.md`  |
| Established new pattern     | Add to `.context/CONVENTIONS.md` |
| Completed task              | Mark [x] in `.context/TASKS.md` |

## Self-Check

Periodically ask yourself:

> "If this session ended right now, would the next session know what happened?"

If no — save a session file or update context files before continuing.

## CLI Commands

If `ctx` is installed, use these commands:

```bash
ctx status        # Context summary and health check
ctx agent         # AI-ready context packet
ctx drift         # Check for stale context
ctx recall list   # Recent session history
```

<!-- ctx:copilot:end -->

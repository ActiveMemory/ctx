# Project Context

<!-- ctx:copilot -->
<!-- DO NOT REMOVE: This marker indicates ctx-managed content -->

## Context System

This project uses Context (`ctx`) for persistent AI context
management. Your memory is NOT ephemeral: it lives in `.context/` files.

## On Session Start

Read these files **in order** before starting any work:

1. `.context/CONSTITUTION.md`: Hard rules, NEVER violate
2. `.context/TASKS.md`: Current work items
3. `.context/CONVENTIONS.md`: Code patterns and standards
4. `.context/ARCHITECTURE.md`: System structure
5. `.context/DECISIONS.md`: Architectural decisions with rationale
6. `.context/LEARNINGS.md`: Gotchas, tips, lessons learned
7. `.context/GLOSSARY.md`: Domain terms and abbreviations
8. `.context/AGENT_PLAYBOOK.md`: How to use this context system

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
# Session: YYYY-MM-DD - Brief Topic Description

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

| Event                       | Action                           |
|-----------------------------|----------------------------------|
| Made architectural decision | Add to `.context/DECISIONS.md`   |
| Discovered gotcha/bug       | Add to `.context/LEARNINGS.md`   |
| Established new pattern     | Add to `.context/CONVENTIONS.md` |
| Completed task              | Mark [x] in `.context/TASKS.md`  |

## Self-Check

Periodically ask yourself:

> "If this session ended right now, would the next session know what happened?"

If no: save a session file or update context files before continuing.

## CLI Commands

If `ctx` is installed, use these commands:

```bash
ctx status        # Context summary and health check
ctx agent         # AI-ready context packet
ctx drift         # Check for stale context
ctx journal source   # Recent session history
```

## MCP Tools (Preferred)

When an MCP server named `ctx` is available, **always prefer MCP tools
over terminal commands** for context operations. MCP tools provide
validation, session tracking, and boundary checks automatically.

| MCP Tool                    | Purpose                              |
|-----------------------------|--------------------------------------|
| `ctx_status`                | Context summary and health check     |
| `ctx_add`                   | Add task, decision, learning, or convention |
| `ctx_complete`              | Mark a task as done                  |
| `ctx_drift`                 | Check for stale or drifted context   |
| `ctx_recall`                | Query session history                |
| `ctx_next`                  | Get the next task to work on         |
| `ctx_compact`               | Archive completed tasks              |
| `ctx_watch_update`          | Write entry and queue for review     |
| `ctx_check_task_completion` | Match recent work to open tasks      |
| `ctx_session_event`         | Signal session start or end          |
| `ctx_remind`                | List pending reminders               |

**Rule**: Do NOT run `ctx` in the terminal when the equivalent MCP tool
exists. MCP tools enforce boundary validation and track session state.
Terminal fallback is only for commands without an MCP equivalent (e.g.,
`ctx agent`, `ctx recall list`).

## Governance: When to Call Tools

The MCP server tracks session state and appends warnings to tool
responses when governance actions are overdue. Follow this protocol:

### Session Lifecycle

1. **BEFORE any work**: call `ctx_session_event(type="start")`, then
   `ctx_status()` to load context.
2. **Before ending**: call `ctx_session_event(type="end")` to flush
   pending state.

### During Work

- **After making a decision or discovering a gotcha**: call `ctx_add()`
  to persist it immediately — not at session end.
- **After completing a task**: call `ctx_complete()` or
  `ctx_check_task_completion()`.
- **Every 10–15 tool calls or 15 minutes**: call `ctx_drift()` to
  check for stale context.
- **Before git commit**: call `ctx_status()` to verify context health.

### Responding to Warnings

When a tool response contains a `⚠` warning, act on it in your next
action. Do not ignore governance warnings — they indicate context
hygiene actions that are overdue.

When a tool response contains a `🚨 CRITICAL` warning, **stop current
work immediately** and address the violation. These indicate dangerous
commands, sensitive file access, or policy violations detected by the
VS Code extension. Review the action, revert if unintended, and explain
what happened before continuing.

### Detection Ring

The VS Code extension monitors terminal commands and file access in
real time. The following actions are flagged as violations:

- **Dangerous commands**: `sudo`, `rm -rf /`, `git push`, `git reset
  --hard`, `curl`, `wget`, `chmod 777`
- **hack/ scripts**: Direct execution of `hack/*.sh` — use `make`
  targets instead
- **Sensitive files**: Editing `.env`, `.pem`, `.key`, or files
  matching `credentials` or `secret`

Violations are recorded and surfaced as CRITICAL warnings in your next
MCP tool response. The user also sees a VS Code notification.

<!-- ctx:copilot:end -->

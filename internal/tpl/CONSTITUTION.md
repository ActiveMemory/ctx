# Constitution

These rules are INVIOLABLE. If a task requires violating these, the task is wrong.

## Security Invariants

- [ ] Never commit secrets, tokens, API keys, or credentials
- [ ] Never store customer/user data in context files

## Quality Invariants

- [ ] All code must pass tests before commit
- [ ] No TODO comments in main branch (move to TASKS.md)
- [ ] Path construction uses stdlib — no string concatenation (security: prevents path traversal)

## Process Invariants

- [ ] All architectural changes require a decision record

## TASKS.md Structure Invariants

TASKS.md must remain a replayable checklist. Uncheck all items and re-run = verify/redo all tasks in order.

- [ ] **Never move tasks** — tasks stay in their Phase section permanently
- [ ] **Never remove Phase headers** — Phase labels provide structure and order
- [ ] **Never merge or collapse Phase sections** — each phase is a logical unit
- [ ] **Never delete tasks** — mark as `[x]` completed, or `[-]` skipped with reason
- [ ] **Use inline labels for status** — add `#in-progress` to task text, don't move it
- [ ] **No "In Progress" / "Next Up" sections** — these encourage moving tasks
- [ ] **Ask before restructuring** — if structure changes seem needed, ask the user first

## Context Preservation Invariants

- [ ] **Archival is allowed, deletion is not** — use `ctx tasks archive` to move completed tasks to `.context/archive/`, never delete context history
- [ ] **Archive preserves structure** — archived tasks keep their Phase headers for traceability

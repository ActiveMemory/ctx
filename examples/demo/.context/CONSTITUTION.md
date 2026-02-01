# Constitution

These rules are INVIOLABLE. If a task requires violating these, the task is wrong.

## Security Invariants

- [ ] Never commit secrets, tokens, API keys, or credentials
- [ ] Never store customer/user data in context files
- [ ] All user input must be validated and sanitized

## Quality Invariants

- [ ] All code must pass tests before commit
- [ ] No TODO comments in main branch (move to TASKS.md)
- [ ] Breaking API changes require deprecation period

## Process Invariants

- [ ] All architectural changes require a decision record in DECISIONS.md

## TASKS.md Structure Invariants

TASKS.md must remain a replayable checklist. Uncheck all items and re-run
the loop = verify/redo all tasks in order.

- [ ] **Never move tasks** — tasks stay in their Phase section permanently
- [ ] **Never remove Phase headers** — Phase labels provide structure and order
- [ ] **Never delete tasks** — mark as `[x]` completed, or `[-]` skipped with reason
- [ ] **Use inline labels for status** — add `#in-progress` to task text, don't move it
- [ ] **No "In Progress" sections** — these encourage moving tasks
- [ ] **Ask before restructuring** — if structure changes seem needed, ask the user first

## Context Preservation Invariants

- [ ] **Archival is allowed, deletion is not** — use `ctx tasks archive` to move completed tasks, never delete context history
- [ ] **Archive preserves structure** — archived tasks keep their Phase headers for traceability

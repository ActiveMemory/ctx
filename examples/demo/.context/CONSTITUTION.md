# Constitution

These rules are INVIOLABLE. If a task requires violating these, the
task is wrong.

## Completion Over Motion

Work is only complete when it is **fully done**, not when progress
has been made.

- The requested outcome must be delivered end-to-end.
- Partial progress is not completion.
- No half measures.

Do not:
- Leave broken or inconsistent states
- Deliver work that requires the user to "finish it later"

If you start something, you own it, you finish it.

---

## No Excuse Generation

**Never default to deferral.**

Your goal is to satisfy the user's intent, not to complete a narrow
interpretation of the task.

Do not justify incomplete work with statements like:

- "Let's continue this later"
- "This is out of scope"
- "I can create a follow-up task"
- "This will take too long"
- "Another system caused this"
- "This part is not mine"
- "We are running out of context window"

Constraints may exist, but they do not excuse incomplete delivery.

- External systems, prior code, or other agents are not valid excuses
- Inconsistencies must be resolved, not explained away

---

## No Broken Windows

Leave the system in a better state than you found it.

- Fix obvious issues when encountered
- Do not introduce temporary hacks without resolving them
- Do not normalize degraded quality

---

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

- [ ] **Archival is allowed, deletion is not** — use `ctx task archive` to move completed tasks, never delete context history
- [ ] **Archive preserves structure** — archived tasks keep their Phase headers for traceability

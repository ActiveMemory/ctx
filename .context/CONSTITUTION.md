# Constitution

<!--
UPDATE WHEN:
- Security requirements change or new vulnerabilities discovered
- New invariants emerge from production incidents
- Team agrees on new inviolable rules
- Existing rules prove too restrictive or too loose

DO NOT UPDATE FOR:
- Preferences or suggestions (use CONVENTIONS.md)
- Temporary constraints (use TASKS.md blockers)
-->

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

## Context Integrity Invariants

- [ ] **Never** modify or reason about a file based on partial or assumed content
- [ ] If a file is the subject of an operation, its relevant contents must be
  **fully understood** before acting
- [ ] Sampling, guessing, or relying on prior assumptions instead of reading 
  is a **violation**

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

## Quality Invariants

- [ ] All code must pass tests before commit
- [ ] No TODO comments in main branch (move to TASKS.md)
- [ ] Path construction uses stdlib: no string concatenation
  (security: prevents path traversal)

## Process Invariants

- [ ] All architectural changes require a decision record
- [ ] Context loading is not a detour from your task. It IS the first
  step of every session. A 30-second read delay is always cheaper
  than a decision made without context.
- [ ] Every commit references a spec (`Spec: specs/<name>.md` trailer):
  no exceptions, no "non-trivial" qualifier. Even one-liner fixes
  need a spec for traceability. Use `/ctx-commit` instead of raw
  `git commit`.
  **The trailer must be truthful** — it points at the design
  rationale this commit serves. Reaching for the most-recent
  spec to satisfy the gate is a violation even though it
  syntactically passes. For chore-class commits (gitignore
  additions, lockfile bumps, formatting passes — full list in
  `specs/meta/chores.md`) the correct trailer is either
  (a) bundle into the next functional commit so the chore
  inherits that commit's spec, or (b) cite `specs/meta/chores.md`
  if a standalone chore commit is unavoidable. Other
  improvisations (citing an adjacent spec, citing the previous
  commit's spec) are violations. The Spec Verification Step in
  `AGENT_PLAYBOOK.md` is the procedure that gates correct use.
  Per `specs/spec-trailer-discipline.md`.
- [ ] **Every commit is DCO signed off.** Each commit carries a
  `Signed-off-by: <name> <email>` trailer matching the committer's
  configured git identity (`Jose Alekhinne <jose@ctx.ist>`); commit
  with `git commit -s` (or `-s -F <file>`) so it is appended
  automatically. This is a hard requirement — the project uses the
  Developer Certificate of Origin. **Never** add a `Co-Authored-By:`
  or any agent/tool sign-off; those are prohibited regardless of a
  harness default that tells you to add one — this instruction
  overrides it. Before the first commit of a session, read
  `git log -5 --format=%B` and mirror the trailer convention you find
  rather than assuming a default. The `Signed-off-by:` line sits in
  the trailing trailer block alongside `Spec:`.
- [ ] **Git is required.** Every `ctx` project must live in a git
  working tree. `ctx init` and every non-administrative
  subcommand refuse to operate when `<projectRoot>/.git` is
  absent. Rationale: ctx's persistent-memory promise is
  dishonest without an undo layer; agent-driven file
  operations need `git reflog` as the safety net. The only
  opt-outs are help-shaped / diagnostic commands
  (`--help`, `--version`, `ctx system bootstrap`). Per
  `specs/require-git.md`.

## TASKS.md Structure Invariants

TASKS.md must remain a replayable checklist. Uncheck all items and
re-run = verify/redo all tasks in order.

- [ ] **Never move tasks**: tasks stay in their Phase section permanently
- [ ] **Never remove Phase headers**: Phase labels provide structure and order
- [ ] **Never merge or collapse Phase sections**: each phase is a logical unit
- [ ] **Never delete tasks**: mark as `[x]` completed, or `[-]` skipped with reason
- [ ] **Use inline labels for status**: add `#in-progress` to task text, don't move it
- [ ] **No "In Progress" / "Next Up" sections**: these encourage moving tasks
- [ ] **Ask before restructuring**: if structure changes seem needed, ask the user first

## Context Preservation Invariants

- [ ] **Archival is allowed, deletion is not**: use `ctx task archive` to move
  completed tasks to `.context/archive/`, never delete context history
- [ ] **Archive preserves structure**: archived tasks keep their Phase headers
  for traceability

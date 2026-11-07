---
name: ctx-implement
description: "Execute a plan step-by-step with verification. Use when you have a plan document — canonically specs/plans/<milestone>.md from /ctx-task-out — and need disciplined, checkpointed implementation."
---

Take a plan — canonically `specs/plans/<milestone>.md` as written
by `/ctx-task-out`, though inline text, another file path, or a
plan from the conversation also work — and execute it
step-by-step with build/test verification between steps.

## When to Use

- After `/ctx-task-out` has decomposed a spec into
  `specs/plans/<milestone>.md` (the canonical input)
- When the user provides a plan document or file and says
  "implement this"
- When a multi-step task has been planned and needs disciplined
  execution
- When the user wants checkpointed progress with verification
  at each step
- After `/ctx-brainstorm` or plan mode produces an approved plan

## When NOT to Use

- For single-step tasks: just do them directly
- When handed a bare multi-milestone spec instead of a plan:
  suggest `/ctx-task-out --spec <path> --milestone <first>`
  first; decomposing on the fly is what it exists to prevent
- When the plan is vague or incomplete: use `/ctx-brainstorm`
  first to refine it
- When the user wants to explore or discuss, not execute
- When changes are trivial (typo fix, config tweak)

## Usage Examples

```text
/ctx-implement
/ctx-implement specs/plans/m0a.md
/ctx-implement path/to/plan.md
/ctx-implement (the plan from our discussion above)
```

## Process

### 1. Load the plan

- If a file path is provided, read it
- If the file is a multi-milestone spec rather than a plan (no
  task breakdown, no acceptance criteria, spans milestones),
  redirect: suggest `/ctx-task-out --spec <path> --milestone
  <first>` and stop rather than improvising a decomposition
- If the plan's header shows `Status: Blocked`, stop: a
  deferrable TBD graduated to blocking mid-milestone. Route its
  resolution (a spec edit or DECISIONS.md entry) and a
  `/ctx-task-out` amendment run before executing further tasks
- If inline text is provided, use it directly
- If neither, look back in the conversation for the most
  recent plan or approved design
- If no plan can be found, ask the user for one

### 2. Break into steps

Parse the plan into discrete, checkable steps. Each step
should be:
- **Atomic**: one logical change (a file, a function, a test)
- **Verifiable**: has a clear pass/fail check
- **Ordered**: dependencies respected (create before use,
  test after implement)

Present the step list to the user for confirmation:

> **Implementation plan** (N steps):
>
> 1. [Step description] - verify: [check]
> 2. [Step description] - verify: [check]
> 3. ...
>
> Ready to start?

### 3. Execute step-by-step

For each step:

1. **Announce** what you're doing (one line)
2. **Think through** the change before writing code: what does
   it touch, what could break, what's the simplest correct path?
3. **Implement** the change
4. **Verify** with the appropriate check:
   - Task from a task-out plan → its acceptance criterion,
     verbatim (in addition to the map below)
   - Go code changed → `CGO_ENABLED=0 go build -o /dev/null ./cmd/ctx`
   - Tests affected → `CGO_ENABLED=0 go test ./...`
   - Config/template changed → build to verify embeds
   - Docs only → no verification needed
5. **Report** step result: pass or fail
6. **If failed**: stop, diagnose, fix, re-verify before
   moving to the next step

Verify after every individual step before proceeding to the next.

### 4. Checkpoint progress

After every 3-5 steps (or after a significant milestone):
- Summarize what has been completed
- If executing `specs/plans/<milestone>.md`, update the execution
  ledger (see Ledger Duties below)
- Note any deviations from the plan
- Ask the user if they want to continue, adjust, or stop

### 5. Wrap up

After all steps complete:
- Run a final full verification (`make check` or
  `CGO_ENABLED=0 go build && go test ./...`)
- Summarize what was implemented
- Note any deviations from the original plan
- Suggest context to persist (decisions, learnings, tasks)

## Ledger Duties (plans from /ctx-task-out)

A task-out plan is the execution ledger — the only record of
milestone progress. Executing one carries four bookkeeping duties:

- **`st` is the record.** Flip a task's `st` cell to `[x]` only
  when its acceptance criterion has demonstrably passed — the
  command ran, the test is green, the behavior was observed.
  Tasks obsoleted by amendment become `[o]`. `st` never moves
  backwards silently; a regression is a deviation to report.
- **DoD is not yours to derive.** Scope & DoD checkboxes are
  confirmed by measurement or by the user — never checked because
  the tasks that "cover" them are done. The rolling-wave gate
  reads DoD only; deriving it from task completion defeats the
  gate.
- **Project epics outward.** TASKS.md epics carry disjoint
  task-id ranges (`Plan: specs/plans/<milestone>.md (Txx–Tyy)`).
  When every task in a range is `[x]` or `[o]`, mark that epic
  `[x]`. Sync is one-way, plan → TASKS.md; never track task
  state in TASKS.md directly.
- **Amendments, not edits.** Never edit a task's acceptance
  criterion in place; a criterion change goes back through
  `/ctx-task-out` (amendment mode). When a measurement gate
  fires (Risks & measurement gates), stop and route the outcome
  through an amendment before executing dependent tasks.

## Step Verification Map

| Change type        | Verification command                              |
|--------------------|---------------------------------------------------|
| Go source code     | `CGO_ENABLED=0 go build -o /dev/null ./cmd/ctx`   |
| Test files         | `CGO_ENABLED=0 go test ./...`                     |
| Templates/embeds   | `CGO_ENABLED=0 go build -o /dev/null ./cmd/ctx`   |
| Makefile           | Run the new/changed target                        |
| Skill files        | Build (to verify embed) + check live copy matches |
| Docs/Markdown only | None required                                     |
| Shell scripts      | `bash -n script.sh` (syntax check)                |

## Handling Failures

When a step fails verification:

1. **Don't panic**: read the error output carefully
2. **Reason through** the failure step-by-step before attempting
   a fix; understand the cause, not just the symptom
3. **Fix** the issue in the current step
4. **Re-verify** the fix
5. **Only then** move to the next step
6. If the fix changes the plan, note the deviation

If a step fails repeatedly (3+ attempts), stop and ask the
user for guidance rather than thrashing.

## Output Format

Progress updates should be concise:

```
Step 1/6: Create ctx-next skill directory .......... OK
Step 2/6: Write SKILL.md template .................. OK
Step 3/6: Copy to live skill directory ............. OK
Step 4/6: Build to verify template embeds .......... OK
Step 5/6: Run tests ................................ OK
Step 6/6: Mark task in TASKS.md .................... OK

All 6 steps complete. Build and tests pass.
```

## Examples

### Good Implementation

> **Step 3/8**: Add `check` target to Makefile
> Added `check: build audit` after the `audit` target.
> Verify: `make check` ... build OK, audit OK.
> **Result**: PASS

### Bad Implementation

> "I'll implement the whole plan now"
> *[makes all changes at once without verification]*
> "Done! Everything should work."

(No step-by-step, no verification, no checkpoints: this
defeats the purpose of the skill.)

## Quality Checklist

Before starting, verify:
- [ ] Plan exists and is clear enough to execute
- [ ] Steps are broken down and presented to the user
- [ ] User confirmed readiness to proceed

During execution, verify:
- [ ] Each step is verified before moving on
- [ ] Failures are fixed in place, not deferred
- [ ] Checkpoints happen every 3-5 steps
- [ ] Task-out plans: `st` flipped only on demonstrated
      acceptance; epics projected to TASKS.md when their range
      completes; DoD boxes left to measurement or the user

After completion, verify:
- [ ] Final full verification passes
- [ ] Deviations from plan are noted
- [ ] Summary of what was implemented is presented
- [ ] Context persistence is suggested if warranted

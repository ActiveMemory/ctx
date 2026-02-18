---
name: ctx-implement
description: "Execute a plan step-by-step with verification. Use when you have a plan document and need disciplined, checkpointed implementation."
---

Take a plan (inline text, file path, or from the conversation)
and execute it step-by-step with build/test verification between
steps.

## When to Use

- When the user provides a plan document or file and says
  "implement this"
- When a multi-step task has been planned and needs disciplined
  execution
- When the user wants checkpointed progress with verification
  at each step
- After `/brainstorm` or plan mode produces an approved plan

## When NOT to Use

- For single-step tasks — just do them directly
- When the plan is vague or incomplete — use `/brainstorm`
  first to refine it
- When the user wants to explore or discuss, not execute
- When changes are trivial (typo fix, config tweak)

## Usage Examples

```text
/ctx-implement
/ctx-implement path/to/plan.md
/ctx-implement (the plan from our discussion above)
```

## Process

### 1. Load the plan

- If a file path is provided, read it
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
> 1. [Step description] — verify: [check]
> 2. [Step description] — verify: [check]
> 3. ...
>
> Ready to start?

### 3. Execute step-by-step

For each step:

1. **Announce** what you're doing (one line)
2. **Think through** the change before writing code — what does
   it touch, what could break, what's the simplest correct path?
3. **Implement** the change
3. **Verify** with the appropriate check:
   - Go code changed → `CGO_ENABLED=0 go build -o /dev/null ./cmd/ctx`
   - Tests affected → `CGO_ENABLED=0 go test ./...`
   - Config/template changed → build to verify embeds
   - Docs only → no verification needed
4. **Report** step result: pass or fail
5. **If failed**: stop, diagnose, fix, re-verify before
   moving to the next step

Do NOT skip verification. Do NOT batch multiple steps
without checking between them.

### 4. Checkpoint progress

After every 3-5 steps (or after a significant milestone):
- Summarize what has been completed
- Note any deviations from the plan
- Ask the user if they want to continue, adjust, or stop

### 5. Wrap up

After all steps complete:
- Run a final full verification (`make check` or
  `CGO_ENABLED=0 go build && go test ./...`)
- Summarize what was implemented
- Note any deviations from the original plan
- Suggest context to persist (decisions, learnings, tasks)

## Step Verification Map

| Change type          | Verification command                          |
|----------------------|-----------------------------------------------|
| Go source code       | `CGO_ENABLED=0 go build -o /dev/null ./cmd/ctx` |
| Test files           | `CGO_ENABLED=0 go test ./...`                 |
| Templates/embeds     | `CGO_ENABLED=0 go build -o /dev/null ./cmd/ctx` |
| Makefile             | Run the new/changed target                    |
| Skill files          | Build (to verify embed) + check live copy matches |
| Docs/markdown only   | None required                                 |
| Shell scripts        | `bash -n script.sh` (syntax check)            |

## Handling Failures

When a step fails verification:

1. **Don't panic** — read the error output carefully
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

(No step-by-step, no verification, no checkpoints — this
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

After completion, verify:
- [ ] Final full verification passes
- [ ] Deviations from plan are noted
- [ ] Summary of what was implemented is presented
- [ ] Context persistence is suggested if warranted

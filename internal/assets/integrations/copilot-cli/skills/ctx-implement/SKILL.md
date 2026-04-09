---
name: ctx-implement
description: "Execute a plan step-by-step with verification. Use when you have a plan document and need disciplined, checkpointed implementation."
tools: [bash, read, write, edit, glob, grep]
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
- After `ctx-brainstorm` or plan mode produces an approved plan

## When NOT to Use

- For single-step tasks: just do them directly
- When the plan is vague or incomplete: use `ctx-brainstorm`
  first to refine it
- When the user wants to explore or discuss, not execute
- When changes are trivial (typo fix, config tweak)

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
> 1. [Step description] - verify: [check]
> 2. [Step description] - verify: [check]
> 3. ...
>
> Ready to start?

### 3. Execute step-by-step

For each step:

1. **Announce** what you're doing (one line)
2. **Think through** the change before writing code
3. **Implement** the change
4. **Verify** with the appropriate check:
   - Go code changed → `go build ./cmd/ctx/...`
   - Tests affected → `go test ./...`
   - Config/template changed → build to verify embeds
   - Docs only → no verification needed
5. **Report** step result: pass or fail
6. **If failed**: stop, diagnose, fix, re-verify before
   moving to the next step

### 4. Checkpoint progress

After every 3-5 steps (or after a significant milestone):
- Summarize what has been completed
- Note any deviations from the plan
- Ask the user if they want to continue, adjust, or stop

### 5. Wrap up

After all steps complete:
- Run a final full verification
- Summarize what was implemented
- Note any deviations from the original plan
- Suggest context to persist (decisions, learnings, tasks)

## Step Verification Map

| Change type        | Verification command           |
|--------------------|--------------------------------|
| Go source code     | `go build ./cmd/ctx/...`       |
| Test files         | `go test ./...`                |
| Templates/embeds   | `go build ./cmd/ctx/...`       |
| Makefile           | Run the changed target         |
| Skill files        | Build to verify embed          |
| Docs/markdown only | None required                  |
| Shell scripts      | `bash -n script.sh`            |

## Handling Failures

When a step fails verification:

1. Read the error output carefully
2. Reason through the failure before attempting a fix
3. Fix the issue in the current step
4. Re-verify the fix
5. Only then move to the next step

If a step fails repeatedly (3+ attempts), stop and ask the
user for guidance.

## Quality Checklist

Before starting:
- [ ] Plan exists and is clear enough to execute
- [ ] Steps are broken down and presented to the user
- [ ] User confirmed readiness to proceed

During execution:
- [ ] Each step is verified before moving on
- [ ] Failures are fixed in place, not deferred
- [ ] Checkpoints happen every 3-5 steps

After completion:
- [ ] Final full verification passes
- [ ] Deviations from plan are noted
- [ ] Context persistence is suggested if warranted

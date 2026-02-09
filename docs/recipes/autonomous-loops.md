---
title: Running an Unattended AI Agent
icon: lucide/repeat
---

![ctx](../images/ctx-banner.png)

## Problem

You have a project with a clear list of tasks and you want an AI agent to
work through them autonomously: overnight, **unattended**, without you sitting
at the keyboard. 

Each iteration needs to **remember** what the previous one did, mark tasks 
as completed, and know **when** to stop.

Without persistent memory, every iteration starts fresh and the loop
collapses. With `ctx`, each iteration picks up exactly where the last one
left off: **but only if the agent proactively persists its context**. 

This is the key insight: unattended operation works because the agent treats
context persistence as part of the work itself, not as an afterthought.

## Commands and Skills Used

| Tool                    | Type    | Purpose                                                            |
|-------------------------|---------|--------------------------------------------------------------------|
| `ctx init --ralph`      | Command | Initialize project for unattended operation (no human in the loop) |
| `ctx loop`              | Command | Generate the loop shell script                                     |
| `ctx watch --auto-save` | Command | Monitor AI output and persist context updates                      |
| `ctx load`              | Command | Display assembled context (for debugging)                          |
| `/ctx-loop`             | Skill   | Generate loop script from inside Claude Code                       |
| `/ctx-implement`        | Skill   | Execute a plan step-by-step with verification                      |
| `/ctx-context-monitor`  | Skill   | Automated context capacity alerts during long sessions             |

## The Workflow

### Step 1: Initialize for Unattended Operation

Start by creating a `.context/` directory configured so the agent can
work without human input. The `--ralph` flag sets up `PROMPT.md` so the
agent makes its own decisions rather than asking clarifying questions.

```bash
ctx init --ralph
```

This creates `.context/` with all template files, `PROMPT.md` configured
for autonomous iteration, `IMPLEMENTATION_PLAN.md`, and `.claude/` hooks
and skills for Claude Code. Without `--ralph`, the agent pauses to ask
questions when requirements are unclear. For unattended runs, you want it
to make reasonable choices and document them in `DECISIONS.md` instead.

### Step 2: Populate TASKS.md with Phased Work

Open `.context/TASKS.md` and organize your work into phases. The agent
works through these systematically, top to bottom, from the highest priority
task first.

```markdown
# Tasks

## Phase 1: Foundation

- [ ] Set up project structure and build system `#priority:high`
- [ ] Configure testing framework `#priority:high`
- [ ] Create CI pipeline `#priority:medium`

## Phase 2: Core Features

- [ ] Implement user registration `#priority:high`
- [ ] Add email verification `#priority:high`
- [ ] Create password reset flow `#priority:medium`

## Phase 3: Hardening

- [ ] Add rate limiting to API endpoints `#priority:medium`
- [ ] Improve error messages `#priority:low`
- [ ] Write integration tests `#priority:medium`
```

Phased organization matters because it gives the agent **natural
boundaries**. Phase 1 tasks should be completable without Phase 2 code
existing yet.

### Step 3: Configure PROMPT.md

The `--ralph` flag generates a `PROMPT.md` that instructs the agent to
operate autonomously:

1. Read `.context/CONSTITUTION.md` first (*hard rules, never violated*)
2. Load context from `.context/` files
3. Pick ONE task per iteration
4. Complete the task and **proactively update context files**
5. Commit changes (including `.context/`)
6. Signal status with a completion signal

You can customize `PROMPT.md` for your project. The critical parts are
the one-task-per-iteration discipline, proactive context persistence,
and the completion signals at the end:

```markdown
## Signal Status

End your response with exactly ONE of:

- `SYSTEM_CONVERGED` — All tasks in TASKS.md are complete
- `SYSTEM_BLOCKED` — Cannot proceed, need human input (explain why)
- (no signal) — More work remains, continue to next iteration
```

### Step 4: Generate the Loop Script

Use `ctx loop` to generate a `loop.sh` tailored to your AI tool:

```bash
# Generate for Claude Code with a 10-iteration cap
ctx loop --tool claude --max-iterations 10

# Generate for Aider
ctx loop --tool aider --max-iterations 10

# Custom prompt and output file
ctx loop --tool claude --prompt TASKS.md --output my-loop.sh
```

The generated script reads `PROMPT.md`, pipes it to the AI tool, checks
for completion signals, and loops until done or the cap is reached. You
can also use the `/ctx-loop` skill from inside Claude Code.

!!! tip "Shell Loop is Best Practice"
    The shell while loop is the recommended approach for autonomous runs.
    Each iteration spawns a **fresh AI process**, so the only state that
    carries between iterations is what lives in `.context/` and git. This
    is "pure ralph": memory is explicit, not accidental.

    Claude Code's built-in `/loop` command runs iterations within the
    same session, which means context window state leaks between
    iterations — the agent "remembers" things from earlier iterations
    that were never persisted. This is convenient for short explorations
    (2-5 iterations) but less reliable for long unattended runs.
    See [Autonomous Loops: Shell Loop vs Built-in Loop](../autonomous-loop.md#quick-start-shell-while-loop-recommended)
    for details.

### Step 5: Run with Watch Mode

Open two terminals. In the first, run the loop. In the second, run
`ctx watch` to automatically process context updates from the AI output.

```bash
# Terminal 1: Run the loop
./loop.sh 2>&1 | tee /tmp/loop.log

# Terminal 2: Watch for context updates
ctx watch --log /tmp/loop.log --auto-save
```

The `--auto-save` flag periodically saves session snapshots to
`.context/sessions/`. The watch command parses XML context-update
commands from the AI output and applies them:

```xml
<context-update type="complete">user registration</context-update>
<context-update type="learning">Email verification needs SMTP configured</context-update>
```

### Step 6: Completion Signals End the Loop

The loop terminates when the agent emits one of these signals:

| Signal               | Meaning                        | What Happens                       |
|----------------------|--------------------------------|------------------------------------|
| `SYSTEM_CONVERGED`   | All tasks in TASKS.md are done | Loop exits successfully            |
| `SYSTEM_BLOCKED`     | Agent cannot proceed           | Loop exits, you review the blocker |
| `BOOTSTRAP_COMPLETE` | Initial scaffolding done       | Loop exits after setup phase       |

When you return in the morning, check the log and the context files:

```bash
# See what happened
tail -100 /tmp/loop.log

# Check task progress
ctx status

# Load full context to see decisions and learnings
ctx load
```

### Step 7: Use /ctx-implement for Plan Execution

Within each iteration, the agent can use `/ctx-implement` to execute
multi-step plans with verification between each step. This is especially
useful for complex tasks that involve multiple files.

The skill breaks a plan into atomic, verifiable steps:

```text
Step 1/6: Create user model .................. OK
Step 2/6: Add database migration ............. OK
Step 3/6: Implement registration handler ..... OK
Step 4/6: Write unit tests ................... OK
Step 5/6: Run test suite ..................... FAIL
  → Fixed: missing test dependency
  → Re-verify ............................ OK
Step 6/6: Update TASKS.md .................... OK
```

Each step is verified (build, test, syntax check) before moving to the
next. Failures are fixed in place, not deferred.

## Putting It Together

The full sequence for an overnight unattended run:

```bash
# 1. Set up the project for unattended operation
ctx init --ralph

# 2. Edit TASKS.md with your phased work items
# 3. Review and customize PROMPT.md

# 4. Generate the loop
ctx loop --tool claude --max-iterations 20

# 5. Start watch mode in background
ctx watch --log /tmp/loop.log --auto-save &

# 6. Run the loop
./loop.sh 2>&1 | tee /tmp/loop.log

# 7. Next morning: review results
ctx status
ctx load
```

## Why Autonomous Loops Work: Proactive Context Persistence

The autonomous loop pattern works **because the agent is proactive about
persisting context**. Without proactive behavior, the loop degrades into
disconnected iterations that repeat work, forget decisions, and lose track
of progress. The agent cannot rely on a human to prompt it — it must treat
context persistence as part of every task, not as a separate step.

### The Agent Playbook's Self-Check Drives Autonomy

The Agent Playbook defines a self-check prompt that the agent applies
continuously:

> "If this session ended right now, would the next session know what
> happened?"

In interactive sessions, a human might remind the agent to save context.
When the agent runs unattended, **no one is there to remind it**. The self-check
becomes the agent's internal trigger: after every meaningful action, the
agent evaluates whether the context files reflect reality and acts
immediately if they don't.

### What the Agent Does Proactively Between Iterations

At each milestone within an iteration, the agent persists without waiting
for instructions:

| Milestone | What the Agent Does Proactively |
|---|---|
| Completes a task | Marks it `[x]` in TASKS.md with `#done` timestamp |
| Discovers a gotcha | Adds it to LEARNINGS.md immediately |
| Makes a design choice | Records it in DECISIONS.md with rationale |
| Identifies follow-up work | Creates new tasks in TASKS.md with `#added` timestamp |
| Finishes a phase | Saves a session snapshot to `.context/sessions/` |
| Hits an unexpected failure | Documents the root cause before moving on |

### Example: What Proactive Persistence Looks Like

Here is what the agent does during a single iteration, beyond just writing
code:

```text
Iteration 4:
  1. Read TASKS.md → pick "Add email verification" (#priority:high)
  2. Add #started:2026-01-25-030012 to the task
  3. Implement the feature (code, tests, etc.)
  4. Tests pass → mark task [x], add #done:2026-01-25-031544
  5. Add learning: "SMTP config must be set before verification
     handler registers — order matters in init()"
  6. Add decision: "Use token-based verification links (not codes)
     because links work better in automated test environments"
  7. Notice rate limiting is needed → create new task:
     "- [ ] Add rate limiting to verification endpoint #added:2026-01-25-031602 #priority:medium"
  8. Commit all changes INCLUDING .context/ files
  9. No signal emitted → loop continues to iteration 5
```

Steps 2, 4, 5, 6, and 7 are **proactive context persistence**. The agent
was not asked to do any of them. It does them because the playbook's
Work-Reflect-Persist cycle is internalized: after completing meaningful
work, reflect on what happened, then persist before moving on.

### Session Snapshots at Milestones

For long autonomous runs, the agent saves session snapshots at natural
boundaries — typically at phase transitions or after completing a cluster
of related tasks. These snapshots give you a narrative of what happened
overnight, not just a list of commits:

```text
.context/sessions/
  2026-01-25-020000-phase1-foundation.md
  2026-01-25-040000-phase2-core-features.md
  2026-01-25-060000-phase3-hardening.md
```

Each snapshot summarizes what was accomplished, what decisions were made,
and what the agent plans to do next. If the loop crashes at 4 AM, the
snapshot from 4 AM tells the next session (or you) exactly where to resume.

### The Persistence Contract

The autonomous loop has an implicit contract:

1. **Every iteration reads context** — TASKS.md, DECISIONS.md, LEARNINGS.md
2. **Every iteration writes context** — task updates, new learnings, decisions
3. **Every commit includes `.context/`** — so the next iteration sees changes
4. **Context is always current** — if the loop stopped right now, nothing is lost

Break any part of this contract and the loop degrades: iterations repeat
work, contradict earlier decisions, or lose track of what's done. The
agent's proactive discipline is what holds the loop together.

## Tips

- **Start with a small iteration cap.** Use `--max-iterations 5` for
  your first run to verify the loop behaves correctly before leaving
  it unattended.

- **Keep tasks atomic.** Each task should be completable in a single
  iteration. "Build the entire authentication system" is too broad;
  break it into registration, login, password reset, etc.

- **Use CONSTITUTION.md for guardrails.** Add rules like "never delete
  production data" or "always run tests before committing" to prevent
  the agent from making dangerous mistakes at 3 AM.

- **Check for signal discipline.** If the loop runs forever, the agent
  is not emitting `SYSTEM_CONVERGED` or `SYSTEM_BLOCKED`. Add explicit
  instructions to PROMPT.md reminding it to signal after every task.

- **Commit after context updates.** The order matters: complete the
  coding work, update context files (`ctx complete`, `ctx add`),
  commit everything including `.context/`, then signal. If context
  updates are not committed, the next iteration loses them.

- **Use `/ctx-context-monitor` for long sessions.** In Claude Code,
  the context checkpoint hook fires automatically and alerts you when
  context capacity is running low, so the agent can save its work
  before hitting limits.

## See Also

- [Autonomous Loops](../autonomous-loop.md): Full documentation of
  the loop pattern, PROMPT.md templates, and troubleshooting
- [CLI Reference: ctx loop](../cli-reference.md#ctx-loop): Command
  flags and options
- [CLI Reference: ctx watch](../cli-reference.md#ctx-watch): Watch
  mode details
- [CLI Reference: ctx init](../cli-reference.md#ctx-init): Init flags
  including `--ralph` for unattended operation
- [The Complete Session](session-lifecycle.md): Interactive workflow
  (the human-attended counterpart)
- [Tracking Work Across Sessions](task-management.md): How to
  structure TASKS.md effectively

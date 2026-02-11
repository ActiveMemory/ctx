---
title: Running an Unattended AI Agent
icon: lucide/repeat
---

![ctx](../images/ctx-banner.png)

## Problem

You have a project with a clear list of tasks, and you want an AI agent to work
through them autonomously: overnight, unattended, without you sitting at the
keyboard.

Each iteration needs to remember what the previous one did, mark tasks as
completed, and know when to stop.

Without persistent memory, every iteration starts fresh and the loop collapses.
With `ctx`, each iteration can pick up where the last one left off, but only if
the agent persists its context as part of the work.

This is the key insight: unattended operation works because the agent treats
context persistence as a first-class deliverable, not an afterthought.

## Commands and Skills Used

| Tool                    | Type    | Purpose                                                            |
|-------------------------|---------|--------------------------------------------------------------------|
| `ctx init --ralph`      | Command | Initialize project for unattended operation (no human in the loop) |
| `ctx loop`              | Command | Generate the loop shell script                                     |
| `ctx watch`             | Command | Monitor AI output and persist context updates                      |
| `ctx load`              | Command | Display assembled context (for debugging)                          |
| `/ctx-loop`             | Skill   | Generate loop script from inside Claude Code                       |
| `/ctx-implement`        | Skill   | Execute a plan step-by-step with verification                      |
| `/ctx-context-monitor`  | Skill   | Automated context capacity alerts during long sessions             |

## The Workflow

### Step 1: Initialize for Unattended Operation

Start by creating a `.context/` directory configured so the agent can work
without human input. The `--ralph` flag sets up `PROMPT.md` so the agent makes
reasonable choices instead of asking clarifying questions.

```bash
ctx init --ralph
````

This creates `.context/` with the template files, a `PROMPT.md` configured for
autonomous iteration, and Claude Code hooks and skills under `.claude/`.

Without `--ralph`, the agent will often pause when requirements are unclear.
For unattended runs, you want it to choose a default and document the trade-off
in `DECISIONS.md` instead.

### Step 2: Populate TASKS.md with Phased Work

Open `.context/TASKS.md` and organize your work into phases. The agent works
through these systematically, top to bottom, using priority tags to break ties.

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

Phased organization matters because it gives the agent natural boundaries.
Phase 1 tasks should be completable without Phase 2 code existing yet.

### Step 3: Configure PROMPT.md

The `--ralph` flag generates a `PROMPT.md` that instructs the agent to operate
autonomously:

1. Read `.context/CONSTITUTION.md` first (hard rules, never violated)
2. Load context from `.context/` files
3. Pick one task per iteration
4. Complete the task and update context files
5. Commit changes (including `.context/`)
6. Signal status with a completion signal

You can customize `PROMPT.md` for your project. The critical parts are the
one-task-per-iteration discipline, proactive context persistence, and completion
signals at the end:

```markdown
## Signal Status

End your response with exactly ONE of:

- `SYSTEM_CONVERGED` — All tasks in TASKS.md are complete (this is the
  signal the loop script detects by default)
- `SYSTEM_BLOCKED` — Cannot proceed, need human input (explain why)
- (no signal) — More work remains, continue to next iteration

Note: the loop script only checks for `SYSTEM_CONVERGED` by default.
`SYSTEM_BLOCKED` is a convention for the human reviewing the log.
```

### Step 4: Configure Permissions

An unattended agent needs permission to use tools without prompting. By default,
Claude Code asks for confirmation on file writes, bash commands, and other
operations, which stops the loop and waits for a human who is not there.

There are two approaches.

#### Option A: Explicit Allowlist (*Recommended*)

Grant only the permissions the agent needs. In `.claude/settings.local.json`:

```json
{
  "permissions": {
    "allow": [
      "Bash(make:*)",
      "Bash(go:*)",
      "Bash(git:*)",
      "Bash(ctx:*)",
      "Read",
      "Write",
      "Edit"
    ]
  }
}
```

Adjust the `Bash` patterns for your project's toolchain. The agent can run
`make`, `go`, `git`, and `ctx` commands but cannot run arbitrary shell commands.

This is recommended even in sandboxed environments because it limits blast
radius.

#### Option B: Skip All Permission Checks

Claude Code supports a `--dangerously-skip-permissions` flag that disables all
permission prompts:

```bash
claude --dangerously-skip-permissions -p "$(cat .context/PROMPT.md)"
```

!!! danger "This flag means what it says"
    With `--dangerously-skip-permissions`, the agent can execute any shell
    command, write to any file, and make network requests without
    confirmation.

    Only use this on a sandboxed machine: ideally a virtual machine with
    no access to host credentials, no SSH keys, and no access to
    production systems.

    If you would not give an untrusted intern `sudo` on this machine, do
    not use this flag.

#### Enforce Isolation at the OS Level

!!! danger "Do not skip this section"
    This is not optional hardening. An unattended agent with unrestricted
    OS access is an unattended shell with unrestricted OS access. The
    allowlist above is a strong first layer, but do not rely on a single
    runtime boundary.

The only controls an agent cannot override are the ones enforced by the
operating system, the container runtime, or the hypervisor.

For unattended runs, enforce isolation at the infrastructure level:

| Layer             | What to enforce                                                                                                                                                                                                                                               |
|-------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| User account      | Run the agent as a dedicated unprivileged user with no `sudo` access and no membership in privileged groups (`docker`, `wheel`, `adm`).                                                                                                                       |
| Filesystem        | Restrict the project directory via POSIX permissions or ACLs. The agent should have no access to other users' files or system directories.                                                                                                                    |
| Container         | Run inside a Docker/Podman sandbox. Mount only the project directory. Drop capabilities (`--cap-drop=ALL`). Disable network if not needed (`--network=none`). Never mount the Docker socket and do not run privileged containers. Prefer rootless containers. |
| Virtual machine   | Prefer a dedicated VM with no shared folders, no host passthrough, and no keys to other machines.                                                                                                                                                             |
| Network           | If the agent does not need the internet, disable outbound access entirely. If it does, restrict to specific domains via firewall rules.                                                                                                                       |
| Resource limits   | Apply CPU, memory, and disk limits (cgroups/container limits). A runaway loop should not fill disk or consume all RAM.                                                                                                                                        |
| Self-modification | Make instruction files read-only. `CLAUDE.md`, `.claude/settings.local.json`, `.claude/hooks/`, and `.context/CONSTITUTION.md` should not be writable by the agent user.                                                                                      |

A minimal Docker setup for overnight runs:

```bash
docker run --rm \
  --network=none \
  --cap-drop=ALL \
  --memory=4g \
  --cpus=2 \
  -v /path/to/project:/workspace \
  -w /workspace \
  your-dev-image \
  ./loop.sh 2>&1 | tee /tmp/loop.log
```

!!! tip "Defense in depth"
    Use multiple layers together: OS-level isolation (*the boundary the
    agent cannot cross*), a permission allowlist (*what Claude Code will do
    within that boundary*), and `CONSTITUTION.md` (*a soft nudge for the
    common case*).

### Step 5: Generate the Loop Script

Use `ctx loop` to generate a `loop.sh` tailored to your AI tool:

```bash
# Generate for Claude Code with a 10-iteration cap
ctx loop --tool claude --max-iterations 10

# Generate for Aider
ctx loop --tool aider --max-iterations 10

# Custom prompt file and output filename (default prompt is PROMPT.md in project root)
ctx loop --tool claude --prompt my-prompt.md --output my-loop.sh
```

The generated script reads `PROMPT.md`, runs the tool, checks for completion
signals, and loops until done or the cap is reached.

You can also use the `/ctx-loop` skill from inside Claude Code.

!!! tip "A shell loop is best practice"
    The shell loop approach spawns a fresh AI process each iteration, so
    the only state that carries between iterations is what lives in
    `.context/` and git.

    Claude Code's built-in `/loop` runs iterations within the same
    session, which can allow context window state to leak between
    iterations. This can be convenient for short runs, but it is less
    reliable for unattended loops. 

    See [Shell Loop vs Built-in
    Loop](../autonomous-loop.md#quick-start-shell-while-loop-recommended)
    for details.

### Step 6: Run with Watch Mode

Open two terminals. In the first, run the loop. In the second, run `ctx watch`
to process context updates from the AI output.

```bash
# Terminal 1: Run the loop
./loop.sh 2>&1 | tee /tmp/loop.log

# Terminal 2: Watch for context updates
ctx watch --log /tmp/loop.log
```

The watch command parses XML context-update commands from the AI output and
applies them:

```xml
<context-update type="complete">user registration</context-update>
<context-update type="learning">Email verification needs SMTP configured</context-update>
```

### Step 7: Completion Signals End the Loop

The generated script checks for **one** completion signal per run. By
default this is `SYSTEM_CONVERGED`. You can change it with the
`--completion` flag:

```bash
ctx loop --tool claude --completion BOOTSTRAP_COMPLETE --max-iterations 5
```

The following signals are conventions used in `PROMPT.md`:

| Signal               | Convention                     | How the script handles it                          |
|----------------------|--------------------------------|----------------------------------------------------|
| `SYSTEM_CONVERGED`   | All tasks in TASKS.md are done | Detected by default (`--completion` default value) |
| `SYSTEM_BLOCKED`     | Agent cannot proceed           | Only detected if you set `--completion` to this    |
| `BOOTSTRAP_COMPLETE` | Initial scaffolding done       | Only detected if you set `--completion` to this    |

The script uses `grep -q` on the agent's output, so any string works as a
signal. If you need to detect multiple signals in one run, edit the
generated `loop.sh` to add additional `grep` checks.

When you return in the morning, check the log and the context files:

```bash
tail -100 /tmp/loop.log
ctx status
ctx load
```

### Step 8: Use `/ctx-implement` for Plan Execution

Within each iteration, the agent can use `/ctx-implement` to execute multi-step
plans with verification between steps. This is useful for complex tasks that
touch multiple files.

The skill breaks a plan into atomic, verifiable steps:

```text
Step 1/6: Create user model .................. OK
Step 2/6: Add database migration ............. OK
Step 3/6: Implement registration handler ..... OK
Step 4/6: Write unit tests ................... OK
Step 5/6: Run test suite ..................... FAIL
  -> Fixed: missing test dependency
  -> Re-verify ................................ OK
Step 6/6: Update TASKS.md .................... OK
```

Each step is verified (build, test, syntax check) before moving to the next.

## Putting It All Together

A typical overnight run:

```bash
ctx init --ralph
# Edit TASKS.md and PROMPT.md

ctx loop --tool claude --max-iterations 20

./loop.sh 2>&1 | tee /tmp/loop.log &
ctx watch --log /tmp/loop.log

# Next morning:
ctx status
ctx load
```

## Why Autonomous Loops Work: Proactive Context Persistence

The autonomous loop pattern works because the agent persists context as part of
the job.

Without proactive persistence, the loop degrades into disconnected iterations
that repeat work, contradict decisions, and lose track of progress. The agent
cannot rely on a human to prompt it. It must treat context updates as part of
every task, not as a separate step.

### The Agent Playbook Self-Check Drives Autonomy

The Agent Playbook defines a self-check prompt the agent applies continuously:

> "If this session ended right now, would the next session know what happened?"

In interactive sessions, a human might remind the agent to save context. When
the agent runs unattended, nobody is there to remind it.

The self-check becomes the trigger: after meaningful work, the agent evaluates
whether the context files reflect reality and updates them immediately if they
do not.

### What the Agent Does Proactively Between Iterations

At milestones within an iteration, the agent persists without waiting for
instructions:

| Milestone                  | What the agent does proactively                       |
|----------------------------|-------------------------------------------------------|
| Completes a task           | Marks it `[x]` in TASKS.md with `#done` timestamp     |
| Discovers a gotcha         | Adds it to LEARNINGS.md                               |
| Makes a design choice      | Records it in DECISIONS.md with rationale             |
| Identifies follow-up work  | Creates new tasks in TASKS.md with `#added` timestamp |
| Hits an unexpected failure | Documents the root cause before moving on             |

### Example: What Proactive Persistence Looks Like

Here is what the agent does during a single iteration beyond writing code:

```text
Iteration 4:
  1. Read TASKS.md -> pick "Add email verification" (#priority:high)
  2. Add #started:2026-01-25-030012 to the task
  3. Implement the feature (code, tests, docs if needed)
  4. Tests pass -> mark task [x], add #done:2026-01-25-031544
  5. Add learning: "SMTP config must be set before verification handler registers. Order matters in init()."
  6. Add decision: "Use token-based verification links (not codes) because links work better in automated tests."
  7. Create follow-up task: "Add rate limiting to verification endpoint" #added:...
  8. Commit all changes including `.context/`
  9. No signal emitted -> loop continues to iteration 5
```

Steps 2, 4, 5, 6, and 7 are proactive context persistence. The agent was not
asked to do any of them.

### Context Persistence at Milestones

For long autonomous runs, the agent persists context at natural boundaries,
often at phase transitions or after completing a cluster of related tasks.
It updates TASKS.md, DECISIONS.md, and LEARNINGS.md as it goes.

If the loop crashes at 4 AM, the context files tell you exactly where to
resume. You can also use `ctx recall list` to review the session transcripts.

### The Persistence Contract

The autonomous loop has an implicit contract:

1. Every iteration reads context: `TASKS.md`, `DECISIONS.md`, `LEARNINGS.md`
2. Every iteration writes context: task updates, new learnings, decisions
3. Every commit includes `.context/` so the next iteration sees changes
4. Context stays current: if the loop stopped right now, nothing important is lost

Break any part of this contract and the loop degrades.

## Tips

!!! warning "Markdown is not enforcement"
    Your real guardrails are permissions and isolation, not markdown.
    `CONSTITUTION.md` can nudge the agent, but it is probabilistic. The
    permission allowlist and OS isolation are deterministic. For
    unattended runs, trust the sandbox and the allowlist, not the prose.

* Start with a small iteration cap. Use `--max-iterations 5` on your first run.
* Keep tasks atomic. Each task should be completable in a single iteration.
* Check signal discipline. If the loop runs forever, the agent is not emitting
  `SYSTEM_CONVERGED` or `SYSTEM_BLOCKED`. Make the signal requirement explicit
  in `PROMPT.md`.
* Commit after context updates. Finish code, update `.context/`, commit including
  `.context/`, then signal.
* Use `/ctx-context-monitor` for long runs. It can warn when context capacity is
  running low so the agent saves before hitting limits.

## Next Up

**[Turning Activity into Content](publishing.md)**:
Generate blog posts and changelogs from your project activity.

## See Also

* [Autonomous Loops](../autonomous-loop.md): loop pattern, PROMPT.md templates, troubleshooting
* [CLI Reference: ctx loop](../cli-reference.md#ctx-loop): flags and options
* [CLI Reference: ctx watch](../cli-reference.md#ctx-watch): watch mode details
* [CLI Reference: ctx init](../cli-reference.md#ctx-init): init flags including `--ralph`
* [The Complete Session](session-lifecycle.md): interactive workflow
* [Tracking Work Across Sessions](task-management.md): structuring TASKS.md

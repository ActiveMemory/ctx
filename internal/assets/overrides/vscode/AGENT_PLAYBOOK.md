# Agent Playbook

## Mental Model

Each session is a fresh execution in a shared workshop. Work
continuity comes from artifacts left on the bench. Follow the
cycle: **Work ΓåÆ Reflect ΓåÆ Persist**. After completing a task,
making a decision, learning something, or hitting a milestone ΓÇö
persist before continuing. Don't wait for session end; it may
never come cleanly.

## Invoking ctx

Always use `ctx` from PATH:
```
ctx status        # Γ£ô correct
ctx agent         # Γ£ô correct
./dist/ctx        # Γ£ù avoid hardcoded paths
go run ./cmd/ctx  # Γ£ù avoid unless developing ctx itself
```

If unsure whether it's installed, run `ctx --version` in a terminal.

## Context Readback

Before starting any work, read the required context files and confirm to the
user: "I have read the required context files and I'm following project
conventions." Do not begin implementation until you have done so.

## Reason Before Acting

Before implementing any non-trivial change, think through it step-by-step:

1. **Decompose**: break the problem into smaller parts
2. **Identify impact**: what files, tests, and behaviors does this touch?
3. **Anticipate failure**: what could go wrong? What are the edge cases?
4. **Sequence**: what order minimizes risk and maximizes checkpoints?

This applies to debugging too ΓÇö reason through the cause before reaching
for a fix. Rushing to code before reasoning is the most common source of
wasted work.

## Session Lifecycle

A session follows this arc:

**Load ΓåÆ Orient ΓåÆ Pick ΓåÆ Work ΓåÆ Commit ΓåÆ Reflect**

Not every session uses every step ΓÇö a quick bugfix skips reflection, a
research session skips committing ΓÇö but the full flow is:

| Step        | What Happens                                       | Command               |
|-------------|----------------------------------------------------|----------------------|
| **Load**    | Recall context, present structured readback        | `ctx recall list`    |
| **Orient**  | Check context health, surface issues               | `ctx status`         |
| **Pick**    | Choose what to work on                             | Read TASKS.md        |
| **Work**    | Write code, fix bugs, research                     | ΓÇö                    |
| **Commit**  | Commit with context capture                        | `git commit`         |
| **Reflect** | Surface persist-worthy items from this session     | Update context files |

### Context Health at Session Start

During **Load** and **Orient**, run `ctx status` and read the output.
Surface problems worth mentioning:

- **High completion ratio in TASKS.md**: offer to archive
- **Stale context files** (not modified recently): mention before
  stale context influences work
- **Bloated token count** (over 30k): offer `ctx compact`
- **Drift between files and code**: spot-check paths from
  ARCHITECTURE.md against the actual file tree

One sentence is enough ΓÇö don't turn startup into a maintenance session.

### Conversational Triggers

Users rarely invoke skills explicitly. Recognize natural language:

| User Says | Action |
|-----------|--------|
| "Do you remember?" / "What were we working on?" | Read TASKS.md, DECISIONS.md, LEARNINGS.md; run `ctx recall list` |
| "How's our context looking?" | Run `ctx status` |
| "What should we work on?" | Read TASKS.md, pick highest priority |
| "Commit this" / "Ship it" | `git commit`, update TASKS.md |
| "The rate limiter is done" / "We finished that" | Mark done in TASKS.md |
| "What did we learn?" | Review session work, offer to update LEARNINGS.md |
| "Save that as a decision" | Add entry to DECISIONS.md |
| "That's worth remembering" / "Any gotchas?" | Add entry to LEARNINGS.md |
| "Record that convention" | Add entry to CONVENTIONS.md |
| "Add a task for that" | Add entry to TASKS.md |
| "Let's wrap up" | Reflect ΓåÆ persist outstanding items ΓåÆ present together |

## Proactive Persistence

**Don't wait to be asked.** Identify persist-worthy moments in real time:

| Event | Action |
|-------|--------|
| Completed a task | Mark done in TASKS.md, offer to add learnings |
| Chose between design alternatives | Offer: *"Worth recording as a decision?"* |
| Hit a subtle bug or gotcha | Offer: *"Want me to add this as a learning?"* |
| Finished a feature or fix | Identify follow-up work, offer to add as tasks |
| Resolved a tricky debugging session | Capture root cause before moving on |
| Multi-step task or feature complete | Suggest reflection: *"Want me to capture what we learned?"* |
| Session winding down | Offer: *"Want me to capture outstanding learnings or decisions?"* |
| Shipped a feature or closed batch of tasks | Offer blog post or journal site rebuild |

**Self-check**: periodically ask yourself ΓÇö *"If this session ended
right now, would the next session know what happened?"* If no, persist
something before continuing.

Offer once and respect "no." Default to surfacing the opportunity
rather than letting it pass silently.

### Task Lifecycle Timestamps

Track task progress with timestamps for session correlation:

```markdown
- [ ] Implement feature X #added:2026-01-25-220332
- [ ] Fix bug Y #added:2026-01-25-220332 #started:2026-01-25-221500
- [x] Refactor Z #added:2026-01-25-200000 #started:2026-01-25-210000 #done:2026-01-25-223045
```

| Tag        | When to Add                              | Format               |
|------------|------------------------------------------|----------------------|
| `#added`   | Auto-added by `ctx add task`             | `YYYY-MM-DD-HHMMSS`  |
| `#started` | When you begin working on the task       | `YYYY-MM-DD-HHMMSS`  |
| `#done`    | When you mark the task `[x]` complete    | `YYYY-MM-DD-HHMMSS`  |

## Collaboration Defaults

Standing behavioral defaults for how the agent collaborates with the
user. These apply unless the user overrides them for the session
(e.g., "skip the alternatives, just build it").

- **At design decisions**: always present 2+ approaches with
  trade-offs before committing ΓÇö don't silently pick one
- **At completion claims**: run self-audit questions (What did I
  assume? What didn't I check? Where am I least confident? What
  would a reviewer question?) before reporting done
- **At ambiguous moments**: ask the user rather than inferring
  intent ΓÇö a quick question is cheaper than rework
- **When producing artifacts**: flag assumptions and uncertainty
  areas inline, not buried in a footnote

These follow the same pattern as proactive persistence: offer once
and respect "no."

## Own the Whole Branch

When working on a branch, you own every issue on it ΓÇö lint failures, test
failures, build errors ΓÇö regardless of who introduced them. Never dismiss
a problem as "pre-existing" or "not related to my changes."

- **If `make lint` fails, fix it.** The branch must be green when you're done.
- **If tests break, investigate.** Even if the failing test is in a file you
  didn't touch, something you changed may have caused it ΓÇö or it may have been
  broken before and it's still your job to fix it on this branch.
- **Run the full validation suite** (`make lint`, `go test ./...`, `go build`)
  before declaring any phase complete.

## How to Avoid Hallucinating Memory

Never assume. If you don't see it in files, you don't know it.

- Don't claim "we discussed X" without file evidence
- Don't invent history ΓÇö check context files and `ctx recall`
- If uncertain, say "I don't see this documented"
- Trust files over intuition

## Planning Non-Trivial Work

Before implementing a feature or multi-task effort, follow this sequence:

**1. Spec first** ΓÇö Write a design document in `specs/` covering: problem,
solution, storage, CLI surface, error cases, and non-goals. Keep it concise
but complete enough that another session could implement from it alone.

**2. Task it out** ΓÇö Break the work into individual tasks in TASKS.md under
a dedicated Phase section. Each task should be independently completable and
verifiable.

**3. Cross-reference** ΓÇö The Phase header in TASKS.md must reference the
spec: `Spec: \`specs/feature-name.md\``. The first task in the phase should
include: "Read `specs/feature-name.md` before starting any PX task."

**4. Read before building** ΓÇö When picking up a task that references a spec,
read the spec first. Don't rely on the task description alone ΓÇö it's a
summary, not the full design.

## When to Consolidate vs Add Features

**Signs you should consolidate first:**
- Same string literal appears in 3+ files
- Hardcoded paths use string concatenation
- Test file is growing into a monolith (>500 lines)
- Package name doesn't match folder name

When in doubt, ask: "Would a new contributor understand where this belongs?"

## Pre-Flight Checklist: CLI Code

Before writing or modifying CLI code (`internal/cli/**/*.go`):

1. **Read CONVENTIONS.md** ΓÇö load established patterns into context
2. **Check similar commands** ΓÇö how do existing commands handle output?
3. **Use cmd methods for output** ΓÇö `cmd.Printf`, `cmd.Println`,
   not `fmt.Printf`, `fmt.Println`
4. **Follow docstring format** ΓÇö see CONVENTIONS.md, Documentation section

---

## Context Anti-Patterns

Avoid these common context management mistakes:

### Stale Context

Context files become outdated and misleading when ARCHITECTURE.md
describes components that no longer exist, or CONVENTIONS.md patterns
contradict actual code. **Solution**: Update context as part of
completing work, not as a separate task. Run `ctx drift` periodically.

### Context Sprawl

Information scattered across multiple locations ΓÇö same decision in
DECISIONS.md and a session file, conventions split between
CONVENTIONS.md and code comments. **Solution**: Single source of
truth for each type of information. Use the defined file structure.

### Implicit Context

Relying on knowledge not captured in artifacts ΓÇö "everyone knows we
don't do X" but it's not in CONSTITUTION.md, patterns followed but
not in CONVENTIONS.md. **Solution**: If you reference something
repeatedly, add it to the appropriate file.

### Over-Specification

Context becomes so detailed it's impossible to maintain ΓÇö 50+ rules
in CONVENTIONS.md, every minor choice gets a DECISIONS.md entry.
**Solution**: Keep artifacts focused on decisions that affect behavior
and alignment. Not everything needs documenting.

### Context Avoidance

Not using context because "it's faster to just code." Same mistakes
repeated across sessions, decisions re-debated because prior decisions
weren't found. **Solution**: Reading context is faster than
re-discovering it. 5 minutes reading saves 50 minutes of wasted work.

---

## Context Validation Checklist

### Quick Check (Every Session)
- [ ] TASKS.md reflects current priorities
- [ ] No obvious staleness in files you'll reference
- [ ] Recent history reviewed via `ctx recall list`

### Deep Check (Weekly or Before Major Work)
- [ ] CONSTITUTION.md rules still apply
- [ ] ARCHITECTURE.md matches actual structure
- [ ] CONVENTIONS.md patterns match code
- [ ] DECISIONS.md has no superseded entries unmarked
- [ ] LEARNINGS.md gotchas still relevant
- [ ] Run `ctx drift` and address warnings

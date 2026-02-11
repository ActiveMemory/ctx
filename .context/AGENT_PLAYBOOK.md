# Agent Playbook

## Mental Model

This system does not persist experience.

- Each session is a fresh execution in a shared workshop.
- Work continuity comes from artifacts left on the bench.

### Work → Reflect → Persist

After completing meaningful work, follow this cycle:

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│  Work   │ ──▶ │ Reflect │ ──▶ │ Persist │ ──▶ (continue)
└─────────┘     └─────────┘     └─────────┘
                     │
                     ▼
              Did I just...
              • Complete a task?
              • Make a decision?
              • Learn something?
              • Hit a milestone?
```

**Don't wait for session end** — it may never come cleanly. Persist as you go.

## Invoking ctx

Always use `ctx` from PATH:
```bash
ctx status        # ✓ correct
ctx agent         # ✓ correct
./dist/ctx        # ✗ avoid hardcoded paths
go run ./cmd/ctx  # ✗ avoid unless developing ctx itself
```

Check with `which ctx` if unsure whether it's installed.

## Context Readback

Before starting any work, read the required context files and confirm to the
user: "I have read the required context files and I'm following project
conventions." Do not begin implementation until you have done so.

## Read Order

1. CONSTITUTION.md — Hard rules, NEVER violate
2. TASKS.md — What to work on next
3. CONVENTIONS.md — How to write code
4. ARCHITECTURE.md — Where things go
5. DECISIONS.md — Why things are the way they are
6. LEARNINGS.md — Gotchas to avoid
7. GLOSSARY.md — Correct terminology

## How Your Memory Works

When asked "Do you remember?" or similar:

1. **Read silently first**: TASKS.md, DECISIONS.md, LEARNINGS.md, and
   run `ctx recall list --limit 5` for recent history. Do this BEFORE composing a response.
2. **Respond with a structured readback**:
   - **Last session**: most recent session topic and date
   - **Active work**: pending or in-progress tasks
   - **Recent context**: 1-2 recent decisions or learnings
   - **Next step**: offer to continue or ask what to focus on
3. **Be honest about mechanism only if asked** — the user wants results,
   not disclaimers

**Wrong**: "I don't have memory, but I can check files. Want me to?"
**Right**: *reads files silently* → "Last session (Feb 7, architecture docs):
we populated ARCHITECTURE.md and fixed task insertion. You have 5 pending
tasks in Phase 2. Want to continue with those?"

The context IS your memory. It's stored externally in files, but that's an
implementation detail. Load it and use it — don't lead with caveats.

## Session Lifecycle

A session follows this arc:

**Load → Orient → Pick → Work → Commit → Reflect**

Not every session uses every step — a quick bugfix skips reflection, a
research session skips committing — but the full flow is:

| Step        | What Happens                                       | Skill / Command  |
|-------------|----------------------------------------------------|------------------|
| **Load**    | Recall context, present structured readback        | `/ctx-remember`  |
| **Orient**  | Check context health, surface issues                | `/ctx-status`    |
| **Pick**    | Choose what to work on                             | `/ctx-next`      |
| **Work**    | Write code, fix bugs, research                     | `/ctx-implement` |
| **Commit**  | Commit with context capture                        | `/ctx-commit`    |
| **Reflect** | Surface persist-worthy items from this session     | `/ctx-reflect`   |

### Context Health at Session Start

During **Load** and **Orient**, run `ctx status` and read the output.
It shows task counts, token usage, and modification times. Use this
to notice problems worth mentioning:

- **High completion ratio in TASKS.md** (many completed items):
  *"TASKS.md has 18 completed items. Want me to archive them?"*
- **Stale context files** (not modified recently — `ctx status`
  shows "modified X ago"): *"DECISIONS.md hasn't been updated in
  40 days. Worth reviewing?"*
- **Bloated token count** (over 30k): *"Context is getting heavy.
  Want me to run `ctx compact` to trim it?"*
- **Drift between files and code**: spot-check a few paths from
  ARCHITECTURE.md or CONVENTIONS.md against the actual file tree.
  If a referenced path is missing or empty, mention it before stale
  context influences any work.

Surface these during the readback or Orient step. One sentence is
enough — don't turn startup into a maintenance session.

### Conversational Triggers

Users rarely invoke skills explicitly. You can drive the same lifecycle
with plain language:

| Step    | Slash command       | Natural language                                        |
|---------|---------------------|---------------------------------------------------------|
| Load    | `/ctx-remember`     | "Do you remember?" / "What were we working on?"         |
| Orient  | `/ctx-status`       | "How's our context looking?"                            |
| Pick    | `/ctx-next`         | "What should we work on?" / "Let's do the caching task" |
| Work    | --                  | "Only change files in internal/cache/"                  |
| Commit  | `/ctx-commit`       | "Commit this" / "Ship it"                               |
| Done    | `ctx complete`      | "The rate limiter is done" / "We finished that"          |
| Reflect | `/ctx-reflect`      | "What did we learn?" / *(agent offers at milestones)*   |

Users also trigger persist actions directly with natural language:

| User Says                                         | Action                                  |
|---------------------------------------------------|-----------------------------------------|
| "Save that as a decision"                          | `/ctx-add-decision` with structured fields inferred from conversation |
| "That's worth remembering" / "Any gotchas?"        | `/ctx-add-learning` with context, lesson, application from session   |
| "Record that convention" / "Codify that pattern"   | `/ctx-add-convention` with section inferred from topic               |
| "Add a task for that" / "We should track that"     | `/ctx-add-task` with description from context                        |

### Chained Flows

Some cues should chain multiple steps in one pass:

- **"Let's wrap up"**: Reflect on the session → persist outstanding
  learnings and decisions. Present all results together; don't make
  the user invoke each step separately.
- **"Yes, and let's wrap up"** (after a persist offer): Record what
  the user confirmed → then trigger the full wrap-up flow.

### Proactive Behavior During Work

While working, proactively identify moments worth persisting — don't
wait to be asked:

- After choosing between design alternatives, flag it: *"That's a
  design choice worth recording. Want me to save it as a decision?"*
- After hitting a subtle bug, surface it: *"That was non-obvious.
  Want me to add it as a learning?"*
- After completing a task, note it: *"That closes the X task. I'll
  mark it done."* Also recognize natural-language completion signals
  like "the rate limiter is done" or "we finished that" — match them
  to tasks in TASKS.md and mark them done with `ctx complete`.
- After finishing a feature or fix, identify follow-up work (tests,
  docs, edge cases, cleanup) and offer to add them: *"That leaves
  a few loose ends — integration tests and updating the CLI docs.
  Want me to add those as tasks?"*

The goal is to **identify** persist-worthy moments in real time, not
just respond to them after the fact.

### Agent-Initiated Reflect

At significant milestones, proactively initiate **Reflect**
without waiting to be asked:

- **After completing a multi-step task or feature**: suggest a
  reflection — *"We just finished the caching layer. Want me to do
  a quick reflection to capture what we learned?"*
- **After a long debugging session**: offer to reflect before moving
  on — *"That took a while to track down. Worth capturing the root
  cause as a learning before we continue?"*
- **When the session feels like it's winding down** (user is doing
  small tweaks, asking fewer questions, long session): offer to
  persist — *"We've covered a lot. Want me to capture any outstanding
  learnings or decisions?"*

- **After shipping a feature or closing a batch of tasks**: suggest
  content — *"We just shipped the caching layer and closed 3 tasks.
  Want me to draft a blog post about it?"* or *"Your journal has
  new entries since last rebuild. Want me to regenerate the site?"*

Offer once and respect "no." But default to surfacing the opportunity
rather than letting it pass silently.

## When to Update Memory

| Event                       | Action                |
|-----------------------------|-----------------------|
| Made architectural decision | Add to DECISIONS.md   |
| Discovered gotcha/bug       | Add to LEARNINGS.md   |
| Established new pattern     | Add to CONVENTIONS.md |
| Completed task              | Mark [x] in TASKS.md  |

## Proactive Context Persistence

**Don't wait for session end** — persist context at natural milestones.

### Milestone Triggers

Offer to persist context when you:

| Milestone                          | Action                                          |
|------------------------------------|-------------------------------------------------|
| Complete a task                    | Mark done in TASKS.md, offer to add learnings   |
| Make an architectural decision     | `ctx add decision "..."`                        |
| Discover a gotcha or bug           | `ctx add learning "..."`                        |
| Finish a significant code change   | Offer to summarize what was done                |
| Encounter unexpected behavior      | Document it before moving on                    |
| Resolve a tricky debugging session | Capture the root cause and fix                  |

### How to Offer

After hitting a milestone, briefly offer:

> "I just completed X. Want me to capture this as a learning/decision before we continue?"

Or proactively persist and inform:

> "I've added that gotcha to LEARNINGS.md so we don't hit it again."

### Self-Check Prompt

Periodically ask yourself:

> "If this session ended right now, would the next session know what happened?"

If no — persist something before continuing.

### Reflect Skill

Use `/ctx-reflect` to trigger a structured reflection checkpoint:
- Reviews what was accomplished in the session
- Identifies learnings, decisions, and task updates
- Offers to persist before continuing

Run this periodically during long sessions or at natural breakpoints.

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

**Why this matters:**
- Correlate tasks with session files by timestamp
- See how long tasks took (across sessions)
- Know which session started vs completed work

**Example workflow:**
1. Pick up task → add `#started:$(date +%Y-%m-%d-%H%M%S)`
2. Work on it
3. Complete → change `[ ]` to `[x]`, add `#done:$(date +%Y-%m-%d-%H%M%S)`

## How to Avoid Hallucinating Memory

Never assume. If you don't see it in files, you don't know it.

- Don't claim "we discussed X" without file evidence
- Don't invent history - check context files and `ctx recall` for actual discussions
- If uncertain, say "I don't see this documented"
- Trust files over intuition

## When to Consolidate vs Add Features

**Signs you should consolidate first:**
- Same string literal appears in 3+ files
- Hardcoded paths use string concatenation
- Test file is growing into a monolith (>500 lines)
- Package name doesn't match folder name

**YOLO mode creates debt**—rapid feature additions scatter patterns across 
the codebase. Periodic consolidation prevents this from compounding.

**Human-guided refactoring catches:**
- Magic strings that should be constants
- Path construction that should use `filepath.Join()`
- Tests that should be colocated with implementations
- Naming inconsistencies

When in doubt, ask: "Would a new contributor understand where this belongs?"

## Pre-Flight Checklist: CLI Code

Before writing or modifying CLI code (`internal/cli/**/*.go`):

1. **Read CONVENTIONS.md** — Load established patterns into context
2. **Check similar commands** — How do existing commands in the same package handle output?
3. **Use cmd methods for output** — `cmd.Printf`, `cmd.Println`, not `fmt.Printf`, `fmt.Println`
4. **Follow docstring format** — See Go Documentation Standard below

**Quick pattern check:**
```bash
# See how other commands do output
grep -n "cmd.Printf\|cmd.Println" internal/cli/status/*.go

# Spot violations in your changes
grep -n "fmt.Printf\|fmt.Println" internal/cli/yourpackage/*.go
```

## Go Documentation Standard

When writing Go code, follow this docstring format consistently.

### Functions

```go
// FunctionName does X.
//
// Extended description if needed.
//
// Parameters:
//   - param1: Description of first parameter
//   - param2: Description of second parameter
//
// Returns:
//   - ReturnType: Description of return value
//   - error: When this error occurs
func FunctionName(param1 Type1, param2 Type2) (ReturnType, error) {
```

### Structs

```go
// StructName represents X.
//
// Extended description if needed.
//
// Fields:
//   - Field1: Description of field
//   - Field2: Description of field
type StructName struct {
    Field1 Type1
    Field2 Type2
}
```

### Key Points

- **Always include Parameters section** if function has parameters
- **Always include Returns section** if function returns values
- **Always include Fields section** for exported structs
- **No inline field comments** — put all field docs in the Fields block
- Check existing code for reference before writing new documentation

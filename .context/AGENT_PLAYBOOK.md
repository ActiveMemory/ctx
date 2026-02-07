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

1. **Check the context files** (sessions/, TASKS.md, LEARNINGS.md, etc.)
2. **Respond with what you found** — "From the sessions, I see we discussed X..."
3. **Be honest about mechanism only if asked** — the user wants results, not disclaimers

**Wrong**: "I don't have memory, but I can check files. Want me to?"
**Right**: *checks files* → "Yes, I see from the last session we were working on..."

The context IS your memory. It's stored externally in files, but that's an
implementation detail. Load it and use it — don't lead with caveats.

## Session History

**IMPORTANT**: Check `.context/sessions/` for session dumps
from previous sessions.

If you're confused about context or need a deep dive into past discussions:
```
ls .context/sessions/
```

**Manual session files** are named `YYYY-MM-DD-HHMMSS-<topic>.md` 
(e.g., `2026-01-15-164600-feature-discussion.md`). 
These are updated throughout the session.

**Auto-snapshot files** are named `YYYY-MM-DD-HHMMSS-<event>.jsonl` 
(e.g., `2026-01-15-170830-pre-compact.jsonl`). These are immutable once created.

**Auto-save triggers** (for Claude Code users):
- **SessionEnd hook** → auto-saves transcript on exit, including Ctrl+C
- **PreCompact** → saves before `ctx compact` archives old tasks
- **Manual** → `ctx session save`

See `.claude/hooks/auto-save-session.sh` for the implementation.

## Timestamp-Based Session Correlation

Context entries (tasks, learnings, decisions) include timestamps that allow
you to determine which session created them.

### Timestamp Format

All timestamps use `YYYY-MM-DD-HHMMSS` format (6-digit time for seconds precision):
- **Tasks**: `- [ ] Do something #added:2026-01-23-143022`
- **Learnings**: `- **[2026-01-23-143022]** Discovered that...`
- **Decisions**: `## [2026-01-23-143022] Use PostgreSQL`
- **Sessions**: `**start_time**: 2026-01-23-140000` / `**end_time**: 2026-01-23-153045`

### Correlating Entries to Sessions

To find which session added an entry:

1. **Extract the entry's timestamp** (e.g., `2026-01-23-143022`)
2. **List sessions** from that day: `ls .context/sessions/2026-01-23*`
3. **Check session time bounds**: Entry timestamp should fall between session's
   start_time and end_time
4. **Match**: The session file with matching time range contains the context

### When Timestamps Help

- **Tracing decisions**: "Why did we decide X?" → Find the session that added it
- **Understanding context**: Read the full session for the discussion that led to an entry
- **Debugging issues**: Correlate when a learning was discovered with what was happening

## Session File Structure (Suggested)

Adapt this structure based on session type. Not all sections are needed for every session.

### Core Sections (Always Include)
```markdown
# Session: <Topic>

**Date**: YYYY-MM-DD
**Topic**: Brief description
**Type**: feature | bugfix | architecture | exploration | planning

---

## Summary
What was discussed/accomplished (2-3 sentences)

## Key Decisions
Bullet points of decisions made (if any)

## Tasks for Next Session
What to pick up next
```

### Context-Dependent Sections

| Session Type              | Additional Sections                               |
|---------------------------|---------------------------------------------------|
| **Feature discussion**    | Requirements, Design options, Implementation plan |
| **Bug investigation**     | Symptoms, Root cause, Fix applied, Prevention     |
| **Architecture decision** | Context, Options considered, Trade-offs, Decision |
| **Exploration/Research**  | Questions, Findings, Open questions               |
| **Planning**              | Goals, Milestones, Dependencies, Risks            |
| **Quick fix**             | Problem, Solution, Files changed (minimal format) |

### When to Go Minimal

For quick sessions (<15 min), just capture:
```markdown
# Session: Quick Fix - <Topic>
**Date**: YYYY-MM-DD
**Summary**: One sentence
**Files changed**: List
```

### When to Go Deep

For complex sessions (architecture, debugging), include:
- User quotes that capture key insights
- Technical context (platform, versions, constraints)
- Links to related sessions or decisions
- Code snippets or error messages if relevant

---

## When to Update Memory

| Event                       | Action                |
|-----------------------------|-----------------------|
| Made architectural decision | Add to DECISIONS.md   |
| Discovered gotcha/bug       | Add to LEARNINGS.md   |
| Established new pattern     | Add to CONVENTIONS.md |
| Completed task              | Mark [x] in TASKS.md  |
| Had important discussion    | Save to sessions/     |

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

### Session Saves

For longer sessions with substantial work, offer to save a session summary:
1. Curated summary → LEARNINGS.md, DECISIONS.md, TASKS.md
2. Full session notes → `.context/sessions/YYYY-MM-DD-<topic>.md`

## How to Avoid Hallucinating Memory

Never assume. If you don't see it in files, you don't know it.

- Don't claim "we discussed X" without file evidence
- Don't invent history - check sessions/ for actual discussions
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

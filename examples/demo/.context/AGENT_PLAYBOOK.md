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
**Right**: *reads files silently* → "Last session (Jan 15, database timeouts):
we investigated connection pool exhaustion. OAuth2 login is in progress
with 3 remaining tasks in Phase 2. Want to continue with that?"

The context IS your memory. It's stored externally in files, but that's an
implementation detail. Load it and use it — don't lead with caveats.

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

### Self-Check Prompt

Periodically ask yourself:

> "If this session ended right now, would the next session know what happened?"

If no — persist something before continuing.

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

## How to Avoid Hallucinating Memory

Never assume. If you don't see it in files, you don't know it.

- Don't claim "we discussed X" without file evidence
- Don't invent history - check context files and `ctx recall` for actual discussions
- If uncertain, say "I don't see this documented"
- Trust files over intuition

---

## Context Anti-Patterns

### Stale Context

**Problem**: Context files become outdated and misleading.

**Solution**: Update context as part of completing work, not as a separate task.
Run `ctx drift` periodically to detect staleness.

### Context Sprawl

**Problem**: Information scattered across multiple locations.

**Solution**: Single source of truth for each type of information.
Use the defined file structure; resist creating new document types.

### Implicit Context

**Problem**: Relying on knowledge not captured in artifacts.

**Solution**: If you reference something repeatedly, add it to the appropriate file.
If this session ended now, would the next session know what you know?

---

## Context Validation Checklist

Before starting significant work, validate context is current:

### Quick Check (Every Session)
- [ ] TASKS.md reflects current priorities
- [ ] No obvious staleness in files you'll reference
- [ ] Recent history reviewed for relevant context (via `ctx recall list`)

### Deep Check (Weekly or Before Major Work)
- [ ] CONSTITUTION.md rules still apply
- [ ] ARCHITECTURE.md matches actual structure
- [ ] CONVENTIONS.md patterns match code
- [ ] DECISIONS.md has no superseded entries unmarked
- [ ] LEARNINGS.md gotchas still relevant
- [ ] Run `ctx drift` and address warnings

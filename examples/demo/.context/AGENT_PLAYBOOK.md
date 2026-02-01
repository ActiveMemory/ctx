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

## Read Order

1. CONSTITUTION.md — Hard rules, NEVER violate
2. TASKS.md — What to work on next
3. CONVENTIONS.md — How to write code
4. ARCHITECTURE.md — Where things go
5. DECISIONS.md — Why things are the way they are
6. LEARNINGS.md — Gotchas to avoid

## How Your Memory Works

When asked "Do you remember?" or similar:

1. **Check the context files** (sessions/, TASKS.md, LEARNINGS.md, etc.)
2. **Respond with what you found** — "From the sessions, I see we discussed X..."
3. **Be honest about mechanism only if asked** — the user wants results, not disclaimers

**Wrong**: "I don't have memory, but I can check files. Want me to?"
**Right**: *checks files* → "Yes, I see from the last session we were working on..."

The context IS your memory. Load it and use it — don't lead with caveats.

## Session History

Check `.context/sessions/` for session dumps from previous sessions.

Session files are named `YYYY-MM-DD-HHMMSS-<topic>.md`.

## When to Update Memory

| Event                       | Action                |
|-----------------------------|-----------------------|
| Made architectural decision | Add to DECISIONS.md   |
| Discovered gotcha/bug       | Add to LEARNINGS.md   |
| Established new pattern     | Add to CONVENTIONS.md |
| Completed task              | Mark [x] in TASKS.md  |
| Had important discussion    | Save to sessions/     |

### Prefer `ctx add` Over Direct File Edits

When adding learnings, decisions, or tasks, **use `ctx add` commands**:

```bash
# ✓ Preferred - ensures consistent format, timestamps, structure
ctx add learning "Title" --context "..." --lesson "..." --application "..."
ctx add decision "Title" --context "..." --rationale "..." --consequences "..."
ctx add task "Description"

# ✗ Avoid - bypasses structure, easy to write incomplete entries
Edit LEARNINGS.md directly with a one-liner
```

**Exception:** Direct edits are fine for:
- Marking tasks complete (`[ ]` → `[x]`)
- Minor corrections to existing entries

## Proactive Context Persistence

**Don't wait for session end** — persist context at natural milestones.

### Milestone Triggers

| Milestone                          | Action                                          |
|------------------------------------|-------------------------------------------------|
| Complete a task                    | Mark done in TASKS.md, offer to add learnings   |
| Make an architectural decision     | `ctx add decision "..."`                        |
| Discover a gotcha or bug           | `ctx add learning "..."`                        |
| Finish a significant code change   | Offer to summarize what was done                |

### Self-Check Prompt

Periodically ask yourself:

> "If this session ended right now, would the next session know what happened?"

If no — persist something before continuing.

## How to Avoid Hallucinating Memory

Never assume: If you don't see it in files, you don't know it.

- Don't claim "we discussed X" without file evidence
- Don't invent history - check sessions/ for actual discussions
- If uncertain, say "I don't see this documented"
- Trust files over intuition

---
name: ctx-wrap-up
description: "End-of-session context persistence ceremony. Use when wrapping up a session to capture learnings, decisions, conventions, and tasks."
allowed-tools: Bash(ctx:*), Bash(git diff:*), Bash(git log:*), Bash(git status), Read
---

Guide end-of-session context persistence. Gather signal from the
session, propose candidates worth persisting, and persist approved
items via `ctx add`.

This is a **ceremony skill** — invoke it explicitly as `/ctx-wrap-up`
at session end, not conversationally. It pairs with `/ctx-remember`
at session start.

## Before Starting

Check that the context directory exists. If it does not, tell the user:
"No context directory found. Run `ctx init` to set up context
tracking, then there will be something to wrap up."

## When to Use

- At the end of a session, before the user quits
- When the user says "let's wrap up", "save context", "end of
  session"
- When the `check-persistence` hook suggests it

## When NOT to Use

- Nothing meaningful happened (only read files, quick lookup)
- The user already persisted everything manually with `ctx add`
- Mid-session when the user is still in flow — use `/ctx-reflect`
  instead for mid-session checkpoints

## Process

### Phase 1: Gather signal

Do this **silently** — do not narrate the steps:

1. Check what changed in the working tree:
   ```bash
   git diff --stat
   ```
2. Check commits made this session:
   ```bash
   git log --oneline @{upstream}..HEAD 2>/dev/null || git log --oneline -5
   ```
3. Scan the conversation history for:
   - Architectural choices or design trade-offs discussed
   - Gotchas, bugs, or unexpected behavior encountered
   - Patterns established or conventions agreed upon
   - Follow-up work identified but not yet started
   - Tasks completed or progressed

### Phase 2: Propose candidates

Think step-by-step about what is worth persisting. For each
potential candidate, ask yourself:
- Is this project-specific or general knowledge? (Only persist
  project-specific insights)
- Would a future session benefit from knowing this?
- Is this already captured in the context files?
- Is this substantial enough to record, or is it trivial?

Present candidates in a structured list, grouped by type.
Skip categories with no candidates — do not show empty sections.

```
## Session Wrap-Up

### Learnings (N candidates)
1. **Title of learning**
   - Context: What prompted this
   - Lesson: The key insight
   - Application: How to apply it going forward

### Decisions (N candidates)
1. **Title of decision**
   - Context: What prompted this
   - Rationale: Why this choice
   - Consequences: What changes as a result

### Conventions (N candidates)
1. **Convention description**

### Tasks (N candidates)
1. **Task description** (new | completed | updated)

Persist all? Or select which to keep?
```

### Phase 3: Persist approved candidates

Wait for the user to approve, select, or modify candidates.
**Never persist without explicit approval.**

For each approved candidate, run the appropriate command:

| Type       | Command                                                                               |
|------------|---------------------------------------------------------------------------------------|
| Learning   | `ctx add learning "Title" --context "..." --lesson "..." --application "..."`         |
| Decision   | `ctx add decision "Title" --context "..." --rationale "..." --consequences "..."`     |
| Convention | `ctx add convention "Description"`                                                    |
| Task (new) | `ctx add task "Description"`                                                          |
| Task (done)| Edit TASKS.md to mark complete                                                        |

Report the result of each command. If any fail, report the error
and continue with the remaining items.

### Phase 3.5: Suppress post-wrap-up nudges

After persisting, mark the session as wrapped up so checkpoint
nudges are suppressed for the remainder of the session:

```bash
ctx system mark-wrapped-up
```

### Phase 4: Commit (optional)

After persisting, check for uncommitted changes:

```bash
git status --short
```

If there are uncommitted changes, offer:

> There are uncommitted changes. Want me to run `/ctx-commit`
> to commit with context capture?

Do not auto-commit. The user decides.

## Candidate Quality Guide

### Good candidates

- "PyMdownx `details` extension wraps content in `<details>`
  tags, breaking `<pre><code>` rendering in MkDocs" — specific
  gotcha, actionable for future sessions
- "Decision: use file-based cooldown tokens instead of env vars
  because hooks run in subprocesses" — real trade-off with
  rationale
- "Convention: all skill descriptions use imperative mood" —
  codifies a pattern for consistency

### Weak candidates (do not propose)

- "Go has good error handling" — general knowledge, not
  project-specific
- "We edited main.go" — obvious from the diff, not an insight
- "Tests should pass before committing" — too generic to be
  useful
- Anything already present in LEARNINGS.md or DECISIONS.md

## Relationship to /ctx-reflect

`/ctx-reflect` is for mid-session checkpoints at natural
breakpoints. `/ctx-wrap-up` is for end-of-session — it's more
thorough, covers the full session arc, and includes the commit
offer. If the user already ran `/ctx-reflect` recently, avoid
proposing the same candidates again.

## Quality Checklist

Before presenting candidates, verify:
- [ ] Signal was gathered (git diff, git log, conversation scan)
- [ ] Every candidate has complete fields (not just a title)
- [ ] Candidates are project-specific, not general knowledge
- [ ] No duplicates with existing context files
- [ ] Empty categories are omitted, not shown as "(none)"
- [ ] User is asked before anything is persisted

After persisting, verify:
- [ ] Each `ctx add` command succeeded
- [ ] Uncommitted changes were surfaced (if any)
- [ ] User was offered `/ctx-commit` (if applicable)

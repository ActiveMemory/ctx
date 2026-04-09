---
name: ctx-wrap-up
description: "End-of-session context persistence ceremony. Use when wrapping up a session to capture learnings, decisions, conventions, and tasks."
tools: [bash, read, write, edit]
---

Guide end-of-session context persistence. Gather signal from the
session, propose candidates worth persisting, and persist approved
items via `ctx add`.

## When to Use

- At the end of a session, before the user quits
- When the user says "let's wrap up", "save context", "end of
  session"

## When NOT to Use

- Nothing meaningful happened (only read files, quick lookup)
- The user already persisted everything manually
- Mid-session: use `ctx-reflect` instead

## Process

### Phase 1: Gather signal

Do this **silently**:

1. Check what changed:
   ```bash
   git diff --stat
   ```
2. Check commits made this session:
   ```bash
   git log --oneline -5
   ```
3. Scan the conversation for:
   - Architectural choices or trade-offs
   - Gotchas or unexpected behavior
   - Patterns established or conventions agreed
   - Follow-up work identified
   - Tasks completed or progressed

### Phase 2: Propose candidates

Think step-by-step about what is worth persisting. For each
candidate ask:
- Is this project-specific or general knowledge?
- Would a future session benefit from knowing this?
- Is this already captured in context files?

Present candidates grouped by type. Skip empty categories.

```
## Session Wrap-Up

### Learnings (N candidates)
1. **Title** — Context, Lesson, Application

### Decisions (N candidates)
1. **Title** — Context, Rationale, Consequence

### Conventions (N candidates)
1. **Convention description**

### Tasks (N candidates)
1. **Task description** (new | completed | updated)

Persist all? Or select which to keep?
```

### Phase 3: Persist approved candidates

Wait for user approval. For each approved item:

| Type        | Command                                                              |
|-------------|----------------------------------------------------------------------|
| Learning    | `ctx add learning "Title" --context "..." --lesson "..." --application "..."` |
| Decision    | `ctx add decision "Title" --context "..." --rationale "..." --consequence "..."` |
| Convention  | `ctx add convention "Description"`                                   |
| Task (new)  | `ctx add task "Description"`                                         |
| Task (done) | Edit TASKS.md to mark complete                                       |

### Phase 4: Commit (optional)

After persisting, check for uncommitted changes:

```bash
git status --short
```

If there are uncommitted changes, offer to commit with
`ctx-commit`.

## Candidate Quality Guide

### Good candidates

- Specific gotchas with actionable lessons
- Real trade-offs with rationale
- Patterns codified for consistency

### Weak candidates (do not propose)

- General programming knowledge
- Obvious facts from the diff
- Things already in context files

## Quality Checklist

Before presenting:
- [ ] Signal was gathered (git diff, git log, conversation scan)
- [ ] Every candidate has complete fields
- [ ] Candidates are project-specific
- [ ] No duplicates with existing context
- [ ] Empty categories are omitted
- [ ] User is asked before persisting

After persisting:
- [ ] Each `ctx add` command succeeded
- [ ] Uncommitted changes were surfaced

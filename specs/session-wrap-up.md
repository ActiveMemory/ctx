# Session Wrap-Up Skill

## Problem

Every session ends with the same ritual:

1. Review what was done.
2. Identify learnings, decisions, conventions worth persisting.
3. Run `ctx add` commands.
4. Optionally commit.

This is mechanical enough to automate but nuanced enough that the AI should
drive it (not a shell script). Currently this depends on the user remembering
to ask, and on the agent knowing the right `ctx add` flags and structure.

## Skill: `/ctx-wrap-up`

A user-invocable skill that guides end-of-session context persistence.

### Trigger

- User says "let's wrap up", "end of session", "save context", or runs
  `/ctx-wrap-up`.
- The existing `check-persistence` hook could suggest running it when a
  session has been long and no context has been persisted.

### Behavior

**Phase 1: Gather signal**

1. Run `git diff --stat` to see files changed in the working tree.
2. Run `git log --oneline @{upstream}..HEAD 2>/dev/null || git log --oneline -5`
   to see commits made this session.
3. Scan conversation history for themes: architectural choices, gotchas
   encountered, patterns established, bugs fixed.

**Phase 2: Propose candidates**

Present a structured list of candidates, grouped by type:

```
## Session Wrap-Up

### Learnings (2 candidates)
1. [learning] Blog posts benefit from periodic enrichment passes
   - Context: ...
   - Lesson: ...
   - Application: ...

2. [learning] ...

### Decisions (1 candidate)
1. [decision] Knowledge scaling: archive path for decisions and learnings
   - Context: ...
   - Rationale: ...
   - Consequences: ...

### Conventions (1 candidate)
1. [convention] Update admonitions for historical blog content

### Tasks (0 candidates)
(none identified)

Persist all? Or select which to keep?
```

**Phase 3: Persist**

For each approved candidate, run the appropriate `ctx add` command.
Report results.

**Phase 4: Commit (optional)**

If there are uncommitted changes, offer to run `/ctx-commit`.

### What It Does NOT Do

- Does not persist anything without user approval.
- Does not invent learnings/decisions that weren't part of the session.
- Does not replace manual `ctx add` — it's a convenience layer.
- Does not auto-commit. It offers; the user decides.

## Skill File

```
.claude/skills/ctx-wrap-up/SKILL.md
```

Standard ctx skill structure: description, when to use, instructions,
examples.

## Integration with Existing Hooks

The `check-persistence` hook (`ctx system check-persistence`) already nudges
the user when no context has been persisted in a long session. It could be
updated to suggest `/ctx-wrap-up` instead of generic advice.

## Implementation Notes

- The skill is pure prompt — no Go code needed.
- It uses existing `ctx add` commands, `git` for signal gathering, and
  conversation context for candidate identification.
- The quality of candidates depends on the AI's ability to distinguish
  session-specific insights from routine work. The skill instructions should
  include examples of good vs. weak candidates.

## Session Ceremonies: Explicit Invocation

Most ctx skills encourage a conversational approach ("jot down: check DNS
after deploy" instead of `ctx pad add "check DNS after deploy"`). The session
bookend skills are the exception.

`/ctx-remember` and `/ctx-wrap-up` should be documented as **explicit
slash-command invocations**, not conversational triggers. Reasons:

1. **Well-defined moments**: They happen at session boundaries, not mid-flow.
   A slash command marks the boundary clearly.
2. **Ambiguity**: "Do you remember?" could mean many things. `/ctx-remember`
   means exactly one thing: load context and present a structured readback.
3. **Completeness**: Conversational triggers risk partial execution — the
   agent might skip steps. The slash command runs the full ceremony.
4. **Muscle memory**: Typing `/ctx-remember` at session start and
   `/ctx-wrap-up` at session end becomes a habit, like opening and closing
   braces.

### Documentation Deliverables

- **New recipe**: `docs/recipes/session-ceremonies.md` — dedicated recipe
  covering the two bookend skills as explicit rituals. Structure:
  - Why explicit invocation (not conversational) for these two skills
  - Session start: `/ctx-remember` — what it does, what to expect
  - Session end: `/ctx-wrap-up` — what it does, approval flow
  - Quick reference card (copy-paste slash commands)
  - When to skip (trivial sessions, quick lookups)

- **Update existing docs**:
  - `docs/skills.md`: add a "Session Ceremonies" grouping note explaining
    these two are explicitly invoked, unlike other skills
  - `docs/recipes/session-lifecycle.md`: cross-link to the new ceremony
    recipe, note the explicit invocation pattern
  - `docs/prompting-guide.md`: add a callout distinguishing ceremony skills
    (explicit) from workflow skills (conversational)
  - `docs/first-session.md`: mention `/ctx-remember` as the recommended
    session start

## Testing

- Manual: run `/ctx-wrap-up` at end of a real session, verify candidates
  are relevant and `ctx add` commands succeed.
- Regression: ensure the skill doesn't propose candidates when nothing
  meaningful happened (e.g., a session that only read files).

---
title: "Session Ceremonies"
icon: lucide/bookmark
---

![ctx](../images/ctx-banner.png)

## The Problem

Sessions have two critical moments: the **start** and the **end**.

At the start, you need the agent to load context and confirm it knows what
is going on. At the end, you need to capture whatever the session produced
before the conversation disappears.

Most ctx skills work conversationally: "jot down: check DNS after deploy"
is as good as `/ctx-pad add "check DNS after deploy"`. But session
boundaries are different. They are well-defined moments with specific
requirements, and partial execution is costly.

If the agent only half-loads context at the start, it works from stale
assumptions. If it only half-persists at the end, learnings and decisions
are lost.

**Session ceremonies** are the two bookend skills that mark these
boundaries. They are the exception to the conversational rule:
invoke them explicitly as slash commands.

!!! tip "TL;DR"
    **Start**: `/ctx-remember` — load context, get a structured readback.

    **End**: `/ctx-wrap-up` — review session, propose candidates, persist approved items.

    Use the slash commands, not conversational triggers, for completeness.

## Commands and Skills Used

| Tool             | Type  | Purpose                                          |
|------------------|-------|--------------------------------------------------|
| `/ctx-remember`  | Skill | Load context and present structured readback     |
| `/ctx-wrap-up`   | Skill | Gather session signal, propose and persist context|
| `/ctx-commit`    | Skill | Commit with context capture (offered by wrap-up) |
| `ctx agent`      | CLI   | Load token-budgeted context packet               |
| `ctx recall list`| CLI   | List recent sessions                             |
| `ctx add`        | CLI   | Persist learnings, decisions, conventions, tasks  |

## Why Explicit Invocation

Most ctx skills encourage natural language. These two are different:

**Well-defined moments.** Sessions have clear boundaries. A slash command
marks the boundary unambiguously.

**Ambiguity risk.** "Do you remember?" could mean many things.
`/ctx-remember` means exactly one thing: load context and present a
structured readback.

**Completeness.** Conversational triggers risk partial execution. The
agent might load some files but skip the session history, or persist one
learning but forget to check for uncommitted changes. The slash command
runs the full ceremony.

**Muscle memory.** Typing `/ctx-remember` at session start and
`/ctx-wrap-up` at session end becomes a habit, like opening and closing
braces.

## Session Start: /ctx-remember

Invoke at the beginning of every session:

```
/ctx-remember
```

The skill silently:

1. Loads the context packet via `ctx agent --budget 4000`
2. Reads TASKS.md, DECISIONS.md, LEARNINGS.md
3. Checks recent sessions via `ctx recall list --limit 3`

Then presents a **structured readback** with four sections:

- **Last session**: topic, date, what was accomplished
- **Active work**: pending and in-progress tasks
- **Recent context**: 1-2 relevant decisions or learnings
- **Next step**: suggestion or question about what to focus on

The readback should feel like recall, not a file system tour. If the
agent says "Let me check if there are files..." instead of a confident
summary, the skill is not working correctly.

!!! note "What about 'Do you remember?'"
    The conversational trigger still works. But `/ctx-remember` guarantees
    the full ceremony runs: context packet, file reads, session history,
    and all four readback sections. The conversational version may cut
    corners.

## Session End: /ctx-wrap-up

Invoke before ending a session where meaningful work happened:

```
/ctx-wrap-up
```

The skill runs four phases:

### Phase 1: Gather signal

Silently checks `git diff --stat`, recent commits, and scans the
conversation for themes: architectural choices, gotchas, patterns
established, follow-up work identified.

### Phase 2: Propose candidates

Presents a structured list grouped by type:

```markdown
## Session Wrap-Up

### Learnings (2 candidates)
1. **PyMdownx details extension breaks pre/code rendering**
   - Context: Journal site showed broken code blocks inside details tags
   - Lesson: details extension wraps content in <details> HTML, which
     interferes with <pre><code> rendering
   - Application: Use fenced code blocks instead of indented code inside
     admonitions when details extension is active

2. **Hook subprocesses cannot propagate env vars**
   - Context: Set env var in PreToolUse hook, invisible in main session
   - Lesson: Hooks execute in child processes; env changes don't propagate
   - Application: Use tombstone files for hook-to-session communication

### Decisions (1 candidate)
1. **File-based cooldown tokens over env vars**
   - Context: Need session-scoped cooldown for ctx agent auto-loading
   - Rationale: File tokens survive across processes, simpler than IPC
   - Consequences: Tombstone files accumulate in /tmp; need TTL cleanup

Persist all? Or select which to keep?
```

Each candidate has complete structured fields, not just a title.
Empty categories are omitted.

### Phase 3: Persist

After you approve (all, some, or modified), the skill runs the
appropriate `ctx add` commands and reports results.

### Phase 4: Commit offer

If there are uncommitted changes, offers to run `/ctx-commit`.
Does not auto-commit.

## When to Skip

Not every session needs ceremonies.

**Skip `/ctx-remember`** when:

- You are doing a quick one-off lookup (reading a file, checking a value)
- Context was already loaded this session via `/ctx-agent`
- You are continuing immediately after a previous session and context is
  still fresh

**Skip `/ctx-wrap-up`** when:

- Nothing meaningful happened (only read files, answered a question)
- You already persisted everything manually during the session
- The session was trivial (typo fix, quick config change)

A good heuristic: if the session produced something a future session
should know about, run `/ctx-wrap-up`. If not, just close.

## Quick Reference

```
# Session start
/ctx-remember

# ... do work ...

# Session end
/ctx-wrap-up
```

That is the complete ceremony. Two commands, bookending your session.

## Relationship to Other Skills

| Skill          | When                          | Purpose                        |
|----------------|-------------------------------|--------------------------------|
| `/ctx-remember`| Session start                 | Load and confirm context       |
| `/ctx-reflect` | Mid-session breakpoints       | Checkpoint at milestones       |
| `/ctx-wrap-up` | Session end                   | Full session review and persist|
| `/ctx-commit`  | After completing work         | Commit with context capture    |

`/ctx-reflect` is for mid-session checkpoints. `/ctx-wrap-up` is for
end-of-session: it is more thorough, covers the full session arc, and
includes the commit offer. If you already ran `/ctx-reflect` recently,
`/ctx-wrap-up` avoids proposing the same candidates again.

## Tips

**Make it a habit.** The value of ceremonies compounds over sessions.
Each `/ctx-wrap-up` makes the next `/ctx-remember` richer.

**Trust the candidates.** The agent scans the full conversation. It
often catches learnings you forgot about.

**Edit before approving.** If a proposed candidate is close but not
quite right, tell the agent what to change. Do not settle for a vague
learning when a precise one is possible.

**Do not force empty ceremonies.** If `/ctx-wrap-up` finds nothing
worth persisting, that is fine. A session that only read files and
answered questions does not need artificial learnings.

## Next Up

**[Browsing and Enriching Past Sessions →](session-archaeology.md)**:
Export session history to a browsable journal and enrich entries with
metadata.

## See Also

* [The Complete Session](session-lifecycle.md): the full session workflow
  that ceremonies bookend
* [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md):
  deep dive on what gets persisted during wrap-up
* [Detecting and Fixing Drift](context-health.md): keeping context files
  accurate between ceremonies

---
title: "Persisting Decisions, Learnings, and Conventions"
icon: lucide/brain
---

![ctx](../images/ctx-banner.png)

## Problem

You debug a subtle issue, discover the root cause, and move on. Three weeks
later, a different session hits the same issue. The knowledge existed briefly
in one session's memory but was never written down. Architectural decisions
suffer the same fate: you weigh trade-offs, pick an approach, and six sessions
later the AI suggests the alternative you already rejected.

How do you make sure important context survives across sessions?

## Commands and Skills Used

| Tool                    | Type    | Purpose                                       |
|-------------------------|---------|-----------------------------------------------|
| `ctx add decision`      | Command | Record an architectural decision              |
| `ctx add learning`      | Command | Record a gotcha, tip, or lesson               |
| `ctx add convention`    | Command | Record a coding pattern or standard           |
| `ctx decisions reindex` | Command | Rebuild the quick-reference index             |
| `ctx learnings reindex` | Command | Rebuild the quick-reference index             |
| `/ctx-add-decision`     | Skill   | AI-guided decision capture with validation    |
| `/ctx-add-learning`     | Skill   | AI-guided learning capture with validation    |
| `/ctx-reflect`          | Skill   | Surface items worth persisting at breakpoints |

## The Workflow

### Step 1: Understand What to Persist

Three context files serve different purposes:

**Decisions** (DECISIONS.md) answer "why is it this way?" They record
trade-offs between alternatives with structured fields: context, rationale,
and consequences. Decisions prevent re-debating settled questions.

**Learnings** (LEARNINGS.md) answer "what did we discover the hard way?"
They record gotchas and debugging insights specific to this project with
structured fields: context, lesson, and application. Learnings prevent
repeating past mistakes.

**Conventions** (CONVENTIONS.md) answer "how do we do things here?" They
record patterns and standards. No structured fields required -- just a name,
a rule, and an example. Conventions keep code consistent across sessions.

The decision point: if you chose between alternatives, it is a decision.
If you discovered something surprising, it is a learning. If you are
codifying a repeated pattern, it is a convention.

### Step 2: Record Decisions with Structured Fields

Decisions require three structured fields: `--context`, `--rationale`,
and `--consequences`.

```bash
ctx add decision "Use file-based cooldown tokens instead of env vars" \
  --context "Hook subprocesses cannot persist env vars to parent shell" \
  --rationale "File tokens survive across processes. Simpler than IPC. Cleanup is automatic via TTL." \
  --consequences "Tombstone files accumulate in /tmp. Cannot share state across machines."
```

Include rejected alternatives in the rationale when multiple options were
considered:

```bash
ctx add decision "Use Cobra for CLI framework" \
  --context "Go CLI project needs subcommand support and shell completion" \
  --rationale "Mature routing and built-in completion. Chose over urfave/cli (weaker subcommands) and kong (smaller community)." \
  --consequences "More boilerplate per command. Completion scripts generated automatically."
```

The `/ctx-add-decision` skill guides you through the fields interactively.
For quick decisions, it supports a Y-statement: "In the context of
[situation], facing [constraint], we decided for [choice] and against
[alternatives], to achieve [benefit], accepting that [trade-off]."

### Step 3: Record Learnings with Structured Fields

Learnings require three structured fields: `--context`, `--lesson`, and
`--application`.

```bash
ctx add learning "Claude Code hooks run in a subprocess" \
  --context "Set env var in PreToolUse hook, but it was not visible in the main session" \
  --lesson "Hook scripts execute in a child process. Env changes do not propagate to parent." \
  --application "Use tombstone files for hook-to-session communication. Never rely on hook env vars."
```

The `/ctx-add-learning` skill applies three filters: (1) Could someone
Google this in 5 minutes? (2) Is it specific to this codebase? (3) Did it
take real effort to discover? All three must pass. Learnings capture
principles and heuristics, not code snippets.

### Step 4: Record Conventions

Conventions do not require structured fields. Use `--section` to place them
in the right category:

```bash
ctx add convention "Use kebab-case for all CLI flag names" --section "Naming"
ctx add convention "Integration tests use t.TempDir() for .context/ isolation" --section "Testing"
ctx add convention "One command per file in cmd/ directory" --section "Structure"
```

Conventions work best for rules that come up repeatedly. Codify a pattern
the third time you see it, not the first.

### Step 5: Reindex After Manual Edits

DECISIONS.md and LEARNINGS.md maintain a quick-reference index at the top --
a compact table of date and title for each entry. The index updates
automatically via `ctx add`, but falls out of sync after hand-edits.

```bash
ctx decisions reindex
ctx learnings reindex
```

Run reindex after any manual edit. The index lets AI tools scan all entries
without reading the full file, which matters when token budgets are tight.

### Step 6: Use /ctx-reflect to Surface What to Capture

At natural breakpoints -- after completing a feature, fixing a bug, or
before ending a session -- use `/ctx-reflect` to identify items worth
persisting.

```text
/ctx-reflect
```

The skill walks through learnings, decisions, tasks, and session notes,
skipping categories with nothing to report. The output includes specific
commands for each suggested persist:

> This session implemented file-based cooldown for `ctx agent` and
> discovered that hook subprocesses cannot set env vars in the parent.
>
> I'd suggest persisting:
> - **Learning**: Hook subprocesses cannot propagate env vars
>   `ctx add learning "..." --context "..." --lesson "..." --application "..."`
> - **Decision**: File-based cooldown tokens over env vars
>   `ctx add decision "..." --context "..." --rationale "..." --consequences "..."`
>
> Want me to persist any of these?

The skill always asks before persisting.

## Putting It Together

```bash
# Decision: record the trade-off
ctx add decision "Use PostgreSQL over SQLite" \
  --context "Need concurrent multi-user access" \
  --rationale "SQLite locks on writes; Postgres handles concurrency" \
  --consequences "Requires a database server; team needs Postgres training"

# Learning: record the gotcha
ctx add learning "SQL migrations must be idempotent" \
  --context "Deploy failed when migration ran twice after rollback" \
  --lesson "CREATE TABLE without IF NOT EXISTS fails on retry" \
  --application "Always use IF NOT EXISTS guards in migrations"

# Convention: record the pattern
ctx add convention "API handlers return structured errors" --section "API"

# Reindex after manual edits
ctx decisions reindex
ctx learnings reindex

# Reflect at breakpoints (from AI assistant)
# /ctx-reflect
```

## Tips

- **Record decisions at the moment of choice.** The alternatives you
  considered and the reasons you rejected them fade quickly. Capture
  trade-offs while they are fresh.
- **Learnings should fail the Google test.** If someone could find it
  in a 5-minute search, it does not belong in LEARNINGS.md.
- **Conventions earn their place through repetition.** Add a convention
  the third time you see a pattern, not the first.
- **Use `/ctx-reflect` at every natural breakpoint.** The skill's
  checklist catches items you might otherwise lose.
- **Keep entries self-contained.** Each entry should make sense on its
  own. A future session may load only one due to token budget constraints.
- **Reindex after every hand-edit.** It takes less than a second. A stale
  index causes AI tools to miss entries.
- **Prefer the structured fields.** The verbosity forces clarity. A
  decision without rationale is just a fact; a learning without application
  is just a story.

## See Also

- [Tracking Work Across Sessions](task-management.md) -- managing the tasks
  that decisions and learnings support
- [The Complete Session](session-lifecycle.md) -- full session lifecycle
  including reflection and context persistence
- [Detecting and Fixing Drift](context-health.md) -- keeping knowledge
  files accurate as the codebase evolves
- [CLI Reference](../cli-reference.md) -- full documentation for `ctx add`,
  `ctx decisions`, `ctx learnings`
- [Context Files](../context-files.md) -- format and conventions for
  DECISIONS.md, LEARNINGS.md, and CONVENTIONS.md

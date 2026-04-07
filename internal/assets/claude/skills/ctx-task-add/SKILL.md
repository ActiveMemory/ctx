---
name: ctx-task-add
description: "Add a task. Use when follow-up work is identified or when breaking down complex work into subtasks."
allowed-tools: Bash(ctx:*)
---

Add a task to TASKS.md.

## Before Recording

Three questions: if any answer is "no", don't record:

1. **"Is this actionable?"** → If it's a vague wish, clarify first
2. **"Would someone else know what to do?"** → If not, add more detail
3. **"Is this tracked elsewhere?"** → If yes, don't duplicate

Tasks should describe **what to do and why**, not just a topic.

## When to Use

- When follow-up work is identified during a session
- When breaking down a complex task into subtasks
- When the user mentions something that should be tracked

## When NOT to Use

- Vague ideas without clear scope (discuss first, then add)
- Work already completed (mark existing tasks done instead)
- One-line fixes you can do right now (just do it)

## Gathering Information

If the user provides only a topic, ask:

1. "What specifically needs to happen?" → Scope the work
2. "Why does this matter?" → Capture motivation
3. "Is this high, medium, or low priority?" → Set priority

## Execution

```bash
ctx add task "Task description" \
  --session-id SESSION --branch BRANCH --commit HASH \
  [--priority high|medium|low] [--section "Phase N"]
```

Provenance flags (`--session-id`, `--branch`, `--commit`) are **required**.
Get these values from the hook-relayed provenance line in your context
(e.g., `Session: abc12345 | Branch: main @ 68fbc00a`).

**Prefer this skill over raw `ctx add task`**: the conversational
approach lets you automatically pick up session ID, branch, and commit
from the provenance line already in your context window.

**Placement**: Without `--section`, the task is inserted before the
first unchecked task in TASKS.md. Use `--section` only when you need
a specific section (e.g., `--section "Maintenance"`).

**Example: specific and actionable:**
```bash
ctx add task "Add --cooldown flag to ctx agent to suppress repeated output within a time window. Use tombstone file per session for isolation." \
  --session-id abc12345 --branch main --commit 68fbc00a \
  --priority medium
```

**Example: with context for why:**
```bash
ctx add task "Investigate ctx init overwriting user-generated content in context files. Commit a9df9dd wiped 18 decisions from DECISIONS.md. Need guard to prevent reinit from destroying user data." \
  --session-id abc12345 --branch main --commit 68fbc00a \
  --priority high
```

**Example: scoped subtask:**
```bash
ctx add task "Add topic-based navigation to blog when post count reaches 15+" \
  --session-id abc12345 --branch main --commit 68fbc00a \
  --priority low
```

**Bad examples (too shallow):**
```bash
ctx add task "Fix bug"              # What bug? Where?
ctx add task "Improve performance"  # Of what? How?
ctx add task "Authentication"       # That's a topic, not a task
# Also bad: missing --session-id, --branch, --commit
```

## Quality Checklist

Before recording, verify:
- [ ] Task starts with a verb (Add, Fix, Implement, Investigate, Update)
- [ ] Someone unfamiliar with the session could act on it
- [ ] Not a duplicate of an existing task in TASKS.md (check first)
- [ ] Priority set if the user indicated urgency

Confirm the task was added.

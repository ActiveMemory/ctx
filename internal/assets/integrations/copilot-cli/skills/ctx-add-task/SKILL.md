---
name: ctx-add-task
description: "Add a task. Use when follow-up work is identified or when breaking down complex work into subtasks."
tools: [bash]
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

## Execution

```bash
ctx add task "Task description" [--priority high|medium|low] [--section "Phase N"]
```

**Good examples:**
```bash
ctx add task "Add --cooldown flag to ctx agent" --priority medium
ctx add task "Investigate ctx init overwriting user content" --priority high
```

**Bad examples (too shallow):**
```bash
ctx add task "Fix bug"              # What bug? Where?
ctx add task "Improve performance"  # Of what? How?
```

## Quality Checklist

- [ ] Task starts with a verb (Add, Fix, Implement, Investigate, Update)
- [ ] Someone unfamiliar with the session could act on it
- [ ] Not a duplicate of an existing task in TASKS.md (check first)
- [ ] Priority set if the user indicated urgency

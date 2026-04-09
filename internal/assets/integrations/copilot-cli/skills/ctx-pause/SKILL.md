---
name: ctx-pause
description: "Pause context nudge/reminder hooks for the session. Security hooks remain active."
tools: [bash]
---

Temporarily pause context nudge and reminder hooks while keeping
security hooks active.

## When to Use

- When hooks are too noisy during focused work
- When doing rapid iteration and nudges interrupt flow
- When the user says "pause hooks" or "too many reminders"

## When NOT to Use

- At session start (hooks haven't fired yet)
- When the user wants to disable security hooks (not supported)

## Process

```bash
ctx system pause-hooks
```

This suppresses:
- Ceremony checks (remember, wrap-up)
- Persistence nudges
- Task completion checks
- Journal reminders

This does NOT suppress:
- Dangerous command blocking
- Context load gate
- Version checks

## Resuming

Use `ctx-resume` to re-enable hooks, or they automatically
resume at next session start.

## Quality Checklist

- [ ] User confirmed they want to pause
- [ ] Security hooks remain active
- [ ] Informed user how to resume

---
name: ctx-resume
description: "Resume context hooks after a pause."
tools: [bash]
---

Re-enable context hooks that were paused with `ctx-pause`.

## When to Use

- After a focused work period where hooks were paused
- When the user is ready for nudges and reminders again
- When the user says "resume hooks"

## When NOT to Use

- Hooks are not currently paused
- At session start (hooks auto-resume)

## Process

```bash
ctx system resume-hooks
```

This re-enables all non-security hooks:
- Ceremony checks
- Persistence nudges
- Task completion checks
- Journal reminders

## Quality Checklist

- [ ] Hooks were actually paused before resuming
- [ ] Confirmed hooks are active again

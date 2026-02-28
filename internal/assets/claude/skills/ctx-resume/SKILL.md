---
name: ctx-resume
description: "Resume context hooks after a pause. Use when ready to restore normal hook behavior."
allowed-tools: Bash(ctx:*)
---

Resume all context hooks after a `/ctx-pause`. Restores normal nudge,
reminder, and ceremony behavior.

## When to Use

- User says "resume ctx", "resume context", "unpause"
- User says "turn nudges back on"
- Session has evolved from a quick task into real project work
- Before running `/ctx-wrap-up` (wrap-up needs hooks active)

## When NOT to Use

- Session is not paused (resume is a silent no-op, but don't confuse the user)
- User wants to restart or reset the session (just start a new session)

## Execution

Run the resume command:

```bash
ctx resume
```

Then confirm to the user:

> Context hooks resumed. Nudges, reminders, and ceremonies are active again.

## Important Notes

- **Silent no-op if not paused**: safe to run even if hooks aren't paused
- **Turn counter resets**: the graduated reminder counter starts fresh if
  you pause again later

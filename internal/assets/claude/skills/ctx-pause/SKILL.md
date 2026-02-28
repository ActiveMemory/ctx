---
name: ctx-pause
description: "Pause context hooks for this session. Use when context nudges aren't needed for the current task."
allowed-tools: Bash(ctx:*)
---

Pause all context nudge and reminder hooks for the current session.
Security hooks (dangerous command blocking) still fire.

## When to Use

- User says "pause ctx", "pause context", "quiet mode"
- User says "stop the nudges", "too many reminders"
- Quick investigation or one-off task that doesn't need ceremonies
- User explicitly asks to reduce context overhead

## When NOT to Use

- User wants to silence a specific hook (use `ctx system message edit` to
  customize or silence individual hooks)
- User wants to permanently disable hooks (edit `.claude/settings.local.json`)
- Session involves real project work that benefits from persistence nudges

## Execution

Run the pause command:

```bash
ctx pause
```

Then confirm to the user:

> Context hooks paused for this session. Nudges, reminders, and ceremony
> prompts are silenced. Security hooks still fire.
>
> Resume anytime with `/ctx-resume`.

## Important Notes

- **Session-scoped**: only affects the current session, not other terminals
- **Hooks still fire silently**: they check the pause flag and no-op
- **Graduated reminder**: a minimal `ctx:paused` indicator appears in hook
  output so the state is never invisible
- **Resume before wrap-up**: if the session evolves into real work, resume
  hooks before wrapping up to capture learnings and decisions
- **Initial context load is unaffected**: the ~8k token startup injection
  happens before any command runs — pause only affects subsequent hooks

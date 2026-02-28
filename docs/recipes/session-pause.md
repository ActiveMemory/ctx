---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "Pausing Context Hooks"
icon: lucide/pause
---

![ctx](../images/ctx-banner.png)

## The Problem

Not every session needs the full ceremony. Quick investigations, one-off
questions, small fixes unrelated to active project work — these tasks
don't benefit from persistence nudges, ceremony reminders, or knowledge
checks. Every hook still fires, consuming tokens and attention on work
that won't produce learnings or decisions worth capturing.

## TL;DR

| Command | What it does |
|---------|--------------|
| `ctx pause` or `/ctx-pause` | Silence all nudge hooks for this session |
| `ctx resume` or `/ctx-resume` | Restore normal hook behavior |

Pause is **session-scoped** — it only affects the current session.
Other sessions (same project, different terminal) are unaffected.

## What Gets Paused

All nudge and reminder hooks go silent:

- Context size checkpoints
- Ceremony adoption nudges
- Persistence reminders
- Journal maintenance reminders
- Knowledge growth nudges
- Map staleness nudges
- Version update nudges
- Resource pressure warnings
- QA reminders
- Post-commit nudges
- Specs nudges
- Backup age warnings
- Context load gate
- Pending reminders relay

## What Still Fires

**Security hooks** always run, even when paused:

- `block-non-path-ctx` — prevents `./ctx` invocations
- `block-dangerous-commands` — blocks `sudo`, force push, etc.

**Housekeeping hooks** also run:

- `cleanup-tmp` — removes stale temp files at session end

## Workflow

```bash
# 1. Session starts — context loads normally (~8k tokens, unavoidable)

# 2. You realize this is a quick task
ctx pause

# 3. Work without interruption — hooks are silent

# 4. Session evolves into real work? Resume first
ctx resume

# 5. Now wrap up normally
# /ctx-wrap-up
```

## Graduated Reminder

Paused hooks aren't completely invisible. A minimal indicator appears so
you always know the state:

| Paused turns | What you see |
|--------------|--------------|
| 1–5 | `ctx:paused` |
| 6+ | `ctx:paused (N turns) — resume with /ctx-resume` |

This prevents the "forgot I paused" problem during long sessions.

## Tips

- **Resume before wrapping up.** If your quick task turns into real work,
  resume hooks before running `/ctx-wrap-up`. The wrap-up ceremony needs
  active hooks to capture learnings properly.

- **Initial context load is unaffected.** The ~8k token startup injection
  (CLAUDE.md, playbook, constitution) happens before any command runs.
  Pause only affects hooks that fire *during* the session.

- **Use for quick investigations.** Debugging a stack trace? Checking a
  git log? Answering a colleague's question? Pause, do the work, close
  the session. No ceremony needed.

- **Don't use for real work.** If you're implementing features, fixing
  bugs, or making decisions — keep hooks active. The nudges exist to
  prevent context loss.

## See Also

*See also: [Session Ceremonies](session-ceremonies.md) — the bookend
rituals that pause lets you skip when they aren't needed.*

*See also: [Customizing Hook Messages](customizing-hook-messages.md) —
if you want to change what hooks say rather than silencing them entirely.*

*See also: [The Complete Session](session-lifecycle.md) — the full
session workflow that pause shortcuts for quick tasks.*

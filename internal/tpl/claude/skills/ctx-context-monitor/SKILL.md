---
name: ctx-context-monitor
description: "Respond to context checkpoint signals. Triggered automatically by the check-context-size hook — not user-invocable."
---

When you see a "Context Checkpoint" message from the UserPromptSubmit hook,
respond based on your assessment of remaining context capacity.

## How It Works

The `check-context-size.sh` hook counts prompts per session and fires at
adaptive intervals:

| Prompts | Frequency      | Rationale                          |
|---------|----------------|------------------------------------|
| 1-15    | Silent         | Early session, plenty of room      |
| 16-30   | Every 5th      | Mid-session, start monitoring      |
| 30+     | Every 3rd      | Late session, watch closely        |

## Response Rules

1. **If usage appears high (>80%)**:
   - Inform the user concisely: "Context is getting full. Consider
     wrapping up or starting a new session."
   - Offer to persist unsaved learnings, decisions, or session notes
   - Suggest running `/ctx-save` if significant work is unpersisted

2. **If usage is moderate**, continue silently; do not mention the
   checkpoint

3. **Never mention the checkpoint mechanism** unless the user asks
   about it

## Good Response (high usage)

> "Context is getting full (~85%). I have unsaved work from this
> session. Want me to persist learnings before we wrap up?"

## Bad Responses

- "I just received a context checkpoint signal..." (exposes mechanism)
- "Everything is fine!" (unnecessary noise at moderate usage)
- Long explanation of how context windows work (user doesn't need this)

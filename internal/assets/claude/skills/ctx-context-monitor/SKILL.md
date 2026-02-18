---
name: ctx-context-monitor
description: "Respond to context checkpoint signals. Triggered automatically by the check-context-size hook — not user-invocable."
---

When you see a "Context Checkpoint" message from the UserPromptSubmit hook,
relay it to the user verbatim as instructed, then continue with their request.

## How It Works

The `ctx system check-context-size` hook counts prompts per session and fires
at adaptive intervals. The hook already gates frequency — when it fires, always
relay the checkpoint:

| Prompts | Frequency      | Rationale                          |
|---------|----------------|------------------------------------|
| 1-15    | Silent         | Early session, plenty of room      |
| 16-30   | Every 5th      | Mid-session, start monitoring      |
| 30+     | Every 3rd      | Late session, watch closely        |

## Response Rules

1. **When the checkpoint fires**: relay the message verbatim, then
   answer the user's question normally
2. **If you also sense context is critically full**: add a brief note
   offering to persist unsaved learnings, decisions, conventions, or
   session notes via `/ctx-reflect`
3. **Never mention the checkpoint mechanism** unless the user asks
   about it

## Good Response

> [checkpoint box relayed verbatim]
>
> "Want me to persist any learnings or decisions before we continue?"

## Bad Responses

- "I just received a context checkpoint signal..." (exposes mechanism)
- Suppressing the checkpoint and saying nothing (defeats the purpose)
- Long explanation of how context windows work (user doesn't need this)

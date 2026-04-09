---
name: ctx-remind
description: "Manage session reminders. Use when the user says 'remind me to...' or asks about pending reminders."
tools: [bash]
---

Manage session-scoped reminders via `ctx remind` commands.

## When to Use

- User says "remind me to..." or "remind me about..."
- User asks "what reminders do I have?"
- User wants to dismiss or clear reminders

## When NOT to Use

- For structured tasks with status tracking (use `ctx add task`)
- For sensitive values or quick notes (use `ctx pad`)
- Create a reminder only when the user explicitly says "remind me"

## Command Mapping

| User intent                      | Command                                    |
|----------------------------------|--------------------------------------------|
| "remind me to refactor swagger"  | `ctx remind "refactor swagger"`            |
| "remind me tomorrow to check CI" | `ctx remind "check CI" --after YYYY-MM-DD` |
| "what reminders do I have?"      | `ctx remind list`                          |
| "dismiss reminder 3"             | `ctx remind dismiss 3`                     |
| "clear all reminders"            | `ctx remind dismiss --all`                 |

## Natural Language Date Handling

The CLI only accepts `YYYY-MM-DD` for `--after`. Convert natural
language dates to this format:

| User says      | You run                                         |
|----------------|-------------------------------------------------|
| "next session" | `ctx remind "..."` (no `--after`)               |
| "tomorrow"     | `ctx remind "..." --after YYYY-MM-DD`           |
| "next week"    | `ctx remind "..." --after YYYY-MM-DD` (+7 days) |

## Important Notes

- Reminders fire **every session** until dismissed: no throttle
- The `--after` flag gates when a reminder starts appearing
- Reminders are stored in `.context/reminders.json` (committed to git)

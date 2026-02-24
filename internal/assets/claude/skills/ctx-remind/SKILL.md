---
name: ctx-remind
description: "Manage session reminders. Use when the user says 'remind me to...' or asks about pending reminders."
allowed-tools: Bash(ctx:*)
---

Manage session-scoped reminders via `ctx remind` commands using
natural language. Translate what the user says into the right
command.

## When to Use

- User says "remind me to..." or "remind me about..."
- User asks "what reminders do I have?"
- User wants to dismiss or clear reminders
- User mentions reminders surfaced at session start

## When NOT to Use

- For structured tasks with status tracking (use `ctx add task`)
- For sensitive values or quick notes (use `ctx pad`)
- For architectural decisions (use `ctx add decision`)
- Do NOT auto-create reminders — only when the user explicitly says "remind me"

## Command Mapping

| User intent | Command |
|---|---|
| "remind me to refactor swagger" | `ctx remind "refactor swagger"` |
| "remind me tomorrow to check CI" | `ctx remind "check CI" --after YYYY-MM-DD` |
| "remind me next week to review auth" | `ctx remind "review auth" --after YYYY-MM-DD` |
| "what reminders do I have?" | `ctx remind list` |
| "dismiss reminder 3" | `ctx remind dismiss 3` |
| "clear all reminders" | `ctx remind dismiss --all` |

## Execution

**Add a reminder:**
```bash
ctx remind "refactor the swagger definitions"
```

**Add with date gate:**
```bash
ctx remind "check CI after the deploy" --after 2026-02-25
```

**List reminders:**
```bash
ctx remind list
```

**Dismiss by ID:**
```bash
ctx remind dismiss 3
```

**Dismiss all:**
```bash
ctx remind dismiss --all
```

## Natural Language Date Handling

The CLI only accepts `YYYY-MM-DD` for `--after`. You must convert
natural language dates to this format.

| User says | You run |
|---|---|
| "remind me next session" | `ctx remind "..."` (no `--after`) |
| "remind me tomorrow" | `ctx remind "..." --after YYYY-MM-DD` (tomorrow's date) |
| "remind me next week" | `ctx remind "..." --after YYYY-MM-DD` (7 days from now) |
| "remind me about X" | `ctx remind "X"` (no `--after`, immediate) |
| "remind me after Friday" | `ctx remind "..." --after YYYY-MM-DD` (next Saturday) |

If the date is ambiguous (e.g., "after the release"), ask the user
for a specific date.

## Important Notes

- Reminders fire **every session** until dismissed — no throttle
- The `--after` flag gates when a reminder starts appearing, not when
  it expires
- IDs are never reused — after dismissing ID 3, the next gets ID 4+
- Reminders are stored in `.context/reminders.json` (committed to git)
- After creating or dismissing, show the command output so the user
  can confirm the action

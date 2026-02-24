---
title: "Session Reminders"
icon: lucide/bell
---

![ctx](../images/ctx-banner.png)

## The Problem

You're deep in a session and realize: "I need to refactor the swagger
definitions next time." You could add a task, but this isn't a work item —
it's a note to future-you. You could jot it on the scratchpad, but
scratchpad entries don't announce themselves.

**How do you leave a message that your next session opens with?**

## TL;DR

```bash
ctx remind "refactor the swagger definitions"
ctx remind list
ctx remind dismiss 1
```

Reminders surface automatically at session start — verbatim, every
session, until you dismiss them.

## Commands and Skills Used

| Tool                    | Type        | Purpose                                        |
|-------------------------|-------------|------------------------------------------------|
| `ctx remind`            | CLI command | Add a reminder (default action)                |
| `ctx remind list`       | CLI command | Show all pending reminders                     |
| `ctx remind dismiss`    | CLI command | Remove a reminder by ID (or `--all`)           |
| `/ctx-remind`           | Skill       | Natural language interface to reminders         |

## The Workflow

### Step 1: Leave a Reminder

Tell your agent what to remember, or run it directly:

```text
You: "remind me to refactor the swagger definitions"

Agent: [runs ctx remind "refactor the swagger definitions"]
       "Reminder set:
         + [1] refactor the swagger definitions"
```

Or from the terminal:

```bash
ctx remind "refactor the swagger definitions"
```

### Step 2: Set a Date Gate (Optional)

If the reminder shouldn't fire until a specific date:

```text
You: "remind me to check the deploy logs after Tuesday"

Agent: [runs ctx remind "check the deploy logs" --after 2026-02-25]
       "Reminder set:
         + [2] check the deploy logs  (after 2026-02-25)"
```

The reminder stays silent until that date, then fires every session.

The agent converts natural language dates ("tomorrow", "next week",
"after the release on Friday") to `YYYY-MM-DD`. If it's ambiguous,
it asks.

### Step 3: Start a New Session

Next session, the reminder appears automatically before anything else:

```text
┌─ Reminders ──────────────────────────────────────
│  [1] refactor the swagger definitions
│
│ Dismiss: ctx remind dismiss <id>
│ Dismiss all: ctx remind dismiss --all
└──────────────────────────────────────────────────
```

No action needed — the `check-reminders` hook fires on
`UserPromptSubmit` and the agent relays the box verbatim.

### Step 4: Dismiss When Done

After you've acted on a reminder (or decided to skip it):

```text
You: "dismiss reminder 1"

Agent: [runs ctx remind dismiss 1]
       "Dismissed:
         - [1] refactor the swagger definitions"
```

Or clear everything:

```bash
ctx remind dismiss --all
```

### Step 5: Check What's Pending

```bash
ctx remind list
```

```
  [1] refactor the swagger definitions
  [3] review auth token expiry logic
  [4] check deploy logs  (after 2026-02-25, not yet due)
```

Date-gated reminders that haven't reached their date show
`(not yet due)`.

## Using `/ctx-remind` in a Session

Invoke the `/ctx-remind` skill, then describe what you want:

```text
You: /ctx-remind remind me to update the API docs
You: /ctx-remind what reminders do I have?
You: /ctx-remind dismiss reminder 3
```

| You say (after `/ctx-remind`)               | What the agent does                             |
|---------------------------------------------|-------------------------------------------------|
| "remind me to update the API docs"          | `ctx remind "update the API docs"`              |
| "remind me next week to check staging"      | `ctx remind "check staging" --after 2026-03-02` |
| "what reminders do I have?"                 | `ctx remind list`                               |
| "dismiss reminder 3"                        | `ctx remind dismiss 3`                          |
| "clear all reminders"                       | `ctx remind dismiss --all`                      |

## Reminders vs Scratchpad vs Tasks

| You want to...                              | Use                |
|---------------------------------------------|--------------------|
| Leave a note that announces itself next session | **`ctx remind`**   |
| Jot down a quick value or sensitive token   | **`ctx pad`**      |
| Track work with status and completion       | **TASKS.md**       |
| Record a decision or lesson for all sessions | **Context files**  |

**Decision guide:**

* If it should **announce itself** at session start → `ctx remind`
* If it's a **quiet note** you'll check manually → `ctx pad`
* If it's a **work item** you'll mark done → `TASKS.md`

!!! tip "Reminders Are Sticky Notes, Not Tasks"
    A reminder has no status, no priority, no lifecycle. It's a
    message to future-you that fires until dismissed. If you need
    tracking, use a task.

## Tips

* **Reminders fire every session.** Unlike nudges (which throttle to
  once per day), reminders repeat until you dismiss them. This is
  intentional — you asked to be reminded.
* **Date gating is session-scoped, not clock-scoped.** `--after
  2026-02-25` means "don't show until sessions on or after Feb 25."
  It does not mean "alarm at midnight on Feb 25."
* **The agent handles date parsing.** Say "next week" or "after
  Friday" — the agent converts it to `YYYY-MM-DD`. The CLI only
  accepts the explicit date format.
* **Reminders are committed to git.** They travel with the repo.
  If you switch machines, your reminders follow.
* **IDs never reuse.** After dismissing reminder 3, the next reminder
  gets ID 4 (or higher). No confusion from recycled numbers.

## Next Up

**[Using the Scratchpad →](scratchpad-with-claude.md)**: For quiet notes
and sensitive values that don't need session-start announcements.

## See Also

* [CLI Reference: ctx remind](../reference/cli-reference.md): full
  command syntax and flags
* [The Complete Session](session-lifecycle.md): how reminders fit into
  the session lifecycle
* [Managing Tasks](task-management.md): for work items that need status
  tracking

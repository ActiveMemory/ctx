# Session Reminders (`ctx remind`)

## Problem

You're mid-session and think: "I need to refactor the swagger definitions
next time." You could add a task, but this isn't a task — it's a sticky
note to future-you. Tasks have status, priority, lifecycle. This is a
one-liner you want relayed verbatim at the next session start, then
dismissed.

Calendar apps can't help either — they don't know when your coding
sessions start.

## Approach

A lightweight `ctx remind` command with:

- **Session-scoped triggers** — reminders fire at session start, not at
  clock times. ctx isn't a daemon and doesn't pretend to be.
- **Verbatim relay** — the message is stored exactly as typed and
  delivered exactly as stored. No summarization, no categorization.
- **Dismiss-on-acknowledge** — the agent or user dismisses after reading.
  No status tracking, no workflows.
- **Optional date gating** — `--after 2026-02-25` suppresses a reminder
  until that date passes. Still session-triggered, just time-gated.

### What this is NOT

- Not a task tracker (use TASKS.md).
- Not a calendar (use your calendar for "3pm tomorrow").
- Not a recurring reminder system (no cron, no repeat).
- Not encrypted (reminders are low-sensitivity; use `ctx pad` for secrets).

## Command

```
ctx remind "refactor the swagger definitions"
ctx remind "check CI after the deploy" --after 2026-02-25
ctx remind list
ctx remind dismiss <id>
ctx remind dismiss --all
```

### Subcommands

| Subcommand | Alias | Args | Description |
|------------|-------|------|-------------|
| (default) | `add` | TEXT | Add a reminder |
| `list` | `ls` | — | Show all pending reminders |
| `dismiss` | `rm` | ID | Dismiss one or all reminders |

The default action (no subcommand) is `add`, so `ctx remind "text"` works
without typing `ctx remind add "text"`.

## Behavior

### Adding a reminder

```
$ ctx remind "refactor the swagger definitions"
  + [1] refactor the swagger definitions
```

```
$ ctx remind "check CI after the deploy" --after 2026-02-25
  + [2] check CI after the deploy  (after 2026-02-25)
```

The reminder is appended to `.context/reminders.json`. The ID is an
auto-incrementing integer, monotonic within the file. IDs are never reused
within a file's lifetime (the next ID is `max(existing IDs) + 1`).

### Listing reminders

```
$ ctx remind list
  [1] refactor the swagger definitions
  [2] check CI after the deploy  (after 2026-02-25, not yet due)
```

Date-gated reminders that aren't yet due show `(after DATE, not yet due)`.
Due reminders show no annotation.

If no reminders exist: `No reminders.`

### Dismissing reminders

```
$ ctx remind dismiss 1
  - [1] refactor the swagger definitions

$ ctx remind dismiss --all
  - [1] refactor the swagger definitions
  - [2] check CI after the deploy
Dismissed 2 reminders.
```

Dismissing removes the entry from `reminders.json`.

Dismissing an ID that doesn't exist:
`No reminder with ID 1.` (exit 1)

## Flags

| Flag | Short | Applies to | Default | Description |
|------|-------|------------|---------|-------------|
| `--after` | `-a` | add | (none) | Don't surface until this date (YYYY-MM-DD) |
| `--all` | | dismiss | false | Dismiss all reminders |

## Storage

### File: `.context/reminders.json`

```json
[
  {
    "id": 1,
    "message": "refactor the swagger definitions",
    "created": "2026-02-23T14:30:00Z",
    "after": null
  },
  {
    "id": 2,
    "message": "check CI after the deploy",
    "created": "2026-02-23T14:31:00Z",
    "after": "2026-02-25"
  }
]
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | int | always | Auto-incremented, never reused |
| `message` | string | always | Verbatim reminder text |
| `created` | string | always | UTC RFC3339 timestamp |
| `after` | string | nullable | Date gate (YYYY-MM-DD), null if immediate |

JSON (not markdown) because:
- Structured fields (id, date) need reliable parsing.
- No human editing expected — the CLI is the interface.
- Small file, no performance concern.

### Gitignore

`reminders.json` is **committed** (not gitignored). Reminders travel with
the repo — useful when switching machines or sharing context with a
collaborator.

## Session-Start Integration

### Hook: `ctx system check-reminders`

A new system hook that fires on `UserPromptSubmit`, alongside existing
hooks like `check-ceremonies` and `check-persistence`.

**Logic:**

1. Check initialized (skip if no `.context/`).
2. Read `.context/reminders.json`.
3. Filter to due reminders: `after` is null OR `after <= today`.
4. If no due reminders → exit silently.
5. Emit VERBATIM relay with due reminders.
6. **No throttle** — reminders fire every session until dismissed. This is
   intentional: unlike nudges (which teach habits), reminders are explicit
   user requests. Suppressing them defeats the purpose.

### Output

```
IMPORTANT: Relay these reminders to the user VERBATIM before answering their question.

┌─ Reminders ──────────────────────────────────────
│  [1] refactor the swagger definitions
│  [3] review the auth token expiry logic
│
│ Dismiss: ctx remind dismiss <id>
│ Dismiss all: ctx remind dismiss --all
└──────────────────────────────────────────────────
```

If only date-gated (not yet due) reminders exist, the hook is silent.

### Hook Registration

Add to `system.go` command tree:

```go
systemCmd.AddCommand(checkRemindersCmd())
```

Add to `internal/assets/claude/hooks.json` under `UserPromptSubmit`:

```json
{
    "type": "command",
    "command": "ctx system check-reminders"
}
```

### Notification

If `ctx notify` is configured, the hook fires a `nudge` event:

```go
notify.Send(notify.Event{
    Name:    "nudge",
    Message: fmt.Sprintf("You have %d pending reminders", len(due)),
})
```

## Skill: `/ctx-remind`

A thin skill so the agent can create and dismiss reminders from
conversation without the user needing to drop to the terminal.

### Skill file: `internal/assets/claude/skills/ctx-remind/SKILL.md`

The skill instructs the agent to:

1. **Create**: When the user says "remind me to X", run
   `ctx remind "X"`. If they mention a date ("next week", "on Monday"),
   parse it to `--after YYYY-MM-DD`. If ambiguous, ask.
2. **Dismiss**: When the user acknowledges a reminder, run
   `ctx remind dismiss <id>`.
3. **List**: When the user asks "what reminders do I have?", run
   `ctx remind list`.
4. **Never invent**: Don't auto-create reminders. Only create when the
   user explicitly says "remind me".

### Natural language date handling

The agent (not the CLI) is responsible for parsing natural language dates.
The CLI only accepts `YYYY-MM-DD`. This keeps the Go code simple and
leverages what LLMs are good at.

| User says | Agent runs |
|-----------|------------|
| "remind me next session" | `ctx remind "..."` (no `--after`) |
| "remind me tomorrow" | `ctx remind "..." --after 2026-02-24` |
| "remind me next week" | `ctx remind "..." --after 2026-03-02` |
| "remind me about X" | `ctx remind "X"` |

## Implementation

### New files

```
internal/cli/remind/
  remind.go          Cmd() with add (default), list, dismiss subcommands
  store.go           readReminders(), writeReminders(), nextID()
  remind_test.go     Tests

internal/cli/system/
  check_reminders.go checkRemindersCmd(), runCheckReminders()

internal/assets/claude/skills/ctx-remind/
  SKILL.md           Agent skill instructions
```

### Modified files

- `internal/bootstrap/bootstrap.go` — Register `remind.Cmd`
- `internal/cli/system/system.go` — Register `checkRemindersCmd()`
- `internal/assets/claude/hooks.json` — Add `check-reminders` hook

### Core types

```go
// internal/cli/remind/store.go

type Reminder struct {
    ID      int     `json:"id"`
    Message string  `json:"message"`
    Created string  `json:"created"`
    After   *string `json:"after"`  // nullable YYYY-MM-DD
}

func readReminders() ([]Reminder, error) {
    path := remindersPath()
    data, err := os.ReadFile(path)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return nil, nil
        }
        return nil, err
    }
    var reminders []Reminder
    if err := json.Unmarshal(data, &reminders); err != nil {
        return nil, fmt.Errorf("parse reminders: %w", err)
    }
    return reminders, nil
}

func writeReminders(reminders []Reminder) error {
    data, err := json.MarshalIndent(reminders, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(remindersPath(), data, 0o644)
}

func nextID(reminders []Reminder) int {
    max := 0
    for _, r := range reminders {
        if r.ID > max {
            max = r.ID
        }
    }
    return max + 1
}

func remindersPath() string {
    return filepath.Join(rc.ContextDir(), "reminders.json")
}
```

### Core function: runAdd

```go
func runAdd(cmd *cobra.Command, message, after string) error {
    reminders, err := readReminders()
    if err != nil {
        return err
    }

    r := Reminder{
        ID:      nextID(reminders),
        Message: message,
        Created: time.Now().UTC().Format(time.RFC3339),
    }
    if after != "" {
        // Validate YYYY-MM-DD format.
        if _, err := time.Parse("2006-01-02", after); err != nil {
            return fmt.Errorf("invalid date %q (expected YYYY-MM-DD)", after)
        }
        r.After = &after
    }

    reminders = append(reminders, r)
    if err := writeReminders(reminders); err != nil {
        return err
    }

    suffix := ""
    if r.After != nil {
        suffix = fmt.Sprintf("  (after %s)", *r.After)
    }
    cmd.Printf("  + [%d] %s%s\n", r.ID, r.Message, suffix)
    return nil
}
```

### Core function: runCheckReminders

```go
// internal/cli/system/check_reminders.go

func checkRemindersCmd() *cobra.Command {
    return &cobra.Command{
        Use:    "check-reminders",
        Short:  "Surface pending reminders at session start",
        Hidden: true,
        RunE: func(cmd *cobra.Command, _ []string) error {
            return runCheckReminders(cmd)
        },
    }
}

func runCheckReminders(cmd *cobra.Command) error {
    if !isInitialized() {
        return nil
    }

    reminders, err := remind.ReadReminders()
    if err != nil {
        return nil // non-fatal: don't break session start
    }

    today := time.Now().Format("2006-01-02")
    var due []remind.Reminder
    for _, r := range reminders {
        if r.After == nil || *r.After <= today {
            due = append(due, r)
        }
    }

    if len(due) == 0 {
        return nil
    }

    cmd.Println("IMPORTANT: Relay these reminders to the user VERBATIM before answering their question.")
    cmd.Println()
    cmd.Println("┌─ Reminders ──────────────────────────────────────")
    for _, r := range due {
        cmd.Printf("│  [%d] %s\n", r.ID, r.Message)
    }
    cmd.Println("│")
    cmd.Println("│ Dismiss: ctx remind dismiss <id>")
    cmd.Println("│ Dismiss all: ctx remind dismiss --all")
    cmd.Println("└──────────────────────────────────────────────────")

    notify.Send(notify.Event{
        Name:    "nudge",
        Message: fmt.Sprintf("You have %d pending reminders", len(due)),
    })

    return nil
}
```

## Tests

| Test | Scenario |
|------|----------|
| `TestAdd_Basic` | Add reminder, verify JSON written with correct fields |
| `TestAdd_WithAfter` | Add with `--after`, verify date stored |
| `TestAdd_InvalidDate` | `--after garbage` → error |
| `TestAdd_IDIncrement` | Add 3, dismiss middle, add another → ID is max+1 |
| `TestList_Empty` | No file → "No reminders." |
| `TestList_Mixed` | Due and not-yet-due reminders, verify annotations |
| `TestDismiss_ByID` | Dismiss one, verify removed from file |
| `TestDismiss_NotFound` | Dismiss nonexistent ID → error |
| `TestDismiss_All` | `--all` clears file |
| `TestCheckReminders_NoDue` | All reminders date-gated in future → silent |
| `TestCheckReminders_Due` | Mix of due and not-due → only due in output |
| `TestCheckReminders_NoFile` | No reminders.json → silent |
| `TestCheckReminders_NullAfter` | After is null → always due |

## Design Decisions

- **No throttle on the session hook.** Unlike nudges (which teach habits
  and self-silence), reminders are explicit user requests. Throttling them
  would mean "I asked to be reminded but ctx decided not to." Reminders
  fire every session until dismissed.

- **JSON not markdown.** Reminders have structured fields (id, date).
  Markdown would require fragile parsing. The file is small and
  machine-managed — no reason for human-readable format.

- **No encryption.** Reminders are "refactor swagger" and "check CI",
  not secrets. If you need encrypted notes, use `ctx pad`. Adding
  encryption here would complicate the implementation for no real benefit.

- **Agent handles natural language dates.** The CLI accepts only
  `YYYY-MM-DD`. Parsing "next week" or "after the release" is exactly
  what the LLM is good at. Keeps Go code trivial.

- **Committed to git.** Reminders are project context, not personal
  state. They should travel with the repo (machine switching, pair
  sessions). If someone wants private reminders, that's a future
  `.gitignore` flag.

- **No recurring reminders.** This is sticky notes, not cron. "Remind me
  every Monday" is a calendar. If it comes up, we'll reconsider.

- **IDs never reuse.** Avoids confusion when the user sees "dismissed
  reminder 3" and later a new reminder 3 appears. The ID ceiling only
  grows. For a sticky-note system this is fine — you'd need thousands
  before the numbers get unwieldy.

## Non-Goals

- **Timer-based reminders** ("in 10 minutes"). ctx isn't a daemon.
- **Recurring reminders**. Use a calendar.
- **Priority or categories**. This is sticky notes, not Jira.
- **Markdown rendering** in reminder text. Stored and relayed as plain text.
- **Encryption**. Use `ctx pad` for sensitive content.
- **Interactive dismiss during session start**. The hook relays; the user
  dismisses when ready. No blocking prompts.

---

## Documentation

### 1. CLI Reference (`docs/reference/cli-reference.md`)

Insert between `ctx pad` and `ctx system` (alphabetical). Full entry:

```markdown
### `ctx remind`

Session-scoped reminders that surface at session start. Reminders are
stored verbatim and relayed verbatim — no summarization, no categories.

When invoked with a text argument and no subcommand, adds a reminder.

```bash
ctx remind "text"
ctx remind <subcommand>
```

#### `ctx remind add`

Add a reminder. This is the default action — `ctx remind "text"` and
`ctx remind add "text"` are equivalent.

```bash
ctx remind "refactor the swagger definitions"
ctx remind add "check CI after the deploy" --after 2026-02-25
```

**Arguments**:

- `text`: The reminder message (verbatim)

**Flags**:

| Flag      | Short | Description                                |
|-----------|-------|--------------------------------------------|
| `--after` | `-a`  | Don't surface until this date (YYYY-MM-DD) |

**Examples**:

```bash
ctx remind "refactor the swagger definitions"
ctx remind "check CI after the deploy" --after 2026-02-25
```

#### `ctx remind list`

List all pending reminders. Date-gated reminders that aren't yet due
are annotated with `(after DATE, not yet due)`.

```bash
ctx remind list
```

**Aliases**: `ls`

#### `ctx remind dismiss`

Remove a reminder by ID, or remove all reminders with `--all`.

```bash
ctx remind dismiss <id>
ctx remind dismiss --all
```

**Arguments**:

- `id`: Reminder ID (shown in `list` output)

**Flags**:

| Flag    | Description              |
|---------|--------------------------|
| `--all` | Dismiss all reminders    |

**Aliases**: `rm`

**Examples**:

```bash
ctx remind dismiss 3
ctx remind dismiss --all
```

---
```

### 2. User-Facing Docs

#### Scratchpad recipe update (`docs/recipes/scratchpad-with-claude.md`)

Update the "When to Use Scratchpad vs Context Files" table. Add a row
for reminders and update the decision guide:

```markdown
| Situation                                                  | Use                |
|------------------------------------------------------------|--------------------|
| Temporary reminders ("*check X after deploy*")             | **Scratchpad**     |
| Session-start reminders ("*remind me next session*")       | **`ctx remind`**   |
| Working values during debugging (ports, endpoints, counts) | **Scratchpad**     |
| Sensitive tokens or API keys (short-term storage)          | **Scratchpad**     |
| Quick notes that don't fit anywhere else                   | **Scratchpad**     |
| Work items with completion tracking                        | **TASKS.md**       |
| Trade-offs between alternatives with rationale             | **DECISIONS.md**   |
| Reusable lessons with context/lesson/application           | **LEARNINGS.md**   |
| Codified patterns and standards                            | **CONVENTIONS.md** |
```

Add to the decision guide:

```markdown
* If you want a message relayed verbatim at the next session start,
  it belongs in `ctx remind`.
```

#### Session lifecycle recipe (`docs/recipes/session-lifecycle.md`)

Mention reminders in the "Load" phase description. Reminders surface
automatically via the `check-reminders` hook — users don't need to do
anything. Add a note like:

```markdown
If you have pending reminders (`ctx remind`), they are relayed
automatically at the start of each session.
```

#### About page (`docs/home/about.md`)

Add `ctx remind` to the feature list if one exists.

### 3. Recipe: Session Reminders (`docs/recipes/session-reminders.md`)

New recipe file:

```markdown
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
```

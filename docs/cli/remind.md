---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Remind
icon: lucide/bell-ring
---

![ctx](../images/ctx-banner.png)

## `ctx remind`

Session-scoped reminders that surface at session start.
Reminders are stored verbatim and relayed verbatim: no
summarization, no categories.

When invoked with a text argument and no subcommand, adds a
reminder.

```bash
ctx remind "text"
ctx remind <subcommand>
```

### `ctx remind add`

Add a reminder. This is the default action: `ctx remind "text"`
and `ctx remind add "text"` are equivalent.

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

### `ctx remind list`

List all pending reminders. Date-gated reminders that aren't
yet due are annotated with `(after DATE, not yet due)`.

**Examples**:

```bash
ctx remind list
ctx remind ls            # alias
```

**Aliases**: `ls`

### `ctx remind dismiss`

Remove one or more reminders by ID, or remove all with
`--all`. Supports individual IDs and ranges.

```bash
ctx remind dismiss <id> [id...]
ctx remind dismiss --all
```

**Arguments**:

- `id`: One or more reminder IDs (e.g., `3`, `3 5-7`)

**Flags**:

| Flag    | Description              |
|---------|--------------------------|
| `--all` | Dismiss all reminders    |

**Aliases**: `rm`

**Examples**:

```bash
ctx remind dismiss 3
ctx remind dismiss 3 5-7
ctx remind dismiss --all
```

### `ctx remind normalize`

Reassign reminder IDs as a contiguous sequence 1..N, closing
any gaps left by dismissals.

**Examples**:

```bash
ctx remind normalize
```

**See also**: [Session Reminders recipe](../recipes/session-reminders.md).

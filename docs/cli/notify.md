---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Notify
icon: lucide/bell
---

![ctx](../images/ctx-banner.png)

## `ctx notify`

Send fire-and-forget webhook notifications from skills, loops,
and hooks.

```bash
ctx notify --event <name> [--session-id <id>] "message"
```

**Flags**:

| Flag           | Short | Description                    |
|----------------|-------|--------------------------------|
| `--event`      | `-e`  | Event name (required)          |
| `--session-id` | `-s`  | Session ID (optional)          |

**Behavior**:

- No webhook configured: silent no-op (exit 0)
- Webhook set but event not in `events` list: silent no-op (exit 0)
- Webhook set and event matches: fire-and-forget HTTP POST
- HTTP errors silently ignored (no retry)

**Examples**:

```bash
ctx notify --event loop "Loop completed after 5 iterations"
ctx notify -e nudge -s session-abc "Context checkpoint at prompt #20"
```

### `ctx notify setup`

Configure the webhook URL interactively. The URL is encrypted
with AES-256-GCM using the encryption key and stored in
`.context/.notify.enc`.

**Examples**:

```bash
ctx notify setup
```

The encrypted file is safe to commit. The key (`~/.ctx/.ctx.key`)
lives outside the project and is never committed.

### `ctx notify test`

Send a test notification and report the HTTP response status.

**Examples**:

```bash
ctx notify test
```

**Payload format** (JSON POST):

```json
{
  "event": "loop",
  "message": "Loop completed after 5 iterations",
  "session_id": "abc123-...",
  "timestamp": "2026-02-22T14:30:00Z",
  "project": "ctx"
}
```

| Field        | Type   | Description                           |
|--------------|--------|---------------------------------------|
| `event`      | string | Event name from `--event` flag        |
| `message`    | string | Notification message                  |
| `session_id` | string | Session ID (omitted if empty)         |
| `timestamp`  | string | UTC RFC3339 timestamp                 |
| `project`    | string | Project directory name                |

**See also**: [Webhook Notifications recipe](../recipes/webhook-notifications.md).

---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Hook
icon: lucide/webhook
---

![ctx](../images/ctx-banner.png)

### `ctx hook`

Manage hook-related settings: messages, notifications, pause/resume,
and event log.

```bash
ctx hook <subcommand> [flags]
```

## Subcommands

| Subcommand                        | Description                                    |
|-----------------------------------|------------------------------------------------|
| `ctx hook message list`           | Show all hook messages with override status     |
| `ctx hook message show <h> <v>`   | Print the effective message template            |
| `ctx hook message edit <h> <v>`   | Copy default to `.context/` for editing         |
| `ctx hook message reset <h> <v>`  | Delete user override, revert to default         |
| `ctx hook notify [message]`       | Send a webhook notification                     |
| `ctx hook notify setup`           | Configure and encrypt webhook URL               |
| `ctx hook notify test`            | Send a test notification                        |
| `ctx hook pause`                  | Pause all context hooks for this session        |
| `ctx hook resume`                 | Resume paused context hooks                     |
| `ctx hook event`                  | Query the local hook event log                  |

## Examples

```bash
# View and manage hook messages
ctx hook message list
ctx hook message show qa-reminder gate
ctx hook message edit qa-reminder gate

# Webhook notifications
ctx hook notify setup
ctx hook notify --event loop "Loop completed"

# Pause/resume hooks
ctx hook pause
ctx hook resume

# Browse event log
ctx hook event --last 20
ctx hook event --hook qa-reminder --json
```

**See also**:
[Customizing Hook Messages](../recipes/customizing-hook-messages.md) |
[Webhook Notifications](../recipes/webhook-notifications.md) |
[Pausing Context Hooks](../recipes/session-pause.md) |
[System Hooks Audit](../recipes/system-hooks-audit.md)

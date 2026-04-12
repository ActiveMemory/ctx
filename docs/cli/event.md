---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Event
icon: lucide/scroll-text
---

![ctx](../images/ctx-banner.png)

### `ctx event`

Query the local hook event log. Requires `event_log: true` in
`.ctxrc`. Reads events from `.context/state/events.jsonl` and outputs
them in a human-readable table or raw JSONL format.

All filter flags combine with AND logic.

```bash
ctx event [flags]
```

**Flags**:

| Flag        | Description                                |
|-------------|--------------------------------------------|
| `--hook`    | Filter by hook name                        |
| `--session` | Filter by session ID                       |
| `--event`   | Filter by event type (`relay`, `nudge`)    |
| `--last`    | Show last N events (default: 50)           |
| `--json`    | Output raw JSONL (for piping to `jq`)      |
| `--all`     | Include rotated log file                   |

**Examples**:

```bash
ctx event                                        # recent events
ctx event --hook check-context-size --last 10    # one hook, last 10
ctx event --json | jq '.hook'                    # pipe to jq
ctx event --session abc123                       # filter by session
```

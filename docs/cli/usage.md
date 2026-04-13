---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Usage
icon: lucide/bar-chart-3
---

![ctx](../images/ctx-banner.png)

### `ctx usage`

Display per-session token usage statistics from the local stats
JSONL files written by the `heartbeat` hook. By default, shows the
last 20 entries across all sessions. Use `--follow` to stream new
entries as they arrive (like `tail -f`).

```bash
ctx usage [flags]
```

**Flags**:

| Flag              | Description                                    |
|-------------------|------------------------------------------------|
| `-f`, `--follow`  | Stream new entries as they arrive              |
| `-s`, `--session` | Filter by session ID (prefix match)            |
| `-n`, `--last`    | Show last N entries (default: 20)              |
| `-j`, `--json`    | Output raw JSONL                               |

**Examples**:

```bash
ctx usage                     # Last 20 entries across all sessions
ctx usage --follow            # Live stream (like tail -f)
ctx usage --session abc123    # Filter to one session
ctx usage --last 100 --json   # Last 100 as raw JSONL
```

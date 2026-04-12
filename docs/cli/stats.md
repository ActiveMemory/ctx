---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Stats
icon: lucide/bar-chart-3
---

![ctx](../images/ctx-banner.png)

### `ctx stats`

Display per-session token usage statistics from the local stats
JSONL files written by the `heartbeat` hook. By default, shows the
last 20 entries across all sessions. Use `--follow` to stream new
entries as they arrive (like `tail -f`).

The name stays plural by convention: `stats` is an idiomatic
abbreviation of "statistics" in CLI tooling, not a countable noun.

```bash
ctx stats [flags]
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
ctx stats                     # Last 20 entries across all sessions
ctx stats --follow            # Live stream (like tail -f)
ctx stats --session abc123    # Filter to one session
ctx stats --last 100 --json   # Last 100 as raw JSONL
```

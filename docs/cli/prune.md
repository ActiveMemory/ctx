---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Prune
icon: lucide/trash-2
---

![ctx](../images/ctx-banner.png)

### `ctx prune`

Remove per-session state files from `.context/state/` that are older
than the specified age. Session state files are identified by UUID
suffixes (`context-check-<session-id>`, `heartbeat-<session-id>`, and
similar). Global files without session IDs (`events.jsonl`,
`memory-import.json`, and other non-per-session markers) are always
preserved.

```bash
ctx prune [flags]
```

**Flags**:

| Flag        | Description                                      |
|-------------|--------------------------------------------------|
| `--days`    | Prune files older than this many days (default: 7) |
| `--dry-run` | Show what would be pruned without deleting      |

**Examples**:

```bash
ctx prune                 # Prune files older than 7 days
ctx prune --days 3        # Prune files older than 3 days
ctx prune --dry-run       # Preview without deleting
```

See [State maintenance](../recipes/state-maintenance.md) for the
recommended cadence and automation recipe.

---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Resume
icon: lucide/play
---

![ctx](../images/ctx-banner.png)

## `ctx hook resume`

Resume context hooks after a pause. Silent no-op if not paused.

```bash
ctx hook resume [flags]
```

**Flags**:

| Flag           | Description                         |
|----------------|-------------------------------------|
| `--session-id` | Session ID (overrides stdin)        |

**Example**:

```bash
ctx hook resume
```

**See also**:

- [`ctx hook pause`](pause.md) — the matching pause command
- [Pausing Context Hooks recipe](../recipes/session-pause.md)

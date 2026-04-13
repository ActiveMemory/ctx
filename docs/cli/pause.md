---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Pause
icon: lucide/pause
---

![ctx](../images/ctx-banner.png)

## `ctx hook pause`

Pause all context nudge and reminder hooks for the current
session. Security hooks (dangerous command blocking) and
housekeeping hooks still fire.

```bash
ctx hook pause [flags]
```

**Flags**:

| Flag           | Description                         |
|----------------|-------------------------------------|
| `--session-id` | Session ID (overrides stdin)        |

**Example**:

```bash
# Pause hooks for a quick investigation
ctx hook pause

# Resume when ready
ctx hook resume
```

**See also**:

- [`ctx hook resume`](resume.md) — the matching resume command
- [Pausing Context Hooks recipe](../recipes/session-pause.md)

---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Sysinfo
icon: lucide/cpu
---

![ctx](../images/ctx-banner.png)

### `ctx sysinfo`

Display a snapshot of system resources (memory, swap, disk, load)
with threshold-based alert severities. Mirrors what the
`check-resource` hook plumbing monitors in the background — but this
command prints the full report at any severity level, not only at
DANGER.

```bash
ctx sysinfo [flags]
```

**Flags**:

| Flag     | Description             |
|----------|-------------------------|
| `--json` | Output in JSON format   |

**Alert thresholds**:

| Resource | WARNING | DANGER |
|----------|---------|--------|
| Memory   | ≥ 75%   | ≥ 90%  |
| Swap     | ≥ 50%   | ≥ 75%  |
| Disk     | ≥ 85%   | ≥ 95%  |
| Load     | ≥ 1.0x CPUs | ≥ 1.5x CPUs |

**Examples**:

```bash
ctx sysinfo                  # Human-readable table
ctx sysinfo --json           # Structured output
```

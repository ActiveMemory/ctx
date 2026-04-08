---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Hub
icon: lucide/network
---

## `ctx hub`

Manage a running shared context hub cluster.

### `ctx hub status`

Show cluster status: role, entry count, and connected projects.

```bash
ctx hub status
```

### `ctx hub peer`

Add or remove peers from the cluster at runtime.

```bash
ctx hub peer add host2:9901
ctx hub peer remove host2:9901
```

### `ctx hub stepdown`

Transfer leadership to another node gracefully. Use before
taking a node offline for maintenance.

```bash
ctx hub stepdown
```

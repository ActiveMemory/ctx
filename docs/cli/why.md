---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Why
icon: lucide/help-circle
---

![ctx](../images/ctx-banner.png)

## `ctx why`

Read `ctx`'s philosophy documents directly in the terminal.

```bash
ctx why [DOCUMENT]
```

**Documents**:

| Name         | Description                                  |
|--------------|----------------------------------------------|
| `manifesto`  | The `ctx` Manifesto: creation, not code      |
| `about`      | About `ctx`: what it is and why it exists    |
| `invariants` | Design invariants: properties that must hold |

**Examples**:

```bash
# Interactive numbered menu
ctx why

# Show a specific document
ctx why manifesto
ctx why about
ctx why invariants

# Pipe to a pager
ctx why manifesto | less
```

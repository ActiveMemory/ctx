---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Guide
icon: lucide/book-open
---

![ctx](../images/ctx-banner.png)

## `ctx guide`

Quick-reference cheat sheet for common ctx commands and skills.

```bash
ctx guide [flags]
```

**Flags**:

| Flag         | Description                  |
|--------------|------------------------------|
| `--skills`   | Show available skills        |
| `--commands` | Show available CLI commands  |

**Example**:

```bash
# Show the full cheat sheet
ctx guide

# Skills only
ctx guide --skills

# Commands only
ctx guide --commands
```

Works without initialization (no `.context/` required). Useful
for a printable one-pager when onboarding to a project.

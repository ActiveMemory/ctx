---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Bootstrap
icon: lucide/compass
---

![ctx](../images/ctx-banner.png)

### `ctx bootstrap`

Print the resolved context directory path so AI agents can anchor
their session. The default output lists the context directory, the
tracked context files, and a short health snapshot. `--quiet` prints
just the path; `--json` produces structured output for automation.

This is the command agents are instructed to run first in their
session-start procedure — it is the authoritative answer to "where
does this project's context live?".

```bash
ctx bootstrap [flags]
```

**Flags**:

| Flag              | Description                                      |
|-------------------|--------------------------------------------------|
| `-q`, `--quiet`   | Output only the context directory path          |
| `--json`          | Output in JSON format                             |

**Examples**:

```bash
ctx bootstrap                 # Text output for agents
ctx bootstrap -q              # Just the context directory path
ctx bootstrap --json          # Structured output for automation
```

**Scripting tip**: `CTX_DIR=$(ctx bootstrap -q)` is the canonical way
for skills and scripts to find the project's context directory
without hardcoding `.context/`.

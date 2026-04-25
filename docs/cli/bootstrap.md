---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: System Bootstrap
icon: lucide/compass
---

![ctx](../images/ctx-banner.png)

### `ctx system bootstrap`

Print the resolved context directory path so AI agents can anchor
their session. The default output lists the context directory, the
tracked context files, and a short health snapshot. `--quiet` prints
just the path; `--json` produces structured output for automation.

This is a hidden, agent-only command that agents are instructed to
run first in their session-start procedure; it is the authoritative
answer to "where does this project's context live?".

```bash
ctx system bootstrap [flags]
```

**Flags**:

| Flag              | Description                                      |
|-------------------|--------------------------------------------------|
| `-q`, `--quiet`   | Output only the context directory path          |
| `--json`          | Output in JSON format                             |

**Examples**:

```bash
ctx system bootstrap                 # Text output for agents
ctx system bootstrap -q              # Just the context directory path
ctx system bootstrap --json          # Structured output for automation
```

**Note**: `-q` prints just the resolved directory path. See
[Activating a Context Directory](../recipes/activating-context.md)
if you hit a "*no context directory specified*" error.

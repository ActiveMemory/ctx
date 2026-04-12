---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Loop
icon: lucide/refresh-cw
---

![ctx](../images/ctx-banner.png)

## `ctx loop`

Generate a shell script for running an autonomous loop.

An autonomous loop continuously runs an AI assistant with the
same prompt until a completion signal is detected, enabling
iterative development where the AI builds on its previous work.

```bash
ctx loop [flags]
```

**Flags**:

| Flag                    | Short | Description                              | Default            |
|-------------------------|-------|------------------------------------------|--------------------|
| `--tool <tool>`         | `-t`  | AI tool: `claude`, `aider`, or `generic` | `claude`           |
| `--prompt <file>`       | `-p`  | Prompt file to use                       | `.context/loop.md` |
| `--max-iterations <n>`  | `-n`  | Maximum iterations (0 = unlimited)       | `0`                |
| `--completion <signal>` | `-c`  | Completion signal to detect              | `SYSTEM_CONVERGED` |
| `--output <file>`       | `-o`  | Output script filename                   | `loop.sh`          |

**Examples**:

```bash
# Generate loop.sh for Claude Code
ctx loop

# Generate for Aider with custom prompt
ctx loop --tool aider --prompt TASKS.md

# Limit to 10 iterations
ctx loop --max-iterations 10

# Output to custom file
ctx loop -o my-loop.sh
```

**Running the generated loop**:

```bash
ctx loop
chmod +x loop.sh
./loop.sh
```

**See also**: [Autonomous Loops](../operations/autonomous-loop.md)
for the full workflow.

---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Setup
icon: lucide/toy-brick
---

![ctx](../images/ctx-banner.png)

## `ctx setup`

Generate AI tool integration configuration.

```bash
ctx setup <tool> [flags]
```

**Flags**:

| Flag      | Short | Description                                                                 |
|-----------|-------|-----------------------------------------------------------------------------|
| `--write` | `-w`  | Write the generated config to disk (e.g. `.github/copilot-instructions.md`) |

**Supported tools**:

| Tool          | Description                                  |
|---------------|----------------------------------------------|
| `claude-code` | Redirects to plugin install instructions     |
| `cursor`      | Cursor IDE                                   |
| `kiro`        | Kiro IDE                                     |
| `cline`       | Cline (VS Code extension)                    |
| `aider`       | Aider CLI                                    |
| `copilot`     | GitHub Copilot                               |
| `windsurf`    | Windsurf IDE                                 |

!!! note "Claude Code uses the plugin system"
    Claude Code integration is now provided via the `ctx` plugin.
    Running `ctx setup claude-code` prints plugin install instructions.

**Examples**:

```bash
# Print hook instructions to stdout
ctx setup cursor
ctx setup aider

# Generate and write .github/copilot-instructions.md
ctx setup copilot --write

# Generate MCP config and sync steering files
ctx setup kiro --write
ctx setup cursor --write
ctx setup cline --write
```

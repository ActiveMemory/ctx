---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Backup
icon: lucide/archive
---

![ctx](../images/ctx-banner.png)

### `ctx backup`

Create timestamped tar.gz archives of project context and/or global
Claude Code data. Optionally copies archives to an SMB share via GVFS.

```bash
ctx backup [flags]
```

**Flags**:

| Flag      | Description                                        |
|-----------|----------------------------------------------------|
| `--scope` | Backup scope: `project`, `global`, or `all` (default: `all`) |
| `--json`  | Output results as JSON                             |

**Scopes**:

| Scope     | What's archived                                |
|-----------|-----------------------------------------------|
| `project` | `.context/`, `.claude/`, `ideas/`, `~/.bashrc` |
| `global`  | `~/.claude/` (excludes `todos/`)              |
| `all`     | Both project and global (default)             |

**Environment**:

| Variable                | Purpose                                          |
|-------------------------|--------------------------------------------------|
| `CTX_BACKUP_SMB_URL`    | SMB share URL (e.g. `smb://host/share`)          |
| `CTX_BACKUP_SMB_SUBDIR` | Subdirectory on share (default: `ctx-sessions`) |

**Examples**:

```bash
ctx backup                       # Back up everything (default: all)
ctx backup --scope project       # Project context only
ctx backup --scope global        # Global Claude data only
ctx backup --scope all --json    # Both, JSON output
```

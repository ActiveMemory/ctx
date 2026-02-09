---
name: backup
description: "Backup project context and global Claude data to SMB share. Use before risky operations, at end of session, or on request."
allowed-tools: Bash(make backup*), Bash(ls /tmp/ctx-backup*)
---

Backup `.context/`, `.claude/`, `ideas/`, and `~/.claude/` to
the configured SMB share.

## When to Use

- Before risky operations (major refactors, dependency upgrades)
- At the end of a productive session
- When the user explicitly asks for a backup
- Before switching branches with uncommitted context changes

## When NOT to Use

- When `CTX_BACKUP_SMB_URL` is not configured (the script will
  error — tell the user to set it up)
- Multiple times in quick succession with no changes in between

## Usage Examples

```text
/backup
/backup project
/backup global
/backup all
```

## Arguments

| Argument  | What it backs up                              |
|-----------|-----------------------------------------------|
| (none)    | Same as `all`                                 |
| `project` | Project context only (`.context/`, `.claude/`, `ideas/`) |
| `global`  | Global Claude data only (`~/.claude/`)        |
| `all`     | Both project and global                       |

## Execution

Based on the argument, run the appropriate make target:

```bash
# For "project"
make backup

# For "global"
make backup-global

# For "all" or no argument
make backup-all
```

## Process

1. Parse the argument (default to `all` if none provided)
2. Run the appropriate `make` target
3. Report the archive path and size from the output
4. Confirm success to the user

## Quality Checklist

- [ ] The make target completed without errors
- [ ] Archive size is reported to the user
- [ ] If the SMB share was not mounted, the error is clearly
      communicated

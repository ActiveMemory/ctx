# Spec: Convert Shell Hook Scripts to `ctx system` Subcommands

**Status**: Implementing
**Added**: 2026-02-23

## Problem

The `.claude/hooks/` directory contains two shell scripts —
`block-dangerous-commands.sh` and `check-backup-age.sh` — that are the last
non-Go hooks in the project. They depend on `jq`, are harder to test, and
don't benefit from shared state infrastructure (`readInput()`,
`secureTempDir()`, `isDailyThrottled()`, `notify.Send()`).

## Solution

Port both scripts to Go `ctx system` subcommands following the same pattern as
`block-non-path-ctx`, `check-version`, and `check-resources`. Update
`settings.local.json` to invoke the Go commands. Delete the shell scripts.

## Scope

- **Project-local hooks only** — they stay in `settings.local.json`, never
  added to `internal/assets/claude/hooks/hooks.json`
- No user-facing docs needed (project-internal)
- The absolute-path-to-ctx regex in `block-dangerous-commands.sh` is dropped
  (already handled by `block-non-path-ctx`)

## Commands

### `ctx system block-dangerous-commands`

PreToolUse (Bash) hook. Regex safety net for commands that the deny-list
cannot express:

| Pattern | Reason |
|---------|--------|
| Mid-command `sudo` after `&&`, `\|\|`, `;` | No password access |
| Mid-command `git push` after `&&`, `\|\|`, `;` | Requires explicit approval |
| `cp`/`mv` to bin directories | Agent must not install binaries |
| `cp`/`install` to `~/.local/bin` | Overrides system ctx in /usr/local/bin |

### `ctx system check-backup-age`

UserPromptSubmit hook. Warns when SMB backup is stale or unmounted:

- Daily throttled via `isDailyThrottled()` + `touchFile()`
- Check 1: If `CTX_BACKUP_SMB_URL` set, derive GVFS mount path and check
  directory existence
- Check 2: Check `~/.local/state/ctx-last-backup` marker age against 2-day
  threshold
- Output: VERBATIM relay with warning box

## Files

| Action | File |
|--------|------|
| New | `internal/cli/system/block_dangerous_commands.go` |
| New | `internal/cli/system/block_dangerous_commands_test.go` |
| New | `internal/cli/system/check_backup_age.go` |
| New | `internal/cli/system/check_backup_age_test.go` |
| Modify | `internal/cli/system/system.go` |
| Modify | `internal/cli/system/doc.go` |
| Modify | `.claude/settings.local.json` |
| Delete | `.claude/hooks/block-dangerous-commands.sh` |
| Delete | `.claude/hooks/check-backup-age.sh` |

# Deprecate and Remove `ctx backup`

## Problem

`ctx backup` is an environment-specific feature masquerading as a
core command. It assumes SMB/GVFS mounts, fires nag hooks for
users who never configured it, and solves a problem that belongs
to the OS/infrastructure layer, not the application layer.

Since its introduction, `ctx hub` has landed as the real answer to
"how do I make my context survive across machines." Meanwhile,
`ctx backup` has accumulated tech debt (the Broadcom mirror repo
broke `sync-to-asgard`, GVFS is Linux-only, the check_backup_age
hook fires even in projects that never use backup) without
delivering proportional value.

No user other than the project maintainer has ever configured it.

## Decision

Remove `ctx backup` entirely. Replace with a documentation runbook
that covers backup strategies for users who need file-level
backups beyond what hub provides.

## Reasoning

### 1. Hub already solves distributed persistence

`ctx hub` replicates decisions, learnings, conventions, and tasks
across machines. For the data that matters most (knowledge files),
hub is the persistence mechanism. Backup of `.context/` is
redundant if you run a hub.

### 2. Backup is inherently environment-specific

SMB, NFS, S3, rsync, Time Machine, Borg, restic: every user has a
different backup story. `ctx backup` picked SMB via GVFS, which is
a narrow choice that does not generalize. Adding
`CTX_BACKUP_MOUNT_PATH` (the original task) would just patch one
assumption with another. You would chase mount strategies forever.

### 3. Wrong layer

Backup is an OS/infrastructure concern. `ctx` manages context
files; the user's backup tool backs up files. Asking `ctx` to also
be a backup tool is scope creep that adds maintenance burden for a
feature that `rsync`, `cp`, or a cron job handles better.

### 4. Active maintenance cost

- The Broadcom mirror repo issue forced disabling `sync-to-asgard`
  (reminder [12] in the project)
- GVFS dependency is Linux-only, breaking macOS/Windows users
- `check_backup_age` hook fires for everyone, even users who never
  configured backup, creating noise
- SMB mount code (`internal/exec/gio/`) is dead weight on macOS

## Scope

### What gets removed

**CLI command** (3 files):
- `internal/cli/backup/cmd.go`
- `internal/cli/backup/run.go`
- `internal/cli/backup/doc.go`

**Hook** (3 files):
- `internal/cli/system/cmd/check_backup_age/cmd.go`
- `internal/cli/system/cmd/check_backup_age/run.go`
- `internal/cli/system/cmd/check_backup_age/doc.go`
- Hook message template: `internal/assets/hooks/messages/check-backup-age/`
- Hook registration in `internal/cli/system/system.go`

**Core archive/SMB** (5 files):
- `internal/cli/system/core/archive/backup.go`
- `internal/cli/system/core/archive/smb.go`
- `internal/cli/system/core/archive/types.go`
- `internal/cli/system/core/archive/archive.go`
- `internal/cli/system/core/archive/doc.go`

**GIO/mount integration**:
- `internal/exec/gio/mount.go` (entire package if only used by backup)

**Config constants** (4 files):
- `internal/config/archive/backup.go`
- `internal/config/archive/var.go`
- `internal/config/archive/archive.go`
- `internal/config/archive/doc.go`

**Env vars** (in `internal/config/env/env.go`):
- `CTX_BACKUP_SMB_URL`
- `CTX_BACKUP_SMB_SUBDIR`

**Embed keys** (scattered across `internal/config/embed/`):
- `cmd/backup.go` (UseBackup, DescKeyBackup)
- `cmd/system.go` (UseSystemCheckBackupAge, DescKeySystemCheckBackupAge)
- `flag/backup.go` (DescKeyBackupJson, DescKeyBackupScope)
- `text/backup.go` (all DescKeyBackup* and DescKeyWriteBackup* keys)
- `text/err_backup.go` (all DescKeyErrBackup* keys)

**Error package**:
- `internal/err/backup/` (entire package)

**Write/output package**:
- `internal/write/backup/` (entire package)

**Entity types** (in `internal/entity/system.go`):
- `ArchiveEntry` struct
- `BackupResult` struct

**Bootstrap registration** (in `internal/bootstrap/group.go`):
- Remove `backup.Cmd` from `GroupRuntime`
- Remove import

**YAML text files** (remove backup-related entries):
- `internal/assets/commands/commands.yaml`
- `internal/assets/commands/text/write.yaml`
- Various error/text YAML files

**Skill**:
- `.claude/skills/_ctx-backup/` (entire directory)

**Documentation**:
- `docs/cli/backup.md` (remove entirely)
- `zensical.toml` (remove backup from nav)
- `docs/home/contributing.md` (remove backup references)
- `docs/home/common-workflows.md` (remove backup workflow)
- `docs/cli/index.md` (remove backup from command index)
- `docs/cli/system.md` (remove check-backup-age reference)
- `docs/recipes/hook-sequence-diagrams.md` (remove backup sequence)
- `docs/recipes/customizing-hook-messages.md` (remove backup example)
- `docs/recipes/hook-output-patterns.md` (remove backup example)
- Various hub docs (remove backup tier references, update to
  recommend external backup tools)

### What gets created

**Runbook**: `docs/operations/runbooks/backup-strategy.md`

Contents:
- What to back up: `.context/`, `.claude/`, `~/.ctx/`
- How hub reduces backup needs (knowledge files are replicated)
- What hub does NOT back up (journal, scratchpad, session logs)
- Example strategies:
  - cron + rsync to NAS/external drive
  - cron + cp to cloud-synced directory
  - Time Machine / system backup (macOS)
  - Borg/restic for versioned backups
- When you still need file-level backup even with hub

### What stays

- `internal/cli/initialize/core/backup/file.go`: this creates
  `.bak` copies during `ctx init --force`. It is NOT part of
  `ctx backup`; it is init's config backup mechanism. Keep it.
- Hub backup/restore documentation in `docs/operations/hub.md`:
  this is about backing up the hub data directory, not `ctx backup`.

## Implementation Order

1. Create the backup-strategy runbook first (so users have a
   migration path before the feature disappears)
2. Remove the `check_backup_age` hook and its registration (this
   is the user-visible annoyance; removing it first gives
   immediate relief)
3. Remove the CLI command, core packages, config constants, error
   package, write package, entity types, bootstrap registration
4. Remove the skill
5. Update all documentation (remove pages, update cross-refs)
6. Remove YAML text entries
7. Run `make audit` to catch any dangling references
8. Update TASKS.md: mark the SMB mount task as skipped, add
   completion note

## Migration

Users who relied on `ctx backup` (currently: only the project
maintainer) should:

1. Replace with `rsync` or `cp` in a cron job
2. Use `ctx hub` for cross-machine knowledge persistence
3. Follow the backup-strategy runbook for file-level needs

## Risks

- **Low**: no external users depend on this feature
- **Medium**: the project maintainer's own backup workflow breaks.
  Mitigated by writing the runbook first and setting up a
  replacement cron job before removing the command.

# Backup Strategy

`ctx backup` was removed. File-level backup is not `ctx`'s
responsibility; your OS or a dedicated backup tool handles it
better and without locking you into a specific mount strategy.

This runbook explains what to back up, how `ctx hub` reduces the
surface, and what options exist for the rest.

## What To Back Up

Per project:

- `.context/`: all context files, journal, state, scratchpad.
- `.claude/`: Claude Code settings, hooks, skills specific to the
  project. Skip this entry when it lives in git; the repo is the
  backup.

Per user:

- `~/.ctx/`: global config, the encryption key (`~/.ctx/.ctx.key`),
  hub data directory (if running a local hub).

## How Hub Reduces Backup Needs

`ctx hub` replicates the knowledge surface across machines:

- `DECISIONS.md`
- `LEARNINGS.md`
- `CONVENTIONS.md`
- `CONSTITUTION.md`
- `ARCHITECTURE.md`
- Task items promoted to hub

If you run `ctx hub` (as a server or by subscribing to someone
else's), the data that matters most survives losing any single
machine.

## What Hub Does *Not* Replicate

Hub is not a file-level backup. The following still live only on
the machine that produced them:

- Journal entries (`.context/journal/*.md`)
- Runtime state (`.context/state/*`)
- Session event log (`.context/events.jsonl`)
- Scratchpad (`.context/.pad`)
- Encrypted notify/webhook config (`.context/.notify.enc`)
- The encryption key itself (`~/.ctx/.ctx.key`)

If you need those to survive a disk failure, use a file-level
backup.

## Example Strategies

### 1. cron + rsync to NAS or External Drive

```cron
# Daily at 03:00, mirror ~/WORKSPACE and ~/.ctx to NAS
0 3 * * * rsync -a --delete \
    --exclude='node_modules' \
    --exclude='dist' \
    --exclude='.context/state' \
    ~/WORKSPACE/ /mnt/nas/backup/workspace/
0 3 * * * rsync -a --delete ~/.ctx/ /mnt/nas/backup/ctx-global/
```

Adjust excludes for the trash you don't want to back up. The
`.context/state/` dir is ephemeral per-session; skip it.

### 2. cron + cp to a Cloud-Synced Directory

iCloud Drive, Dropbox, or any directory watched by a sync client:

```cron
0 3 * * * cp -a ~/WORKSPACE/some-project/.context \
    ~/CloudDrive/ctx-backups/some-project/$(date +\%Y-\%m-\%d)
```

Daily snapshots, cloud provider handles the replication.

### 3. Time Machine (macOS)

If you already run Time Machine, ensure `~/WORKSPACE` and `~/.ctx`
are not in its exclusion list. Time Machine handles versioning;
you get point-in-time recovery for free.

### 4. Borg or restic for Versioned Backups

For deduplicated, versioned, encrypted backups:

```bash
# Borg init (once)
borg init --encryption=repokey /mnt/nas/borg-ctx

# Daily backup
borg create /mnt/nas/borg-ctx::'ctx-{now}' \
    ~/WORKSPACE ~/.ctx \
    --exclude '*/node_modules' \
    --exclude '*/.context/state'
```

Use `restic` if you prefer S3-compatible targets.

## When You Still Need File-Level Backup Even With Hub

- **Journal**: session histories are local-only until exported.
- **Scratchpad**: private notes, encrypted locally.
- **Encryption key**: losing `~/.ctx/.ctx.key` means losing access
  to every encrypted file in every project.
- **Non-hub projects**: projects that never called `ctx hub
  register` have zero cross-machine persistence.

For these, pick one strategy above and forget about it.

## Why ctx No Longer Ships a Backup Command

Backup is inherently environment-specific: SMB, NFS, S3, rsync,
Time Machine, Borg, restic. Every user has a different story. The
previous `ctx backup` picked SMB via GVFS, which was Linux-only and
narrow. Chasing mount strategies would never generalize.

Hub is the right answer for the data `ctx` owns (knowledge). For
everything else, your OS or a dedicated backup tool is the right
layer.

# Pad Undo and Snapshot Safety Net

Every destructive `ctx pad` operation should leave a recoverable
snapshot, and `ctx pad undo` should reverse the most recent
mutation. The pad is encrypted, so accidental deletion currently
has no readable trace and no in-product recovery path.

## Problem

`ctx pad rm <entry>` removes an entry with no confirmation and no
backup. The entry is then gone for good: the pad blob is
re-encrypted in place via `store.WriteEntries`, the previous
ciphertext is overwritten by the atomic rename, and the user has
no in-product recovery surface.

The encrypted nature of the pad makes this worse than a typical
delete: the user cannot peek at an entry's content without
reading it (which is what they were avoiding in the first place
by skipping the read step before `rm`). The accident pattern is
mundane — confusing one entry tag for another, hitting the wrong
slug in a short list, or removing the wrong line during an
`edit` save — and the resulting recovery dance is severe:

1. Locate the most recent off-host backup of the pad blob.
2. Stash the current `.scratchpad.enc` and drop the backup in
   its place.
3. Read and extract the lost content.
4. Restore the stashed current blob.
5. Re-add the missing content as a new entry.
6. Clean up the downloaded backup so it doesn't drift.

That sequence presumes the user has working backup hygiene, can
locate a recent enough backup, and is willing to walk a six-step
ritual for what was a single fat-finger. The same pattern
applies to `edit` (saving over the wrong entry's body) and to
`mv` / `merge` / `normalize` / `tag` / `resolve` — anything that
re-writes the pad blob.

A confirmation prompt was considered and rejected: the user
explicitly prefers the current no-prompt UX. The safety net has
to live downstream of the command, not in front of it.

## Approach

Two mechanisms working together:

1. **Snapshot-on-mutate**: every call to `store.WriteEntries`
   that overwrites an existing pad copies the prior encrypted
   blob to a per-mutation history file under
   `.context/.scratchpad.history/` *before* the new ciphertext
   lands. Snapshot is the existing ciphertext byte-for-byte; no
   re-encryption needed.

2. **`ctx pad undo`**: a new subcommand that restores the most
   recent snapshot, displacing the current pad. The undo
   operation is itself a mutation that snapshots first, so
   `ctx pad undo; ctx pad undo` yields a redo.

The choke point is `internal/cli/pad/core/store/store.go:WriteEntries`.
Every mutating `ctx pad` subcommand routes through it; adding
the safety net there means no per-command changes.

## Behavior

### Happy Path

1. `ctx pad rm urgent-notes` runs.
2. `store.WriteEntries` is invoked with the new entry set.
3. Before the atomic rename lands the new ciphertext,
   `WriteEntries` copies the existing `.scratchpad.enc` to
   `.context/.scratchpad.history/<UTC-timestamp>-rm.enc`.
4. The new pad is written; the old one survives in history.
5. User realizes the mistake, runs `ctx pad undo`.
6. The most recent snapshot is promoted back to
   `.scratchpad.enc` (after itself being snapshotted as
   `<UTC-timestamp>-undo.enc`).
7. The deleted entry is back.

### First-Write Edge Case

If `.scratchpad.enc` does not exist when `WriteEntries` is
called (first ever entry), no snapshot is written. There is
nothing to preserve. The history directory is created lazily
on first real snapshot.

### No-op Edge Case

If `WriteEntries` is called with content that re-encrypts to a
byte-equivalent ciphertext, the snapshot is still written. The
predictability of "every write snapshots" beats the small
storage saving of trying to detect no-ops, and AES-GCM with a
fresh nonce won't produce byte-identical output anyway.

### Empty History

`ctx pad undo` with no snapshots prints `no pad history to
restore` and exits 0. Not an error — the empty case is
expected on a fresh project.

### Encryption Key Loss

If the per-machine key at `~/.ctx/.ctx.key` is gone or rotated,
snapshots are unreadable. This is the same failure mode as the
live pad and outside the scope of this safety net. Document it
in the undo help text.

## Retention Policy

Bounded ring buffer with two ceilings, whichever trips first:

- **Count cap**: keep the most recent 20 snapshots.
- **Age cap**: drop snapshots older than 30 days.

Both defaults exposed via `.ctxrc`:

```toml
[pad.history]
max_snapshots = 20      # 0 disables snapshotting entirely
max_age_days = 30       # 0 disables age-based pruning
```

`0` for `max_snapshots` is the opt-out for users who want the
current behavior. Setting both to high values is supported but
unrecommended (each snapshot is a full pad copy; storage grows
linearly).

Pruning runs at the tail of `WriteEntries`, after the new
snapshot lands. Pruning failure does not block the write; it
logs via `internal/log/warn` and continues.

## Interface

```
ctx pad undo                  # restore the most recent snapshot
ctx pad undo --list           # show available snapshots, newest first
ctx pad undo --to <slot>      # restore a specific snapshot by ID
ctx pad undo --prune          # force-run retention pruning now
ctx pad undo --clear          # delete all snapshots (with confirmation prompt)
```

`--list` output is one line per snapshot:

```
<slot>  <UTC timestamp>  <op>      <entries before>  <entries after>
abc1234 2026-05-24 14:30 rm        7                 6
def5678 2026-05-24 14:12 edit      6                 6
ghi9abc 2026-05-23 09:08 add       5                 6
```

`<slot>` is the first 7 chars of the snapshot file's hex digest
(stable across runs, short enough to type).

Entry counts come from a tiny sidecar `<slot>.meta.json`
written alongside the snapshot — counts only, no entry content
or tags (those live inside the encrypted blob).

The `--clear` confirmation prompt is the *only* prompt in this
feature; the safety net's whole point is no prompts on the hot
path, but mass-deleting the safety net itself deserves one.

## Files to Create / Modify

- `internal/cli/pad/core/store/history.go` — new file:
  snapshot writer, snapshot lister, snapshot restorer, retention
  pruner. All operations encrypted-blob-level (no decryption
  needed for snapshot/restore; decryption only needed for the
  count metadata, and even that can be deferred to the next
  `--list` invocation).
- `internal/cli/pad/core/store/store.go` —
  `WriteEntries` calls `history.SnapshotBefore` before the
  atomic rename and `history.Prune` after.
- `internal/cli/pad/cmd/undo/` — new subcommand package
  mirroring the layout of `cmd/rm/` etc.
- `internal/cli/pad/cmd/root/` — register the new
  subcommand on the `ctx pad` cobra tree.
- `internal/config/pad/history.go` — new file: path
  constants (`HistoryDirName = ".scratchpad.history"`),
  default retention constants, `.ctxrc` key paths.
- `internal/err/pad/pad.go` — new error constructors:
  `errPad.HistoryRead`, `errPad.HistoryWrite`,
  `errPad.HistoryRestore`, `errPad.UnknownSlot`.
- `internal/write/pad/history.go` — user-facing strings:
  "snapshot saved", "no pad history to restore",
  "restored snapshot from <timestamp>", "pruned N old
  snapshots", "clear all N snapshots? this cannot be
  undone".
- `internal/assets/claude/skills/ctx-pad/SKILL.md` —
  document the new `undo` command and the snapshot safety
  net so the agent recommends it.
- `docs/recipes/scratchpad-with-claude.md` — add a "If You
  Delete the Wrong Thing" subsection.
- `docs/cli/pad.md` (if present) — document `ctx pad undo`.

## Testing

### Unit

- `TestSnapshotBefore_FirstWriteWritesNoSnapshot` — empty pad
  state, `WriteEntries` runs, history dir stays empty.
- `TestSnapshotBefore_PreservesExactCiphertext` — seed a pad
  blob, mutate, assert the snapshot file's bytes equal the
  pre-mutation pad bytes.
- `TestPrune_RingBufferCountCap` — write 25 mutations, assert
  exactly 20 snapshots remain, oldest 5 gone.
- `TestPrune_AgeCap` — backdate snapshot mtimes via
  `os.Chtimes`, run prune, assert age-evicted snapshots are
  gone.
- `TestPrune_BothCapsZeroDisablesFeature` — set both caps to
  0 in config, mutate pad, assert no snapshot written.
- `TestUndo_RestoresMostRecent` — mutate, undo, assert pad
  content matches pre-mutation state.
- `TestUndo_IsItselfSnapshotted` — mutate → undo → undo,
  assert the second undo brings the post-mutation state back
  (redo via re-undo).
- `TestUndo_EmptyHistoryPrintsAndExitsZero` — fresh project,
  `undo` runs, exit code 0, expected message.
- `TestUndoTo_UnknownSlotErrors` — `undo --to bogus123`,
  assert `errPad.UnknownSlot`.

### Integration

- `TestPadRmUndoRoundTrip` — `ctx pad add foo "x" && ctx pad
  rm foo && ctx pad undo && ctx pad show foo` returns
  `"x"`.
- `TestPadEditUndoRoundTrip` — same shape for edit
  shrinkage.
- `TestPadUndoListAcrossOperations` — perform add / edit /
  mv / rm in sequence, run `undo --list`, assert four rows
  in newest-first order with correct `op` labels.
- `TestPadUndoClearConfirmation` — invoke `undo --clear`
  with a piped `n`, assert nothing was deleted; with `y`,
  assert history dir is empty.

## Non-Goals

- **Cross-machine snapshot sync.** Snapshots are local-only.
  The `scratchpad-sync` recipe handles encrypted pad sync
  between hosts; replicating per-host history would multiply
  storage and conflict surface. If a user wants
  cross-machine recovery, the existing sync flow plus the
  retention window on the source host is the path.
- **Per-entry undo.** Operations may span multiple entries
  (`merge`, `normalize`, `mv` of tags). Restoring the whole
  pad blob is the simpler and safer primitive; per-entry
  rollback can come later if there's demand.
- **Time-travel beyond the retention window.** Off-host
  backups remain the user's responsibility for recovery
  older than 30 days.
- **Confirmation prompts on destructive subcommands.** Out of
  scope per the user's explicit preference; the safety net
  exists precisely to make confirmations unnecessary.
- **Snapshotting on non-mutating reads** (`show`, `export`,
  `root`). These don't touch `WriteEntries`.

## Configuration

```toml
[pad.history]
max_snapshots = 20      # ring-buffer count cap; 0 disables
max_age_days = 30       # age cap in days; 0 disables age pruning
```

Defaults are written into a fresh `.ctxrc` by `ctx init` but
not required — absence falls back to the defaults above.

## Open Questions

1. **Should `ctx pad undo --to <slot>` reset the ring after
   the restore point, or preserve newer snapshots as
   "redo-able"?** Default leaning: preserve. The newer
   snapshots are already mutations that snapshot-on-mutate
   captured; tossing them on a `--to` restore loses redo
   capability and surprises users who expected non-
   destructive time travel. Decision deferred to first
   implementation; revisit if it complicates the snapshot
   index.
2. **Should the safety net cover `ctx pad export`?** Export
   doesn't mutate the pad but emits plaintext to disk. Out
   of scope for this spec; tracked separately if needed.

## Source

User request, 2026-05-24 session: pad `rm` is silent and
irreversible, and the recovery dance via off-host backups is
disproportionate to the original mistake. Snapshot-on-mutate
plus an `undo` subcommand was the agreed-on shape.

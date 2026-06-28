# Fix Connect Sync Lock TOCTOU Race

`internal/cli/connection/core/sync.loadState` guarded
`ctx connect sync` against concurrent execution with a
check-then-write pattern (`os.Stat` followed by
`io.SafeWriteFile`), not a real lock. Upstream issue:
[ActiveMemory/ctx#93](https://github.com/ActiveMemory/ctx/issues/93).

## Problem

The window between the existence check and the write is
unbounded, and the pattern is racy by construction:

```go
// Acquire lock: fail if another sync is running.
if _, statErr := os.Stat(lockPath); statErr == nil {
    return s, nil, os.ErrExist
}
if writeErr := io.SafeWriteFile(
    lockPath, []byte(cfgHub.LockSentinel), fs.PermFile,
); writeErr != nil {
    return s, nil, writeErr
}
```

Two processes racing past the `os.Stat` (hook-triggered plus
manual invocation; cron plus interactive run; two terminals)
can both decide the lock is free, both write the lock file,
both load the same `LastSequence`, both fetch the same entries
from the hub, and both write duplicate content into
`.context/hub/`. (Issue #93 says `.context/shared/`; that is
the pre-rename path — the renderer joins `cfgHub.DirHub` at
`internal/cli/connection/core/render/render.go:32`.)

The doc comment ("Acquires a lock file to prevent concurrent
access.") overstated what the code did.

The package had no test coverage at all: no test exercised
the contention path, the release path, or the lock location.

## Solution

Replace the stat-then-write with the atomic create-or-fail
primitive that already exists for exactly this purpose:
`io.SafeTryLock` (`O_CREATE|O_EXCL` in a single syscall) and
its counterpart `io.SafeUnlock`. Both landed with the dream
consolidator and have prior art at
`internal/cli/dream/core/pass/pass.go:62-70`.

1. `internal/cli/connection/core/sync/state.go` — swap the
   `os.Stat` / `io.SafeWriteFile` pair for one
   `io.SafeTryLock` call. `acquired == false` maps to the
   existing `os.ErrExist` contract so `Run`'s caller-facing
   behavior is unchanged. The `release` closure delegates to
   `io.SafeUnlock`, keeping the warn-on-failure logging.
2. `internal/config/hub/hub.go` — delete `LockSentinel`. Its
   only consumer was the racy write; `SafeTryLock` creates an
   empty lock file (the lock is the file's existence, not its
   content), so the constant is dead. Leaving it would be a
   broken window.
3. `internal/config/hub/doc.go` — drop `LockSentinel` from
   the package-doc constant inventory.
4. `internal/cli/connection/core/sync/state_test.go` — new;
   regression-pins the contract:
   - `TestLoadState_RejectsConcurrentSyncs`: N goroutines
     race `loadState`; exactly one acquires, the rest observe
     `os.ErrExist`.
   - `TestLoadState_ReleaseRemovesLock`: after `release()`,
     the lock file is gone and a subsequent `loadState`
     succeeds.
   - `TestLoadState_ReleasesLockOnCorruptState`: a corrupt
     `.sync-state.json` makes `loadState` fail *after*
     acquisition; the lock must not leak.
   - `TestLoadState_LockFileLocation`: the lock lives at
     `<ctxDir>/hub/.sync.lock` (pins `DirHub` +
     `FileSyncLock` composition).

## Review Follow-ups (PR #113)

Code review approved the concurrency fix and raised two
hardening points (both non-blocking); addressed here:

1. **Actionable lock-contention error.** Making the lock
   *reliable* raises the stakes of a wedged stale lock: every
   subsequent sync now hard-fails, and `Run` surfaced the raw
   `os.ErrExist` ("file already exists") with no hint. The
   `!acquired` branch now returns `errHub.ConnectSyncLocked(lockPath)`
   (new `err.hub.connect-sync-locked` text key) which names the
   lock path — "remove the stale lock `<ctxDir>/hub/.sync.lock`"
   — and wraps `os.ErrExist` with `%w` so the pre-existing
   `errors.Is(err, os.ErrExist)` contract (and the contention
   test) still holds. This is the cheap half of the deferred
   stale-lock work: a self-documenting wedge ahead of automatic
   PID/age recovery. `TestLoadState_LockedErrorIsActionable`
   pins both the path hint and the `os.ErrExist` wrap.
2. **`SafeTryLock` no longer leaks on close failure.**
   Previously a create-succeeds-but-`Close`-fails returned
   `(true, closeErr)`; callers treat any non-nil error as
   "not acquired" and get no release func, yet the file
   persisted — a leaked lock that would wedge every future
   caller. `SafeTryLock` now removes the freshly created file
   (best-effort) and returns `(false, closeErr)` so the on-disk
   state matches the reported outcome. Both callers (`loadState`
   and `dream/pass`) check the error before `!acquired`, so the
   observable behavior is unchanged except the file no longer
   leaks. The close-failure path is exotic (only ENOSPC-on-close
   / NFS) and not unit-testable without fault injection, so it
   carries no dedicated test.

The remaining review footnotes (NFS `O_EXCL` caveat on ancient
NFSv2; the lock lives in local `.context/hub/`) are unchanged
and remain negligible.

## Out of Scope

- Stale-lock detection (PID-in-lockfile, age-based cleanup).
  A crashed process still leaves a stale lock; the issue
  explicitly defers this to a follow-up. The review-follow-up
  error above makes such a wedge self-documenting, but does not
  auto-recover it.
- `flock`-based locking (issue's Option B). Rejected for now:
  `syscall.Flock` is Unix-only and `SafeTryLock` already
  matches the existing lockfile-as-sentinel model.
- The hubsync hook's silent error suppression
  (ActiveMemory/ctx#100) — adjacent code, separate issue.

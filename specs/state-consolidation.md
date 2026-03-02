# Spec: State Consolidation and Init Guard

## Problem

### Scattered State

ctx writes session-scoped state to `/tmp` (or `$XDG_RUNTIME_DIR`):

- **Agent cooldown tombstones**: `secureTempDir()/ctx-agent-<PID>`
  (prevents repeated `ctx agent` output)
- **Pause markers**: `secureTempDir()/ctx-paused-<sessionID>`
  (suppresses hook output when paused)

This creates three problems:

1. **Split brain**: project-scoped state in `.context/state/`, session
   state in `/tmp` — two directories to reason about
2. **Code duplication**: `secureTempDir()` is copy-pasted identically
   in `agent/cooldown.go` and `system/state.go`
3. **Orphan cleanup**: `cleanup-tmp` SessionEnd hook exists solely to
   sweep stale temp files — infrastructure that shouldn't need to exist

### Missing Init Guard

Not every `ctx` subcommand checks for initialization. Some silently
degrade or produce confusing output when `.context/` is missing or
incomplete. This is analogous to a Kubernetes resource in `Pending`
state — nothing should work until initialization is stable.

## Decision

1. Move all session state from `/tmp` to `.context/state/`
2. Delete `secureTempDir()` from both packages
3. Delete `cleanup-tmp` command and its SessionEnd hook registration
4. Add init guard to all `ctx` subcommands that touch `.context/`

## Changes Required

### Phase 1: Init Guard

Add a shared `requireInit()` check (or Cobra `PersistentPreRunE` on the
root command) that verifies `.context/` exists and contains required files
before any subcommand runs. Commands exempt from the check:

- `ctx init` (creates `.context/`)
- `ctx system bootstrap` (reports context dir, works pre-init)
- `ctx hook` (generates integration configs, works pre-init)
- `ctx version` / `ctx help` (informational)

Error message:

```
ctx: not initialized. Run "ctx init" first.
```

Short, clear, actionable. No partial execution.

### Phase 2: State Relocation

1. **`agent/cooldown.go`**: Replace `secureTempDir()` with
   `filepath.Join(rc.ContextDir(), config.DirState)`. Tombstone files
   stay named `ctx-agent-<PID>`.

2. **`system/state.go`**: Remove `secureTempDir()`. Update
   `pauseMarkerPath()` to use `.context/state/`. Pause marker files
   stay named `ctx-paused-<sessionID>`.

3. **Delete `system/cleanup_tmp.go`**: No longer needed — `.context/state/`
   files are cleaned up by `ctx compact` or manual deletion, and are
   gitignored anyway.

4. **Delete `cleanup-tmp` from hooks.json**: Remove the SessionEnd hook
   registration.

5. **Tests**: Update all tests that reference `secureTempDir()` or
   temp paths.

### Phase 3: Cleanup

- Remove `secureTempDir()` from both packages
- Remove `cleanupTmpCmd()` and its Cobra registration
- Update `ARCHITECTURE.md` hook count (SessionEnd drops from 1 to 0,
  or gets replaced by other SessionEnd hooks)

## Edge Cases

- **Multiple sessions in same project**: Already handled — state files
  are namespaced by session ID or PID
- **Pre-init state writes**: Blocked by init guard — cooldown/pause
  can't fire before `ctx init`
- **Stale state files**: `.context/state/` is gitignored and ephemeral.
  Old tombstones/pause markers are harmless (checked by modtime or
  session validity). Could add cleanup to `ctx compact` if needed.

## Non-Goals

- Global (cross-project) state: if needed in the future, that goes
  under `~/.ctx/` per the dir relocation spec
- Changing state file format — just moving location

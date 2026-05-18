# state.Dir() must not mkdir in uninitialized projects

## Problem

`state.Dir()` calls `SafeMkdirAll` unconditionally on every invocation,
materializing `.context/state/` (mode `0750`) even in projects that
have never been initialized with `ctx init`. This violates the
documented invariant on `state.Initialized()`:

> "Hooks should no-op when this returns false to avoid creating a
> partial state (e.g., logs/) before initialization."

The invariant is unenforceable as long as `Dir()` mkdirs first: any
caller that touches `Dir()` before consulting `Initialized()` leaks.

The leak is reachable in practice today. Cursor's docs
(https://cursor.com/docs/hooks) state that Cursor imports Claude Code
hooks and sets `CLAUDE_PROJECT_DIR` to the workspace root for Claude
compatibility. With the `ctx@activememory-ctx` Claude plugin enabled
globally in `~/.claude/settings.json`, Cursor fires the imported
`UserPromptSubmit` hook chain in every workspace it opens. The chain
includes `ctx system check-reminder`, whose `Run` deliberately calls
`coreCheck.Preamble(stdin)` *before* `state.Initialized()` so that
provenance prints unconditionally:

```
checkreminder.Run
  └─ coreCheck.Preamble                      # before the Initialized gate
       └─ nudge.Paused(sessionID)
            └─ PauseMarkerPath
                 └─ state.Dir()              # ← SafeMkdirAll(.context/state, 0750)
```

Result: opening any non-ctx project in Cursor and submitting a single
prompt deposits a stub `.context/state/` directory into the project
root. Confirmed via leak inspection: mode bits `drwxr-x---` on the
created directory match `fs.PermRestrictedDir = 0750` exactly, which
is uniquely produced by `state.Dir()`.

## Approach

Move the initialization gate inside `state.Dir()` itself, so the
invariant is structural rather than conventional. `Dir()` returns a
new typed error, `errCtx.ErrNotInitialized`, when
`ctxContext.Initialized(ctxDir)` is false. The mkdir runs only on the
initialized path.

This mirrors the existing handling of `errCtx.ErrDirNotDeclared`:
both are "legitimate absence" conditions that callers should treat
as a silent bail rather than a true error. Callers that already have
the `dirErr != nil → return nil` pattern (the dominant one) continue
to work unchanged. Callers that need to distinguish absence from
failure use `errors.Is(err, errCtx.ErrNotInitialized)`.

## Behavior

### Happy Path

1. Hook process starts with `CTX_DIR` set to a properly initialized
   project's `.context/`.
2. Caller invokes `state.Dir()`.
3. `rc.ContextDir()` resolves; `ctxContext.Initialized(ctxDir)`
   returns true.
4. `SafeMkdirAll(ctxDir/state, 0o750)` runs (no-op when present).
5. Returns `(stateDir, nil)`. Existing behavior preserved.

### Uninitialized Path (the new case)

1. Hook process starts with `CTX_DIR` set to a non-ctx project's
   `.context/` (which does not exist on disk).
2. Caller invokes `state.Dir()`.
3. `rc.ContextDir()` resolves; `ctxContext.Initialized(ctxDir)`
   returns false (required files like `AGENT_PLAYBOOK.md` are absent).
4. `Dir()` returns `("", errCtx.ErrNotInitialized)` *without* mkdir.
5. Caller bails silently via its existing `dirErr != nil` branch.
6. Filesystem unchanged — no `.context/` materialized.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `CTX_DIR` unset | Unchanged: `rc.ContextDir()` returns `errCtx.ErrDirNotDeclared`; `Dir()` propagates as before. |
| `CTX_DIR` set, `.context/` does not exist on disk | New: `Initialized()` returns false (required files absent); returns `ErrNotInitialized`. No mkdir of `.context/` or `.context/state/`. |
| `CTX_DIR` set, `.context/` exists but partial (some required files missing) | New: `Initialized()` returns false; returns `ErrNotInitialized`. Critically, this is *also* a leak path today — `state/` would be created inside an otherwise-correct ctx-shaped dir. The fix closes it. |
| `CTX_DIR` set, fully initialized | Unchanged: mkdir runs (no-op when present), returns `(stateDir, nil)`. |
| Concurrent calls during a single hook invocation | Unchanged: `SafeMkdirAll` is idempotent. Initialized check is a stat-only read; no race introduced. |
| `Initialized()` itself errors (resolver failure) | Unchanged: propagate the error so callers' `dirErr != nil` branch fires. Do not silently treat as uninitialized — that would hide real failures. |
| `dirOverride` (test override) is set | Unchanged: bypass the gate entirely, return `(dirOverride, nil)`. Tests that explicitly opt into a state dir continue to work without needing to fake `Initialized()`. |

### Validation Rules

No new input validation. The change is internal to `state.Dir()`.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| `ErrNotInitialized` returned to a hook caller | None — hook bails silently. This is by design: hooks running in non-ctx projects must be invisible. | User runs `ctx init` if they want ctx in this project; otherwise no action needed. |
| `ErrNotInitialized` returned to an interactive command (e.g., `ctx remind list`, `ctx pad show`, `ctx task complete`) that calls `state.Dir()` | Print to stderr: `ctx: this project is not initialized. Run 'ctx init' to set up context here.` Exit non-zero (use the standard CLI error exit code; do not invent a new one). | Run `ctx init` in the project root. |
| `ErrDirNotDeclared` (existing) | Unchanged. | Unchanged. |
| Resolver / mkdir failures | Unchanged: surfaced as today. | Unchanged. |

## Interface

No CLI surface change. This is a library-internal contract change.

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/err/context/errors.go` (or wherever `ErrDirNotDeclared` lives) | Add `ErrNotInitialized` sentinel error. |
| `internal/cli/system/core/state/state.go` | Insert `Initialized()` check between `rc.ContextDir()` and `SafeMkdirAll`. Update package-level docstring on `Dir()` to document the new contract. |
| `internal/cli/system/core/state/state_test.go` (new or existing) | Add unit tests: uninitialized → returns `ErrNotInitialized`, no mkdir occurs; initialized → mkdir runs as before; `dirOverride` bypasses the gate. |
| `internal/cli/system/cmd/checkreminder/run.go` | No code change required — the existing `if dirErr != nil { return nil }` branch in `Preamble`'s call chain absorbs the new error. Add a comment cross-referencing this spec to document why the order is now safe. |
| Other call sites of `state.Dir()` | Two-pass audit during implementation. Pass 1 — hook commands and other non-interactive callers: confirm the existing `dirErr != nil → return nil` branch absorbs `ErrNotInitialized` without warning. Pass 2 — interactive callers (every entry point reachable from a user-typed `ctx ...` subcommand): wrap the call so `errors.Is(err, errCtx.ErrNotInitialized)` produces the stderr message above and a non-zero exit. The full classified list is part of the PR description and must be exhaustive (no "rest as follow-up"). |
| `specs/tests/regression/uninit-no-state-leak/` | Add a test harness that simulates the Cursor scenario end-to-end. See Testing section. |

### Key Functions

```go
// internal/err/context/errors.go
var ErrNotInitialized = errors.New("ctx: project is not initialized")

// internal/cli/system/core/state/state.go
func Dir() (string, error) {
    if dirOverride != "" {
        return dirOverride, nil
    }
    ctxDir, err := rc.ContextDir()
    if err != nil {
        return "", err
    }
    if !ctxContext.Initialized(ctxDir) {
        return "", errCtx.ErrNotInitialized
    }
    d := filepath.Join(ctxDir, dir.State)
    if mkdirErr := ctxIo.SafeMkdirAll(d, fs.PermRestrictedDir); mkdirErr != nil {
        return "", mkdirErr
    }
    return d, nil
}
```

### Helpers to Reuse

- `ctxContext.Initialized(ctxDir)` — already exists at
  `internal/context/validate/validate.go:26`.
- `errCtx.ErrDirNotDeclared` pattern — model `ErrNotInitialized` on it
  for symmetry (same package, same `errors.Is` ergonomics).

## Configuration

None. No new `.ctxrc` keys, env vars, or settings.

## Testing

### Unit

- `state.Dir()` with uninitialized ctxDir → returns `ErrNotInitialized`
  with empty path, no `state/` directory created.
- `state.Dir()` with initialized ctxDir → existing happy-path test
  continues to pass.
- `state.Dir()` with `dirOverride` set → bypasses the gate.
- `errors.Is(err, errCtx.ErrNotInitialized)` works as expected for
  callers that need to distinguish.

### Integration / Regression

The acceptance test for this spec:

```
specs/tests/regression/uninit-no-state-leak/
```

Setup:
1. Create a tempdir; do NOT run `ctx init`.
2. Set `CTX_DIR=<tempdir>/.context`.
3. Invoke `ctx system check-reminder` directly (the leaking entry
   point), feeding it a minimal valid hook JSON envelope on stdin.

Assertions:
- Exit code 0 (hooks must never fail).
- `<tempdir>/.context` does not exist after the call.
- `<tempdir>/.context/state` does not exist after the call.
- stdout/stderr contain no warnings about the missing directory
  (hooks in uninitialized projects are silent by contract).

Repeat the same harness for every UserPromptSubmit hook in
`internal/assets/claude/hooks/hooks.json` to catch any future
regression. Parameterize over the hook list rather than hand-rolling
per hook.

### Edge Case Coverage

- Partial-init case: create `<tempdir>/.context/` with one file but
  not all required files; verify `state/` is still not created.
- `CTX_DIR` unset: verify behavior is identical to before this change.

## Non-Goals

- **Not changing `event.Append`'s mkdir.** That path uses
  `fs.PermExec` (0755) and goes through its own `logFilePath` resolver
  which has its own gate semantics. It is not implicated in this leak
  (mode bits don't match) and is out of scope.
- **Not introducing a top-level "is ctx active in this project?" CLI
  command.** Users diagnose via `ctx system bootstrap` and `ctx
  status`; no new surface needed.
- **Not redesigning the `Preamble` / `FullPreamble` split.** The
  `checkreminder` "provenance-first" ordering is preserved exactly;
  the fix sits below it in the call stack. A future cleanup may
  collapse the two preamble variants, but it is not required to close
  this leak.
- **Not auditing or fixing other places that may pre-create files in
  uninitialized projects.** This spec closes the documented `0750`
  leak. The implementation PR must additionally grep for `SafeMkdirAll`
  / `os.MkdirAll` calls that target paths inside `ctxDir` and verify
  each is reachable only from an initialized-gated path. Any
  additional leak found in that grep is fixed in the same PR
  (per CONSTITUTION's "No Broken Windows" — fix obvious issues when
  encountered). What is excluded from this spec: a wider sweep of
  *non-mkdir* writes (e.g., `notify.Send` webhook payloads), because
  those do not create persistent filesystem state in the project tree.
- **Not changing Cursor's behavior or asking the user to disable the
  Claude-hook import.** The hook firing is intended; only the
  resulting filesystem side-effect is wrong.

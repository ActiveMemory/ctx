---
title: Activate Strict-CWD (No Walk-Up)
status: superseded
date: 2026-05-20
owner: jose
scope: behavioral — `ctx activate` resolver
supersedes-section-of:
  - specs/single-source-context-anchor.md (the `ctx activate` carve-out that
    preserved upward candidate-scanning under "interactive discovery")
superseded-by:
  - specs/cwd-anchored-context.md (deletes `ctx activate` entirely under
    the cwd-anchored model; this spec's resolver tightening is moot once
    the command itself is gone)
related:
  - specs/single-source-context-anchor.md
---

# Spec: Activate Strict-CWD (No Walk-Up)

## Problem

`ctx activate` resolves the target `.context/` by walking upward from
CWD via `rc.ScanCandidates`, picking the innermost candidate and
emitting the rest as "also visible upward" advisories. The walk was
preserved as a deliberate carve-out under
[single-source-context-anchor.md](single-source-context-anchor.md) on
the rationale that activate is "user-invoked discovery, not silent
resolution" and that workspace-shared `.context/` next to per-project
ones is a legitimate layout.

In practice the carve-out leaks the exact failure mode the parent spec
was written to kill:

1. User runs `git init` in a fresh directory under a workspace that
   already has its own `.context/`.
2. `eval $(ctx activate)` walks past the new project's git boundary,
   binds `CTX_DIR` to the workspace `.context/`.
3. `ctx init` then either targets the wrong directory or refuses with
   a confusing "parent context already exists" warning. The user's
   stated intent ("activate ctx in *this* directory") has been
   silently overridden.

Tracked as TASKS.md item line 58.

## Bet

**`ctx activate` is strictly local to `$PWD`. Zero walk, zero
candidate list, zero advisories about parent directories.**

`activate` succeeds if and only if `$PWD/.context/` exists. Otherwise
it bails with a typed error that points at `ctx init`. The
workspace-shared `.context/` use case is preserved by user action
(`cd` to the workspace before activating), not by inferred walk.

### Why the prior justifications no longer hold

- **"Matches `git` / `make` for nested layouts."** `git` walks up for
  *read* commands (`status`, `log`) but refuses to cross repo
  boundaries for *state* commands (`git init` in a subdir of a repo
  warns; `git checkout` cannot reach into a parent repo).
  `ctx activate` is a state command: it exports an env var that
  governs every subsequent ctx invocation in the shell. By git's own
  pattern, state commands are strict.
- **"Workspace-shared `.context/` alongside per-project ones."**
  Preserved without walk-up. To bind the workspace `.context/`, `cd`
  to the workspace dir first. To bind a project's, `cd` to that
  project. The walk only saved one `cd` for a single one-time
  invocation per shell; the surprise tax was recurring.
- **"User-invoked discovery, not silent resolution."** The user did
  invoke the command, but the *resolution* was silent: the user did
  not name the parent path and did not see it before the bind. The
  "also visible upward" advisory on stderr is invisible to the only
  documented invocation form (`eval "$(ctx activate)"`), so the
  multi-candidate path was effectively silent in the use case where
  it mattered most.

### Interaction with TASKS.md item line 63

The pair: `ctx init` should also refuse when `CTX_DIR` is set and
points at a `.context/` whose `realpath` differs from `$PWD/.context`.
Strict-CWD `activate` shrinks the population of users who hit that
bug (the env var no longer gets silently bound to the wrong place),
but does not fix it; the deliberate path ("activated project A, then
cd into project B and ran `ctx init`") still needs the init-side
guard. Out of scope for this spec; tracked separately.

## Implementation

### Resolver

`internal/cli/activate/core/resolve/internal.go::scan()`:

```go
func scan() (string, error) {
    cwd, cwdErr := os.Getwd()
    if cwdErr != nil {
        return "", cwdErr
    }
    candidate := filepath.Join(cwd, dir.Context)
    info, statErr := os.Stat(candidate)
    if statErr != nil {
        if errors.Is(statErr, os.ErrNotExist) {
            return "", errActivate.NoLocalContext(cwd)
        }
        return "", statErr
    }
    if !info.IsDir() {
        return "", errActivate.NoLocalContext(cwd)
    }
    return candidate, nil
}
```

`rc.ScanCandidates` stays in the codebase (still used by
`rc.Require` and `ctx system bootstrap` for error-message context).
The activate path no longer calls it.

### Public surface

- `resolve.Selected()` → `(string, error)`. The `[]string others`
  return value is removed. Sole caller is `cmd/root/run.go`.
- `writeActivate.AlsoVisible` is deleted. No other callers.
- `errActivate.NoCandidates()` → `errActivate.NoLocalContext(cwd)`.
  The new error names `$PWD` and points at `ctx init`, replacing
  the generic "no `.context/` directory found from this location."

### Error message

```
ctx activate: no .context/ at <cwd>
Run `ctx init` here, or `cd` to a directory that has .context/.
See: https://ctx.ist/recipes/activating-context/
```

### Tests

- `TestActivate_NoArgs_NoCandidates` → still passes (renamed
  internally to `TestActivate_NoLocalContext_Bails`).
- `TestActivate_NoArgs_OneCandidate` → still passes (local
  `.context/` → success).
- `TestActivate_NoArgs_ManyCandidates` → **deleted**. Replaced by
  `TestActivate_DeepSubdir_WithParentContext_Bails`: CWD is a
  sub-sub-directory of a tree where a parent has `.context/`;
  activate now bails rather than walking up.
- `TestActivate_RejectsArgs`, `TestActivate_StaleReplacementComment`,
  `TestActivate_NoStaleCommentOnFirstActivate`, `TestActivate_ShellFlag`
  → unchanged.

### Documentation

- `internal/cli/activate/core/resolve/{resolve,internal,doc}.go`
  rewrite the package commentary: drop "walks via
  `rc.ScanCandidates`," drop "innermost-first," drop the
  multi-candidate paragraph. Replace with the strict-CWD contract.
- `internal/cli/activate/cmd/root/{run,doc,cmd}.go`: drop the
  "also visible upward" stderr channel from the documented output
  shape.

## Out of scope

- TASKS.md item line 63 (`ctx init` env-vs-CWD mismatch guard).
- Any change to `rc.Require` or `ctx system bootstrap` resolution.
  These continue to use `CTX_DIR` strictly per the parent spec.
- Any change to the basename guard, hook injection, or
  `check-anchor-drift` machinery from the parent spec.

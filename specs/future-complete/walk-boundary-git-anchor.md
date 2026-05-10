# Walk Boundary: Git-Anchored Context Resolution

## Problem

When a non-ctx-initialized repo lives inside a ctx-initialized parent
workspace, `walkForContextDir` walks upward, finds the parent's
`.context`, and returns it. The boundary check in `PersistentPreRunE`
then rejects it because the resolved path is outside CWD's project
root. The user sees:

```
Error: context directory "/parent/.context" resolves outside project root "/parent/child"
Use --allow-outside-cwd to override this check
```

This affects `ctx init` (can't initialize) and hooks (fire in the
child repo, walk finds parent's `.context`, boundary rejects it).

Root cause: the walk has no concept of project boundaries — it walks
until it finds *any* `.context`, even one belonging to a different
project.

## Approach

Add a git-root anchor to `walkForContextDir`. After the walk finds a
candidate `.context`, check whether it falls within the git root of
the current working directory. If it doesn't, discard the candidate
and use the git root as the anchor for the context directory instead.

Git is a **hint**, not a hard requirement. If no `.git` is found
anywhere, the walk falls through to the existing CWD fallback.

### Resolution chain (updated)

1. CLI `--context-dir` → use it (no walk, existing behavior)
2. `.ctxrc` configured absolute path → use it (no walk, existing behavior)
3. Walk finds `.context` within git root → use it
4. Walk finds `.context` outside git root → discard, use `git-root/<name>`
5. Walk finds `.context`, no git root found → discard, use `cwd/<name>`
   (existing fallback, unchanged)
6. Walk finds nothing, git root found → use `git-root/<name>`
7. Walk finds nothing, no git root → use `cwd/<name>` (existing fallback)

## Behavior

### Happy Path

1. User has `workspace/my-project/` with `.git` and `.context/`
2. User runs `ctx status` from `workspace/my-project/src/pkg/`
3. Walk finds `workspace/my-project/.context`
4. Git walk finds `workspace/my-project/.git` → git root = `my-project/`
5. `.context` is within git root → use it ✓

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Parent workspace has `.context`, child has `.git` but no `.context` | Walk finds parent's `.context`, git root = child → `.context` outside git root → use `child/.context` |
| No `.git` anywhere, parent has `.context` | Walk finds parent's `.context`, no git to confirm → fall through to `cwd/.context` |
| Monorepo: single `.git` at root, `.context` at root, CWD in subpackage | Walk finds root `.context`, git root = repo root → `.context` within git root → use it ✓ |
| Git worktree (`.git` is a file, not a directory) | `os.Stat` matches files too — `.git` file detected, parent = git root ✓ |
| CWD is the git root itself, no `.context` exists yet | Walk finds nothing, git root = CWD → `cwd/.context` (same as before) |
| `.ctxrc` sets `context_dir: /shared/team-context` (absolute) | Absolute path skips walk entirely → existing behavior unchanged |
| Nested git repos (submodule has own `.git`) | Inner `.git` found first → inner git root used → parent's `.context` rejected correctly |

### Validation Rules

No new validation. The boundary check in `PersistentPreRunE` remains
as a safety net. The walk now makes correct decisions before the
boundary check runs.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| `os.Getwd()` fails in git walk | Fall through to existing behavior (CWD fallback) | No user action needed |
| `.git` exists but is unreadable | Treated as "no `.git` found" — fall through | No user action needed |

No new user-facing errors introduced. The fix *removes* an error
path (the boundary violation on init in nested repos).

## Interface

No CLI changes. No new flags. Internal refactor of `walkForContextDir`.

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/rc/walk.go` | Add `findGitRoot` helper; update `walkForContextDir` to validate candidates against git root |
| `internal/rc/walk_test.go` | New file: unit tests for git-anchored walk behavior |
| `internal/rc/rc_test.go` | Update `TestContextDir_UpwardWalkFromSubdir` if needed |

### Key Functions

```go
// findGitRoot walks upward from cwd looking for a .git entry
// (directory or file, to support worktrees). Returns the parent
// directory of the .git entry, or "" if none is found.
func findGitRoot(cwd string) string

// walkForContextDir — updated to call findGitRoot when the walk
// finds a candidate outside CWD, using git root as anchor.
func walkForContextDir(name string) string
```

### Helpers to Reuse

- `filepath.Dir`, `os.Stat` — same pattern already used in walk loop

## Configuration

No new `.ctxrc` keys or environment variables.

## Testing

- Unit: `TestWalkForContextDir_GitAnchor` — parent has `.context`,
  child has `.git`, walk returns `child/.context`
- Unit: `TestWalkForContextDir_NoGit` — no `.git` anywhere, walk
  falls back to `cwd/.context`
- Unit: `TestWalkForContextDir_SameGitRoot` — `.context` and CWD
  share same git root, walk returns the found `.context`
- Unit: `TestWalkForContextDir_GitWorktreeFile` — `.git` is a file
  (worktree), still detected as git root
- Integration: existing `TestContextDir_UpwardWalkFromSubdir`
  continues to pass (same git root scenario)

## Non-Goals

- Making git a hard dependency for ctx
- Detecting project roots via `go.mod`, `package.json`, etc.
- Changing the boundary check in `PersistentPreRunE`
- Changing `rc.ContextDir()` resolution for explicit overrides
  (CLI flag, `.ctxrc`, env var)

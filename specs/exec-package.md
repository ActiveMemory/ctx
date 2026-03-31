# Spec: Centralize exec.Command under internal/exec

## Problem

`exec.Command` is called from 15 files across 10 packages. Each
call site independently handles:

- LookPath checks (some do, some don't)
- gosec G204 suppression (`//nolint:gosec`)
- Input validation (some validate, some trust callers)
- Error wrapping (inconsistent across call sites)

This creates duplicated patterns and scattered `//nolint:gosec`
markers. Centralizing under `internal/exec` gives one place to
audit command execution, validate inputs, and satisfy gosec.

## Non-Goals

- Mocking exec for unit tests. Callers that need testability
  should use interfaces at their own level.
- Wrapping test exec.Command calls (`cli_test.go`,
  `compliance_test.go`). Test binaries are inherently dynamic.

## Inventory

### Git operations (6 call sites → `internal/exec/git`)

| Call site | Operation | Args |
|-----------|-----------|------|
| `config/core/core.go:105` | rev-parse --show-toplevel | literal |
| `change/core/scan/scan.go:133` | log --since (time) | validated |
| `system/cmd/post_commit/score.go:33` | log -1 --format=%B | literal |
| `system/cmd/post_commit/score.go:66` | diff-tree HEAD | literal |
| `system/core/health/map_staleness.go:66` | log --oneline --since | validated |
| `journal/parser/git.go:42` | remote get-url origin | literal |

All use `cfgGit.Binary` constant. All return `[]byte` via
`.Output()`. Several duplicate the LookPath check.

**Proposed API:**

```go
package git

// Run executes a git command with literal arguments and
// returns its output.
func Run(args ...string) ([]byte, error)

// Root returns the repository root (rev-parse --show-toplevel).
func Root() (string, error)

// RemoteURL returns the origin remote URL for a directory.
func RemoteURL(dir string) (string, error)

// LogSince runs git log with a --since filter.
func LogSince(since time.Time, extraArgs ...string) ([]byte, error)

// LastCommitMessage returns the most recent commit message.
func LastCommitMessage() ([]byte, error)

// DiffTreeHead returns the list of changed files in HEAD.
func DiffTreeHead() ([]byte, error)
```

`Run` handles LookPath once. Specific functions compose `Run`
with validated arguments. gosec nolint lives in one place.

### Platform commands (3 call sites → keep in sysinfo)

| Call site | Binary | Platform |
|-----------|--------|----------|
| `sysinfo/memory_darwin.go:29` | sysctl | darwin |
| `sysinfo/memory_darwin.go:42` | vm_stat | darwin |
| `sysinfo/load_darwin.go:27` | sysctl | darwin |

These are platform-specific with build tags. They use literal
constant args and don't benefit from centralization. **Leave
in sysinfo** — the build tag isolation is more important than
the exec centralization.

### Zensical (1 call site → already in `exec/zensical`)

Already lives in `internal/exec/zensical`. No changes needed.

### Dependency listing (3 call sites → `internal/exec/dep`)

| Call site | Binary | Operation |
|-----------|--------|-----------|
| `dep/core/go.go:84` | go | list -m all |
| `dep/core/rust.go:122` | cargo | metadata |
| `dep/core/rust.go:147` | cargo | metadata (alt) |

These call external toolchain binaries. Move exec calls to
`internal/exec/dep` with `GoListModules()` and
`CargoMetadata()`.

### Archive (1 call site → `internal/exec/gio`)

| Call site | Binary | Operation |
|-----------|--------|-----------|
| `system/core/archive/smb.go:79` | gio | mount |

Single call site. Move to `internal/exec/gio` with `Mount()`.

### Validate (1 call site → uses LookPath only)

| Call site | Binary | Operation |
|-----------|--------|-----------|
| `initialize/core/validate/validate.go:32` | ctx | LookPath |

Only uses `exec.LookPath`, not `exec.Command`. No change.

### Test files (7 call sites → leave as-is)

| Call site | Binary |
|-----------|--------|
| `cli/cli_test.go` (6 calls) | ctx binary |
| `compliance/compliance_test.go` (3 calls) | go, golangci-lint |

Test exec calls are inherently dynamic (build paths, test
binaries). Leave as-is.

## Package structure

```
internal/exec/
├── git/          # Git operations (6 callers → 1)
│   ├── doc.go
│   └── git.go
├── dep/          # Dependency toolchains (3 callers → 1)
│   ├── doc.go
│   ├── go.go
│   └── rust.go
├── gio/          # GIO/mount operations (1 caller)
│   ├── doc.go
│   └── mount.go
└── zensical/     # Already exists
    └── zensical.go
```

## Migration strategy

1. Build `exec/git` with all 6 git operations
2. Migrate callers one package at a time, verify tests
3. Build `exec/dep` with go/cargo operations
4. Build `exec/gio` with mount operation
5. Delete orphaned `//nolint:gosec` markers
6. Verify: `grep -rn 'os/exec' internal/` returns only
   `exec/`, `sysinfo/`, test files, and `validate/`

## Outcome

- `os/exec` imports drop from 15 files to ~8 (exec/, sysinfo/,
  tests, validate)
- gosec G204 nolint markers drop from 6 to 0 in non-exec code
- LookPath checks consolidated (no more "sometimes checked")
- One audit surface for command execution

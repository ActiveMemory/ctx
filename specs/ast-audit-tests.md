---
title: AST-based audit tests to replace lint-drift.sh
date: 2026-03-23
status: ready
---

# AST-Based Audit Tests

## Problem

`hack/lint-drift.sh` uses regex/grep to detect code-level drift. This
approach has three structural weaknesses:

1. **No type awareness** — it cannot distinguish `Use*` constants from
   `DescKey*` constants, causing false positives (71 in the cmd↔YAML
   check before the fix on 2026-03-23).
2. **Fragile exclusions** — when a constant moves from `token.go` to
   `whitespace.go`, the exclusion glob breaks silently.
3. **Ceiling on detection** — checks that require understanding call
   sites, import graphs, or type relationships are impossible in shell.

## Solution

Replace the five shell checks with Go tests using `go/ast` and
`go/packages`. The tests run as part of `go test ./...` — no separate
script, no separate CI step.

## Package Location

`internal/audit/` — a dedicated package for codebase-level invariant
tests. All files are `_test.go` so the package produces no binary
output and is not importable.

## Check Migration

### From lint-drift.sh (replace)

| # | Shell check | Go test | Improvement |
|---|-------------|---------|-------------|
| 1 | Literal `"\n"` grep | `TestNoLiteralNewline` — walk AST, find string literals `== "\n"`, skip constant definition sites | No exclusion filenames to maintain |
| 2 | `cmd.Printf`/`cmd.PrintErrf` grep | `TestNoPrintfCalls` — find call expressions matching `cmd.Printf` / `cmd.PrintErrf` | Distinguishes calls from comments/strings |
| 3 | Magic dir strings in `filepath.Join` | `TestNoMagicDirStrings` — find `filepath.Join` calls with string literal args matching known `Dir*` constant values | Also catches `path.Join`, ignores test code properly |
| 4 | Literal `".md"` grep | `TestNoLiteralExtMarkdown` — same pattern as check 1 | Same benefit |
| 5 | DescKey ↔ YAML linkage | `TestDescKeyYAMLLinkage` — load `DescKey*` constants from each embed sub-package, load corresponding YAML, diff key sets | Type-aware, no regex on Go source |

### New checks (not possible in shell)

| Test | What it catches |
|------|----------------|
| `TestUseConstantsOnlyInCobraUse` | `Use*` constant appears only in cobra `Use:` struct field assignments |
| `TestDescKeyOnlyInLookupCalls` | `DescKey*` passed only to `assets.CommandDesc()` / `assets.FlagDesc()` / `desc.Text()` |
| `TestNoDeadConstants` | Constants defined in `embed/cmd`, `embed/flag`, `embed/text` but never referenced outside their definition |
| `TestNoWrongNamespaceLookup` | `TextDescKey` not passed to `CommandDesc()`, `FlagDescKey` not passed to `Text()`, etc. |
| `TestNoStringConcatPaths` | `filepath.Join` arguments do not contain `+` string concatenation |

## Implementation Approach

### AST scanning pattern

Each test follows the same structure:

```go
func TestNoLiteralNewline(t *testing.T) {
    pkgs := loadPackages(t, "github.com/ActiveMemory/ctx/internal/...")
    violations := []string{}
    for _, pkg := range pkgs {
        for _, file := range pkg.Syntax {
            ast.Inspect(file, func(n ast.Node) bool {
                // check node, append to violations
                return true
            })
        }
    }
    for _, v := range violations {
        t.Error(v)
    }
}
```

### Shared helpers

- `loadPackages(t, pattern)` — wraps `go/packages.Load` with test
  caching. Load once per test run via `sync.Once`.
- `isTestFile(filename)` — skip `_test.go` files.
- `isConstantDef(node)` — detect `const ( X = "..." )` so literal
  checks can skip definition sites.
- `posString(fset, pos)` — format `file:line` for error messages.

### Performance

`go/packages.Load` with `NeedSyntax` for `internal/...` takes ~2-3s.
Cache the loaded packages across tests in the same package via
package-level `sync.Once`. Total overhead is comparable to the shell
script (~1s for grep vs ~3s for AST load + walk).

## Migration Plan

1. Create `internal/audit/` with `doc.go` and shared helpers
2. Port checks 1-4 (literal detection) as AST tests
3. Port check 5 (DescKey ↔ YAML linkage) as a type-aware test
4. Add new checks (dead constants, wrong namespace, Use-site validation)
5. Remove `hack/lint-drift.sh` and its call from `Makefile`
6. Update `Makefile` audit target to rely solely on `go test ./...`

Steps 1-3 can ship together. Steps 4-5 can follow incrementally.

## Non-Goals

- Replacing `golangci-lint` — it handles standard Go linting well
- Replacing `go vet` — same
- Runtime checking — these are compile-time/test-time invariants only
- Checking non-Go files (YAML structure, Markdown formatting) — those
  stay in shell or dedicated tools

## Open Questions

- Should the AST tests live in `internal/audit/` (dedicated) or be
  colocated in `internal/config/embed/` (closer to the constants)?
  Recommendation: `internal/audit/` because the checks span multiple
  packages — they are codebase-level, not package-level.
- Should `hack/lint-drift.sh` be removed immediately or kept as a
  fallback during migration? Recommendation: remove once all 5 checks
  are ported — no parallel maintenance.

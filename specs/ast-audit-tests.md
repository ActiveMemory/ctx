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
| 6 | Inline `fmt.Errorf`/`errors.New` outside `internal/err/` | `TestNoNakedErrors` — find `fmt.Errorf` and `errors.New` call expressions outside `internal/err/**` packages | Distinguishes calls from comments/strings; covers all error construction patterns |
| 7 | `err.go` files outside `internal/err/` | `TestNoStrayErrFiles` — filesystem check for `err.go` outside `internal/err/` | Same as shell but unified in Go test runner |
| 8 | `strings.Join` with inline separator | `TestNoInlineSeparators` — find `strings.Join` calls with string literal second argument, outside `internal/config/token/` | Distinguishes calls from comments/strings |

### New checks (not possible in shell)

| Test | What it catches |
|------|----------------|
| `TestUseConstantsOnlyInCobraUse` | `Use*` constant appears only in cobra `Use:` struct field assignments |
| `TestDescKeyOnlyInLookupCalls` | `DescKey*` passed only to `assets.CommandDesc()` / `assets.FlagDesc()` / `desc.Text()` |
| `TestNoDeadConstants` | Constants defined in `embed/cmd`, `embed/flag`, `embed/text` but never referenced outside their definition |
| `TestNoWrongNamespaceLookup` | `TextDescKey` not passed to `CommandDesc()`, `FlagDescKey` not passed to `Text()`, etc. |
| `TestNoStringConcatPaths` | `filepath.Join` arguments do not contain `+` string concatenation |
| `TestNoStutteryFunctions` | Function names that redundantly include their package name as a PascalCase word boundary (e.g., `write.WriteJournal` → `write.Journal`, `write.journalWriteSilent` → `write.journalSilent`). Covers both exported and unexported functions. Identity functions like `write.Write` / `write.write` are excluded. |
| `TestNoRawLogging` | Direct logging to stderr or log files must only appear in `internal/log/**`. All other packages must use `log/warn.Warn` (stderr warnings) or `log/event.Append` (structured event log). Flags `fmt.Fprintf(os.Stderr, ...)`, `fmt.Fprintln(os.Stderr, ...)`, `os.Stderr.Write*`, and stdlib `log.Print*`/`log.Fatal*`/`log.Panic*` calls outside `internal/log/`. Test files exempt. |
| `TestNoRawFileIO` | Direct `os.ReadFile`, `os.WriteFile`, `os.Open`, `os.OpenFile`, `os.Create`, `os.MkdirAll` calls must only appear in `internal/io/`. All other packages must use the `Safe*` wrappers (`SafeReadFile`, `SafeWriteFile`, `SafeOpenUserFile`, etc.) which centralize path validation, sanitization, and `nolint:gosec` suppression. Test files exempt. |
| `TestNoFlagBindOutsideFlagbind` | Direct cobra flag registration (`.Flags().StringVar`, `.Flags().BoolVarP`, etc.) must only appear in `internal/flagbind/`. No exceptions — missing helpers must be added to `flagbind` before this check can pass. Prerequisite: extend `flagbind` with `IntFlag`, `DurationFlag`, `DurationFlagP`, `StringP`, `BoolP` and migrate all ~50 call sites. Test files exempt. |
| `TestNoExecOutsideExecPkg` | `exec.Command` and `exec.CommandContext` calls must only appear in `internal/exec/**` packages, which centralize `nolint:gosec` and per-command sanitization. Test files exempt. Detects `*ast.CallExpr` with `exec.Command*` selectors, flags any outside `internal/exec/`. |
| `TestNoCmdPrintOutsideWrite` | `cmd.Println`, `cmd.Print`, `cmd.Printf`, `cmd.PrintErr*` calls must only appear in `internal/write/**` packages. All other packages must delegate output through the corresponding `write/` subpackage. Test files exempt. Detects `*ast.CallExpr` with `cmd.Print*` selectors, flags any outside `internal/write/`. |
| `TestNoErrorsAs` | Flags calls to `errors.As()` which should use the generic `errors.AsType()` (available since Go 1.23). Detects `*ast.CallExpr` with selector `errors.As`. |
| `TestNoMagicValues` | Flags magic string and numeric literals in non-test Go files under `internal/`. Walks `ast.BasicLit` nodes and checks parent context. **String exceptions**: empty string `""`, single space `" "`, `const`/`var` definition sites, struct tags, import paths. **Numeric exceptions**: `0`, `1`, `-1`, strconv radix/bitsize arguments (`10`, `32`, `64` when parent is a `strconv.Parse*`/`Format*` call), `const` definition sites. User-facing text must go through `internal/assets`, configuration values through `internal/config`. Does not apply to `editors/` (TypeScript). |
| `TestPackageDocQuality` | Every package under `internal/` must have a `doc.go` — including packages that exist only as parents for subpackages. The package doc comment (`ast.File.Doc`) must have at least 8 lines of meaningful text (excluding blank comment lines and the `// Package X` opener). Flags lazy one-liners and file-list patterns (`// - foo.go`, `// Source files:`) which are maintenance-fragile. Existence via `os.ReadDir`, quality via AST. |
| `TestDocComments` | All functions (exported and unexported), structs, and package-level variables must have a doc comment (`ast.FuncDecl.Doc`, `ast.TypeSpec.Doc`, `ast.ValueSpec.Doc` non-nil). Test files (`_test.go`) are exempt. Aligns with the godoc format convention in CONVENTIONS.md. |
| `TestNoRawPermissions` | Octal file permission literals (`0644`, `0755`, `0600`, `0700`, `0750`) must not appear outside `internal/config/fs/` and `internal/io/`. All other packages must use `config/fs` constants (`fs.PermFile`, `fs.PermExec`, `fs.PermSecret`, etc.). Detects octal `ast.BasicLit` nodes in `os.WriteFile`, `os.Mkdir*`, `os.OpenFile`, `os.Chmod` call arguments. Test files exempt. |
| `TestNoInlineRegexpCompile` | `regexp.MustCompile` and `regexp.Compile` calls must appear at package level (`var` declarations), never inside function bodies. The project centralizes compiled patterns in `internal/config/regex/` as package-level vars. Detects `*ast.CallExpr` for `regexp.*Compile` where the parent is not a `*ast.ValueSpec` at file scope. Prevents per-call recompilation and enforces the single-definition pattern. |
| `TestNoVariableShadowing` | Detects variable shadowing: (a) **error variables** — multiple `:=` assignments to bare `err` in the same function scope; the convention requires descriptive names (`readErr`, `writeErr`, `parseErr`); (b) **general shadowing** — inner-scope `:=` declarations that shadow an outer-scope variable of the same name. Walk `*ast.AssignStmt` (`:=`) nodes and track identifier scopes per function body. Test files exempt. |
| `TestNoRawTimeFormats` | Raw time layout strings in `time.Parse`, `time.Format`, and `time.AppendFormat` calls must use `config/time` constants (`time.DateFormat`, `time.DateTimeFmt`, `time.CompactTimestamp`, etc.) instead of inline format strings. Detects string literal arguments to these calls outside `internal/config/time/`. Test files exempt. |
| `TestCLICmdStructure` | Enforces `internal/cli/` cmd directory conventions: (a) `cmd/$sub/` dirs contain only `cmd.go`, `run.go`, `doc.go`, and test files — no stray `.go` files; (b) `cmd.go` declares only a `Cmd()` function — no structs, global vars, or extra functions; (c) `run.go` declares only `Run()` or `Run*()` functions — no structs, global vars, or extra functions; (d) `cmd/` dirs themselves contain only subdirectories, `doc.go`, and test files. Filesystem checks via `os.ReadDir`, shape checks via AST `File.Decls`. |

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

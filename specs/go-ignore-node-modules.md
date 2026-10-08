# Keep node_modules Out of the Go Module

The repo holds four npm packages (`editors/vscode`, `ctx-desktop`,
`tools/typecheck/opencode`, `tools/typecheck/pi`) inside the root
Go module. npm dependencies may ship Go sources, and nothing kept
Go tooling from treating them as part of ctx.

## Problem

`flatted` (pulled in by eslint through `flat-cache`) ships
`flatted/golang/pkg/flatted/flatted.go`. After `npm install` in
`editors/vscode`:

- `go list ./...` reports
  `github.com/ActiveMemory/ctx/editors/vscode/node_modules/flatted/golang/pkg/flatted`
  as a package of this module, so `go build`, `go vet`, and
  `go test ./...` compile it.
- `make lint` fails with two `govet` findings inside that file.
- `make test` fails in `TestGolangciLint`, which shells out to
  golangci-lint.

CI never saw it, because the Go jobs don't run `npm install`. Any
contributor who works on the extension and then runs `make lint`
does.

`./...` skips only directories starting with `.` or `_`,
`testdata`, and nested modules; `node_modules` is none of those.
golangci-lint's `exclusions.paths` only hides reported issues: the
package is still loaded and type-checked.

## Solution

1. **`go.mod`: `ignore node_modules`** (Go 1.25+). The bare form,
   without a `./` prefix, matches a `node_modules` directory at
   any depth (`./node_modules` would match only the root, which
   was verified to leave the package in `go list ./...`). This
   removes the directory from `./...` for every Go command,
   golangci-lint included, since it loads packages through the
   go command.
2. **`.golangci.yml`: `node_modules` in `exclusions.paths`**,
   beside `vendor`, `dist`, and `site`. With (1) in place this is
   belt-and-braces; it states in the linter's own config that
   node code is not in scope, and holds if golangci-lint is ever
   pointed at an explicit path.
3. **Regression guard:**
   `internal/compliance/go_mod_ignore_test.go` reads go.mod via
   `go mod edit -json` and fails if the directive is missing.

The compliance tests' own file walkers (`allGoFiles`,
`allSourceFiles`) already skip `node_modules`; they need no
change.

## Verification

With `editors/vscode/node_modules` installed:

- `go list ./... | grep node_modules` prints nothing.
- `make lint`: 0 issues. `make test`: passes.
- `TestGoModIgnoresNodeModules` passes, and fails after
  `go mod edit -dropignore=node_modules`.

## Non-Goals

- `tools/ctxctl` is a separate module with no npm package inside
  it; its go.mod is unchanged.

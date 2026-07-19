# go-idioms-and-structure

## [2026-06-07-170003] Error handling: sentinels, unwrapping, and silent discards (consolidated)

**Consolidated from**: 6 entries (2026-03-06 to 2026-06-02)

- os.IsNotExist does NOT unwrap — it is false on any fmt.Errorf("…%w…")
  error; prefer errors.Is(err, os.ErrNotExist). But errors.Is only holds if the
  wrap carries %w at runtime, and a wrap whose format string comes from the
  text/i18n registry only carries %w when that registry is initialized (so it
  behaves differently in prod vs a bare test binary; go vet can't see it). To
  detect file absence reliably, stat directly: os.Stat returns an unwrapped
  *fs.PathError so errors.Is(statErr, os.ErrNotExist) is dependable everywhere.
- An error-discard catalogue (grep + name/regex classification) is an inventory
  of candidates, not findings. Name-inference produces false positives (a
  discarded bool mistaken for an error; a value type that can't nil-deref; an
  already-failed cleanup-close path). Read the callee signature and enclosing
  control flow before assigning return-error vs logWarn vs annotate.
- Canonical sentinel shape: a typed zero-data struct (or fielded struct for
  parameterised errors) whose Error() resolves text via
  desc.Text(text.DescKey…) lazily at call time — never `var ErrX =
  errors.New("english")` and never an ErrMsg* string-const layer. Empty-struct
  values are comparable and errors.Is finds them through %w wraps. Reference:
  internal/err/context/context.go.
- fmt.Fprintf to strings.Builder silently discards errors (Write never fails) so
  errcheck allows it, but project convention forbids any silent discard —
  TestNoUncheckedFmtWrite enforces `if _, err := fmt.Fprintf(...)`.
- A path-returning (string, error) function must never return ('', nil):
  filepath.Join('', rel) yields rel as a CWD-relative path, causing orphan
  writes at project root. Sentinel errors force callers to gate. Audit any
  path-returner with a historic ('', nil) shortcut (fixed: state.Dir,
  rc.ContextDir).
- Package-local err.go files in CLI packages invite agents to duplicate error
  constructors (errFileWrite, errMkdir repeated). Centralize in internal/err; no
  err.go files in CLI packages.

---

## [2026-06-07-170009] Constant placement & helper smells (consolidated)

**Consolidated from**: 6 entries (2026-03-07 to 2026-03-23)

- A constant used by only one domain (agent scoring, budget %, cooldowns)
  belongs in that domain's config package, not a god-object file.go. Check
  callers before placing.
- Before adding any constant to internal/config, grep by VALUE (".jsonl") not
  just name — camelCase vs ALLCAPS variants hide duplicates (ExtJsonl vs
  existing ExtJSONL).
- Project-root files created by `ctx init` (Makefile) are scaffolding
  (config/file), NOT context files loaded via ReadOrder (config/ctx). Check
  ReadOrder membership before moving a file constant.
- SafeReadFile / validation.SafeReadFile take (baseDir, filename) separately —
  split full paths with filepath.Dir + filepath.Base when adapting os.ReadFile
  calls.
- One-liner method wrappers that just forward a struct field to a stdlib/pkg
  function (checkBoundary → validation.ValidateBoundary with h.ContextDir)
  obscure the real dependency — inline them.
- A param-struct field that is a function pointer where all callers pass thin
  wrappers varying only by a text key (MergeParams.UpdateFn) is "data in
  disguise" — replace the callback with the key and let the consumer dispatch.

---

## [2026-06-07-170011] Go toolchain, gofmt & build-tag pitfalls (consolidated)

**Consolidated from**: 5 entries (2026-03-16 to 2026-05-10)

- gofmt strips bare `//` padding lines as unnecessary whitespace, so
  programmatic Go generation must produce substantive content lines; always run
  gofmt after any scripted Go-file generation.
- Agents reliably introduce gofmt issues during bulk renames (75+ files, 12
  broken); run `gofmt -l` (then `-w`) as a standard step after any agent-driven
  bulk edit before trusting the build.
- The "compile version X does not match go tool version Y" error comes from the
  CACHED toolchain (~/go/pkg/mod/golang.org/toolchain@…), not the system Go
  — reinstalling Go does nothing. Diagnose via `go env GOROOT`; fix by
  deleting the cached dir, bumping go.mod, or GOTOOLCHAIN=go<system>. `go clean
  -cache` and GOTOOLCHAIN=local don't help.
- `make test` exit code is unreliable: the -cover flag can fail with "no such
  tool covdata" even when every package passes. Fall back to `go test ./...` (no
  -cover) and tally ^ok/^FAIL.
- AST checks via go/packages only see files matching the current GOOS —
  darwin-only (_darwin.go) violations are invisible on Linux. Fix violations
  regardless; note coverage is platform-dependent (need multi-GOOS CI or a
  go/parser fallback).

---

## [2026-06-07-170018] Go test isolation & patterns (consolidated)

**Consolidated from**: 4 entries (2026-01-19 to 2026-04-25)

- Any code using os.UserHomeDir() / user-level paths (~/.ctx/, ~/.config/) needs
  t.Setenv("HOME", tmpDir) in tests — especially shared setup helpers. Under
  parallel `make test`, fourteen test files invoking initialize.Cmd().Execute()
  raced on read-modify-write of ~/.claude/settings.json, surfacing as flaky
  "FAIL coverage: [no statements]"; testctx.Declare now sets HOME alongside
  CTX_DIR (centralized fix).
- Go testing patterns: `go build ./...` misses test-file callsite breaks —
  always `go test ./...` after signature changes. Consume all runCmd() returns
  (`_, _ = runCmd(...)`) for errcheck. Disable ANSI via color.NoColor=true in
  package init for string assertions. Recall tests isolate via t.Setenv("HOME",
  tmpDir) with .claude/projects/. formatDuration takes an interface with
  Minutes() (use a stubDuration). CI needs CTX_SKIP_PATH_CHECK=1 (init checks
  PATH). CGO_ENABLED=0 for ARM64 Linux.
- Converting PersistentPreRun → PersistentPreRunE changes exit behavior:
  errors propagate through Cobra Execute() return with no os.Exit.
  Subprocess-based tests expecting exit codes must convert to direct error
  assertions.

---
## [2026-05-28-201400] A non-root Go module nested under the main module's path CAN import its internal/ packages

**Context**: While designing the ctxctl module split, the initial spec (and a
lot of online consensus) claimed a separate `go.mod` cannot import the parent
module's `internal/` packages, which would have forced relocating or duplicating
~25 foundation packages (`rc`, `desc`, `nudge`, `config/*`, …). The "obvious"
reading made same-module the only viable option.

**Lesson**: Go's internal-import rule is **lexical on import paths, not
module-scoped**. A separate module whose path is
`github.com/<owner>/<main>/tools/<x>` CAN import
`github.com/<owner>/<main>/internal/...` — verified by an empirical build
experiment this session. An outsider path (`example.com/...`) is rejected with
`use of internal package … not allowed`. The rule fires on the import-path
prefix relative to the `internal/` directory's parent, not on module boundaries.

**Application**: For monorepo splits (maintainer-only tooling, isolated
experiments, ancillary CLIs), choose a module path nested under the main module
so the new module reuses the parent's foundations via the lexical-internal
allowance. Full self-containment of a maintainer module would be a DRY
catastrophe; the lexical allowance is the correct shape. Prove it with a
throwaway `go build` against a representative `internal/` import before
designing around the *wrong* constraint.

---

## [2026-05-17-061500] `_helpers.go` / `_utils.go` filenames are project anti-pattern; use domain nouns

**Context**: During Phase KB / Phase RG audit cleanup, the first file split I
attempted to satisfy the mixed-visibility audit named the new file
`read_helpers.go`. The user vetoed on sight: "utils; helpers, etcs are ALL lazy
naming; I will veto them the moment I see them; find proper domain objects."

**Lesson**: ctx's per-package file layout follows domain nouns, not
visibility-suffixed catch-alls. The canonical reference shape is
`internal/journal/parser/` which splits 18 files by domain (envelope, markdown,
parse, validate, claude, copilot, ...). The mixed-visibility audit demands a
split, but the split target must be a real noun: `frontmatter.go` (YAML
parsing/validation), `markdown.go` (rendering), `filename.go` (filename
derivation), `provenance.go` (sha/branch resolution), `parse.go` (one-shot
parser), `cursor.go` (latest-pointer logic). Never `_helpers.go`.

**Application**: When splitting a file to satisfy `mixed_visibility_test`, name
the new file for what the helpers ARE about, not for what visibility they have.
If you can't name it cleanly, the split itself may be wrong and the funcs may
belong in a different package entirely.

---

## [2026-04-03-180000] Import cycles and package splits (consolidated)

**Consolidated from**: 5 entries (2026-03-06 to 2026-03-22)

- Types in god-object files (e.g. hook/types.go with 15+ types from 8 domains)
  create circular dependencies — move types to their owning domain package
- Tests in parent package X cannot import X/sub packages that import X back —
  move tests to the sub-package they exercise
- Variable shadowing causes cascading failures after splits: `dir`, `file`,
  `entry` are common Go variable names that collide with new sub-package names
  — run `go test ./...` before committing splits
- When moving constants between packages, change imports and all references in a
  single atomic write so the linter never sees an inconsistent state
- Import cycle rule: the package providing implementation logic must own the
  shared types; the facade package aliases them (e.g. `entry.Params` aliases
  `add/core.EntryParams`)

---

## [2026-03-18-133457] Lazy sync.Once per-accessor is a code smell for static embedded data

**Context**: assets package had 4 sync.Once guards, 4 exported maps, 4 Load*()
functions, and a wrapper desc package — all to lazily load YAML from embed.FS
that never mutates. Every accessor call went through sync.Once + global map +
wrapper indirection.

**Lesson**: When data is static and loaded from embedded bytes, scatter-loading
with per-accessor sync.Once is over-engineering. A single Init() called eagerly
at startup is simpler, and one sync.Once on Init() itself provides the test
safety net. Exported maps that exist only for wrapper packages to reach are a
sign the abstraction boundary is wrong.

**Application**: Prefer eager Init() in main.go for static embedded data. Keep
maps unexported. Accessors do plain map lookups. If a wrapper package exists
solely to break a cycle caused by exported state, delete the wrapper and
unexport the state.

---


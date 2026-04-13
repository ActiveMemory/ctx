---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Code Conventions
icon: lucide/scroll-text
---
![ctx](../images/ctx-banner.png)


# Code Conventions: Common Patterns and Fixes

This guide documents the code conventions enforced by `internal/audit/`
AST tests. Each section shows the violation pattern, the fix, and the
rationale. When a test fails, find the matching section below.

All tests skip `_test.go` files. The patterns apply only to production
code under `internal/`.

---

## Variable Shadowing (bare `err :=` reuse)

**Test:** `TestNoVariableShadowing`

When a function has multiple `:=` assignments to `err`, each shadows
the previous one. This makes it impossible to tell which error a later
`if err != nil` is checking.

**Before:**

```go
func Run(cmd *cobra.Command) error {
    data, err := os.ReadFile(path) 
    if err != nil {
        return err
    }

    result, err := json.Unmarshal(data)  // shadows first err
    if err != nil {
        return err
    }

    err = validate(result)  // shadows again
    return err
}
```

**After:**

```go
func Run(cmd *cobra.Command) error {
    data, readErr := os.ReadFile(path)
    if readErr != nil {
        return readErr
    }

    result, parseErr := json.Unmarshal(data)
    if parseErr != nil {
        return parseErr
    }

    validateErr := validate(result)
    return validateErr
}
```

**Rule:** Use descriptive error names (`readErr`, `writeErr`, `parseErr`,
`walkErr`, `absErr`, `relErr`) so each error site is independently
identifiable.

---

## Import Name Shadowing

**Test:** `TestNoImportNameShadowing`

When a local variable has the same name as an imported package, the
import becomes inaccessible in that scope.

**Before:**

```go
import "github.com/ActiveMemory/ctx/internal/session"

func process(session *entity.Session) {  // param shadows import
    // session package is now unreachable here
}
```

**After:**

```go
import "github.com/ActiveMemory/ctx/internal/session"

func process(sess *entity.Session) {
    // session package still accessible
}
```

**Rule:** Parameters, variables, and return values must not reuse
imported package names. Common renames: `session` -> `sess`,
`token` -> `tok`, `config` -> `cfg`, `entry` -> `ent`.

---

## Magic Strings

**Test:** `TestNoMagicStrings`

String literals in function bodies are invisible to refactoring tools
and cause silent breakage when the value changes in one place but not
another.

**Before (string literals):**

```go
func loadContext() {
    data := filepath.Join(dir, "TASKS.md")
    if strings.HasSuffix(name, ".yaml") {
        // ...
    }
}
```

**After:**

```go
func loadContext() {
    data := filepath.Join(dir, config.FilenameTask)
    if strings.HasSuffix(name, config.ExtYAML) {
        // ...
    }
}
```

**Before (format verbs — also caught):**

```go
func EntryHash(text string) string {
    h := sha256.Sum256([]byte(text))
    return fmt.Sprintf("%x", h[:8])
}
```

**After:**

```go
func EntryHash(text string) string {
    h := sha256.Sum256([]byte(text))
    return hex.EncodeToString(h[:cfgFmt.HashPrefixLen])
}
```

**Before (URL schemes — also caught):**

```go
if strings.HasPrefix(target, "https://") ||
    strings.HasPrefix(target, "http://") {
    return target
}
```

**After:**

```go
if strings.HasPrefix(target, cfgHTTP.PrefixHTTPS) ||
    strings.HasPrefix(target, cfgHTTP.PrefixHTTP) {
    return target
}
```

**Exempt from this check:**

- Empty string `""`, single space `" "`, indentation strings
- Regex capture references (`$1`, `${name}`)
- `const` and `var` definition sites (that's where constants live)
- Struct tags
- Import paths
- Packages under `internal/config/`, `internal/assets/tpl/`

**Rule:** If a string is used for comparison, path construction, or
appears in 3+ files, it belongs in `internal/config/` as a constant.
Format strings belong in `internal/config/` as named constants
(e.g., `cfgGit.FlagLastN`, `cfgTrace.RefFormat`). User-facing prose
belongs in `internal/assets/` YAML files accessed via `desc.Text()`.

**Common fix for `fmt.Sprintf` with format verbs:**

| Pattern | Fix |
|---------|-----|
| `fmt.Sprintf("%d", n)` | `strconv.Itoa(n)` |
| `fmt.Sprintf("%d", int64Val)` | `strconv.FormatInt(int64Val, 10)` |
| `fmt.Sprintf("%x", bytes)` | `hex.EncodeToString(bytes)` |
| `fmt.Sprintf("%q", s)` | `strconv.Quote(s)` |
| `fmt.Sscanf(s, "%d", &n)` | `strconv.Atoi(s)` |
| `fmt.Sprintf("-%d", n)` | `fmt.Sprintf(cfgGit.FlagLastN, n)` |
| `"https://"` | `cfgHTTP.PrefixHTTPS` |
| `"&lt;"` | config constant in `config/html/` |

---

## Direct Printf Calls

**Test:** `TestNoPrintfCalls`

`cmd.Printf` and `cmd.PrintErrf` bypass the write-package formatting
pipeline and scatter user-facing text across the codebase.

**Before:**

```go
func Run(cmd *cobra.Command, args []string) {
    cmd.Printf("Found %d tasks\n", count)
}
```

**After:**

```go
func Run(cmd *cobra.Command, args []string) {
    write.TaskCount(cmd, count)
}
```

**Rule:** All formatted output goes through `internal/write/` which
uses `cmd.Print`/`cmd.Println` with pre-formatted strings from
`desc.Text()`.

---

## Raw Time Format Strings

**Test:** `TestNoRawTimeFormats`

Inline time format strings (`"2006-01-02"`, `"15:04:05"`) drift when
one call site is updated but others are missed.

**Before:**

```go
func formatDate(t time.Time) string {
    return t.Format("2006-01-02")
}
```

**After:**

```go
func formatDate(t time.Time) string {
    return t.Format(cfgTime.DateFormat)
}
```

**Rule:** All time format strings must use constants from
`internal/config/time/`.

---

## Direct Flag Registration

**Test:** `TestNoFlagBindOutsideFlagbind`

Direct cobra flag calls (`.Flags().StringVar()`, etc.) scatter flag
wiring across dozens of `cmd.go` files. Centralizing through
`internal/flagbind/` gives one place to audit flag names, defaults,
and description key lookups.

**Before:**

```go
func Cmd() *cobra.Command {
    var output string
    c := &cobra.Command{Use: cmd.UseStatus}
    c.Flags().StringVarP(&output, "output", "o", "",
        "output format")
    return c
}
```

**After:**

```go
func Cmd() *cobra.Command {
    var output string
    c := &cobra.Command{Use: cmd.UseStatus}
    flagbind.StringFlagShort(c, &output, flag.Output,
        flag.OutputShort, cmd.DescKeyOutput)
    return c
}
```

**Rule:** All flag registration goes through `internal/flagbind/`.
If the helper you need doesn't exist, add it to `flagbind/flag.go`
before using it.

---

## TODO Comments

**Test:** `TestNoTODOComments`

TODO, FIXME, HACK, and XXX comments in production code are invisible
to project tracking. They accumulate silently and never get addressed.

**Before:**

```go
// TODO: handle pagination
func listEntries() []Entry {
```

**After:**

Remove the comment and add a task to `.context/TASKS.md`:

```
- [ ] Handle pagination in listEntries (internal/task/task.go)
```

**Rule:** Deferred work lives in TASKS.md, not in source comments.

---

## Dead Exports

**Test:** `TestNoDeadExports`

Exported symbols with zero references outside their definition file
are dead weight. They increase API surface, confuse contributors, and
cost maintenance.

**Fix:** Either delete the export (preferred) or demote it to
unexported if it's still used within the file.

If the symbol existed for historical reasons and might be needed again,
move it to `quarantine/deadcode/` with a `.dead` extension. This
preserves the code in git without polluting the live codebase:

```
quarantine/deadcode/internal/config/flag/flag.go.dead
```

Each `.dead` file includes a header:

```go
// Dead exports quarantined from internal/config/flag/flag.go
// Quarantined: 2026-04-02
// Restore from git history if needed.
```

**Rule:** If a test-only allowlist entry is needed (the export exists
only for test use), add the fully qualified symbol to `testOnlyExports`
in `dead_exports_test.go`. Keep this list small — prefer eliminating
the export.

---

## Core Package Structure

**Test:** `TestCoreStructure`

`core/` directories under `internal/cli/` must contain only `doc.go`
and test files at the top level. All domain logic lives in subpackages.
This prevents `core/` from becoming a god package.

**Before:**

```
internal/cli/dep/core/
    go.go           # violation — logic at core/ level
    python.go       # violation
    node.go         # violation
    types.go        # violation
```

**After:**

```
internal/cli/dep/core/
    doc.go          # package doc only
    golang/
        golang.go
        golang_test.go
        doc.go
    python/
        python.go
        python_test.go
        doc.go
    node/
        node.go
        node_test.go
        doc.go
```

**Rule:** Extract each logical unit into its own subpackage under
`core/`. Each subpackage gets a `doc.go`. The subpackage name should
match the domain concept (`golang`, `check`, `fix`, `store`), not a
generic label (`util`, `helper`).

---

## Cross-Package Types

**Test:** `TestCrossPackageTypes`

When a type defined in one package is used from a different module
(e.g., `cli/doctor` importing a type from `cli/notify`), the type
has crossed its module boundary. Cross-cutting types belong in
`internal/entity/` for discoverability.

**Before:**

```go
// internal/cli/notify/core/types.go
type NotifyPayload struct { ... }

// internal/cli/doctor/core/check/check.go
import "github.com/ActiveMemory/ctx/internal/cli/notify/core"
func check(p core.NotifyPayload) { ... }
```

**After:**

```go
// internal/entity/notify.go
type NotifyPayload struct { ... }

// internal/cli/doctor/core/check/check.go
import "github.com/ActiveMemory/ctx/internal/entity"
func check(p entity.NotifyPayload) { ... }
```

**Exempt:** Types inside `entity/`, `proto/`, `core/` subpackages,
and `config/` packages. Same-module usage (e.g., `cli/doctor/cmd/`
using `cli/doctor/core/`) is not flagged.

---

## Type File Convention

**Test:** `TestTypeFileConvention`, `TestTypeFileConventionReport`

Exported types in `core/` subpackages should live in `types.go` (the
convention from CONVENTIONS.md), not scattered across implementation
files. This makes type definitions discoverable.
`TestTypeFileConventionReport` generates a diagnostic summary of all
type placements for triage.

**Exception:** `entity/` organizes by domain (`task.go`, `session.go`),
`proto/` uses `schema.go`, and `err/` packages colocate error types
with their domain context.

---

## DescKey / YAML Linkage

**Test:** `TestDescKeyYAMLLinkage`

Every DescKey constant must have a corresponding key in the YAML asset
files, and every YAML key must have a corresponding DescKey constant.
Orphans in either direction mean dead text or runtime panics.

**Fix for orphan YAML key:** Delete the YAML entry, or add the
corresponding `DescKey` constant in `config/embed/{text,cmd,flag}/`.

**Fix for orphan DescKey:** Delete the constant, or add the
corresponding entry in the YAML file under
`internal/assets/commands/text/`, `cmd/`, or `flag/`.

If the orphan YAML entry was once valid but the feature was removed,
move the YAML entry to a `.dead` file in `quarantine/deadcode/`.

---

## Package Doc Quality

**Test:** `TestPackageDocQuality`

Every package under `internal/` must have a `doc.go` with a meaningful
package doc comment (at least 8 lines of real content). One-liners and
file-list patterns (`// - foo.go`, `// Source files:`) are flagged
because they drift as files change.

**Template:**

```go
//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mypackage does X.
//
// It handles Y by doing Z. The main entry point is [FunctionName]
// which accepts A and returns B.
//
// Configuration is read from [config.SomeConstant]. Output is
// written through [write.SomeHelper].
//
// This package is used by [parentpackage] during the W lifecycle
// phase.
package mypackage
```

---

## Inline Regex Compilation

**Test:** `TestNoInlineRegexpCompile`

`regexp.MustCompile` and `regexp.Compile` inside function bodies
recompile the pattern on every call. Compiled patterns belong at
package level.

**Before:**

```go
func parse(s string) bool {
    re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
    return re.MatchString(s)
}
```

**After:**

```go
// In internal/config/regex/regex.go:
// DatePattern matches ISO date format (YYYY-MM-DD).
var DatePattern = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

// In calling package:
func parse(s string) bool {
    return regex.DatePattern.MatchString(s)
}
```

**Rule:** All compiled regexes live in `internal/config/regex/` as
package-level `var` declarations. Two tests enforce this:
`TestNoInlineRegexpCompile` catches function-body compilation, and
`TestNoRegexpOutsideRegexPkg` catches package-level compilation
outside `config/regex/`.

---

## Doc Comments

**Test:** `TestDocComments`

All functions (exported and unexported), structs, and package-level
variables must have a doc comment. Config packages allow group doc
comments for `const` blocks.

**Before:**

```go
func buildIndex(entries []Entry) map[string]int {
```

**After:**

```go
// buildIndex maps entry names to their position in the
// ordered slice for O(1) lookup during reconciliation.
//
// Parameters:
//   - entries: ordered slice of entries to index
//
// Returns:
//   - map[string]int: name-to-position mapping
func buildIndex(entries []Entry) map[string]int {
```

**Rule:** Every function, struct, and package-level `var` gets a doc
comment in godoc format. Functions include `Parameters:` and
`Returns:` sections. Structs with 2+ fields document every field.
See CONVENTIONS.md for the full template.

---

## Line Length

**Test:** `TestLineLength`

Lines in non-test Go files must not exceed 80 characters. This is a
hard check, not a suggestion.

**Before:**

```go
_ = trace.Record(fmt.Sprintf(cfgTrace.RefFormat, cfgTrace.RefTypeTask, matchedNum), state.Dir())
```

**After:**

```go
ref := fmt.Sprintf(
    cfgTrace.RefFormat, cfgTrace.RefTypeTask, matchedNum,
)
_ = trace.Record(ref, state.Dir())
```

**Rule:** Break at natural points: function arguments, struct fields,
chained calls. Long strings (URLs, struct tags) are the rare
acceptable exception.

---

## Literal Whitespace

**Test:** `TestNoLiteralWhitespace`

Bare whitespace string and byte literals (`"\n"`, `"\r\n"`, `"\t"`)
must not appear outside `internal/config/token/`. All other packages
use the token constants.

**Before:**

```go
output := strings.Join(lines, "\n")
```

**After:**

```go
output := strings.Join(lines, token.Newline)
```

**Rule:** Whitespace literals are defined once in
`internal/config/token/`. Use `token.Newline`, `token.Tab`,
`token.CRLF`, etc.

---

## Magic Numeric Values

**Test:** `TestNoMagicValues`

Numeric literals in function bodies need constants, with narrow
exceptions.

**Before:**

```go
if len(entries) > 100 {
    entries = entries[:100]
}
```

**After:**

```go
if len(entries) > config.MaxEntries {
    entries = entries[:config.MaxEntries]
}
```

**Exempt:** `0`, `1`, `-1`, `2`–`10`, strconv radix/bitsize args
(`10`, `32`, `64` in `strconv.Parse*`/`Format*`), octal permissions
(caught separately by `TestNoRawPermissions`), and `const`/`var`
definition sites.

---

## Inline Separators

**Test:** `TestNoInlineSeparators`

`strings.Join` calls must use token constants for their separator
argument, not string literals.

**Before:**

```go
result := strings.Join(parts, ", ")
```

**After:**

```go
result := strings.Join(parts, token.CommaSep)
```

**Rule:** Separator strings live in `internal/config/token/`.

---

## Stuttery Function Names

**Test:** `TestNoStutteryFunctions`

Function names must not redundantly include their package name as a
PascalCase word boundary. Go callers already write `pkg.Function`,
so `pkg.PkgFunction` stutters.

**Before:**

```go
// In package write
func WriteJournal(cmd *cobra.Command, ...) {
```

**After:**

```go
// In package write
func Journal(cmd *cobra.Command, ...) {
```

**Exempt:** Identity functions like `write.Write` / `write.write`.

---

## Predicate Naming (no `Is`/`Has`/`Can` prefix)

**Test:** None (manual review convention)

Exported methods that return `bool` must not use `Is`, `Has`, or
`Can` prefixes. The predicate reads more naturally without them,
especially at call sites where the package name provides context.

**Before:**

```go
func IsCompleted(t *Task) bool { ... }
func HasChildren(n *Node) bool { ... }
func IsExemptPackage(path string) bool { ... }
```

**After:**

```go
func Completed(t *Task) bool { ... }
func Children(n *Node) bool { ... }  // or: ChildCount > 0
func ExemptPackage(path string) bool { ... }
```

**Rule:** Drop the prefix. Private helpers may use prefixes when it
reads more naturally (`isValid` in a local context is fine). This
convention applies to exported methods and package-level functions.
See CONVENTIONS.md "Predicates" section.

This is not yet enforced by an AST test — it requires semantic
understanding of return types and naming intent that makes automated
detection fragile. Apply during code review.

---

## Mixed Visibility

**Test:** `TestNoMixedVisibility`

Files with exported functions must not also contain unexported
functions. Public API and private helpers live in separate files.

**Before:**

```
load.go
    func Load() { ... }        // exported
    func parseHeader() { ... } // unexported — violation
```

**After:**

```
load.go
    func Load() { ... }        // exported only
parse.go
    func parseHeader() { ... } // private helper
```

**Exempt:** Files with exactly one function, `doc.go`, test files.

---

## Stray err.go Files

**Test:** `TestNoStrayErrFiles`

`err.go` files must only exist under `internal/err/`. Error
constructors anywhere else create a broken-window pattern where
contributors add local error definitions when they see a local
`err.go`.

**Fix:** Move the error constructor to `internal/err/<domain>/`.

---

## CLI Cmd Structure

**Test:** `TestCLICmdStructure`

Each `cmd/$sub/` directory under `internal/cli/` may contain only
`cmd.go`, `run.go`, `doc.go`, and test files. Extra `.go` files
(helpers, output formatters, types) belong in the corresponding
`core/` subpackage.

**Before:**

```
internal/cli/doctor/cmd/root/
    cmd.go
    run.go
    format.go   # violation — helper in cmd dir
```

**After:**

```
internal/cli/doctor/cmd/root/
    cmd.go
    run.go
internal/cli/doctor/core/format/
    format.go
    doc.go
```

---

## DescKey Namespace

**Test:** `TestUseConstantsOnlyInCobraUse`,
`TestDescKeyOnlyInLookupCalls`,
`TestNoWrongNamespaceLookup`

Three tests enforce DescKey/Use constant discipline:

1. `Use*` constants appear only in cobra `Use:` struct field
   assignments — never as arguments to `desc.Text()` or elsewhere.
2. `DescKey*` constants are passed only to `assets.CommandDesc()`,
   `assets.FlagDesc()`, or `desc.Text()` — never to cobra `Use:`.
3. No cross-namespace lookups — `TextDescKey` must not be passed to
   `CommandDesc()`, `FlagDescKey` must not be passed to `Text()`, etc.

---

## YAML Examples / Registry Linkage

**Test:** `TestExamplesYAMLLinkage`, `TestRegistryYAMLLinkage`

Every key in `examples.yaml` and `registry.yaml` must match a known
entry type constant. Prevents orphan entries that are never rendered.

**Fix:** Delete the orphan YAML entry, or add the corresponding
constant in `config/entry/`.

---

## Other Enforced Patterns

These tests follow the same fix approach — extract the operation to
its designated package:

| Test | Violation | Fix |
|------|-----------|-----|
| `TestNoNakedErrors` | `fmt.Errorf`/`errors.New` outside `internal/err/` | Add error constructor to `internal/err/<domain>/` |
| `TestNoRawFileIO` | Direct `os.ReadFile`, `os.Create`, etc. | Use `io.SafeReadFile`, `io.SafeWriteFile`, etc. |
| `TestNoRawLogging` | Direct `fmt.Fprintf(os.Stderr, ...)` | Use `log/warn.Warn()` or `log/event.Append()` |
| `TestNoExecOutsideExecPkg` | `exec.Command` outside `internal/exec/` | Add command to `internal/exec/<domain>/` |
| `TestNoCmdPrintOutsideWrite` | `cmd.Print*` outside `internal/write/` | Add output helper to `internal/write/<domain>/` |
| `TestNoRawPermissions` | Octal literals (`0644`, `0755`) | Use `config/fs.PermFile`, `config/fs.PermExec`, etc. |
| `TestNoErrorsAs` | `errors.As()` | Use `errors.AsType()` (generic, Go 1.23+) |
| `TestNoStringConcatPaths` | `dir + "/" + file` | Use `filepath.Join(dir, file)` |

---

## General Fix Workflow

When an audit test fails:

1. **Read the error message.** It includes `file:line` and a
   description of the violation.
2. **Find the matching section above.** The test name maps directly
   to a section.
3. **Apply the pattern.** Most fixes are mechanical: extract to the
   right package, rename a variable, or replace a literal with a
   constant.
4. **Run `make test` before committing.** Audit tests run as part of
   `go test ./internal/audit/`.
5. **Don't add allowlist entries as a first resort.** Fix the code.
   Allowlists exist only for genuinely unfixable cases (test-only
   exports, config packages that are definitionally exempt).

# code style

## Naming

- **Constants use semantic prefixes**: Group related constants with prefixes
  - `Dir*` for directories (`DirContext`, `DirArchive`)
  - `File*` for file paths (`FileSettings`, `FileClaudeMd`)
  - `Filename*` for file names only (`FilenameTask`, `FilenameDecision`)
  - `*Type*` for enum-like values (`UpdateTypeTask`, `UpdateTypeDecision`)
- **Package name = folder name**: Go canonical pattern
  - `package initialize` in `initialize/` folder
  - Never `package initcmd` in `init/` folder
- **Go package names: lowercase, no underscores, no
  mixedCaps**: per the [Effective Go](https://go.dev/blog/package-names)
  guidance and the stdlib precedent (`strconv`, `httptest`,
  `bufio`). Apply to the directory too — `internal/flagbind/`,
  not `internal/flag_bind/`. Filenames may use underscores
  (`foo_test.go` is canonical); package names may not. When in
  doubt, find the closest stdlib analogue and copy its shape.
- **Maps reference constants**: Use constants as keys, not literals
  - `map[string]X{ConstKey: value}` not `map[string]X{"literal": value}`

## Casing

- **Proper nouns keep their casing** in comments, strings, and docs
  - `Markdown` not `markdown` (it's a language name)
  - `YAML`, `JSON`, `TOML` — always uppercase
  - `GitHub`, `JavaScript`, `PostgreSQL` — match official casing
  - Exception: code fence language identifiers are lowercase (`` ```markdown ``)

## Predicates

- **No Is/Has/Can prefixes**: `Completed()` not
  `IsCompleted()`, `Empty()` not `IsEmpty()`
- Applies to all bool-returning funcs and methods, exported or
  not (`topicNested`, not `hasNestedTopic`)

## Line Width

- **Target ~80 characters**: Highly encouraged, not a hard limit
  - Some lines will naturally exceed it (long strings,
    struct tags, URLs) — that's fine
  - Drift accumulates silently, especially in test code
  - Break at natural points: function arguments, struct fields, chained calls

## Patterns

- **Centralize magic strings**: All repeated literals
  belong in a `config` or `constants` package
  - If a string appears in 3+ files, it needs a constant
  - If a string is used for comparison, it needs a constant
- **Path construction**: Always use stdlib path joining
  - Go: `filepath.Join(dir, file)`
  - Python: `os.path.join(dir, file)`
  - Node: `path.join(dir, file)`
  - Never: `dir + "/" + file`
- **Constants reference constants**: Self-referential definitions
  - `FileType[UpdateTypeTask] = FilenameTask` not
    `FileType["task"] = "TASKS.md"`
- **No error variable shadowing**: Use descriptive names
  when multiple errors exist in a function
  - `readErr`, `writeErr`, `indexErr` — not repeated `err` / `err :=`
  - Shadowed `err` silently disconnects from the outer
    variable, causing subtle bugs
- **Colocate related code**: Group by feature, not by type
  - `session/run.go`, `session/types.go`, `session/parse.go`
  - Not: `runners/session.go`, `types/session.go`, `parsers/session.go`

## Duplication

- **Non-test code**: Apply the rule of three — extract
  when a block appears 3+ times
  - Watch for copy-paste during task-focused sessions
    where the agent prioritizes completion over shape
- **Test code**: Some duplication is acceptable for readability
  - When the same setup/assertion block appears 3+ times, extract a test helper
  - Use `t.Helper()` so failure messages point to the caller, not the helper

## Error Handling

- **Zero silent error discard**: Handle every error, never suppress with
  `_ =` or `//nolint:errcheck`. Production: defer-close logs to stderr
  via `log.Warn()`. Test: `t.Fatal(err)` for setup, `t.Log(err)` for
  cleanup. For gosec false positives: fix the code rather than adding
  nolint markers — the goal is zero golangci-lint suppressions
- **Error constructors in internal/err**: Never in per-package err.go
  files — eliminates the broken-window pattern where agents add local
  errors when they see a local err.go exists
- **Identity sentinels are `entity.Sentinel` consts, not
  `errors.New`**: Declare `errors.Is` targets as
  `const ErrX = entity.Sentinel(text.DescKey...)`. The
  user-facing text lives in `commands/text/errors.yaml` keyed by
  `err.<pkg>.<name>`; the sentinel's `Error()` resolves it via
  `desc.Text` at call time. Never write
  `var ErrX = errors.New("english")` — the English leaks into
  `.Error()` output and bypasses localization. Never add an
  `ErrMsg* = "english"` const layer in `internal/config/<pkg>/`
  to back the sentinel; that layer is dead text once the typed
  Sentinel does the lookup itself.
- **Parameterised errors use typed structs**: When the error
  needs to carry fields (path, name, etc.), define a struct in
  `internal/err/<area>/` with a pointer-receiver `Error()` and
  optional `Is(error) bool` for sentinel-compatibility. See
  `internal/err/context.NotFoundError` for the canonical shape.


# Conventions

## Naming

- **Constants use semantic prefixes**: Group related constants with prefixes
  - `Dir*` for directories (`DirContext`, `DirArchive`)
  - `File*` for file paths (`FileSettings`, `FileClaudeMd`)
  - `Filename*` for file names only (`FilenameTask`, `FilenameDecision`)
  - `*Type*` for enum-like values (`UpdateTypeTask`, `UpdateTypeDecision`)
- **Package name = folder name**: Go canonical pattern
  - `package initialize` in `initialize/` folder
  - Never `package initcmd` in `init/` folder
- **Maps reference constants**: Use constants as keys, not literals
  - `map[string]X{ConstKey: value}` not `map[string]X{"literal": value}`

## Predicates

- **No Is/Has/Can prefixes**: `Completed()` not `IsCompleted()`, `Empty()` not `IsEmpty()`
- Applies to exported methods that return bool
- Private helpers may use prefixes when it reads more naturally

## File Organization

- **Public API in main file, private helpers in separate logical files**
  - `loader.go` (exports `Load()`) + `process.go` (unexported helpers)
  - NOT: one file with unexported functions stacked at the bottom
- Reasoning: agent loads only the public API file unless it needs implementation detail

## Patterns

- **Centralize magic strings**: All repeated literals belong in a `config` or `constants` package
  - If a string appears in 3+ files, it needs a constant
  - If a string is used for comparison, it needs a constant
- **Path construction**: Always use stdlib path joining
  - Go: `filepath.Join(dir, file)`
  - Python: `os.path.join(dir, file)`
  - Node: `path.join(dir, file)`
  - Never: `dir + "/" + file`
- **Constants reference constants**: Self-referential definitions
  - `FileType[UpdateTypeTask] = FilenameTask` not `FileType["task"] = "TASKS.md"`
- **Colocate related code**: Group by feature, not by type
  - `session/run.go`, `session/types.go`, `session/parse.go`
  - Not: `runners/session.go`, `types/session.go`, `parsers/session.go`

## Line Width

- **Target ~80 characters**: Highly encouraged, not a hard limit
  - Some lines will naturally exceed it (long strings, struct tags, URLs) — that's fine
  - Drift accumulates silently, especially in test code
  - Break at natural points: function arguments, struct fields, chained calls

## Duplication

- **Non-test code**: Apply the rule of three — extract when a block appears 3+ times
  - Watch for copy-paste during task-focused sessions where the agent prioritizes completion over shape
- **Test code**: Some duplication is acceptable for readability
  - When the same setup/assertion block appears 3+ times, extract a test helper
  - Use `t.Helper()` so failure messages point to the caller, not the helper

## Testing

- **Colocate tests**: Test files live next to source files
  - `foo.go` → `foo_test.go` in same package
  - Not a separate `tests/` folder
- **Test the unit, not the file**: One test file can test multiple related functions
- **Integration tests are separate**: `cli_test.go` for end-to-end binary tests

## Code Change Heuristics

- **Present interpretations, don't pick silently**: If a request has multiple
  valid readings, lay them out rather than guessing
- **Push back when warranted**: If a simpler approach exists, say so
- **"Would a senior engineer call this overcomplicated?"**: If yes, simplify
- **Match existing style**: Even if you'd write it differently in a greenfield
- **Every changed line traces to the request**: If it doesn't, revert it

## Decision Heuristics

- **"Would I start this today?"**: If not, continuing is the sunk cost — evaluate only future value
- **"Reversible or one-way door?"**: Reversible decisions don't need deep analysis
- **"Does the analysis cost more than the decision?"**: Stop deliberating when the options are within an order of magnitude
- **"Order of magnitude, not precision"**: 10x better matters; 10% better usually doesn't

## Refactoring

- **Measure the end state, not the effort**: When refactoring, ask what the
  codebase looks like *after*, not how much work the change is
- **Three questions before restructuring**:
  1. What's the smallest codebase that solves this?
  2. Does the proposed change result in less total code?
  3. What can we delete now that this change makes obsolete?
- **Deletion is a feature**: Writing 50 lines that delete 200 is a net win

## Documentation

- **Godoc format**: Use canonical sections
  ```go
  // FunctionName does X.
  //
  // Longer description if needed.
  //
  // Parameters:
  //   - param1: Description
  //   - param2: Description
  //
  // Returns:
  //   - Type: Description of return value
  func FunctionName(param1, param2 string) error
  ```
- **Package doc in doc.go**: Each package gets a `doc.go` with package-level documentation
- **Copyright headers**: All source files get the project copyright header

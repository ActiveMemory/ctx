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

## Casing

- **Proper nouns keep their casing** in comments, strings, and docs
  - `Markdown` not `markdown` (it's a language name)
  - `YAML`, `JSON`, `TOML` — always uppercase
  - `GitHub`, `JavaScript`, `PostgreSQL` — match official casing
  - Exception: code fence language identifiers are lowercase (`` ```markdown ``)

## Predicates

- **No Is/Has/Can prefixes**: `Completed()` not `IsCompleted()`, `Empty()` not `IsEmpty()`
- Applies to exported methods that return bool
- Private helpers may use prefixes when it reads more naturally

## File Organization

- **Public API in main file, private helpers in separate logical files**
  - `loader.go` (exports `Load()`) + `process.go` (unexported helpers)
  - NOT: one file with unexported functions stacked at the bottom
- Reasoning: agent loads only the public API file unless it needs implementation detail
- **Name files after what they contain, not their role**
  - `format.go`, `sort.go`, `parse.go` — named by responsibility
  - NOT: `util.go`, `utils.go`, `helper.go`, `common.go` — junk drawer names
  - If a file can't be named without a generic label, its contents don't belong together
  - Existing junk drawers should be split as their contents grow

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
- **No error variable shadowing**: Use descriptive names when multiple errors exist in a function
  - `readErr`, `writeErr`, `indexErr` — not repeated `err` / `err :=`
  - Shadowed `err` silently disconnects from the outer variable, causing subtle bugs
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
- **Package doc in doc.go**: Each package gets a `doc.go` with package-level
  documentation describing behavior, not structure. Do NOT include
  `# File Organization` sections listing files — they drift when files are
  added, renamed, or removed, and the filesystem is self-documenting
- **Copyright headers**: All source files get the project copyright header

## Blog Publishing

- **Checklist for ideas/ → docs/blog/ promotion**:
  1. Update date in frontmatter to publish date
  2. Fix relative paths (from `../docs/blog/` to peer references)
  3. Add cross-links to/from companion posts ("See also" sections)
  4. Add "The Arc" section connecting to the series narrative
  5. Update `docs/blog/index.md` with entry (newest first)
  6. Verify all link targets exist
  7. Build and test before commit
- **Arc section**: Every post includes "The Arc" near the end, framing
  where the post sits in the broader blog narrative
- **See also links**: Use italic `*See also: [Title](file) -- one-line
  description connecting the two posts.*` format at the end of posts
- **Frontmatter**: Include copyright header, title, date, author, topics list
- **Blog index order**: Newest post first, with topic tags and 3-4 line summary

- Update admonitions for historical blog content: Use MkDocs admonitions (\!\!\! note "Update") at the top of blog post sections where features have been superseded or installation has changed. Link to current documentation. Keep original content intact below for historical context.

- New CLI subcommand documentation checklist: When adding a new CLI subcommand, update docs in at least three places: (1) Feature page (e.g., docs/scratchpad.md) — commands table, usage section, skill/NL table. (2) CLI reference (docs/cli-reference.md) — full reference entry with args, flags, examples. (3) Relevant recipes — any recipe that covers the feature area. (4) zensical.toml — only if adding a new page.

- Always stage site/ when committing docs/ changes — the generated HTML is tracked in git with no CI build step.

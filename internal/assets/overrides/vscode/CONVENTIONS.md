# Conventions

<!--
UPDATE WHEN:
- New pattern is established and should be followed consistently
- Existing pattern is deprecated or superseded
- Team adopts new tooling that changes workflows
- Code review reveals recurring issues that need a convention

DO NOT UPDATE FOR:
- One-off exceptions (document in code comments)
- Experimental patterns not yet proven
- Personal preferences without team consensus
-->

## Naming

- **Use semantic prefixes for constants**: Group related constants with prefixes
  - `DIR_*` / `Dir*` for directories
  - `FILE_*` / `File*` for file paths
  - `*_TYPE` / `*Type` for enum-like values
- **Module/package name = folder name**: Keep names consistent with the filesystem
- **Avoid magic strings**: Use named constants instead of string literals for comparison

## Patterns

- **Centralize repeated literals**: All repeated literals belong in a constants/config module
  - If a string appears in 3+ files, it needs a constant
  - If a string is used for comparison, it needs a constant
- **Path construction**: Always use your language's standard path joining
  - Python: `os.path.join(dir, file)` or `pathlib.Path(dir) / file`
  - Node/TS: `path.join(dir, file)`
  - Go: `filepath.Join(dir, file)`
  - Rust: `PathBuf::from(dir).join(file)`
  - Never: `dir + "/" + file` (string concatenation)
- **Colocate related code**: Group by feature, not by type
  - `session/run.ext`, `session/types.ext`, `session/parse.ext`
  - Not: `runners/session.ext`, `types/session.ext`, `parsers/session.ext`

## Testing

- **Colocate tests**: Test files live next to source files
  - Not in a separate `tests/` folder (unless the language convention requires it)
- **Test the unit, not the file**: One test file can test multiple related functions
- **Integration tests are separate**: Clearly distinguish unit tests from end-to-end tests

## Documentation

- **Follow language conventions**: Use the standard doc format for your language
  - Python: docstrings (Google/NumPy/Sphinx style)
  - TypeScript/JavaScript: JSDoc or TSDoc
  - Go: Godoc comments
  - Rust: `///` doc comments with Markdown
- **Document public APIs**: Every exported function/class/type gets a doc comment
- **Copyright headers**: All source files get the project copyright header

# ctx fmt

## Problem

Context file entries — especially tasks — are written as single long
lines. Agents produce verbose descriptions that blow past 80 chars,
making files hard to scan in terminals and diffs. No formatting step
exists today; the only workaround is manual line breaks.

## Approach

- Elevate `internal/cli/journal/core/wrap/` to `internal/wrap/`
  so wrapping logic is shared between journal and context formatting
- Add list-aware wrapping: 2-space continuation indent for markdown
  list items (`- [ ]`, `- [x]`, `- [-]`, `- `)
- New `ctx fmt` subcommand that formats all context files in-place
- `make fmt-context` target as convenience alias

## Behavior

### Happy Path

1. User runs `ctx fmt`
2. Command resolves the context directory (same as all ctx commands)
3. Reads TASKS.md, DECISIONS.md, LEARNINGS.md, CONVENTIONS.md
4. Wraps long lines to 80 characters at word boundaries
5. Writes back only files that changed
6. Prints summary: `Formatted 2 of 4 context files`

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Line is exactly 80 chars | Leave it alone |
| No word boundary before 80 chars (long URL) | Don't break mid-word, let it exceed 80 |
| Already-wrapped continuation lines | Idempotent: don't double-wrap |
| Markdown headings (`##`) | Skip, don't wrap |
| Markdown tables (`\|` rows) | Skip, don't wrap |
| Frontmatter (`---` blocks) | Skip, don't wrap |
| Tag clusters (`#priority:high #session:abc`) | Wrap normally at word boundaries; tags are just words |
| Empty context file | No-op, no error |
| Missing context file | Skip with warning, continue with remaining files |
| HTML comment blocks (`<!-- -->`) | Skip, don't wrap |

### Validation Rules

- Context directory must exist (standard ctx system bootstrap check)
- No input validation needed: the command operates on known files

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| No context directory | `context directory not found — run ctx init` | Run `ctx init` |
| File read/write failure | `failed to format {file}: {error}` | Check file permissions |
| All files missing | `no context files found in {dir}` | Run `ctx init` |

## Interface

### CLI

```
ctx fmt [flags]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--width` | int | 80 | Target line width |
| `--check` | bool | false | Check only, exit 1 if files would change (CI mode) |
| `--quiet` | bool | false | Suppress output |

### Make target

```makefile
fmt-context:
	ctx fmt
```

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/wrap/wrap.go` | Move from `internal/cli/journal/core/wrap/wrap.go`, add list-aware wrapping |
| `internal/wrap/wrap_test.go` | Move and extend tests |
| `internal/cli/fmt/cmd/root/cmd.go` | New `ctx fmt` command |
| `internal/cli/fmt/cmd/root/run.go` | Command logic: read-wrap-write loop |
| `internal/cli/fmt/core/format.go` | Context file formatting orchestration |
| `internal/cli/journal/core/wrap/` | Remove; update all journal imports to `internal/wrap/` |
| `Makefile` | Add `fmt-context` target |
| `internal/config/flag/flag.go` | Add `Width` constant |
| `internal/config/embed/flag/fmt.go` | Flag description keys |

### Key Functions

```go
// internal/wrap/wrap.go

// ContextFile wraps long lines in a context file (.context/*.md).
// Handles list continuation with 2-space indent.
func ContextFile(content string) string

// ListItem wraps a markdown list line, using 2-space continuation
// indent for wrapped lines.
func ListItem(line string, width int) []string
```

### Helpers to Reuse

- `wrap.Soft()` — word-boundary breaking (already exists)
- `wrap.Content()` — frontmatter/table skipping (already exists)
- `rc.ContextDir()` — resolve context directory
- `io.SafeReadUserFile()` / `io.SafeWriteFile()` — safe file I/O
- `config/entry.CtxFile()` — file name mapping

## Configuration

- `--width` flag overrides the default 80-char width
- No `.ctxrc` keys needed; 80 is the right default

## Testing

- Unit: `wrap.ListItem()` with checkbox lines, plain list items,
  long URLs, tags, short lines
- Unit: `wrap.ContextFile()` with mixed content (headings, tables,
  frontmatter, list items, paragraphs)
- Unit: idempotency — `ContextFile(ContextFile(input)) == ContextFile(input)`
- Integration: `ctx fmt` on a temp context directory with known
  content, verify output matches expected
- Edge case: file with only headings and tables (nothing to wrap)

## Non-Goals

- Perfect markdown rendering. Good enough for humans and agents.
- Reflowing already-wrapped paragraphs into optimal line lengths.
  Only lines exceeding the width get wrapped.
- Formatting non-context files (README, specs, etc.)
- `--wrap` flag on `ctx add`. Can be added later; `ctx fmt` after
  adding entries achieves the same result.

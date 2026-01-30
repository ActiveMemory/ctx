# Conventions

## Naming

## Patterns

### CLI Output Methods

In CLI command code (`internal/cli/**/*.go`), always use Cobra's command methods
for output instead of raw `fmt` functions:

```go
// ✓ Correct - uses cmd methods
cmd.Printf("%s Done\n", green("✓"))
cmd.Println("Status:")
fmt.Fprintln(cmd.OutOrStdout(), "Output here")

// ✗ Wrong - bypasses Cobra's output handling
fmt.Printf("%s Done\n", green("✓"))
fmt.Println("Status:")
```

**Why this matters:**
- Cobra's `cmd.OutOrStdout()` respects output redirection in tests
- Makes CLI commands testable without capturing os.Stdout
- Consistent pattern across all commands

**Exceptions:**
- Helper functions that don't have access to `*cobra.Command` may use `fmt`
- But prefer passing `io.Writer` or the command down to helpers

## Testing

## Documentation

### Doc-Impact Rule

**When to update docs:** Any change to CLI commands, flags, or behavior 
requires updating the corresponding documentation in `docs/`.

**Mapping:**

| Code Change                         | Doc Update Required              |
|-------------------------------------|----------------------------------|
| `internal/cli/*.go` command changes | `docs/cli-reference.md`          |
| `.context/` file format changes     | `docs/context-files.md`          |
| AI tool integration changes         | `docs/integrations.md`           |
| New features                        | `docs/index.md` (if user-facing) |

**How to remember:** When marking a CLI-related task complete, check if `docs/`
needs updating. Add `#doc-impact` tag to tasks that affect documentation.

**Drift detection:** Run `ctx drift` — it will warn if `internal/cli/` 
is newer than `docs/cli-reference.md`.

- Decision records must include filled-in Context, Rationale, and Consequences sections. Never leave placeholder text.

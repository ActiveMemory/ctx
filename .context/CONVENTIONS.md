# Conventions

## Naming

## Patterns

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

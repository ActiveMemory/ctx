# Conventions

Coding standards and patterns used in this project.

## Naming

- Use camelCase for variables and functions
- Use PascalCase for types and interfaces
- Use SCREAMING_SNAKE_CASE for constants

## Code Style

- Prefer early returns over nested conditionals
- Maximum line length: 100 characters
- One component per file

## Patterns

### Error Handling

Always return errors, never panic in library code:

```go
// ✓ Correct
func ProcessData(data []byte) (Result, error) {
    if len(data) == 0 {
        return Result{}, fmt.Errorf("empty data")
    }
    // ...
}

// ✗ Wrong
func ProcessData(data []byte) Result {
    if len(data) == 0 {
        panic("empty data")  // Never panic in libraries
    }
    // ...
}
```

Wrap errors with context:

```go
if err != nil {
    return fmt.Errorf("processing user %s: %w", userID, err)
}
```

### Configuration

Load order (highest priority first):
1. Environment variables
2. Config file (config.yaml)
3. Default values

Log config source at startup for debuggability.

## Testing

- Test files adjacent to source files (`foo.go` → `foo_test.go`)
- Use table-driven tests for multiple cases
- Mock external dependencies, never call real services in tests

## Git Practices

- Commit messages follow Conventional Commits format
- Feature branches: `feature/<description>`
- Bug fixes: `fix/<description>`
- All PRs require at least one approval

## Documentation

### Doc-Impact Rule

When modifying code that affects user-facing behavior, update the corresponding
documentation:

| Code Change              | Doc Update Required    |
|--------------------------|------------------------|
| API endpoint changes     | `docs/api.md`          |
| CLI command changes      | `docs/cli.md`          |
| Configuration changes    | `docs/configuration.md`|
| New features             | `README.md`            |

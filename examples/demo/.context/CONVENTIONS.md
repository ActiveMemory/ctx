# Conventions

Coding standards and patterns used in this project.

## Code Style

- Use camelCase for variables and functions
- Use PascalCase for types and interfaces
- Prefer early returns over nested conditionals
- Maximum line length: 100 characters

## File Organization

- One component per file
- Group related files in directories
- Test files should be adjacent to source files

## Git Practices

- Commit messages follow Conventional Commits
- Feature branches: `feature/<description>`
- Bug fixes: `fix/<description>`
- All PRs require at least one approval

## Error Handling

- Always return errors, never panic in libraries
- Wrap errors with context using `fmt.Errorf`
- Log errors at the boundary, not in helpers

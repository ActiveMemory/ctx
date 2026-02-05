---
name: update-docs
description: "Check if docs and code conventions are consistent after changes. Use after modifying source code, before committing, or when asked to sync docs."
---

When source code changes, check public docs AND internal code conventions.

## Workflow

1. **Diff** the branch: `git diff main --stat` (or relevant base)
2. **Verify mapping** is current (see Self-Maintenance below)
3. **Map** changed packages to affected docs (see table below)
4. **Read** each affected doc — flag sections that contradict the new code
5. **Update** or flag for the user
6. **Validate**: `mkdocs build` in docs site (if available)

## Code-to-Docs Mapping

| Source Path                        | Likely Affected Docs                               |
|------------------------------------|----------------------------------------------------|
| `cmd/ctx/`, `internal/cli/`        | `docs/cli-reference.md`                            |
| `internal/config/`                 | `docs/context-files.md`                            |
| `internal/context/`                | `docs/context-files.md`, `docs/prompting-guide.md` |
| `internal/drift/`                  | `docs/context-files.md` (drift section)            |
| `internal/recall/`                 | `docs/session-journal.md`                          |
| `internal/bootstrap/`              | `docs/index.md` (getting started)                  |
| `internal/claude/`, `internal/rc/` | `docs/integrations.md`                             |
| `internal/tpl/`                    | `docs/context-files.md` (templates)                |
| `SECURITY.md`                      | `docs/security.md`                                 |
| `.context/` schema changes         | `docs/context-files.md`                            |

## What to Check

- **New CLI flags/commands**: Are they in `docs/cli-reference.md`?
- **Changed file formats**: Does `docs/context-files.md` match?
- **New context files**: Added to both the read order docs and `docs/context-files.md`?
- **Removed features**: Still referenced in docs?
- **Changed defaults**: Do examples in docs use the old defaults?

## Self-Maintenance

This mapping table will drift. Before relying on it:

1. `ls internal/` — any packages not in the table? Add them.
2. `ls docs/*.md` — any doc pages not in the table? Map them.
3. If you update the table, edit this skill file directly.

The skill is its own first test case: if the mapping is stale, the skill
has already failed at its job.

## Internal Code Conventions

Also check that changed code follows project patterns (not Go defaults):

### Godoc Style
Project uses explicit **Parameters/Returns** sections, not standard godoc.
```go
// Good (project style):
// FunctionName does X.
//
// Parameters:
//   - param1: Description
//
// Returns:
//   - Type: Description
func FunctionName(param1 string) error

// Bad (standard godoc — agent corpus drift):
// FunctionName does X with param1.
func FunctionName(param1 string) error
```

Verify that godoc comments match actual parameters and behavior.

### Predicate Naming
Project uses predicates **without** Is/Has/Can prefixes:
- `Completed()` not `IsCompleted()`
- `Empty()` not `IsEmpty()`
- `Exists()` not `DoesExist()`

### File Organization
Public API in the main file, private helpers in **separate logical files**:
- `loader.go` (public `Load()`) + `process.go` (private helpers)
- NOT: everything in one file with unexported functions at the bottom

### Magic Strings
Literals belong in `internal/config/`. If you see a hardcoded string used
in 2+ files, it needs a constant. Check `internal/config/` for existing
constants before introducing new literals.

## Relationship to ctx drift

`ctx drift` checks `.context/` file health (dead paths, staleness).
This skill checks `docs/` ↔ source code alignment and internal conventions.

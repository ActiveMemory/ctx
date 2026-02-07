---
name: update-docs
description: "Check if docs and code conventions are consistent after changes. Use after modifying source code, before committing, or when asked to sync docs."
---

When source code changes, check public docs AND internal code
conventions for consistency.

## When to Use

- After modifying source code that has user-facing behavior
- Before committing changes that touch CLI flags, file formats,
  or defaults
- When the user asks to sync docs with code
- After adding or removing commands, subcommands, or flags

## When NOT to Use

- After changes that are purely internal with no docs impact
  (e.g., renaming a private variable)
- When only docs were changed (no code to drift from)
- When the user explicitly says docs updates are not needed

## Usage Examples

```text
/update-docs
/update-docs (after adding --session flag to ctx agent)
```

## Workflow

1. **Diff** the branch: `git diff main --stat` (or relevant
   base)
2. **Verify mapping** is current (see Self-Maintenance below)
3. **Map** changed packages to affected docs (see table)
4. **Read** each affected doc; flag sections that contradict
   the new code
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
| `internal/tpl/claude/skills/`      | `.claude/skills/` (live versions)                  |
| `SECURITY.md`                      | `docs/security.md`                                 |
| `.context/` schema changes         | `docs/context-files.md`                            |

## What to Check

- **New CLI flags/commands**: are they in
  `docs/cli-reference.md`?
- **Changed file formats**: does `docs/context-files.md`
  match?
- **New context files**: added to both read order docs and
  `docs/context-files.md`?
- **Removed features**: still referenced in docs?
- **Changed defaults**: do examples in docs use old defaults?
- **Skill templates changed**: do live versions in
  `.claude/skills/` match `internal/tpl/claude/skills/`?

## Self-Maintenance

This mapping table will drift. Before relying on it:

1. `ls internal/`: any packages not in the table? Add them.
2. `ls docs/*.md`: any doc pages not in the table? Map them.
3. If you update the table, edit this skill file directly.

The skill is its own first test case: if the mapping is stale,
the skill has already failed at its job.

## Internal Code Conventions

Also check that changed code follows project patterns (not
Go defaults):

### Godoc Style

Project uses explicit **Parameters/Returns** sections, not
standard godoc.

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

// Bad (standard godoc; agent corpus drift):
// FunctionName does X with param1.
func FunctionName(param1 string) error
```

Verify that godoc comments match actual parameters and
behavior.

### Predicate Naming

Project uses predicates **without** Is/Has/Can prefixes:
- `Completed()` not `IsCompleted()`
- `Empty()` not `IsEmpty()`
- `Exists()` not `DoesExist()`

### File Organization

Public API in the main file, private helpers in **separate
logical files**:
- `loader.go` (public `Load()`) + `process.go` (private)
- NOT: everything in one file with unexported functions at
  the bottom

### Magic Strings

Literals belong in `internal/config/`. If you see a hardcoded
string used in 2+ files, it needs a constant. Check
`internal/config/` for existing constants before introducing
new literals.

## Relationship to ctx drift

`ctx drift` checks `.context/` file health (dead paths,
staleness). This skill checks docs-to-source-code alignment
and internal conventions. They are complementary.

## Quality Checklist

Before reporting results, verify:
- [ ] All changed packages were mapped to their affected docs
- [ ] Each flagged doc section was actually read (not assumed)
- [ ] Skill template/live drift was checked if `internal/tpl/`
      was touched
- [ ] Self-maintenance was done (mapping table is current)

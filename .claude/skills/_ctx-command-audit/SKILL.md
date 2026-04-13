---
name: _ctx-command-audit
description: "Audit CLI command surface after renames, moves, or deletions. Use after any namespace change to catch stale references."
---

Verify that a CLI command change (rename, move, deletion, new parent
grouping) is fully propagated across code, assets, docs, and tests.

## Before Running

1. **Build passes**: `go build ./...` must succeed
2. **Know what changed**: list the command changes (old name -> new
   name, or deleted, or moved under a parent)

## When to Use

- After renaming a CLI command (e.g. `ctx stats` -> `ctx usage`)
- After moving a command under a parent (e.g. `ctx pause` -> `ctx hook pause`)
- After deleting a command entirely (e.g. `ctx dep`)
- After promoting or demoting a command (top-level <-> system subcommand)

## When NOT to Use

- For flag-only changes (no command surface change)
- For internal refactors that don't change the user-facing command name
- When only adding a brand-new command (use `_ctx-qa` instead)

## Usage Examples

```text
/_ctx-command-audit
/_ctx-command-audit (after renaming ctx stats to ctx usage)
```

## Checks

For each renamed/moved/deleted command, run these checks in order.

### 1. Cobra --help examples

Verify the examples shown in `--help` output use the correct
command path.

```bash
grep -r "ctx <old-name>" internal/assets/commands/examples.yaml
grep -r "ctx <old-name>" internal/assets/commands/commands.yaml
```

Fix: Update YAML values to use the new command path. For
subcommands under a parent, examples must show the full path
(e.g. `ctx hook event`, not `ctx event`).

### 2. Use/DescKey constants

Verify the cobra Use string and DescKey constant match the new
command name.

```bash
grep -r "Use<OldName>\|DescKey<OldName>" internal/config/embed/cmd/
```

Fix: Rename constants, rename the file if needed (e.g.
`stats.go` -> `usage.go`), update all callers.

### 3. Flag DescKey constants

Verify flag description keys use the new command prefix.

```bash
grep -r "<old-name>\." internal/config/embed/flag/
```

Fix: Rename constants and their string values (e.g.
`stats.follow` -> `usage.follow`), rename the file, update
the corresponding `flags.yaml` keys.

### 4. Text DescKey constants

Verify text/write/error description keys use the new prefix.

```bash
grep -r "<old-name>\." internal/config/embed/text/
```

Fix: Rename constants and YAML keys in
`internal/assets/commands/text/*.yaml`.

### 5. group.go registration

Verify the command is registered in the correct group (or
removed if deleted).

```bash
grep "<old-package>" internal/bootstrap/group.go
```

Fix: Update import path and registration entry. Update the
function's doc comment to list the current commands.

### 6. Doc pages

Verify `docs/cli/<name>.md` exists with correct content, or
is removed for deleted commands.

```bash
ls docs/cli/<old-name>.md  # should not exist
ls docs/cli/<new-name>.md  # should exist
```

Fix: `git mv` for renames, `git rm` for deletions, create new
files for new parent commands.

### 7. CLI index

Verify `docs/cli/index.md` references the correct command names
and links.

```bash
grep "<old-name>" docs/cli/index.md
```

Fix: Update table entries, links, and anchors.

### 8. zensical.toml nav

Verify the site navigation references correct filenames.

```bash
grep "<old-name>" zensical.toml
```

Fix: Update nav entries to match renamed/new doc files.

### 9. Recipes and docs sweep

Broad sweep for any remaining stale references in user-facing
documentation.

```bash
grep -r "ctx <old-name>" docs/ --include="*.md"
grep -r "ctx <old-name>" internal/assets/claude/ --include="*.md"
```

Fix: Update all references. For commands moved under a parent,
every `ctx <name>` becomes `ctx <parent> <name>`.

### 10. Dead export cascade

After deleting a command, its support packages may become
orphaned (config constants, error constructors, write functions,
exec helpers, format constants, regex patterns).

```bash
make test  # TestNoDeadExports will catch these
```

Fix: Delete orphaned packages and constants. Do NOT add
allowlist exceptions — trace the dead export to its root cause
and remove it properly.

### 11. YAML orphan check

After deleting constants, their YAML counterparts may remain
as orphans.

```bash
make test  # TestDescKeyYAMLLinkage will catch these
```

Fix: Remove the orphaned YAML entries from `commands.yaml`,
`examples.yaml`, `flags.yaml`, and `text/*.yaml`.

### 12. Magic string check

If you inlined data that previously lived in a config package,
the magic strings test will catch it.

```bash
make test  # TestNoMagicStrings will catch these
```

Fix: Move the data to an appropriate `internal/config/`
package. Never inline string literals in source files.

### 13. AST test exceptions

Verify no allowlist entries were added to pass tests.

```bash
grep -r "allowlist\|whitelist\|skip\|except" internal/audit/ internal/compliance/ | grep -i "<command-name>"
```

Fix: If you added an exception, remove it and fix the
underlying issue instead.

### 14. Docstring review

Verify doc comments on touched packages are accurate, not
just minimally gate-passing.

Check these files for each renamed/moved command:
- `doc.go` — package description matches new command path
- `cmd.go` — `Cmd()` godoc names the correct command
- `system.go` or parent — subcommand list in comments is current
- `group.go` — function doc comments list current commands

## Output Format

Report as a checklist:

```
## Command Audit: <change description>

- [x] Cobra --help examples updated
- [x] Use/DescKey constants renamed
- [x] Flag DescKeys renamed
- [x] Text DescKeys renamed
- [x] group.go registration updated
- [x] Doc pages renamed/created/deleted
- [x] CLI index updated
- [x] zensical.toml nav updated
- [x] Recipes and docs swept
- [x] No dead exports
- [x] No orphan YAML keys
- [x] No magic strings
- [x] No AST test exceptions added
- [x] Docstrings reviewed

**Result**: PASS / FAIL (N issues)
```

## Relationship to Other Skills

| Skill                  | Scope                              |
|------------------------|------------------------------------|
| `_ctx-qa`              | General build/lint/test gate       |
| `_ctx-audit`           | Codebase convention sweep          |
| `_ctx-command-audit`   | CLI surface change completeness    |
| `_ctx-update-docs`     | Docs-code consistency after edits  |
| `_ctx-alignment-audit` | Agent playbook vs docs alignment   |

## Quality Checklist

Before reporting results, verify:
- [ ] Every check was run for every changed command
- [ ] No stale references remain in docs, recipes, or assets
- [ ] Tests pass with zero exceptions added
- [ ] Docstrings read as if written by someone who understands
      the domain, not as minimal gate-passers

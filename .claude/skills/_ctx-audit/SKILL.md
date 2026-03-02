---
name: _ctx-audit
description: "Detect and fix code-level drift. Use after YOLO sprints, before releases, or when the 3:1 consolidation ratio is due."
---

Run a code-level consolidation pass on the ctx codebase. This
complements `ctx drift` (which checks context-level drift) by
checking **source code** against project conventions.

## Before Consolidating

1. **Check the ratio**: have there been 3+ sessions since the
   last consolidation? If not, it may be too early
2. **Clean working tree**: `git status` should show no
   uncommitted changes; consolidation touches many files and
   you need a clean diff baseline
3. **Run tests first**: `make audit` should pass before you
   start; do not consolidate on top of a broken build

## When to Use

- After 3+ rapid sessions (the 3:1 ratio)
- Before tagging a release
- When a session touched many files
- When you suspect convention drift

## When NOT to Use

- Mid-feature when code is intentionally incomplete
- Immediately after the last consolidation with no new work
  in between
- When the user is focused on shipping and explicitly defers
  cleanup

## Usage Examples

```text
/_ctx-audit
/_ctx-audit (before v0.3.0 release)
/_ctx-audit (after the YOLO sprint this week)
```

## Checks

Before running checks mechanically, reason through which areas
are most likely to have drifted based on recent changes. This
focuses attention where it matters most.

Run each check in order. Report findings per check, then summarize.

### 1. Predicate Naming

Convention: no `Is`/`Has`/`Can` prefixes on exported bool-returning methods.

```bash
rg '^\s*func\s+\([^)]+\)\s+(Is|Has|Can)[A-Z]\w*\(' --type go -l
```

Accepted exceptions (do NOT flag these):
- `IsUser()`, `IsAssistant()` on `Message`: dropping `Is` makes
  these look like getters (`msg.User()` reads as "get user?").
  The prefix earns its keep.

**Fix**: Rename to drop the prefix. `IsPending()` → `Pending()`.
Flag any NEW `Is`/`Has`/`Can` methods not listed as exceptions above.

### 2. Magic Strings

Convention: literals used in 3+ files need a constant.

```bash
# Find repeated string literals across files
rg '"[A-Z][A-Z_]+\.md"' --type go -c | sort -t: -k2 -rn
rg '"\.context/' --type go -c | sort -t: -k2 -rn
```

Check `internal/config/` for existing constants. If a literal is already
defined there but not used everywhere, that's a drift.

**Fix**: Replace literal with the constant from `internal/config/`.

### 3. Hardcoded Permissions

Convention: file permissions should use named constants, not literals.

```bash
rg '0[67][0-7][0-5]' --type go -l
```

**Fix**: Define constants in `internal/config/` if missing, then reference them.

### 4. File Size

Convention: source files > 300 LOC should be evaluated for splitting.

```bash
find . -name '*.go' -not -name '*_test.go' -exec wc -l {} + | sort -rn | head -20
```

Files over 300 LOC: check if they mix public API with private helpers
(convention says split them).

### 5. TODO/FIXME in Source

Constitution: no TODO comments in main branch (move to TASKS.md).

```bash
rg 'TODO|FIXME|HACK|XXX' --type go -n
```

**Fix**: Move the item to `.context/TASKS.md`, delete the comment.

### 6. Path Construction

Constitution: path construction uses stdlib, no string concatenation.

```bash
rg '"\.\./|"/"|"/" \+|+ "/"' --type go -l
```

**Fix**: Replace with `filepath.Join()`.

### 7. Line Width

Highly encouraged: keep lines to ~80 characters. This is not a hard limit —
some lines (long string literals, struct tags, URLs) will exceed it and that's
fine. But drift happens quietly, especially in test code where long assertion
messages and deeply nested structs push lines wide without anyone noticing.

```bash
# Lines exceeding 100 chars (flag the worst offenders, not every 81-char line)
rg '.{101,}' --type go -c | sort -t: -k2 -rn | head -20
```

**Fix**: Break long lines at natural points: function arguments,
struct fields, chained calls. For test code, extract repeated
long values into local variables or constants.

### 8. Duplicate Code Blocks

Drift pattern: copy-paste blocks accumulate when the agent is focused on
getting the task done rather than keeping the code in shape. This is
especially common in test files but also appears in non-test code.

In test code, some duplication is acceptable; test readability matters.
But when the same setup/assertion block appears 3+ times, consider a test
helper (`testutil` or unexported helpers in `_test.go`).

In non-test code, apply the Consolidation Decision Matrix below.

```bash
# Heuristic: find functions with very similar signatures in the same package
# Manual review is more effective here; look for:
#   - Identical error-handling blocks
#   - Repeated struct construction
#   - Copy-paste command setup patterns
```

**Fix (tests)**: Extract a helper function in the same `_test.go` file.
Use `t.Helper()` so failure messages point to the caller.

**Fix (non-test)**: Extract shared logic into a package-level unexported
function, or into a shared internal package if it spans packages.

### 9. Architecture Diagram Drift

After structural changes (new packages, moved files, changed
dependencies), verify `.context/ARCHITECTURE.md` diagrams match
actual code:

```bash
# Compare packages listed in ARCHITECTURE.md to actual packages
ls internal/
# Compare dependency graph claims to actual imports
grep -r '"github.com/ActiveMemory/ctx/internal/' internal/ | \
  sed 's|.*ctx/internal/|internal/|' | sort -u
```

**Fix**: Update the component map table, dependency graph, and file
layout sections in `.context/ARCHITECTURE.md`. Run `ctx drift` to
verify no dead path references remain.

### 10. Dead Exports

Check for exported functions/types with no callers outside their package.

```bash
# Quick heuristic: exported func defined but only used in its own package
```

Use `go vet` and `golangci-lint run --enable=unused` for a more thorough check.

### 11. Package Documentation Drift

Convention: packages with a `doc.go` must stay accurate in two ways:

**a) File Organization listing** — must list every `.go` file in the
package (excluding `_test.go`). Missing or extra entries mean files
were added/removed without updating the doc.

```bash
make lint-docs
```

**b) Package description** — the opening paragraph describes what the
package does. When behavior changes (new subcommands, new
responsibilities, renamed concepts), the description drifts.

Review each `doc.go` manually: does the description still match what
the package actually does today? Check exported symbols, command
`Use`/`Short`/`Long` strings, and the file organization listing for
clues that the scope expanded or shifted.

**Fix (a)**: Add missing files, remove stale entries.
**Fix (b)**: Rewrite the description to match current behavior.

### 12. Dead Doc Links

Documentation links drift when pages are renamed, moved, or deleted.

Invoke the `/ctx-check-links` skill to scan all `docs/` markdown files for:

- **Internal links** pointing to files that don't exist
- **External links** that return errors (reported as warnings, not failures)
- **Image references** to missing files

Internal broken links count as findings to fix. External failures are
informational — network partitions happen.

### 13. ctxrc Schema Drift

Three sources of truth define `.ctxrc` options and must stay in sync:

- **`internal/rc/types.go`** — the `CtxRC` struct (yaml tags are the keys)
- **`.ctxrc`** — the sample config in the project root
- **`docs/configuration.md`** — the "Full Reference" yaml block and the
  "Option Reference" table

Drift patterns:
- A new field is added to the struct but not to the sample or docs
- A field is removed from the struct but lingers in the sample or docs
- A yaml tag is renamed but the sample/docs still use the old key
- Default values diverge between the struct defaults (`rc.go:Default()`)
  and the commented values in the sample / docs table

```bash
# Extract yaml tags from the struct
rg 'yaml:"(\w+)"' internal/rc/types.go -o --replace '$1'

# Extract commented keys from the sample .ctxrc (top-level scalars)
rg '^#?\s*(\w+):' .ctxrc -o --replace '$1' | grep -v '^#'

# Extract keys from the docs Full Reference block
rg '^#?\s*(\w+):' docs/configuration.md -o --replace '$1' | grep -v '^#'
```

**Fix**: Add missing keys to the sample and docs, or remove stale ones.
Ensure default values in the sample comments match `rc.Default()` and the
docs Option Reference table.

### 14. Makefile Target Drift

Skills and docs reference Makefile targets (e.g., `make audit`, `make lint-docs`,
`make smoke`). Targets get renamed or removed without updating references.

```bash
# Extract all .PHONY targets from the Makefile
rg '\.PHONY:' Makefile | tr ' ' '\n' | grep -v '\.PHONY' | sort -u

# Find Makefile target references in skills and docs
rg 'make \w+' .claude/skills/ docs/ --type md -o | sort -u
```

Compare the two lists. Any referenced target that doesn't exist in the
Makefile is a finding.

**Fix**: Update the reference to the current target name, or restore
the target if it was removed accidentally.

### 15. Skill Metadata Drift

Skill directories under `.claude/skills/` and skill entries in
`internal/config/file.go` (`DefaultClaudePermissions`) must stay aligned.
Skills get renamed (e.g., `ctx-borrow` → `absorb`) but stale references
linger in permission lists, docs, or other skills.

```bash
# List skill directories
ls .claude/skills/

# List ctx-plugin skill directories
ls .claude/ctx-skills/ 2>/dev/null

# Extract Skill() entries from DefaultClaudePermissions
rg 'Skill\(' internal/config/file.go -o | sort -u

# Find cross-references to skill names in other skills
rg '/\w[-\w]+' .claude/skills/ --type md -o | sort -u
```

Check:
- Every skill directory has a matching `Skill()` permission entry (or is
  intentionally unpermissioned)
- No `Skill()` entry references a skill that no longer exists
- Cross-references between skills use current names

**Fix**: Update stale names in permission lists and cross-references.

### 16. Config Constant Drift

`internal/config/` defines constants (`DirContext`, `FileReadOrder`,
`DefaultTokenBudget`, etc.) that docs and the sample `.ctxrc` cite as
defaults. When constants change, docs lag behind.

```bash
# Extract key constants and their values
rg 'const|var' internal/config/file.go internal/config/dir.go \
  internal/rc/types.go internal/rc/rc.go --type go -A1

# Check docs claim the same defaults
rg 'default|Default' docs/configuration.md -n
```

Cross-check:
- `config.DirContext` value matches the `context_dir` default in docs
  and sample `.ctxrc`
- `FileReadOrder` entries match the `priority_order` list in sample
  `.ctxrc` and the docs "Default priority order" section
- `DefaultTokenBudget`, `DefaultArchiveAfterDays`, etc. in `rc.go`
  match the commented values in `.ctxrc` and docs Option Reference table

**Fix**: Align the docs/sample values to match the code constants.

### 17. CLI Subcommand Drift

`docs/cli-reference.md` documents commands and flags. When commands are
added, renamed, or removed, the docs drift.

```bash
# Extract registered cobra commands (Use: field)
rg 'Use:\s+"(\w+)' internal/cli/ --type go -o --replace '$1' | sort -u

# Extract commands documented in cli-reference.md
rg '^## `ctx (\w+)' docs/cli-reference.md -o --replace '$1' | sort -u

# Compare
diff <(rg 'Use:\s+"(\w+)' internal/cli/ --type go -o --replace '$1' | sort -u) \
     <(rg '^## `ctx (\w+)' docs/cli-reference.md -o --replace '$1' | sort -u)
```

Also check that global flags documented in `docs/configuration.md`
match the flags actually registered in `internal/bootstrap/cmd.go`.

**Fix**: Add missing commands/flags to the docs, or remove stale entries.

## Consolidation Decision Matrix

Use this to prioritize what to fix:

| Similarity | Instances | Action |
|------------|-----------|--------|
| Exact duplicate | 2+ | Consolidate immediately |
| Same pattern, different args | 3+ | Extract with parameters |
| Similar shape | 5+ | Consider abstraction |
| < 3 instances | Any | Leave it; duplication is cheaper than wrong abstraction |

## Safe Migration Pattern

When consolidating would change public API:

1. Create new function alongside old
2. Deprecate old with `// Deprecated:` godoc comment
3. Migrate callers incrementally
4. Delete old function when no callers remain

Never bulk-rename in a single commit if callers span packages.

## Output Format

After running checks, report:

```
## Consolidation Report

### Findings
- [check name]: N issues (list files)
- [check name]: clean

### Priority
1. [highest impact finding]: [why]
2. [next]: [why]

### Suggested Fixes
- [file:line]: [what to change]
```

## Relationship to Other Skills

| Skill          | Scope                                     |
|----------------|-------------------------------------------|
| `/_ctx-qa`          | Build/test/lint; this checks conventions  |
| `/ctx-verify`       | Confirms claims; use after fixing findings|
| `/_ctx-update-docs` | Syncs docs with code; run after changes   |
| `ctx drift`         | Checks `.context/` files; this checks `.go` |
| `/ctx-check-links`  | Dead doc links; invoked as check #12      |

## Quality Checklist

Before reporting the consolidation results:
- [ ] All 17 checks were run (not skipped)
- [ ] Accepted exceptions were respected (e.g., `IsUser()`)
- [ ] Findings are prioritized (highest impact first)
- [ ] Each finding has a concrete fix suggestion with file path
- [ ] `make audit` still passes after fixes are applied

---
name: consolidate
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
/consolidate
/consolidate (before v0.3.0 release)
/consolidate (after the YOLO sprint this week)
```

## Checks

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

### 9. Dead Exports

Check for exported functions/types with no callers outside their package.

```bash
# Quick heuristic: exported func defined but only used in its own package
```

Use `go vet` and `golangci-lint run --enable=unused` for a more thorough check.

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
| `/qa`          | Build/test/lint; this checks conventions  |
| `/verify`      | Confirms claims; use after fixing findings|
| `/update-docs` | Syncs docs with code; run after changes   |
| `ctx drift`    | Checks `.context/` files; this checks `.go` |

## Quality Checklist

Before reporting the consolidation results:
- [ ] All 9 checks were run (not skipped)
- [ ] Accepted exceptions were respected (e.g., `IsUser()`)
- [ ] Findings are prioritized (highest impact first)
- [ ] Each finding has a concrete fix suggestion with file path
- [ ] `make audit` still passes after fixes are applied

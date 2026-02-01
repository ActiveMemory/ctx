---
description: "Audit documentation for consistency and accuracy"
---

Perform a documentation audit across the codebase. This is an AI-assisted 
review that produces a report for human judgment.

## Audit Scope

### 1. Go Docstring Consistency

Check that exported functions follow the canonical docstring format:

```go
// FunctionName does X.
//
// Extended description of behavior.
//
// Parameters:
//   - param1: Description of first parameter
//   - param2: Description of second parameter
//
// Returns:
//   - error: Non-nil if X fails or Y is invalid
func FunctionName(param1 Type1, param2 Type2) error {
```

**Check for:**
- Missing docstrings on exported functions
- Missing or malformed Parameters/Returns sections
- Inconsistent formatting (e.g., "Args:" vs "Parameters:")
- One-liner docstrings that should have more detail

**Files to scan:** `internal/**/*.go`, `cmd/**/*.go`

### 2. Public Docs vs CLI Reality

Verify `docs/*.md` accurately describes what the CLI actually does.

**Check for:**
- Commands mentioned in docs that don't exist
- Flags documented that don't exist (or vice versa)
- Examples that wouldn't actually work
- Default values that are incorrect

**Method:**
1. Extract all `ctx <command>` references from docs
2. Run `ctx <command> --help` and compare
3. Flag mismatches

### 3. Narrative Consistency

Review prose docs for internal consistency.

**Check for:**
- Terminology drift (e.g., "context file" vs "context document")
- Inconsistent explanations of the same concept
- Outdated references to removed features
- Tone inconsistencies

### 4. Code Pattern Drift

Check that code follows established conventions in CONVENTIONS.md.

**Check for:**
- CLI output methods: `fmt.Print*` instead of `cmd.Print*` in CLI code
- Other patterns documented in CONVENTIONS.md

**Method:**
1. Read CONVENTIONS.md to understand established patterns
2. Grep for violations (e.g., `fmt.Print` in `internal/cli/**/*.go`)
3. Flag files that violate conventions

**Files to scan:** `internal/cli/**/*.go` for CLI patterns

## Output Format

Produce a structured report:

```markdown
# Documentation Audit Report

## Summary
- X Go docstring issues found
- Y CLI/docs mismatches found
- Z narrative inconsistencies flagged
- W code pattern violations found

## Go Docstring Issues

### Missing Docstrings
- `internal/foo/bar.go:42` — `ExportedFunc` has no docstring

### Format Issues
- `internal/cli/watch/run.go:85` — Uses "Args:" instead of "Parameters:"

## CLI/Docs Mismatches

### Missing from Docs
- `ctx recall export` — Command exists but not documented

### Incorrect in Docs
- `docs/cli-reference.md:123` — Says `--verbose` but flag is `--debug`

## Narrative Issues

### Terminology
- "context files" used in index.md, "context documents" in getting-started.md

### Outdated References
- `docs/integrations.md:45` — References removed `--minimal` flag

## Code Pattern Violations

### CLI Output Methods
- `internal/cli/task/run.go:127` — Uses `fmt.Printf` instead of `cmd.Printf`
- `internal/cli/watch/run.go:45` — Uses `fmt.Println` instead of `cmd.Println`
```

## Execution

1. **Read CONVENTIONS.md** to understand established patterns
2. **Scan Go files** for docstring patterns
3. **Parse docs/*.md** for CLI references
4. **Run ctx --help** variants to get actual CLI surface
5. **Check code patterns** against conventions (e.g., CLI output methods)
6. **Compare and report**

Do NOT auto-fix anything. This audit produces a report for human review.

## Notes

- Focus on actionable issues, not style nitpicks
- Prioritize: missing docs > incorrect docs > inconsistent style
- Skip test files and generated code
- For large codebases, can focus on specific packages with $ARGUMENTS

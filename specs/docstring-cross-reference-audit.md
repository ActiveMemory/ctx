# Spec: Docstring Cross-Reference Audit

**Status**: Planned

## Problem

Agent-written docstrings frequently describe the wrong domain. A function
in `write/loop/pad.go` says "used during init to show which template files
were skipped" when it's actually called from `pad/cmd/export`. These
mismatches are invisible without manually tracing every caller.

The codebase has ~125 `write/**` output functions plus additional helpers
across `internal/`. Manual review doesn't scale.

## Solution

A skill (or standalone audit script) that cross-references function
docstrings against their actual callers to surface domain mismatches.

### Algorithm

For each exported function in target packages:

1. **Find callers** — grep for `pkg.FunctionName(` across the tree.
2. **Extract domain keywords** from the docstring — package names, command
   names, feature areas (init, pad, drift, compact, etc.).
3. **Extract caller domains** — the package path of each caller
   (e.g., `cli/pad/cmd/export` → "pad", "export").
4. **Flag mismatches** — docstring mentions domain X, but no caller lives
   in domain X.

### Output

A checklist of suspicious mismatches, one per function:

```
write/loop/pad.go:InfoPathConversionExists
  docstring mentions: init, template
  callers: cli/pad/cmd/export
  → domain mismatch: "init", "template" not in caller path
```

The skill reports findings but does not auto-fix. A human or focused
follow-up prompt rewrites the flagged docstrings.

## Scope

### Phase 1: `write/**` packages

The `write/` tree is the highest-risk area — output functions written by
agents that guessed at context rather than tracing callers.

### Phase 2: All `internal/` exported functions

Extend to the full tree once the approach is validated on `write/`.

## Design Decisions

- **Scope per run**: Process one subpackage per invocation, not the whole
  tree. Small batches produce better results than "audit everything."
- **Checklist output, not auto-fix**: The skill reports suspicious
  mismatches; a human does the actual rewrite.
- **Domain keywords are heuristic**: The list of recognized domains comes
  from the CLI command names and top-level package names. False positives
  are acceptable — the goal is to surface candidates, not prove correctness.

## Non-Goals

- Detecting *semantically* wrong descriptions that use correct domain
  words (e.g., "prints status" when it actually prints a warning). This
  requires understanding intent, not pattern matching.
- Auto-rewriting docstrings.

## Implementation

- Skill file: `.claude/skills/ctx-docstring-audit/SKILL.md`
- Alternatively: a Go test in `internal/compliance/` that runs as part of
  `make lint`, similar to the existing `TestDocGoExists` and
  `TestNoLiteralNewline` compliance checks.
- The Go test approach is preferred — it's deterministic, runs in CI, and
  doesn't depend on agent quality.

---
name: ctx-refactor
description: "Refactor code safely: test-first, one change at a time, preserve behavior. Use when the user says 'refactor this', 'clean this up', or wants structural improvement."
allowed-tools: Read, Grep, Glob, Edit, Write, Bash(go:*), Bash(make:*), Bash(git:*)
---

Refactor the specified code following strict safety rules.
Refactoring changes structure, not outcomes.

## When to Use

- User says "refactor this", "clean this up", "simplify this"
- User wants to extract, rename, split, or reorganize code
- User says "this is messy" or "can we improve this"

## When NOT to Use

- User wants to add new behavior (that's a feature, not a refactor)
- User wants a rename across the codebase (prefer
  `/gitnexus-refactoring` if GitNexus is available; otherwise
  use grep-based search to find all references before renaming)

## Rules

Follow these in order. Do not skip steps.

1. **Write or verify tests first**: confirm existing behavior is
   captured before changing structure.
2. **Preserve all existing behavior**: refactoring changes
   structure, not outcomes. If a step would change observable
   behavior, stop and flag it as a separate task.
3. **Make one structural change at a time**: keep each step
   reviewable and revertible.
4. **Run tests after each step**: catch regressions immediately,
   not at the end.
5. **Check project conventions**: consult `.context/CONVENTIONS.md`
   to ensure the refactored code follows established patterns.

## Execution

1. Read `.context/CONVENTIONS.md` to load project patterns
2. Read the target code and its tests
3. If no tests exist, write them first (confirm with user)
4. Plan the refactoring steps: present to user before starting
5. Execute one step at a time, running tests between each
6. After all steps, run `make lint && make test`

## Output Format

Before starting, present the plan:

```
## Refactoring Plan: <target>

1. <step>: why
2. <step>: why
...

Tests to verify: <list>
```

After each step, report: what changed, tests still passing.

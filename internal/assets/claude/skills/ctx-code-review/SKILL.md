---
name: ctx-code-review
description: "Review code changes for correctness, edge cases, and convention adherence. Use when the user asks to review code, a diff, a PR, or says 'review this'."
allowed-tools: Read, Grep, Glob, Bash(git:*)
---

Review the specified code change focusing on substance over style.

## When to Use

- User says "review this code", "review this change", "code review"
- User asks for feedback on a diff, PR, or set of changes
- User says "what do you think of this?"

## When NOT to Use

- User wants a full PR review with GitHub integration (use
  `/gitnexus-pr-review` instead)
- User wants an architecture-level review (use `/ctx-architecture`)

## Review Checklist

Work through each dimension. Flag issues, don't fix them unless asked.

1. **Correctness**: Does the logic do what it claims? Off-by-one
   errors, nil dereferences, race conditions?
2. **Edge cases**: What happens with empty input, max values,
   concurrent access, or partial failures?
3. **Naming clarity**: Do function, variable, and type names
   communicate intent without needing comments?
4. **Test coverage gaps**: What behavior is untested? What inputs
   would exercise uncovered paths?
5. **Convention adherence**: Does this follow the project patterns
   documented in `.context/CONVENTIONS.md`?

## Execution

1. Read `.context/CONVENTIONS.md` to load project patterns
2. Identify the scope: file(s), diff, or recent changes
3. If no specific target, check `git diff` for unstaged changes
4. Work through each checklist dimension
5. Present findings grouped by severity: bugs > logic gaps >
   conventions > style observations

## Output Format

Lead with the most important finding. Use this structure:

```
## Review: <scope>

### Issues
- **[severity]** file:line — description

### Observations
- Note anything non-obvious but not necessarily wrong

### Verdict
One sentence: ship it, fix N issues first, or needs rethink.
```

Flag but don't fix style issues. Focus review on substance over
formatting.

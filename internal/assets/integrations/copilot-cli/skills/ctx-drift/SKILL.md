---
name: ctx-drift
description: "Detect and fix context drift. Use to find stale paths, broken references, and constitution violations in context files."
tools: [bash, read, write, edit, glob, grep]
---

Detect context drift at two layers: **structural** (stale paths,
missing files, constitution violations) via `ctx drift`, and
**semantic** (outdated conventions, superseded decisions,
irrelevant learnings) via agent analysis.

## When to Use

- At session start to verify context health before working
- After refactors, renames, or major structural changes
- When the user asks "is our context clean?" or "check for drift"
- Before a release or milestone

## When NOT to Use

- When you just ran status and everything looked fine
- Repeatedly in the same session without changes
- Mid-flow when the user is focused on a task

## Execution

### Layer 1: Structural Checks

```bash
ctx drift
```

Catches dead paths, missing files, staleness indicators.

### Layer 2: Semantic Analysis

After structural check, read context files and compare to the
codebase. Check for:

- **Outdated conventions**: patterns the code no longer follows
- **Superseded decisions**: entries overridden by later work
- **Stale architecture**: module descriptions that have changed
- **Irrelevant learnings**: entries about fixed bugs
- **Contradictions**: context files contradicting each other

### Reporting

1. Summarize findings by severity
2. Explain each finding: what file, why it matters
3. Distinguish structural from semantic
4. Offer to auto-fix structural: `ctx drift --fix`
5. Propose specific edits for semantic issues

## Quality Checklist

- [ ] Summarized findings (did not dump raw output)
- [ ] Explained why each finding matters
- [ ] Offered auto-fix before running it
- [ ] Did not run `--fix` without user confirmation

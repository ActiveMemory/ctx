---
name: ctx-drift
description: "Detect context drift. Use to find stale paths, broken references, and outdated context."
---

Detect context drift: stale paths, missing files, constitution
violations, and semantic staleness in context files.

## When to Use

- At session start to verify context health
- After refactors, renames, or major structural changes
- When asked "is our context clean?" or "anything stale?"
- Before a release or milestone

## When NOT to Use

- When you just ran ctx status and everything looked fine
- Repeatedly without changes in between

## Execution

Run the structural drift check:

```bash
ctx drift
```

This catches dead paths, missing files, staleness indicators,
and constitution violations.

After running, also do a semantic check: read the context files
and compare them to what you know about the codebase. Look for:

- Outdated conventions that the code no longer follows
- Decisions whose rationale no longer applies
- Architecture descriptions that have changed
- Learnings about bugs that were since fixed

Report findings as actionable items, not raw output. Propose
specific fixes for each issue found.

---
name: ctx-add-learning
description: "Record a learning. Use when discovering gotchas, bugs, or unexpected behavior that future sessions should know about."
tools: [bash]
---

Record a learning in LEARNINGS.md.

## Before Recording

Three questions: if any answer is "no", don't record:

1. **"Could someone Google this in 5 minutes?"** → If yes, skip it
2. **"Is this specific to this codebase?"** → If no, skip it
3. **"Did it take real effort to discover?"** → If no, skip it

Learnings should capture **principles and heuristics**, not code snippets.

## When to Use

- After discovering a gotcha or unexpected behavior
- When a debugging session reveals root cause
- When finding a pattern that will help future work

## When NOT to Use

- General programming knowledge (not specific to this project)
- One-off workarounds that won't recur
- Things already documented in the codebase

## Execution

```bash
ctx learning add "Title" \
  --context "What were you doing when you discovered this?" \
  --lesson "What's the key insight?" \
  --application "How should we handle this going forward?"
```

## Quality Checklist

- [ ] Context explains what happened (not just what you learned)
- [ ] Lesson is a principle, not a code snippet
- [ ] Application gives actionable guidance for next time
- [ ] Not already in LEARNINGS.md (check first)

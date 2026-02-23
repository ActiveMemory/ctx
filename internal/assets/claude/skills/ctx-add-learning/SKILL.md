---
name: ctx-add-learning
description: "Record a learning. Use when discovering gotchas, bugs, or unexpected behavior that future sessions should know about."
allowed-tools: Bash(ctx:*)
---

Record a learning in LEARNINGS.md.

## Before Recording

Three questions — if any answer is "no", don't record:

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

## Gathering Information

If the user provides only a title, ask:

1. "What were you doing when you discovered this?" → Context
2. "What's the key insight?" → Lesson
3. "How should we handle this going forward?" → Application

## Execution

```bash
ctx add learning "Title" --context "..." --lesson "..." --application "..."
```

**Example — behavioral pattern:**
```bash
ctx add learning "Agent ignores repeated hook output (repetition fatigue)" \
  --context "PreToolUse hook ran ctx agent on every tool use, injecting the same context packet repeatedly. Agent tuned it out and didn't follow conventions." \
  --lesson "Repeated injection causes the agent to ignore the output. A cooldown tombstone emits once per window. A readback instruction creates a behavioral gate harder to skip than silent injection." \
  --application "Use --session \$PPID in hook commands to enable cooldown. Pair context injection with a readback instruction."
```

**Example — technical gotcha:**
```bash
ctx add learning "go:embed only works with files in same or child directories" \
  --context "Tried to embed files from parent directory, got compile error" \
  --lesson "go:embed paths are relative to the source file and cannot use .. to escape the package" \
  --application "Keep embedded files in internal/assets/ or child directories, not project root"
```

**Example — workflow insight:**
```bash
ctx add learning "ctx init overwrites user content without guard" \
  --context "Commit a9df9dd wiped 18 decisions from DECISIONS.md, replacing with empty template" \
  --lesson "Init treats all context files as templates, but after first use they contain user data" \
  --application "Skip existing files by default, only overwrite with --force"
```

## Quality Checklist

Before recording, verify:
- [ ] Context explains what happened (not just what you learned)
- [ ] Lesson is a principle, not a code snippet
- [ ] Application gives actionable guidance for next time
- [ ] Not already in LEARNINGS.md (check first)

Confirm the learning was added.

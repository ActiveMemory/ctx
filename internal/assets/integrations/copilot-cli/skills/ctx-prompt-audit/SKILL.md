---
name: ctx-prompt-audit
description: "Analyze session transcripts to identify vague prompts that caused unnecessary clarification cycles."
tools: [bash, read, write]
---

Analyze session history to find prompts that led to wasted work
due to ambiguity, and suggest improvements.

## When to Use

- After a session with many back-and-forth clarifications
- When improving prompt discipline
- During periodic workflow reviews

## When NOT to Use

- No session history exists
- Sessions were straightforward

## Process

### 1. Load recent sessions

```bash
ctx recall list --limit 5
```

### 2. Scan for patterns

Look for:
- Multiple clarifying questions before work began
- Misunderstood instructions leading to rework
- Vague requests like "fix it" or "make it better"
- Missing context that was discovered mid-task

### 3. Categorize findings

| Pattern | Example | Improvement |
|---------|---------|-------------|
| Vague scope | "Fix the tests" | "Fix TestFoo in internal/cli — it's failing on empty input" |
| Missing context | "Add a feature" | "Add JSON output to ctx status (see spec in specs/)" |
| Ambiguous reference | "Update that file" | "Update internal/config/mcp/tool/tool.go" |

### 4. Present recommendations

Provide actionable suggestions for clearer prompts.

## Quality Checklist

- [ ] At least 3 sessions analyzed
- [ ] Patterns categorized with examples
- [ ] Concrete improvements suggested
- [ ] No session data exposed inappropriately

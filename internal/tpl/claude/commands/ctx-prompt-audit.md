---
description: "Analyze session logs to identify vague prompts and suggest improvements"
---

Analyze recent session transcripts to identify prompts that led to unnecessary 
clarification back-and-forth. This helps the user improve their prompting 
patterns.

## Your Task

1. **Read recent sessions** from `.context/sessions/` (focus on the 3-5 most recent)
2. **Identify vague prompts** - user messages that caused you to ask clarifying questions
3. **Generate a coaching report** with concrete examples and suggestions

## What Makes a Prompt "Vague"

Look for user prompts where Claude's immediate response was to ask clarifying 
questions rather than take action. Signs include:

- **Missing file context**: "fix the bug" without specifying which file or error
- **Ambiguous scope**: "optimize it" without what to optimize or success criteria
- **Undefined targets**: "update the component" when multiple components exist
- **Missing error details**: "it's not working" without symptoms or expected behavior
- **Vague action words**: "make it better", "clean this up", "improve the code"

## Important Nuance

Not every short prompt is vague! Consider context:
- "fix the bug" after discussing a specific error in detail → **NOT vague**
- "fix the bug" as the first message → **VAGUE**
- "optimize it" when working on a single function → probably fine
- "optimize it" in a large codebase with no context → **VAGUE**

## Output Format

Generate a report like this:

```
## Prompt Audit Report

**Sessions analyzed**: 5
**User prompts reviewed**: 47
**Vague prompts found**: 4 (8.5%)

---

### Example 1: Missing File Context

**Your prompt**: "fix the bug"

**What happened**: I had to ask which file and what error you were seeing, a
dding 2 messages of back-and-forth.

**Better prompt**: "fix the authentication error in src/auth/login.ts where 
JWT validation fails with 401"

**Cost**: ~2 extra messages, ~30 seconds

---

### Example 2: Undefined Target

**Your prompt**: "optimize the component"

**What happened**: Multiple components exist. I asked which one and what 
performance issue to address.

**Better prompt**: "optimize UserList in src/components/UserList.tsx to reduce 
re-renders when parent state updates"

**Cost**: ~3 extra messages, ~1 minute

---

## Patterns to Watch

Based on your sessions, you tend to:
1. Skip mentioning file paths (3 occurrences)
2. Use "it" without establishing what "it" refers to (2 occurrences)

## Tips

- Start prompts with the **file path** when discussing specific code
- Include **error messages** when debugging
- Specify **success criteria** for optimization tasks
```

## Guidelines

- Be constructive, not critical - the goal is to help, not shame
- Show the actual prompt from their session (quoted)
- Explain what happened (what you had to ask)
- Provide a concrete better alternative
- Estimate the "cost" in extra messages/time
- Look for patterns across multiple examples
- End with actionable tips based on their specific tendencies

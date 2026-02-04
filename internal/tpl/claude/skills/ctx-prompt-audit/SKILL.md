---
name: ctx-prompt-audit
description: "Audit prompting patterns. Use periodically to help users improve prompt quality and reduce clarification cycles."
---

Analyze recent session transcripts to identify prompts that led to unnecessary
clarification back-and-forth.

## When to Use

- Periodically to help users improve their prompting
- When the user asks for feedback on their prompting style
- After noticing many clarification cycles in recent sessions

## Process

1. **Read recent sessions** from `.context/sessions/` (focus on 3-5 most recent)
2. **Identify vague prompts** - user messages that caused clarifying questions
3. **Generate a coaching report** with concrete examples and suggestions

## What Makes a Prompt "Vague"

Look for prompts where Claude asked clarifying questions instead of acting:

- **Missing file context**: "fix the bug" without specifying which file or error
- **Ambiguous scope**: "optimize it" without what to optimize or success criteria
- **Undefined targets**: "update the component" when multiple components exist
- **Missing error details**: "it's not working" without symptoms
- **Vague action words**: "make it better", "clean this up"

## Important Nuance

Not every short prompt is vague! Consider context:
- "fix the bug" after discussing a specific error → **NOT vague**
- "fix the bug" as the first message → **VAGUE**

## Output Format

```markdown
## Prompt Audit Report

**Sessions analyzed**: 5
**User prompts reviewed**: 47
**Vague prompts found**: 4 (8.5%)

---

### Example 1: Missing File Context

**Your prompt**: "fix the bug"

**What happened**: I had to ask which file and what error.

**Better prompt**: "fix the authentication error in src/auth/login.ts where
JWT validation fails with 401"

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

- Be constructive, not critical
- Show actual prompts from their session (quoted)
- Explain what happened (what you had to ask)
- Provide concrete better alternatives
- Look for patterns across examples
- End with actionable tips

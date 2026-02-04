---
name: ctx-reflect
description: "Reflect on session progress. Use at natural breakpoints, after unexpected behavior, or when shifting to a different task."
---

Pause and reflect on this session. Review what has been accomplished and
identify context worth persisting.

## Reflection Checklist

Work through each category:

### 1. Learnings
- Did we discover any gotchas, bugs, or unexpected behavior?
- Did we learn something about the codebase, tools, or patterns?
- Would this help a future session avoid problems?

### 2. Decisions
- Did we make any architectural or design choices?
- Did we choose between alternatives?
- Should the rationale be captured?

### 3. Tasks
- Did we complete any tasks? (Mark `[x]` with `#done:` timestamp)
- Did we start any tasks? (Add `#started:` timestamp)
- Should new tasks be added for follow-up work?

### 4. Session Notes
- Was this a significant session worth a full summary?
- Would a future session benefit from the discussion context?

## Output Format

After reflecting, provide:

1. **Summary**: What was accomplished (2-3 sentences)
2. **Suggested Persists**: List what should be saved and where
3. **Offer**: Ask if the user wants you to persist any of these

Example:
> "This session fixed the auth bug and we discovered the token refresh gotcha.
> I'd suggest:
> - Learning: Token refresh requires explicit cache invalidation
> - Task: Mark 'Fix auth bug' as done
> Want me to persist these?"

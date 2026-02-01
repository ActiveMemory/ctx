---
description: "Add a decision to DECISIONS.md (requires context, rationale, consequences)"
argument-hint: "\"decision title\""
---

When the user runs /ctx-add-decision, you need to gather the complete ADR 
(Architecture Decision Record) format.

If $ARGUMENTS contains only a title (no flags), ask the user for:
1. **Context**: What prompted this decision?
2. **Rationale**: Why this choice over alternatives?
3. **Consequences**: What changes as a result?

Then run the command with all required flags:

```!
ctx add decision "$ARGUMENTS" --context "..." --rationale "..." --consequences "..."
```

If the user already provided flags in $ARGUMENTS, run the command directly:

```!
ctx add decision $ARGUMENTS
```

Confirm the decision was added and show a brief summary.

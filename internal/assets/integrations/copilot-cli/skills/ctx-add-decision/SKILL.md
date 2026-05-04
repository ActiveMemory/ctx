---
name: ctx-add-decision
description: "Record architectural decision. Use when a trade-off is resolved or a non-obvious design choice is made that future sessions need to know."
tools: [bash]
---

Record an architectural decision in DECISIONS.md.

## When to Use

- After resolving a trade-off between alternatives
- When making a non-obvious design choice
- When the "why" behind a choice needs to be preserved

## When NOT to Use

- Minor implementation details (use code comments instead)
- Routine maintenance or bug fixes
- When there was no real alternative to consider

## Decision Formats

### Quick Format (Y-Statement)

> "In the context of **[situation]**, facing **[constraint]**, we decided
> for **[choice]** and against **[alternatives]**, to achieve
> **[benefit]**, accepting that **[trade-off]**."

### Full Format

Gather: Context, Alternatives, Decision, Rationale, Consequence.

## Execution

```bash
ctx decision add "Use Cobra for CLI framework" \
  --context "Need CLI framework for Go project" \
  --rationale "Better subcommand support, team familiarity" \
  --consequence "More boilerplate, but clearer command structure"
```

## Quality Checklist

- [ ] Context explains the problem clearly
- [ ] At least one alternative was considered
- [ ] Rationale addresses why alternatives were rejected
- [ ] Consequence includes both benefits and trade-offs

---
name: ctx-add-convention
description: "Record a coding convention. Use when a repeated pattern should be codified so all sessions follow it consistently."
tools: [bash]
---

Record a coding convention in CONVENTIONS.md.

## When to Use

- When a pattern has been used 2-3 times and should be standardized
- When establishing a naming, formatting, or structural rule
- When a new contributor would need to know "how we do things here"
- When the user says "codify that" or "make that a convention"

## When NOT to Use

- One-off implementation details (use code comments instead)
- Architectural decisions with trade-offs (use `ctx-add-decision`)
- Debugging insights or gotchas (use `ctx-add-learning`)
- Rules that are already enforced by linters or formatters

## Gathering Information

Conventions are simpler than decisions or learnings. You need:

1. **Name**: What is the convention called?
2. **Rule**: What is the rule? One clear sentence.
3. **Section**: Where does it belong in CONVENTIONS.md?

If the user provides only a description, infer the section from the
topic. Check existing sections in CONVENTIONS.md first to place it
correctly: don't create a new section if an existing one fits.

## Execution

```bash
ctx convention add "Use kebab-case for all CLI flag names" --section "Naming"
```

## Quality Checklist

- [ ] The rule is clear enough that someone unfamiliar could follow it
- [ ] It is specific to this project (not a general rule)
- [ ] It is not already in CONVENTIONS.md (check first)
- [ ] The section matches an existing section, or a new one is needed
- [ ] It describes a pattern, not a one-time choice (that's a decision)

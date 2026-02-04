---
name: ctx-add-decision
description: "Record architectural decision. Use when a trade-off is resolved or a non-obvious design choice is made that future sessions need to know."
allowed-tools: Bash(ctx:*)
---

Record an architectural decision in DECISIONS.md.

## When to Use

- After resolving a trade-off between alternatives
- When making a non-obvious design choice
- When the "why" behind a choice needs to be preserved
- When future sessions need to understand why something is the way it is

## When NOT to Use

- Minor implementation details (use code comments instead)
- Routine maintenance or bug fixes
- Configuration changes that don't affect architecture
- When there was no real alternative to consider

## Decision Formats

### Quick Format (Y-Statement)

For lightweight decisions, use a single statement:

> "In the context of **[situation]**, facing **[constraint]**, we decided for **[choice]** and against **[alternatives]**, to achieve **[benefit]**, accepting that **[trade-off]**."

Example:
> "In the context of needing a CLI framework, facing Go ecosystem options, we decided for Cobra and against urfave/cli, to achieve better subcommand support, accepting that it has more boilerplate."

### Full Format

For significant decisions, gather:

1. **Context**: What situation prompted this decision? What constraints exist?
2. **Alternatives**: What options were considered? (At least 2)
3. **Decision**: What was chosen?
4. **Rationale**: Why this choice over the alternatives?
5. **Consequences**: What changes as a result? (Both positive and negative)

## Gathering Information

If the user provides only a title, ask:

1. "What prompted this decision?" → Context
2. "What alternatives did you consider?" → Options
3. "Why this choice over the alternatives?" → Rationale
4. "What are the consequences (good and bad)?" → Consequences

For quick decisions, offer the Y-statement format instead.

## Cross-Referencing

When a decision **supersedes** an earlier one:
- Mark the old decision as "Superseded by [new decision]"
- Reference the old decision in the new one
- Capture lessons learned from the original decision

When decisions are **related**:
- Note "See also: [related decision]" in consequences

## Execution

**Quick format:**
```bash
ctx add decision "Use Cobra for CLI framework" \
  --context "Need CLI framework for Go project" \
  --rationale "Better subcommand support than urfave/cli, team familiarity" \
  --consequences "More boilerplate, but clearer command structure"
```

**Full format with alternatives:**
```bash
ctx add decision "Use PostgreSQL for primary database" \
  --context "Need ACID-compliant database for e-commerce transactions" \
  --rationale "PostgreSQL offers JSONB, full-text search, and team has experience. Chose over MySQL (weaker JSON) and MongoDB (no multi-doc ACID)." \
  --consequences "Single database handles transactions and search. Team needs PostgreSQL-specific training."
```

## Quality Checklist

Before recording, verify:
- [ ] Context explains the problem clearly
- [ ] At least one alternative was considered
- [ ] Rationale addresses why alternatives were rejected
- [ ] Consequences include both benefits and trade-offs

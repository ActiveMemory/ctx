# Decisions

<!-- INDEX:START -->
| Date | Decision |
|------|--------|
| 2026-02-04 | E/A/R classification as the standard for skill evaluation |
<!-- INDEX:END -->

<!-- DECISION FORMATS

## Quick Format (Y-Statement)

For lightweight decisions, a single statement suffices:

> "In the context of [situation], facing [constraint], we decided for [choice]
> and against [alternatives], to achieve [benefit], accepting that [trade-off]."

## Full Format

For significant decisions:

## [2026-02-04-230933] E/A/R classification as the standard for skill evaluation

**Status**: Accepted

**Context**: Reviewed ~30 external skill/prompt files; needed a systematic way to evaluate what to keep vs delete

**Decision**: E/A/R classification as the standard for skill evaluation

**Rationale**: Expert/Activation/Redundant taxonomy from judge.txt captures the key insight: Good Skill = Expert Knowledge - What Claude Already Knows. Gives a concrete target (>70% Expert, <10% Redundant)

**Consequences**: skill-creator SKILL.md updated with E/A/R as core principle. All future skills evaluated against this framework

---

## [YYYY-MM-DD] Decision Title

**Status**: Accepted | Superseded | Deprecated

**Context**: What situation prompted this decision? What constraints exist?

**Alternatives Considered**:
- Option A: [Pros] / [Cons]
- Option B: [Pros] / [Cons]

**Decision**: What was decided?

**Rationale**: Why this choice over the alternatives?

**Consequences**: What are the implications? (Include both positive and negative)

**Related**: See also [other decision] | Supersedes [old decision]

## When to Record a Decision

✓ Trade-offs between alternatives
✓ Non-obvious design choices
✓ Choices that affect architecture
✓ "Why" that needs preservation

✗ Minor implementation details
✗ Routine maintenance
✗ Configuration changes
✗ No real alternatives existed

-->

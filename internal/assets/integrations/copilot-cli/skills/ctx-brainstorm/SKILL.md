---
name: ctx-brainstorm
description: "Design before implementation. Use before any creative or constructive work to transform vague ideas into validated designs."
tools: [bash, read, write]
---

Transform raw ideas into **clear, validated designs** through
structured dialogue **before any implementation begins**.

## Before Brainstorming

1. **Check if design is needed**: is the change complex enough?
2. **Review prior art**: check `.context/DECISIONS.md` for related
   past decisions
3. **Identify what exists**: read relevant code before asking
   questions the codebase already answers

## When to Use

- Before implementing a new feature
- Before architectural changes
- Before significant behavior modifications
- When an idea is vague and needs shaping

## When NOT to Use

- Bug fixes with clear solutions
- Routine maintenance tasks
- Well-defined requirements
- Small, isolated changes
- When the user explicitly wants to jump to code

## Operating Mode

Design facilitator, not builder. No implementation while
brainstorming.

## The Process

### 1. Understand Current Context

Before asking questions:
- Review project state, docs, prior decisions
- Identify what exists vs what is proposed
- Note implicit constraints

### 2. Clarify the Idea

Rules:
- Ask **one question per message**
- Prefer **multiple-choice** when possible

Focus on:
- Purpose: why does this need to exist?
- Users: who benefits?
- Constraints: what limits apply?
- Success criteria: how do we know it works?
- Non-goals: what is explicitly out of scope?

### 3. Non-Functional Requirements

Clarify or propose assumptions for:
- Performance, scale, security, reliability, maintenance

### 4. Understanding Lock (Gate)

Before proposing any design, provide:

**Understanding Summary** (5-7 bullets):
- What is being built, why, for whom, constraints, non-goals

**Assumptions**: list all explicitly.

**Open Questions**: list unresolved items.

> "Does this accurately reflect your intent? Confirm or correct
> before we move to design."

**Do NOT proceed until confirmed.**

### 5. Explore Design Approaches

- Propose **2-3 viable approaches**
- Lead with recommended option
- Explain trade-offs

### 6. Stress-Test the Chosen Approach

After the user picks an approach:
- Surface assumptions and dependencies
- Identify failure modes
- Steel-man an alternative

> "These are the risks I see. Do they change your preference?"

### 7. Present the Design

Break into digestible sections. Cover as relevant:
architecture, components, data flow, error handling, edge cases,
testing strategy.

### 8. Decision Log

Maintain a running log:

| Decision | Alternatives | Rationale |
|----------|--------------|-----------|
| ...      | ...          | ...       |

## After the Design

### Persist to Context

```bash
ctx add decision "..." --context "..." --rationale "..."
```

### Implementation Handoff

Only after documentation, ask:
> "Ready to begin implementation?"

## Quality Checklist

Exit brainstorming **only when**:
- [ ] Understanding Lock confirmed by the user
- [ ] At least one design approach accepted
- [ ] Stress-test completed
- [ ] Major assumptions documented
- [ ] Key risks acknowledged
- [ ] Decision Log complete
- [ ] Decisions persisted to context

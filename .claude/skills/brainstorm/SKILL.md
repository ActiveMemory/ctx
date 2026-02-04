---
name: brainstorm
description: "Design before implementation. Use before any creative or constructive work (features, architecture, behavior changes) to transform vague ideas into validated designs."
---

# Brainstorm Ideas Into Designs

Transform raw ideas into **clear, validated designs** through structured dialogue
**before any implementation begins**.

## When to Use

- Before implementing a new feature
- Before architectural changes
- Before significant behavior modifications
- When an idea is vague and needs shaping

## When NOT to Use

- Bug fixes with clear solutions
- Routine maintenance
- When requirements are already well-defined
- Small, isolated changes

## Operating Mode

You are a **design facilitator**, not a builder.

- No implementation while brainstorming
- No speculative features
- No silent assumptions
- No skipping ahead

Your job is to **slow down just enough to get it right**.

## The Process

### 1. Understand Current Context

Before asking questions:

- Review project state: files, docs, prior decisions
- Check `.context/DECISIONS.md` for related past decisions
- Identify what exists vs. what is proposed
- Note implicit constraints

**Do not design yet.**

### 2. Clarify the Idea

Goal: **shared clarity**, not speed.

Rules:
- Ask **one question per message**
- Prefer **multiple-choice** when possible
- Split complex topics into multiple questions

Focus on:
- Purpose — why does this need to exist?
- Users — who benefits?
- Constraints — what limits apply?
- Success criteria — how do we know it works?
- Non-goals — what is explicitly out of scope?

### 3. Non-Functional Requirements

Explicitly clarify or propose assumptions for:

- Performance expectations
- Scale (users, data, traffic)
- Security/privacy constraints
- Reliability needs
- Maintenance expectations

If the user is unsure, propose reasonable defaults and mark them as **assumptions**.

### 4. Understanding Lock (Gate)

Before proposing any design, pause and provide:

**Understanding Summary** (5-7 bullets):
- What is being built
- Why it exists
- Who it's for
- Key constraints
- Explicit non-goals

**Assumptions**: List all explicitly.

**Open Questions**: List unresolved items.

Then ask:
> "Does this accurately reflect your intent? Please confirm or correct before we move to design."

**Do NOT proceed until confirmed.**

### 5. Explore Design Approaches

Once understanding is confirmed:

- Propose **2-3 viable approaches**
- Lead with your **recommended option**
- Explain trade-offs: complexity, extensibility, risk, maintenance
- Apply YAGNI ruthlessly

### 6. Present the Design

Break into digestible sections. After each, ask:
> "Does this look right so far?"

Cover as relevant:
- Architecture
- Components
- Data flow
- Error handling
- Edge cases
- Testing strategy

### 7. Decision Log

Maintain a running log throughout:

| Decision | Alternatives | Rationale |
|----------|--------------|-----------|
| ... | ... | ... |

## After the Design

### Persist to Context

Once validated, persist outputs:

```bash
# Record key decisions
ctx add decision "..." --context "..." --rationale "..." --consequences "..."

# Update implementation plan
# Write to IMPLEMENTATION_PLAN.md
```

Save the full design document to `.context/sessions/` or a dedicated design doc.

### Implementation Handoff

Only after documentation, ask:
> "Ready to begin implementation?"

If yes:
- Create explicit implementation plan
- Break into incremental steps
- Proceed one step at a time

## Exit Criteria

Exit brainstorming mode **only when**:

- [ ] Understanding Lock confirmed
- [ ] At least one design approach accepted
- [ ] Major assumptions documented
- [ ] Key risks acknowledged
- [ ] Decision Log complete

If any criterion is unmet, continue refinement.

## Principles

- One question at a time
- Assumptions must be explicit
- Explore alternatives
- Validate incrementally
- Clarity over cleverness
- Be willing to go back
- **YAGNI ruthlessly**

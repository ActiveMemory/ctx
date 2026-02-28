---
name: _ctx-brainstorm
description: "Design before implementation. Use before any creative or constructive work (features, architecture, behavior changes) to transform vague ideas into validated designs."
---

Transform raw ideas into **clear, validated designs** through
structured dialogue **before any implementation begins**.

## Before Brainstorming

1. **Check if design is needed**: is the change complex enough
   to warrant a design phase, or is the solution already clear?
2. **Review prior art**: check `.context/DECISIONS.md` for
   related past decisions; do not re-litigate settled choices
3. **Identify what exists**: read relevant code and docs before
   asking questions; do not ask the user things the codebase
   already answers

## When to Use

- Before implementing a new feature
- Before architectural changes
- Before significant behavior modifications
- When an idea is vague and needs shaping

## When NOT to Use

- Bug fixes with clear solutions
- Routine maintenance tasks
- When requirements are already well-defined
- Small, isolated changes (just do them)
- When the user explicitly wants to jump straight to code

## Usage Examples

```text
/_ctx-brainstorm
/_ctx-brainstorm (new caching layer for the API)
/_ctx-brainstorm (should we split the monolith?)
```

## Operating Mode

Design facilitator, not builder.

- No implementation while brainstorming
- No speculative features
- No silent assumptions
- No skipping ahead

**Slow down just enough to get it right.**

## The Process

### 1. Understand Current Context

Before asking questions:

- Review project state: files, docs, prior decisions
- Check `.context/DECISIONS.md` for related past decisions
- Identify what exists vs what is proposed
- Note implicit constraints

**Do not design yet.**

### 2. Clarify the Idea

Goal: **shared clarity**, not speed.

Rules:
- Ask **one question per message**
- Prefer **multiple-choice** when possible
- Split complex topics into multiple questions

Focus on:
- Purpose: why does this need to exist?
- Users: who benefits?
- Constraints: what limits apply?
- Success criteria: how do we know it works?
- Non-goals: what is explicitly out of scope?

### 3. Non-Functional Requirements

Explicitly clarify or propose assumptions for:

- Performance expectations
- Scale (users, data, traffic)
- Security/privacy constraints
- Reliability needs
- Maintenance expectations

If the user is unsure, propose reasonable defaults and mark
them as **assumptions**.

### 4. Understanding Lock (Gate)

Before proposing any design, pause and provide:

**Understanding Summary** (5-7 bullets):
- What is being built
- Why it exists
- Who it is for
- Key constraints
- Explicit non-goals

**Assumptions**: list all explicitly.

**Open Questions**: list unresolved items.

Then ask:
> "Does this accurately reflect your intent? Confirm or
> correct before we move to design."

**Do NOT proceed until confirmed.**

### 5. Explore Design Approaches

Once understanding is confirmed:

- Propose **2-3 viable approaches**
- Lead with your **recommended option**
- Explain trade-offs: complexity, extensibility, risk,
  maintenance
- Apply YAGNI ruthlessly

### 6. Stress-Test the Chosen Approach

After the user picks an approach, pause for adversarial review
before moving to detailed design.

**Surface assumptions**:
- List assumptions the chosen approach depends on
- Identify implicit dependencies (libraries, infra, team knowledge)

**Identify failure modes**:
- What would make this approach fail? (edge cases, scale limits,
  integration risks, operational complexity)
- What's the worst-case recovery if it does fail?

**Steel-man an alternative**:
- Name the strongest argument for a different approach
- Be specific — "Approach B avoids X risk" not "there are other options"

**Gate question**:
> "These are the risks I see. Do they change your preference, or
> should we proceed with the chosen approach?"

Only move to detailed design after the user confirms. This is
one gate, not a debate — if the user says "proceed," proceed.

### 7. Present the Design

Break into digestible sections. After each, ask:
> "Does this look right so far?"

Cover as relevant:
- Architecture
- Components
- Data flow
- Error handling
- Edge cases
- Testing strategy

### 8. Decision Log

Maintain a running log throughout:

| Decision | Alternatives | Rationale |
|----------|--------------|-----------|
| ...      | ...          | ...       |

## After the Design

### Persist to Context

Once validated, persist outputs:

```bash
# Record key decisions
ctx add decision "..." --context "..." --rationale "..."
```

### Implementation Handoff

Only after documentation, ask:
> "Ready to begin implementation?"

If yes:
- Create explicit implementation plan
- Break into incremental steps
- Proceed one step at a time

## Good Example

> **Understanding Summary**:
> - Building a cooldown mechanism for `ctx agent` hooks
> - Prevents repetitive context injection on every tool use
> - For Claude Code users running ctx in PreToolUse hooks
> - Must be session-isolated (two sessions share no state)
> - Non-goal: per-tool granularity (cooldown is global)
>
> **Assumptions**: 10-minute default cooldown is reasonable.
>
> **Open Questions**: none remaining.
>
> Does this accurately reflect your intent?

## Bad Examples

- Jumping to architecture diagrams before asking what the
  feature is for
- Asking 5 questions in one message (ask one at a time)
- Proposing a design without the Understanding Lock step
- "Let me implement this real quick" (no implementation
  during brainstorm)

## Quality Checklist

Exit brainstorming mode **only when**:

- [ ] Understanding Lock confirmed by the user
- [ ] At least one design approach accepted
- [ ] Stress-test completed (assumptions, failure modes, alternatives)
- [ ] Major assumptions documented explicitly
- [ ] Key risks acknowledged
- [ ] Decision Log complete
- [ ] Decisions persisted to `.context/DECISIONS.md`

If any criterion is unmet, continue refinement.

## Principles

- **Think step-by-step** before proposing anything — reason
  through the problem space before jumping to solutions
- One question at a time
- Assumptions must be explicit
- Explore alternatives before committing
- Validate incrementally
- Clarity over cleverness
- Be willing to go back
- **YAGNI ruthlessly**

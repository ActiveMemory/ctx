---
title: Scrutinizing a Plan
---

# Scrutinizing a Plan

When you have a plan and want it attacked, not validated, the
`/ctx-plan` skill runs an adversarial interview. It surfaces what's
weak, missing, or unexamined *before* you commit.

## When to Use It

- Before starting a multi-day implementation.
- After writing a spec but before opening the first PR.
- When the team aligned suspiciously fast on a complex change.
- When you've drafted something and the obvious next step is to
  hand it to a senior reviewer.

`/ctx-plan` is not a brainstorming tool. Use
[`/ctx-brainstorm`](design-before-coding.md) for that. `/ctx-plan`
assumes the plan exists and asks: **does it survive contact with
reality?**

## Stop Conditions

The skill runs depth-first through assumptions, tradeoffs, and
failure modes. Stop when one of these is true:

- **Plan converges.** The bet is clear, the sequencing is justified,
  and you can articulate what would falsify it.
- **Plan dissolves.** The interview surfaces a load-bearing
  assumption you can't defend; the plan needs rework, not more
  scrutiny.
- **You hit information you don't have.** Pause, gather it,
  resume.

If the agent loops on a single angle (asking variations of the same
question), say so. The skill is designed to walk depth-first, not
exhaust one branch.

## Worked Example

```text
You: /ctx-plan

I want to migrate our queue from RabbitMQ to NATS. I've drafted a
two-week plan: spike NATS in staging, port one consumer, port the
rest, decommission RabbitMQ.

Agent: One question at a time. What breaks if your message ordering
       guarantees aren't equivalent across the two systems?

You: NATS gives us per-subject ordering by default; we don't actually
     need cross-subject ordering, so I think we're fine.

Agent: How would you notice if you were wrong? What's the loudest
       symptom of a stray out-of-order message in production today?

You: ...probably a duplicate notification. Which we already see
     occasionally. So we wouldn't notice immediately.

Agent: Then your migration's correctness check needs to be stronger
       than "no incidents." Before you port consumer two, what
       evidence would convince you that consumer one's ordering
       behavior is correct under load?

You: Right. Add a dedicated assertion test in staging that replays
     a known-out-of-order sequence and verifies the consumer
     dedupes correctly. Otherwise we're just hoping.

Agent: Good. Next: what would change your mind about NATS itself?
       What rejected alternative is closest to coming back?
```

The interview keeps going until the bet is articulated, the
falsifiable conditions are written down, and the rejected
alternatives have a recall trigger.

## Output

`/ctx-plan` produces a clearer plan, not a document. Persist the
deltas via:

- **`/ctx-spec`** if the conclusions belong in a feature spec.
- **`/ctx-decision-add`** if a tradeoff resolved into an
  architectural decision.
- **`/ctx-learning-add`** if you discovered a project-specific
  gotcha during the interview.

The skill itself is in
`internal/assets/claude/skills/ctx-plan/SKILL.md`; the working
contract lives there, the recipe is the on-ramp.

## See Also

- [Design Before Coding](design-before-coding.md): the
  brainstorming counterpart, used *before* a plan exists.
- [`ctx-spec`](../cli/index.md): scaffolds a feature spec from
  the project template.

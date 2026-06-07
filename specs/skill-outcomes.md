# Skill Outcomes (success criteria for ctx skills)

*Promoted from the outcomes sub-idea in `ideas/000-dreams.md` by a
ctx-dream serendipity round (2026-06-07). Seed-stage spec. Inspired by
Anthropic Managed Agents "define outcomes"
(platform.claude.com/docs/en/managed-agents/define-outcomes): agents do
their best work when they know what "good" looks like, and a grader can
check output against an example ideal.*

## Problem

ctx skills describe *what to do* and *when to trigger*, but rarely state
what a **good result looks like**. Without an explicit success criterion,
quality is implicit and uneven: a skill can "run" and produce weak output
(a vague learning, a spec with placeholders, a shallow review) with
nothing to check it against. The agent has no target, and there's no
cheap way to grade whether a skill's output met the bar.

## Approach

Give skills an optional, declared **outcome**: a short success criterion
(and, where useful, an example of an ideal result) the skill's output
should meet. Then either:

- **Self-check** — the skill ends by checking its own output against the
  stated outcome before returning, or
- **Grader pass** — a separate, cheap grader evaluates the output against
  the outcome (the Anthropic "grader agent" shape), surfacing a
  pass/weak/fail signal.

Start with the skills where "good" is most checkable and most often
missed (e.g. `/ctx-learning-add` — is the lesson a principle, not a
snippet?; `/ctx-spec` — are edge cases + error handling present, no
placeholder `...`?; `/ctx-code-review` — are findings evidence-backed?).

## Behavior

### Happy Path

1. A skill declares an `outcome:` (criterion + optional example) in its
   frontmatter or body.
2. On completion, the outcome is evaluated (self-check or grader).
3. A weak/fail result prompts the agent to revise before returning, or
   flags the output for the user.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Skill has no declared outcome | No grading; unchanged behavior (opt-in). |
| Grader unavailable / budget exceeded | Skip grading; never block on it. |
| Subjective outcome (e.g. "brand voice") | Use an example-based criterion; grade leniently, surface not enforce. |

## Interface

- A convention for declaring `outcome:` in a SKILL.md (criterion +
  optional example).
- Optionally a `ctx skill grade <name>` / grader hook (TBD).

## Non-Goals

- Not auto-rejecting skill output (outcomes inform/flag; the human
  decides — consistent with ctx's no-autonomous-gate stance).
- Not a per-task outcome system (this is per-skill success criteria, not
  Anthropic's per-task outcomes verbatim).

## Open Questions

- **Where outcomes live** — SKILL.md frontmatter vs a sibling file vs a
  registry. (TBD)
- **Self-check vs grader** — which skills warrant a separate grader pass
  vs an inline self-check? (TBD)
- **Grader model/budget** — reuse the session model, or a cheap grader?
  (TBD; mirror the ctx-dream summary-model question.)
- **Which skills get outcomes first** — needs a pass over the skill set
  ranking "good is checkable & often missed." (TBD — `/ctx-plan`.)

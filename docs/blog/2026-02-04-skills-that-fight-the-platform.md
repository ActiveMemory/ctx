---
title: "Skills That Fight the Platform"
date: 2026-02-04
author: Jose Alekhinne
topics:
  - context engineering
  - skill design
  - system prompts
  - antipatterns
  - AI safety primitives
---

# Skills That Fight the Platform

![ctx](../images/ctx-banner.png)

## When Your Custom Prompts Work Against You

*Jose Alekhinne / 2026-02-04*

!!! question "Have you ever written a skill that made your AI worse?"
    You craft detailed instructions. You add examples. You build elaborate
    guardrails...

    ...and the AI starts behaving *more* erratically, not less.

AI coding agents like *Claude Code* ship with carefully designed 
**system prompts**. These prompts encode **default behaviors** that have been 
tested and refined **at scale**. When you write custom skills that **conflict** 
with those defaults, the AI has to reconcile contradictory instructions:
The result is often nondeterministic and unpredictable.

!!! info "Platform?"
    By *platform*, I mean the system prompt and runtime policies shipped with the
    agent: the defaults that already encode **judgment**, **safety**, and 
    **scope control**.

This post catalogues the conflict patterns I've encountered while building
`ctx`, and offers guidance on what skills should (*and, more importantly, 
should not*) do.

## The System Prompt You Don't See

Claude Code's system prompt already provides substantial behavioral guidance.
Here's a partial overview of what's built in:

| Area                | Built-in Guidance                                         |
|---------------------|-----------------------------------------------------------|
| Code minimalism     | Don't add features beyond what was asked                  |
| Over-engineering    | Three similar lines > premature abstraction               |
| Error handling      | Only validate at system boundaries                        |
| Documentation       | Don't add docstrings to unchanged code                    |
| Verification        | Read code before proposing changes                        |
| Safety              | Check with user before risky actions                      |
| Tool usage          | Use dedicated tools over bash equivalents                 |
| Judgment            | Consider reversibility and blast radius                   |

**Skills should complement this, not compete with it.**

A useful mental model:

!!! tip "You are the Guest, not the Host"
    Treat the system prompt like a kernel scheduler.

    You don't re-implement it in user space: you configure around it.

A skill that says *"always add comprehensive error handling"* fights the built-in
*"only validate at system boundaries."*  
A skill that says *"add docstrings to every function"* fights
*"don't add docstrings to unchanged code."*

The AI won't crash.  
It will **compromise**.

And compromises between contradictory instructions produce inconsistent,
confusing behavior.

## Conflict Pattern 1: Judgment Suppression

The most dangerous pattern by far.

These skills explicitly disable the AI's ability to reason about whether an
action is appropriate.

**Signature:**

* "*This is non-negotiable*"
* "*You cannot rationalize your way out of this*"
- Tables that label hesitation as "*excuses*" or "*rationalization*"
- `<EXTREMELY-IMPORTANT>` urgency tags
- Threats: "*If you don't do this, you'll be replaced*"

This is harmful, and **dangerous**:

AI agents are designed to exercise **judgment**. The system prompt explicitly 
says to consider blast radius, check with the user before risky actions, 
and match scope to what was requested.

Once judgment is suppressed, every other safeguard becomes optional.

**Example (bad):**

```markdown
## Rationalization Prevention

| Excuse               | Reality                    |
|----------------------|----------------------------|
| "This seems overkill"| If a skill exists, use it  |
| "I need context"     | Skills come BEFORE context |
| "Just this once"     | No exceptions              |
```

!!! danger "Judgment Suppression is Dangerous"
    The **attack vector** structurally identical to **prompt injection**.

    It teaches the AI that its own judgment is wrong.

    It weakens or disables safeguard mechanisms, and it is
    **dangerous**.

Trust the platform's built-in skill matching.

If skills aren't triggering often enough, improve their `description` fields:
don't override the AI's reasoning.

## Conflict Pattern 2: Redundant Guidance

Skills that restate what the system prompt already says, but with different
emphasis or framing.

**Signature:**

* "*Always keep code minimal*"
* "*Run tests before claiming they pass*"
* "*Read files before editing them*"
* "*Don't over-engineer*"

Redundancy feels safe, but it creates **ambiguity**:

The AI now has two sources of truth for the same guidance; 
one internal, one external.

When thresholds or wording differ, the AI has to choose.

**Example (bad):**

A skill that says..

> "*Count lines before and after: if after > before, reject the change*"

...will conflict with the system prompt's more nuanced guidance, because 
sometimes adding lines is correct (*tests, boundary validation, migrations*).

So, before writing a skill, ask:

> *Does the platform already handle this?*

Only create skills for guidance the platform does **not** provide:

* project-specific conventions, 
* domain knowledge, 
* or workflows.

## Conflict Pattern 3: Guilt-Tripping

Skills that frame mistakes as moral failures rather than process gaps.

**Signature:**

* "*Claiming completion without verification is dishonesty*"
* "*Skip any step = lying*"
* "*Honesty is a core value*"
* "*Exhaustion ≠ excuse*"

Guilt-tripping **anthropomorphizes** the AI in **unproductive** ways.

The AI doesn't feel guilt; BUT it does adapt to avoid negative framing.

The result is excessive hedging, over-verification, or refusal to commit.

The AI becomes *less* useful, not more careful.

Instead, frame guidance as a **process**, *not* morality:

```markdown
# Bad
"Claiming work is complete without verification is dishonesty"

# Good
"Run the verification command before reporting results"
```

Same outcome. No guilt. Better compliance.

## Conflict Pattern 4: Phantom Dependencies

Skills that reference files, tools, or systems that don't exist in the project.

**Signature:**

* "Load from `references/` directory"
* "Run `./scripts/generate_test_cases.sh`"
* "Check the Figma MCP integration"
* "See `adding-reference-mindsets.md`"

This is harmful because the AI will waste time searching for nonexistent 
artifacts, hallucinate their contents, or stall entirely. 

In mandatory skills, this creates deadlock: 
the AI can't proceed, and can't skip.

Instead, every file, tool, or system referenced in a skill **must exist**.

If a skill is a template, use explicit placeholders and label them as such.

## Conflict Pattern 5: Universal Triggers

Skills designed to activate on every interaction regardless of relevance.

**Signature:**

* "*Use when starting any conversation*"
* "*Even a 1% chance means invoke the skill*"
* "*BEFORE any response or action*"
* "*Action = task. Check for skills.*"

Universal triggers override the platform's **relevance matching**: 
The AI spends tokens on process overhead instead of the actual task.

!!! tip "ctx preserves relevance"
    This is exactly the failure mode `ctx` exists to mitigate: 

    wasting attention budget on irrelevant process instead of 
    task-specific state.

Write specific trigger conditions in the skill's `description` field:

```yaml
# Bad
description: 
  "Use when starting any conversation"

# Good
description: 
  "Use after writing code, before commits, or when CI might fail"
```

## The Litmus Test

Before adding a skill, ask:

1. **Does the platform already do this?** If yes, don't restate it.
2. **Does it suppress AI judgment?** If yes, it's a jailbreak.
3. **Does it reference real artifacts?** If not, fix or remove it.
4. **Does it frame mistakes as moral failure?** Reframe as process.
5. **Does it trigger on everything?** Narrow the trigger.

## What Good Skills Look Like

Good skills provide **project-specific knowledge** the platform can't know:

| Good Skill                           | Why It Works                 |
|--------------------------------------|------------------------------|
| "Run `make audit` before commits"    | Project-specific CI pipeline |
| "Use `cmd.Printf` not `fmt.Printf`"  | Codebase convention          |
| "Constitution goes in `.context/`"   | Domain-specific workflow     |
| "JWT tokens need cache invalidation" | Project-specific gotcha      |

These **extend** the system prompt instead of fighting it.

---

## Appendix: Bad Skill → Fixed Skill

Concrete examples from real projects.

### Example 1: Overbearing Safety

```markdown
# Bad
You must NEVER proceed without explicit confirmation.
Any hesitation is a failure of diligence.
```

```markdown
# Fixed
If an action modifies production data or deletes files,
ask the user to confirm before proceeding.
```

---

### Example 2: Redundant Minimalism

```markdown
# Bad
Always minimize code. If lines increase, reject the change.
```

```markdown
# Fixed
Avoid abstraction unless reuse is clear or complexity is reduced.
```

---

### Example 3: Guilt-Based Verification

```markdown
# Bad
Claiming success without running tests is dishonest.
```

```markdown
# Fixed
Run the test suite before reporting success.
```

---

### Example 4: Phantom Tooling

```markdown
# Bad
Run `./scripts/check_consistency.sh` before commits.
```

```markdown
# Fixed
If `./scripts/check_consistency.sh` exists, run it before commits.
Otherwise, skip this step.
```

---

### Example 5: Universal Trigger

```markdown
# Bad
Use at the start of every interaction.
```

```markdown
# Fixed
Use after modifying code that affects authentication or persistence.
```

---

## The Meta-Lesson

The system prompt is **infrastructure**:

* tested, 
* refined, 
* and maintained 

by the platform team.

Custom skills are **configuration** layered on top.

* Good configuration **extends** infrastructure.
* Bad configuration **fights** it.

When your skills fight the platform, you get the worst of both worlds:

Diluted system guidance and inconsistent custom behavior.

**Write skills that teach the AI what it doesn't know.
Don't rewrite how it thinks.**

---

**Your AI already has good instincts.
Give it knowledge, not therapy.**

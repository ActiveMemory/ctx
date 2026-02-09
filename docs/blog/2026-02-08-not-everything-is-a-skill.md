---
title: "Not Everything Is a Skill"
date: 2026-02-08
author: Jose Alekhinne
topics:
  - skill design
  - context engineering
  - automation discipline
  - recipes
  - agent teams
---

# Not Everything Is a Skill

![ctx](../images/ctx-banner.png)

## What a Codebase Audit Taught Me About Restraint

*Jose Alekhinne / 2026-02-08*

!!! question "When you find a useful prompt, what do you do with it?"
    My instinct was to make it a skill. 

    I had just spent **three posts** explaining how to build skills that work. 
    Naturally, **the hammer wanted nails**.

    Then I looked at what I was holding and realized: **this is not
    a nail**.

## The Audit

I wanted to understand how I use `ctx`: 
* 
* where the friction is, 
* what works, what drifts, 
* what I keep doing manually that could be automated. 

So I wrote a prompt that spawned eight agents to analyze the codebase from 
different angles:

| Agent | Analysis                                          |
|-------|---------------------------------------------------|
| 1     | Extractable patterns from session history         |
| 2     | Documentation drift (godoc, inline comments)      |
| 3     | Maintainability (large functions, misplaced code) |
| 4     | Security review (CLI-specific surface)            |
| 5     | Blog theme discovery                              |
| 6     | Roadmap and value opportunities                   |
| 7     | User-facing documentation gaps                    |
| 8     | Agent team strategies for future sessions         |

The prompt was **specific**: 

* read-only agents, 
* structured output format,
* concrete file references, 
* ranked recommendations. 

It ran for about  20 minutes and produced **eight Markdown reports**.

The reports were good: Not perfect, but actionable.

What mattered was not the speed. It was that the work could be
explored without committing to any single outcome.

They surfaced a stale `doc.go` referencing a subcommand that was
never built. 

They found **311 build-then-test sequences** I could reduce
to a single `make check`. 

They identified that **42% of my sessions start with "do you remember?"**, 
which is a lot of repetition for something a skill could handle.

**I had findings. I had recommendations. I had the instinct to
automate.**

And then... **I stopped**.

## The Question

The natural next step was to wrap the audit prompt as `/ctx-audit`:
a skill you invoke periodically to get a health check. It fits the
pattern. It has a clear trigger. It produces structured output.

But I had just spent a week writing about what makes skills work,
and the criteria I established argued against it.

From [The Anatomy of a Skill That Works][anatomy-post]:

> "*A skill without boundaries is just a suggestion.*"

From [You Can't Import Expertise][import-post]:

> "*Frameworks travel, expertise doesn't.*"

From [Skills That Fight the Platform][fight-post]:

> "*You are the guest, not the host.*"

The audit prompt fails all three tests:

| Criterion | Audit prompt               | Good skill              |
|-----------|----------------------------|-------------------------|
| Frequency | Quarterly, maybe           | Daily or weekly         |
| Stability | Tweaked every time         | Consistent invocation   |
| Scope     | Bespoke, 8 parallel agents | Single focused action   |
| Trigger   | "I feel like auditing"     | Clear, repeatable event |

Skills are **contracts**. Contracts need **stable terms**. 

A prompt I will rewrite every time I use it is not a contract. 
It is a *conversation starter*.

[anatomy-post]: 2026-02-07-the-anatomy-of-a-skill-that-works.md
[import-post]: 2026-02-05-you-cant-import-expertise.md
[fight-post]: 2026-02-04-skills-that-fight-the-platform.md

## Recipes vs Skills

The distinction that emerged:

|                    | Skill                        | Recipe                       |
|--------------------|------------------------------|------------------------------|
| **Invocation**     | `/slash-command`             | Copy-paste from a doc        |
| **Frequency**      | High (daily, weekly)         | Low (quarterly, ad hoc)      |
| **Stability**      | Fixed contract               | Adapted each time            |
| **Scope**          | One focused action           | Multi-step orchestration     |
| **Audience**       | The agent                    | The human (who then prompts) |
| **Lives in**       | `.claude/skills/`            | `hack/` or `docs/`           |
| **Attention cost** | Loaded into context on match | Zero until needed            |

Recipes can later graduate into skills, but only after repetition
proves stability.

That last row matters. Skills consume the
[attention budget][attention-post] every time the platform considers
activating them. A skill that triggers quarterly but gets evaluated
on every prompt is pure waste: attention spent on something that
will say "When NOT to Use: now" 99% of the time.

Recipes have **zero** attention cost. They sit in a Markdown file until
a human decides to use them. The human provides the **judgment** about
timing. The prompt provides the **structure**.

[attention-post]: 2026-02-03-the-attention-budget.md

!!! tip "The Attention Budget Applies to Skills Too"
    Every skill in `.claude/skills/` is a standing claim on the
    context window. The platform evaluates skill descriptions
    against every user prompt to decide whether to activate.

    Twenty focused skills are fine. Thirty might be fine. But
    each one added reduces the headroom available for actual work.

    **Recipes are skills that opted out of the attention tax.**

## What the Audit Actually Produced

The audit was not wasted. It was a planning exercise that generated
concrete tasks:

| Finding                                         | Action                                                           |
|-------------------------------------------------|------------------------------------------------------------------|
| 42% of sessions start with memory check         | Task: `/ctx-remember` skill (this one *is* a skill; it is daily) |
| Auto-save stubs are empty                       | Task: enhance `/ctx-save` with richer summaries                  |
| 311 raw build-test sequences                    | Task: `make check` target                                        |
| Stale `recall/doc.go` lists nonexistent `serve` | Task: fix the doc.go                                             |
| 120 commit sequences disconnected from context  | Task: `/ctx-commit` workflow                                     |

Some findings became skills. Some became `Makefile` targets. Some
became one-line doc fixes. 

The audit did not prescribe the artifact type. **The findings did**.

**The audit is the input. Skills are one possible output. Not
the only one.**

## The Audit Prompt

Here is the exact prompt I used, for those who are curious.

This is not a template. It worked because it was written against this
codebase, at this moment, with specific goals in mind.

```markdown
I want you to create an agent team to audit this codebase. Save each report as
a separate markdown file under `./ideas/` (or another directory if you prefer).

Use read-only agents (subagent_type: Explore) for all analyses. No code changes.

For each report, use this structure:
- Executive Summary (2-3 sentences + severity table)
- Findings (grouped, with file:line references)
- Ranked Recommendations (high/medium/low priority)
- Methodology (what was examined, how)

Keep reports actionable. Every finding should suggest a concrete fix or next step.

## Analyses to Run

### 1. Extractable Patterns (session mining)
Search session JSONL files, journal entries, and task archives for repetitive
multi-step workflows. Count frequency of bash command sequences, slash command
usage, and recurring user prompts. Identify patterns that could become skills
or scripts. Cross-reference with existing skills to find coverage gaps.
Output: ranked list of automation opportunities with frequency data.

### 2. Documentation Drift (godoc + inline)
Compare every doc.go against its package's actual exports and behavior. Check
inline godoc comments on exported functions against their implementations.
Scan for stale TODO/FIXME/HACK comments. Check that package-level comments match
package names.
Output: drift items ranked by severity with exact file:line references.

### 3. Maintainability
Look for:
- functions longer than 80 lines with clear split points
- switch blocks with more than 5 cases that could be table-driven
- inline comments like "step 1", "step 2" that indicate a block wants to be a function
- files longer than 400 lines
- flat packages that could benefit from sub-packages
- functions that appear misplaced in their file

Do NOT flag things that are fine as-is just because they could theoretically
be different.
Output: concrete refactoring suggestions, not style nitpicks.

### 4. Security Review
This is a CLI app. Focus on CLI-relevant attack surface, not web OWASP:
- file path traversal
- command injection
- symlink following when writing to `.context/`
- permission handling
- sensitive data in outputs

Output: findings with severity ratings and plausible exploit scenarios.

### 5. Blog Theme Discovery
Read existing blog posts for style and narrative voice. Analyze git history,
recent session discussions, and DECISIONS.md for story arcs worth writing about.
Suggest 3-5 blog post themes with:
- title
- angle
- target audience
- key commits or sessions to reference
- a 2-sentence pitch

Prioritize themes that build a coherent narrative across posts.

### 6. Roadmap and Value Opportunities
Based on current features, recent momentum, and gaps found in other analyses,
identify the highest-value improvements. Consider user-facing features,
developer experience, integration opportunities, and low-hanging fruit.
Output: prioritized list with rough effort and impact estimates.

### 7. User-Facing Documentation
Evaluate README, help text, and user docs. Suggest improvements structured as
use-case pages: the problem, how ctx solves it, a typical workflow, and gotchas.
Identify gaps where a user would get stuck without reading source code.
Output: documentation gaps with suggested page outlines.

### 8. Agent Team Strategies
Based on the codebase structure, suggest 2-3 agent team configurations for
upcoming work sessions. For each, include:
- team composition (roles and agent types)
- task distribution strategy
- coordination approach
- the kinds of work it suits
```

!!! warning "Avoid Generic Advice"
    Suggestions that are not grounded in a project's actual structure,
    history, and workflows are **worse than useless**:

    They create **false confidence**.

    If an analysis cannot point to concrete files, commits, 
    sessions, or patterns, it should say "*no finding*" 
    instead of inventing best practices.

## The Deeper Pattern

This is part of a pattern I keep rediscovering: the **urge** to automate
is **not** the same as the **need** to automate:

* The [3:1 ratio][refactor-post] taught me that not every session
should be a YOLO sprint. 
* The [E/A/R framework][import-post] taught me that not every template 
  is worth importing. Now the audit is teaching me that 
  **not every useful prompt is worth institutionalizing**.

[refactor-post]: 2026-02-01-refactoring-with-intent.md

The common thread is **restraint**: Knowing when to stop. Recognizing
that the cost of automation is not just the effort to build it.
It is the **ongoing attention tax** of maintaining it, the context
it consumes, and the false confidence it creates when it drifts.

A recipe in `hack/codebase-audit.md` is honest about what it is:

A prompt I wrote once, improved once, and will adapt again next
time: 

* It does **not** pretend to be a reliable contract. 
* It does **not** claim *attention budget*. 
* It does **not** drift silently.

!!! warning "The Automation Instinct"
    When you find a useful prompt, the instinct is to
    institutionalize it. **Resist**.

    Ask first: **will I use this the same way next time?**

    If yes, it is a skill. If no, it is a recipe. If you are
    not sure, it is a recipe until proven otherwise.

## This Mindset In the Context of `ctx`

`ctx` is a **tool** that gives AI agents persistent memory. Its purpose
is **automation**: reducing the **friction** of context loading, session
recall, decision tracking.

But **automation has boundaries**, and knowing where those boundaries
are is as important as pushing them forward. 

The **skills system** is for high-frequency, stable workflows. 

The **recipes**, the **journal entries**, the **session dumps** in 
`.context/sessions/`: those are for everything else.

**Not everything needs to be a slash command. Some things are
better as Markdown files you read when you need them.**

The goal of `ctx` is **not** to automate everything: It is to automate
the right things and to make the rest easy to find when you need it.

---

!!! quote "If you remember one thing from this post..."
    **The best automation decision is sometimes not to automate.**

    A recipe in a markdown file costs nothing until you use it.
    A skill costs attention on every prompt, whether it fires or not.

    **Automate the daily. Document the periodic. Forget the rest.**

---

*This post was written during the session that produced the codebase
audit reports and distilled the prompt into `hack/codebase-audit.md`.
The audit generated seven tasks, one Makefile target, and zero new
skills. The meta continues.*

---
title: "You Can't Import Expertise"
date: 2026-02-05
author: Jose Alekhinne
topics:
  - skill adaptation
  - E/A/R framework
  - convention drift
  - consolidation
  - project-specific expertise
---

# You Can't Import Expertise

![ctx](../images/ctx-banner.png)

## Why Good Skills Can't Be Copy-Pasted

*Jose Alekhinne / 2026-02-05*

!!! question "Have you ever dropped a well-crafted template into a project and had it do... nothing useful?"
    The template was thorough. The structure was sound. The advice was correct.

    And yet it sat there, inert, while the same old problems kept drifting in.

I found a consolidation skill online. It was well-organized: four files,
ten refactoring patterns, eight analysis dimensions, six report templates.
Professional. Comprehensive. Exactly the kind of thing you'd bookmark and
think *"I'll use this."*

Then I stopped, and applied `ctx`'s own evaluation framework: 

**70% of it was noise!**

This post is about **why**.

> **Templates describe categories of problems.**
> 
> **Expertise encodes which problems actually happen, and how often.**

## The Skill Looked Great on Paper

Here is what the consolidation skill offered:

| File                        | Content                                                      |
|-----------------------------|--------------------------------------------------------------|
| `SKILL.md`                  | Entry point: 8 analysis dimensions, workflow, output formats |
| `analysis-dimensions.md`    | Detailed criteria for duplication, architecture, quality     |
| `consolidation-patterns.md` | 10 refactoring patterns with before/after code               |
| `report-templates.md`       | 6 output templates: executive summary, roadmap, onboarding   |

It had a scoring system (`0-10` per dimension, letter grades `A+` through `F`).
It had severity classifications with color-coded emojis. It had bash commands
for detection. It even had antipattern warnings.

**By any standard template review, this skill passes.**

It looks like something an expert wrote. 

And that's **exactly** the trap.

## Applying E/A/R: The 70-20-10 Split

In a [previous post][skills-post], I described the **E/A/R framework** for
evaluating skills:

[skills-post]: 2026-02-04-skills-that-fight-the-platform.md "Skills That Fight the Platform"

* **Expert**: Knowledge that took years to learn. Keep.
* **Activation**: Useful triggers or scaffolding. Keep if lightweight.
* **Redundant**: Restates what the AI already knows. Delete.

Target: >70% Expert, <10% Redundant.

This skill scored the inverse.

### What Was Redundant (~70%)

Every code example was **Rust**. My project is **Go**.

The analysis dimensions: duplication detection, architectural structure,
code organization, refactoring opportunities... These are things Claude already
does when you ask it to review code. 

The skill restated them with more ceremony but no more insight.

The six report templates were generic scaffolding: Executive Summary,
Onboarding Document, Architecture Documentation. They are useful if you are
writing a consulting deliverable, but not when you are trying to catch
convention drift in a >15K-line Go CLI.

## What Does a `B+` in Code Organization *Actually* Mean?!

The scoring system (*`0-10` per dimension, letter grades*) added ceremony
without actionable insight. 

What is a `B+`? What do I do differently for an `A-`?

**The skill told the AI what it already knew, in more words.**

### What Was Activation (*~10%*)

The consolidation checklist (*semantics preserved? tests pass? docs updated?*)
was useful as a **gate**. But, it's the kind of thing you could inline in three
lines.

The phased roadmap structure was reasonable scaffolding for sequencing work.

### What Was Expert (*~20%*)

Three concepts survived:

1. **The Consolidation Decision Matrix**: A concrete framework mapping
   similarity level and instance count to action. "*Exact duplicate, 2+
   instances: consolidate immediately.*" "*<3 instances: leave it:
   duplication is cheaper than wrong abstraction.*" This is the kind of
   nuance that prevents premature generalization.

2. **The Safe Migration Pattern**: Create the new API alongside old, deprecate,
   migrate incrementally, delete. Straightforward to describe, yet
   forgettable under pressure.

3. **Debt Interest Rate framing**: Categorizing technical debt by how fast
   it compounds (security vulns = daily, missing tests = per-change,
   doc gaps = constant low cost). This changes prioritization.

**Three ideas out of four files and 700+ lines.** The rest was filler
that competed with the AI's built-in capabilities.

## What the Skill Didn't Know

!!! tip "AI Without Context is Just a Corpus"
    LLMs are optimized on insanely large corpora.
    And then they are passed through several layers of
    human-assisted refinement.
    The whole process costs millions of dollars.

    Yet, the uncomfortable truth is that no corpus can "*infer*"
    your project's design, convetions, patterns, habits, history,
    vision, and deliverables.

    Your project is unique: So should your skills be.

Here is the part no template can provide: 

**`ctx`'s actual drift patterns.**

Before evaluating the skill, I did **archaeology**. I read through:

* Blog posts from previous refactoring sessions
* The project's learnings and decisions files
* Session journals spanning **weeks of development**

What I found was specific:

| Drift Pattern                       | Where                  | How Often                  |
|-------------------------------------|------------------------|----------------------------|
| `Is`/`Has`/`Can` predicate prefixes | 5+ exported methods    | Every YOLO sprint          |
| Magic strings instead of constants  | 7+ files               | Gradual accumulation       |
| Hardcoded file permissions (`0755`) | 80+ instances          | Since day one              |
| Lines exceeding 80 characters       | Especially test files  | Every session              |
| Duplicate code blocks               | Test and non-test code | When agent is task-focused |

The generic skill had no check for any of these. It couldn't; because these
patterns are specific to this project's conventions, its Go codebase, and
its development rhythm.

!!! tip "The Insight"
    The skill's analysis dimensions were about *categories of problems*.

    What I needed was *my specific problems*.

## The Adapted Skill

The adapted skill is roughly a quarter of the original's size.
It has nine checks, each targeting a known drift pattern:

1. **Predicate naming**: `rg` for `Is`/`Has`/`Can` prefixes
2. **Magic strings**: literals that should be constants
3. **Hardcoded permissions**: `0755`/`0644` literals
4. **File size**: source files over 300 LOC
5. **TODO/FIXME**: constitution violation (move to TASKS.md)
6. **Path construction**: string concatenation instead of `filepath.Join`
7. **Line width**: lines exceeding ~80 characters
8. **Duplicate blocks**: copy-paste drift, especially in tests
9. **Dead exports**: unused public API

**Every check** has a detection command. **Every check** maps to a specific
convention or constitution rule. **Every check** was discovered through
actual project history; **not** invented from a template.

The three expert concepts from the original survived:

* The decision matrix gates when to consolidate vs. when to leave
  duplication alone
* The safe migration pattern guides public API changes
* The relationship to other skills (`/qa`, `/verify`, `/update-docs`,
  `ctx drift`) prevents overlap

Nothing else made it.

## The Deeper Pattern

This experience crystallized something I've been circling for weeks:

**You can't import expertise. You have to grow it from your project's
own history.**

A skill that says "*check for code duplication*" is **not** expertise: 
It's a *category*. 

Expertise is **knowing**, in the heart of your hearts, that *this* 
project accumulates `Is*` predicate violations during velocity sprints, 
that *this* codebase has 80 hardcoded permission literals because nobody 
made a constant, that **this team**'s test files drift wide because the 
agent prioritizes getting the task done over keeping the code in shape.

!!! note "The Parallel to the 3:1 Ratio"
    In [Refactoring with Intent][refactoring-post], I described the 3:1
    ratio: three YOLO sessions followed by one consolidation session.

    The same ratio applies to skills: you need experience *in* the project
    before you can write effective guidance *for* the project.

    Importing a skill on day one is like scheduling a consolidation session
    before you've written any code.

[refactoring-post]: 2026-02-01-refactoring-with-intent.md "Refactoring with Intent"

## The Template Trap

Templates are seductive because they feel like progress:

* You found something
* It's well-organized
* It covers the topic
* It has concrete examples

**But coverage is not relevance.**

A template that covers eight analysis dimensions with Rust examples
adds zero value to a Go project with five known drift patterns. Worse,
it adds **negative** value: the AI spends **attention** *defending generic
advice* instead of noticing project-specific drift.

This is the [attention budget][attention-post] problem again. Every token
of generic guidance displaces a token of specific guidance. A 700-line
skill that's 70% redundant doesn't just waste 490 lines: it *dilutes*
the 210 lines that matter.

[attention-post]: 2026-02-03-the-attention-budget.md "The Attention Budget"

## The Litmus Test

Before dropping any external skill into your project:

1. **Run E/A/R**: What percentage is expert knowledge vs. what the AI
   already knows? If it's less than 50% expert, it's probably not worth
   the attention cost.

2. **Check the language**: Does it use your stack? Generic patterns in
   the wrong language are noise, not signal.

3. **List your actual drift**: Read your own session history, learnings,
   and post-mortems. What breaks *in practice*? Does the skill check for
   those things?

4. **Measure by deletion**: After adaptation, how much of the original
   survives? If you're keeping less than 30%, you would have been faster
   writing from scratch.

5. **Test against your conventions**: Does every check in the skill map
   to a specific convention or rule in your project? If not, it's
   generic advice wearing a skill's clothing.

## What Good Adaptation Looks Like

The consolidation skill went from:

| Before                   | After                                                    |
|--------------------------|----------------------------------------------------------|
| 4 files, 700+ lines      | 1 file, ~120 lines                                       |
| Rust examples            | Go-specific `rg` commands                                |
| 8 generic dimensions     | 9 project-specific checks                                |
| 6 report templates       | 1 focused output format                                  |
| Scoring system (A+ to F) | Findings + priority + suggested fixes                    |
| "Check for duplication"  | "Check for `Is*` predicate prefixes in exported methods" |

The adapted version is smaller, faster to parse, and catches the things
that **actually** drift in this project.

**That's the difference between a template and a tool.**

---

!!! quote "If you remember one thing from this post..."
    **Frameworks travel. Expertise doesn’t.**

    You can import structures, matrices, and workflows.

    But the checks that matter only grow where the scars are:

    * the conventions that were violated, 
    * the patterns that drifted,
    * and the specific ways **this** codebase accumulates debt.

---

*This post was written during a consolidation session where the
consolidation skill itself became the subject of consolidation.
The meta continues.*
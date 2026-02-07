---
title: Blog
icon: lucide/newspaper
---

![ctx](../images/ctx-banner.png)

Stories, insights, and lessons learned from building and using ctx.

---

## Posts

### [The Anatomy of a Skill That Works](2026-02-07-the-anatomy-of-a-skill-that-works.md)

*Jose Alekhinne / February 7, 2026*

I had 20 skills. Most were well-intentioned stubs: a description, a command,
and a wish for the best. Then I rewrote all of them in a single session.
Seven lessons emerged: quality gates prevent premature execution, negative
triggers are load-bearing, examples set boundaries better than rules, and
skills are contracts, not instructions. The practical companion to the
previous two skill design posts.

**Topics**: skill design, context engineering, quality gates, E/A/R framework,
practical patterns

---

### [You Can't Import Expertise](2026-02-05-you-cant-import-expertise.md)

*Jose Alekhinne / February 5, 2026*

I found a well-crafted consolidation skill: four files, ten patterns, eight
analysis dimensions. Applied my own E/A/R framework: **70% was noise**. The
template was thorough, correct, and almost entirely useless for my project.
This post is about why good skills can't be copy-pasted, and how to grow
them from your project's own drift history instead.

**Topics**: skill adaptation, E/A/R framework, convention drift, consolidation,
project-specific expertise

---

### [Skills That Fight the Platform](2026-02-04-skills-that-fight-the-platform.md)

*Jose Alekhinne / February 4, 2026*

AI coding agents ship with carefully designed system prompts. When custom
skills conflict with those defaults, the AI has to **reconcile contradictory
instructions**: The result is **unpredictable**. This post catalogues five
conflict patterns discovered while building `ctx`: judgment suppression,
redundant guidance, guilt-tripping, phantom dependencies, and universal
triggers.

**Topics**: context engineering, skill design, system prompts, antipatterns,
AI safety primitives

---

### [The Attention Budget: Why Your AI Forgets What You Just Told It](2026-02-03-the-attention-budget.md)

*Jose Alekhinne / February 3, 2026*

Every token you send to an AI consumes a finite resource: **the attention budget**.
Understanding this constraint shaped every design decision in ctx: hierarchical
file structure, explicit budgets, progressive disclosure, and filesystem-as-index.
This post explains the theory and how ctx operationalizes it.

**Topics**: attention mechanics, context engineering, progressive disclosure,
ctx primitives, token budgets

---

### [ctx v0.2.0: The Archaeology Release](2026-02-01-ctx-v0.2.0-the-archaeology-release.md)

*Jose Alekhinne / February 1, 2026*

What if your AI could remember everything? Not just the current session, but
every session. `ctx v0.2.0` introduces the **recall** and **journal** systems:
making 86 commits of history searchable, exportable, and analyzable. 
This post tells the story of why those features exist.

**Topics**: session recall, journal system, structured entries, token budgets,
meta-tools

---

### [Refactoring with Intent: Human-Guided Sessions in AI Development](2026-02-01-refactoring-with-intent.md)

*Jose Alekhinne / February 1, 2026*

The **YOLO mode shipped** **14 commands** **in a week**. But technical debt 
doesn't send invoices:it just waits. This is the story of what happened when 
we stopped auto-accepting everything and started guiding the AI with intent: 
27 commits across 4 days, a major version release, and lessons that apply far 
beyond ctx.

**Topics**: refactoring, code quality, documentation standards, module
decomposition, YOLO versus intentional development

---

### [Building ctx Using ctx: A Meta-Experiment in AI-Assisted Development](2026-01-27-building-ctx-using-ctx.md)

*Jose Alekhinne / January 27, 2026*

What happens when you build a tool designed to give AI memory, using that very
same tool to remember what you're building? **This is the story of `ctx`**:
how `ctx` evolved from a hasty "*YOLO*" experiment to a **disciplined system** 
for persistent AI context, and what we learned along the way.

**Topics**: dogfooding, AI-assisted development, Ralph Loop, session
persistence, architectural decisions

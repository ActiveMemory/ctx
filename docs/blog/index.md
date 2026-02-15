---
title: Blog
icon: lucide/newspaper
---

![ctx](../images/ctx-banner.png)

Stories, insights, and lessons learned from building and using ctx.

---

## Posts

### [ctx v0.3.0: The Discipline Release](2026-02-15-ctx-v0.3.0-the-discipline-release.md)

*Jose Alekhinne / February 15, 2026*

No new headline feature. No architectural pivot. No rewrite. Just 35+
documentation and quality commits against ~15 feature commits — and
somehow, the tool feels like it grew up overnight. This post is about
what a release looks like when the ratio of polish to features is 3:1.

**Topics**: release, skills migration, consolidation, code quality,
E/A/R framework

---

### [Why Zensical](2026-02-15-why-zensical.md)

*Jose Alekhinne / February 15, 2026*

I needed a static site generator for the journal system and the instinct
was Hugo — same language as ctx, fast, well-established. But instinct is
not analysis. The journal entries are standard Markdown with YAML
frontmatter. That is it. No JSX, no shortcodes, no custom templating.
This post is about why zensical — a pre-1.0 tool from the Material for
MkDocs team — was the right choice: thin dependencies, MkDocs-compatible
config, 4-5x faster incremental rebuilds, and zero lock-in.

**Topics**: tooling, static site generators, journal system,
infrastructure decisions, context engineering

---

### [Parallel Agents with Git Worktrees](2026-02-14-parallel-agents-with-worktrees.md)

*Jose Alekhinne / February 14, 2026*

I had 30 open tasks and most of them didn't touch the same files. Running
one agent at a time meant serial execution on work that was fundamentally
parallel. This post is about using git worktrees to partition a backlog by
file overlap, run 3-4 agents simultaneously, and merge the results — the
same attention budget principle applied to execution instead of context.

**Topics**: agent teams, parallelism, git worktrees, context engineering,
task management

---

### [Before Context Windows, We Had Bouncers](2026-02-14-irc-as-context.md)

*Jose Alekhinne / February 14, 2026*

IRC is stateless. You disconnect, you vanish. Modern systems are not much
different: close the tab, lose the scrollback, open a new LLM session, start
from zero. This post traces the line from IRC bouncers like ZNC to context
engineering: stateless protocols require stateful wrappers, volatile interfaces
require durable memory. Before context windows, we had bouncers. Before AI
memory files, we had buffers.

**Topics**: context engineering, infrastructure, IRC, persistence,
state continuity

---

### [How Deep Is Too Deep?](2026-02-12-how-deep-is-too-deep.md)

*Jose Alekhinne / February 12, 2026*

I kept feeling like I should go deeper into ML theory. Then I spent a week
debugging an agent failure that had nothing to do with model architecture and
everything to do with knowing which abstraction was leaking. This post is
about when depth compounds and when it doesn't: why the useful understanding
lives one or two layers below where you work, not at the bottom of the stack.

**Topics**: AI foundations, abstraction boundaries, agentic systems,
context engineering, failure modes

---

### [Defense in Depth: Securing AI Agents](2026-02-09-defense-in-depth-securing-ai-agents.md)

*Jose Alekhinne / February 9, 2026*

I was writing the autonomous loops recipe and realized the security
advice was "use CONSTITUTION.md for guardrails." Then I read that
sentence back and realized: that is wishful thinking. This post
traces five defense layers for unattended AI agents, each with a
bypass, and shows why the strength is in the combination, not any
single layer.

**Topics**: agent security, defense in depth, prompt injection,
autonomous loops, container isolation

---

### [Not Everything Is a Skill](2026-02-08-not-everything-is-a-skill.md)

*Jose Alekhinne / February 8, 2026*

I ran an 8-agent codebase audit and got actionable results. The natural
instinct was to wrap the prompt as a `/ctx-audit` skill. Then I applied
my own criteria from the skill trilogy: **it failed all three tests**.
This post is about the difference between skills and recipes, why the
attention budget applies to your skill library too, and why the best
automation decision is sometimes not to automate.

**Topics**: skill design, context engineering, automation discipline,
recipes, agent teams

---

### [The Anatomy of a Skill That Works](2026-02-07-the-anatomy-of-a-skill-that-works.md)

*Jose Alekhinne / February 7, 2026*

I had 20 skills. Most were well-intentioned stubs: a description, a command,
and a wish for the best. Then I rewrote all of them in a single session.
**Seven lessons emerged**: quality gates prevent premature execution, negative
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

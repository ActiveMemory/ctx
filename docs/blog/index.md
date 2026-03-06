---
title: Blog
icon: lucide/newspaper
---

![ctx](../images/ctx-banner.png)

Stories, insights, and lessons learned from **building** and **using** `ctx`.

---

## Context Engineering

### [Agent Memory Is Infrastructure](2026-03-04-agent-memory-is-infrastructure.md)

*March 4, 2026* — Every AI coding agent starts fresh. The obvious fix is
"memory." But there's a different problem memory doesn't touch: the project
itself accumulates knowledge that has nothing to do with any single session.
This post argues that agent memory is L2 (runtime cache); what's missing is
L3 (project infrastructure).

**Topics**: context engineering, agent memory, infrastructure, persistence,
team knowledge

---

### [Context as Infrastructure](2026-02-17-context-as-infrastructure.md)

*February 17, 2026* — Where does your AI's knowledge live between sessions?
If the answer is "in a prompt I paste at the start," you are treating context
as a consumable. This post argues for treating it as infrastructure instead:
persistent files, separation of concerns, two-tier storage, progressive
disclosure, and the filesystem as the most mature interface available.

**Topics**: context engineering, infrastructure, progressive disclosure,
persistence, design philosophy

---

### [The Attention Budget: Why Your AI Forgets What You Just Told It](2026-02-03-the-attention-budget.md)

*February 3, 2026* — Every token you send to an AI consumes a finite
resource: the attention budget. Understanding this constraint shaped every
design decision in ctx: hierarchical file structure, explicit budgets,
progressive disclosure, and filesystem-as-index.

**Topics**: attention mechanics, context engineering, progressive disclosure,
ctx primitives, token budgets

---

### [Before Context Windows, We Had Bouncers](2026-02-14-irc-as-context.md)

*February 14, 2026* — IRC is stateless. You disconnect, you vanish. Modern
systems are not much different. This post traces the line from IRC bouncers
to context engineering: stateless protocols require stateful wrappers,
volatile interfaces require durable memory.

**Topics**: context engineering, infrastructure, IRC, persistence,
state continuity

---

### [The Last Question](2026-02-28-the-last-question.md)

*February 28, 2026* — In 1956, Asimov wrote a story about a question that
spans the entire future of the universe. A reading of "The Last Question"
through the lens of persistence, substrate migration, and what it means to
build systems where sessions don't reset.

**Topics**: context continuity, long-lived systems, persistence,
intelligence over time, field notes

---

## Agent Behavior and Design

### [The Dog Ate My Homework: Teaching AI Agents to Read Before They Write](2026-02-25-the-homework-problem.md)

*February 25, 2026* — You wrote the playbook. The agent skipped all of it.
Five sessions, five failure modes, and the discovery that observable
compliance beats perfect compliance.

**Topics**: hooks, agent behavior, context engineering, behavioral design,
testing methodology, compliance monitoring

---

### [Skills That Fight the Platform](2026-02-04-skills-that-fight-the-platform.md)

*February 4, 2026* — When custom skills conflict with system prompt defaults,
the AI has to reconcile contradictory instructions. Five conflict patterns
discovered while building ctx.

**Topics**: context engineering, skill design, system prompts, antipatterns,
AI safety primitives

---

### [The Anatomy of a Skill That Works](2026-02-07-the-anatomy-of-a-skill-that-works.md)

*February 7, 2026* — I had 20 skills. Most were well-intentioned stubs. Then
I rewrote all of them. Seven lessons emerged: quality gates prevent premature
execution, negative triggers are load-bearing, examples set boundaries better
than rules.

**Topics**: skill design, context engineering, quality gates, E/A/R framework,
practical patterns

---

### [You Can't Import Expertise](2026-02-05-you-cant-import-expertise.md)

*February 5, 2026* — I found a well-crafted consolidation skill. Applied my
own E/A/R framework: 70% was noise. This post is about why good skills can't
be copy-pasted, and how to grow them from your project's own drift history.

**Topics**: skill adaptation, E/A/R framework, convention drift, consolidation,
project-specific expertise

---

### [Not Everything Is a Skill](2026-02-08-not-everything-is-a-skill.md)

*February 8, 2026* — I ran an 8-agent codebase audit and got actionable
results. The natural instinct was to wrap the prompt as a skill. Then I
applied my own criteria: it failed all three tests.

**Topics**: skill design, context engineering, automation discipline,
recipes, agent teams

---

### [Defense in Depth: Securing AI Agents](2026-02-09-defense-in-depth-securing-ai-agents.md)

*February 9, 2026* — The security advice was "use CONSTITUTION.md for
guardrails." That is wishful thinking. Five defense layers for unattended
AI agents, each with a bypass, and why the strength is in the combination.

**Topics**: agent security, defense in depth, prompt injection,
autonomous loops, container isolation

---

## Development Practice

### [Code Is Cheap. Judgment Is Not.](2026-02-17-code-is-cheap-judgment-is-not.md)

*February 17, 2026* — AI does not replace workers. It replaces unstructured
effort. Three weeks of building ctx with an AI agent proved it: YOLO mode
showed production is cheap, the 3:1 ratio showed judgment has a cadence.

**Topics**: AI and expertise, context engineering, judgment vs production,
human-AI collaboration, automation discipline

---

### [The 3:1 Ratio](2026-02-17-the-3-1-ratio.md)

*February 17, 2026* — AI makes technical debt worse — not because it writes
bad code, but because it writes code so fast that drift accumulates before
you notice. Three feature sessions, one consolidation session.

**Topics**: consolidation, technical debt, development workflow, convention
drift, code quality

---

### [Refactoring with Intent: Human-Guided Sessions in AI Development](2026-02-01-refactoring-with-intent.md)

*February 1, 2026* — The YOLO mode shipped 14 commands in a week. But
technical debt doesn't send invoices. This is the story of what happened
when we started guiding the AI with intent.

**Topics**: refactoring, code quality, documentation standards, module
decomposition, YOLO versus intentional development

---

### [How Deep Is Too Deep?](2026-02-12-how-deep-is-too-deep.md)

*February 12, 2026* — I kept feeling like I should go deeper into ML theory.
Then I spent a week debugging an agent failure that had nothing to do with
model architecture. When depth compounds and when it doesn't.

**Topics**: AI foundations, abstraction boundaries, agentic systems,
context engineering, failure modes

---

## Agent Workflows

### [Parallel Agents, Merge Debt, and the Myth of Overnight Progress](2026-02-17-parallel-agents-merge-debt-and-the-myth-of-overnight-progress.md)

*February 17, 2026* — You discover agents can run in parallel. So you open
ten terminals. It is not progress: it is merge debt being manufactured in
real time. The five-agent ceiling and why role separation beats file locking.

**Topics**: agent workflows, parallelism, verification, context
engineering, engineering practice

---

### [Parallel Agents with Git Worktrees](2026-02-14-parallel-agents-with-worktrees.md)

*February 14, 2026* — I had 30 open tasks that didn't touch the same files.
Using git worktrees to partition a backlog by file overlap, run 3-4 agents
simultaneously, and merge the results.

**Topics**: agent teams, parallelism, git worktrees, context engineering,
task management

---

## Field Notes and Signals

### [When a System Starts Explaining Itself](2026-02-17-when-a-system-starts-explaining-itself.md)

*February 17, 2026* — Every new substrate begins as a private advantage.
Reality begins when other people start describing it in their own language.
"Better than Adderall" is not praise; it is a diagnostic.

**Topics**: field notes, adoption signals, infrastructure vs tools,
context engineering, substrates

---

### [Why Zensical](2026-02-15-why-zensical.md)

*February 15, 2026* — I needed a static site generator for the journal
system. The instinct was Hugo. But instinct is not analysis. Why zensical
was the right choice: thin dependencies, MkDocs-compatible config, and
zero lock-in.

**Topics**: tooling, static site generators, journal system,
infrastructure decisions, context engineering

---

## Releases

### [ctx v0.6.0: The Integration Release](2026-02-16-ctx-v0.6.0-the-integration-release.md)

*February 16, 2026* — ctx is now a Claude Marketplace plugin. Two commands,
no build step, no shell scripts. v0.6.0 replaces six Bash hook scripts with
compiled Go subcommands and ships 25+ Skills as a plugin.

**Topics**: release, plugin system, Claude Marketplace, distribution,
security hardening

---

### [ctx v0.3.0: The Discipline Release](2026-02-15-ctx-v0.3.0-the-discipline-release.md)

*February 15, 2026* — No new headline feature. Just 35+ documentation and
quality commits against ~15 feature commits. What a release looks like when
the ratio of polish to features is 3:1.

**Topics**: release, skills migration, consolidation, code quality,
E/A/R framework

---

### [ctx v0.2.0: The Archaeology Release](2026-02-01-ctx-v0.2.0-the-archaeology-release.md)

*February 1, 2026* — What if your AI could remember everything? Not just
the current session, but every session. ctx v0.2.0 introduces the recall
and journal systems.

**Topics**: session recall, journal system, structured entries, token budgets,
meta-tools

---

### [Building ctx Using ctx: A Meta-Experiment in AI-Assisted Development](2026-01-27-building-ctx-using-ctx.md)

*January 27, 2026* — What happens when you build a tool designed to give AI
memory, using that very same tool to remember what you're building? This is
the story of ctx.

**Topics**: dogfooding, AI-assisted development, Ralph Loop, session
persistence, architectural decisions

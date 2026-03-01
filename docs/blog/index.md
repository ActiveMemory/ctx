---
title: Blog
icon: lucide/newspaper
---

![ctx](../images/ctx-banner.png)

Stories, insights, and lessons learned from **building** and **using** `ctx`.

---

## Posts

### [The Last Question](2026-02-28-the-last-question.md)

*Jose Alekhinne / February 28, 2026*

In 1956, Asimov wrote a story about a question that spans the entire future
of the universe. Everyone remembers the last line. What they forget is how
many times the question had to be asked (*and why*). The answer isn't about
computation. It's about **continuity**: the system that never forgets is the
one that can finally answer. A reading of "The Last Question" through the lens
of persistence, substrate migration, and what it means to build systems where
sessions don't reset.

**Topics**: context continuity, long-lived systems, persistence,
intelligence over time, field notes

---

### [The Dog Ate My Homework: Teaching AI Agents to Read Before They Write](2026-02-25-the-homework-problem.md)

*Jose Alekhinne / February 25, 2026*

You wrote the playbook. You organized the files. You even put "**CRITICAL,
not optional**" in **bold**. The agent skipped all of it. Five sessions, five
failure modes, and the discovery that **observable compliance** beats
perfect compliance: A relay block that works as a compliance canary,
a feedback loop where agents read their own prior failure analysis,
and the principle that a 30-second read delay is always cheaper than
a decision made without context.

**Topics**: hooks, agent behavior, context engineering, behavioral design,
testing methodology, compliance monitoring

---

### [When a System Starts Explaining Itself](2026-02-17-when-a-system-starts-explaining-itself.md)

*Jose Alekhinne / February 17, 2026*

Every new substrate begins as a private advantage. Reality begins when
other people start describing it in their own language: Not accurately,
not consistently, but **involuntarily**.
"*Better than Adderall*" is not praise; it is a **diagnostic**:
Executive function has been externalized.
This post collects the first external signals from real users and asks what it
means when a system starts explaining itself without its creator present.

**Topics**: field notes, adoption signals, infrastructure vs tools,
context engineering, substrates

---

### [Code Is Cheap. Judgment Is Not.](2026-02-17-code-is-cheap-judgment-is-not.md)

*Jose Alekhinne / February 17, 2026*

AI does not replace workers. It replaces **unstructured effort**. Three weeks
of building ctx with an AI agent proved it: YOLO mode showed production
is cheap, the 3:1 ratio showed judgment has a cadence, the attention
budget showed framing is scarce, and the skill trilogy showed taste is
load-bearing. This post threads every previous blog entry into one
argument, and ends with a personal note on why that's empowering,
not threatening.

**Topics**: AI and expertise, context engineering, judgment vs production,
human-AI collaboration, automation discipline

---

### [Context as Infrastructure](2026-02-17-context-as-infrastructure.md)

*Jose Alekhinne / February 17, 2026*

Where does your AI's knowledge live between sessions? If the answer is
"*in a prompt I paste at the start*," you are treating context as a consumable.
This post argues for treating it as infrastructure instead: persistent files,
separation of concerns, two-tier storage, progressive disclosure, and the
filesystem as the most mature interface available. The architectural companion
to **The Attention Budget** and the second piece of the trilogy completed by
**Code Is Cheap**.

**Topics**: context engineering, infrastructure, progressive disclosure,
persistence, design philosophy

---

### [The 3:1 Ratio](2026-02-17-the-3-1-ratio.md)

*Jose Alekhinne / February 17, 2026*

AI-assisted development makes technical debt worse: **Not** because the AI
writes bad code, **but** because it writes code so fast that drift accumulates
before you notice. Three weeks of git history prove a rhythm: three feature
sessions, one consolidation session. This post shows the **evidence**, the
concrete **drift** that accumulated, and the **decision matrix** for when to clean
up versus leave things alone.

**Topics**: consolidation, technical debt, development workflow, convention
drift, code quality

---

### [Parallel Agents, Merge Debt, and the Myth of Overnight Progress](2026-02-17-parallel-agents-merge-debt-and-the-myth-of-overnight-progress.md)

*Jose Alekhinne / February 17, 2026*

You discover agents can run in parallel. So you open ten terminals.
Then twenty. The fans spin, tokens burn, the screen looks like progress.
It is **not** progress: It is merge debt being manufactured in real time.
This post is about the five-agent ceiling, why role separation beats
file locking, why overnight loops are an illusion, and why the only
metric that matters is time from **idea to verified**, mergeable change.

**Topics**: agent workflows, parallelism, verification, context
engineering, engineering practice

---

### [ctx v0.6.0: The Integration Release](2026-02-16-ctx-v0.6.0-the-integration-release.md)

*Jose Alekhinne / February 16, 2026*

`ctx` is now a **Claude Marketplace** plugin. Two commands, no build step,
no shell scripts. `v0.6.0` replaces six Bash hook scripts with compiled
Go subcommands, ships **25+ Skills** as a plugin served directly from
source, and closes three medium-severity security findings. The version
jumped from `0.3.0` to `0.6.0` because this is the release that turns a
developer tool into a **distributable product**.

**Topics**: release, plugin system, Claude Marketplace, distribution,
security hardening

---

### [ctx v0.3.0: The Discipline Release](2026-02-15-ctx-v0.3.0-the-discipline-release.md)

*Jose Alekhinne / February 15, 2026*

No new headline feature. No architectural pivot. No rewrite. Just **35+
documentation and quality commits against ~15 feature commits**: And
somehow, the tool feels like it grew up overnight. This post is about
what a release looks like when the ratio of polish to features is 3:1.

**Topics**: release, skills migration, consolidation, code quality,
E/A/R framework

---

### [Why Zensical](2026-02-15-why-zensical.md)

*Jose Alekhinne / February 15, 2026*

I needed a static site generator for the journal system, and the instinct
was Hugo:  same language as ctx, fast, well-established. But instinct is
not analysis. The journal entries are standard Markdown with YAML
frontmatter. That is it. **No JSX, no shortcodes, no custom templating**.
This post is about why `zensical` (*a pre-1.0 tool from the Material for
MkDocs team*) was the right choice: thin dependencies, MkDocs-compatible
config, 4-5x faster incremental rebuilds, and zero lock-in.

**Topics**: tooling, static site generators, journal system,
infrastructure decisions, context engineering

---

### [Parallel Agents with Git Worktrees](2026-02-14-parallel-agents-with-worktrees.md)

*Jose Alekhinne / February 14, 2026*

I had 30 open tasks, and most of them didn't touch the same files. Running
one agent at a time meant serial execution on work that was fundamentally
parallel. This post is about using **git worktrees** to **partition** a backlog by
file overlap, run 3-4 agents simultaneously, and merge the results: The
same **attention budget** principle applied to **execution** instead of context.

**Topics**: agent teams, parallelism, git worktrees, context engineering,
task management

---

### [Before Context Windows, We Had Bouncers](2026-02-14-irc-as-context.md)

*Jose Alekhinne / February 14, 2026*

IRC is **stateless**. You disconnect, you vanish. Modern systems are not much
different: close the tab, lose the scrollback, open a new LLM session, start
from zero. This post traces the line from IRC bouncers like ZNC to context
engineering: stateless protocols require stateful wrappers, volatile interfaces
require durable memory. Before context windows, we had **bouncers**. Before AI
memory files, we had buffers.

**Topics**: context engineering, infrastructure, IRC, persistence,
state continuity

---

### [How Deep Is Too Deep?](2026-02-12-how-deep-is-too-deep.md)

*Jose Alekhinne / February 12, 2026*

I kept feeling like I should go deeper into **ML theory**. Then I spent a week
debugging an agent failure that had nothing to do with model architecture and
everything to do with knowing which abstraction was leaking. This post is
about when depth **compounds** and when it doesn't: why the useful understanding
lives one or two layers below where you work, **not** at the bottom of the stack.

**Topics**: AI foundations, abstraction boundaries, agentic systems,
context engineering, failure modes

---

### [Defense in Depth: Securing AI Agents](2026-02-09-defense-in-depth-securing-ai-agents.md)

*Jose Alekhinne / February 9, 2026*

I was writing the autonomous loops recipe and realized the security
advice was "*use `CONSTITUTION.md` for guardrails*." Then I read that
sentence back and realized: That is **wishful thinking**. This post
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

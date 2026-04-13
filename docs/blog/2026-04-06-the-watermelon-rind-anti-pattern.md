---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "The Watermelon-Rind Anti-Pattern:
  Why Smarter Tools Make Shallower Agents"
date: 2026-04-06
author: Jose Alekhinne
reviewed_and_finalized: true
topics:
  - architecture
  - code intelligence
  - agent behavior
  - design patterns
  - field notes
---

# The Watermelon-Rind Anti-Pattern

## Why Smarter Tools Make Shallower Agents

![ctx](../images/ctx-banner.png)

**Give an agent a graph query tool, and it will tell you everything
about your codebase except what actually matters.**

*Jose Alekhinne / April 6, 2026*

## A Turkish Proverb Walks into a Codebase

There's a Turkish idiom: *esegin aklina karpuz kabugu sokmak*
(*literally, "to put watermelon rind into a donkey's mind." It means
to plant an idea in someone's head that they wouldn't have come up
with on their own*) usually one that leads them astray.

In English, let's call this a "**watermelon metric**": a project management term 
for something that's green on the outside and red on the inside: all
dashboards passing, reality crumbling.

Both halves of this metaphor showed up in a single experiment. And
the result changed how we design architecture analysis in
[ctx][ctx].

[ctx]: https://ctx.ist

---

## The Experiment

We ran three sessions analyzing the same large codebase (~34,000
symbols) using the same architecture skill, varying only what tools
the agent had access to.

| Session | Tools Available   | Output (lines) | Character                 |
|---------|-------------------|----------------|---------------------------|
| 1       | None (MCP broken) | 5,866          | Deep, intimate            |
| 2       | Full graph MCP    | 1,124          | Structural, correct       |
| 3       | Enrichment pass   | +verified data | Additive, not restorative |

Session 1 was an accident. The MCP server that provides code
intelligence queries was broken, so the agent couldn't ask the
graph anything. It had to read code. Line by line. File by file.

It produced 5,866 lines of architecture analysis: per-controller
data flows, scale math, startup sequences, timeout defaults, edge
cases that only surface when you actually look at the
implementation.

Session 2 had working tools. Same skill, same codebase. The agent
produced 1,124 lines (**5.2x less**). Structurally correct. Valid
symbol references. Proper call chains.

And **hollow**.

---

## The Rind

The Session 2 output was a **watermelon rind**: the right shape, the
right color, the right texture on the outside. But the substance
(*the operational details, the defaults nobody documents, the
scale math that tells you when a component will fall over*) was
missing.

Not wrong. Not broken. Just... thin.

The agent had answered every question correctly. The problem was
that it never discovered the questions it should have asked. When
you can query a graph for "*what calls this function?*", you don't
stumble into the retry loop that silently swallows errors three
layers down. When you can ask for the dependency tree, you don't
notice that two packages share a mutable state through a global
variable that isn't in any interface.

**The tool answered the question asked but prevented the discovery of
answers to questions never asked.**

Here's what that looks like concretely: the graph tells you that
`ReconcileDeployment` calls `SyncPods`. It does not tell you that
`SyncPods` retries three times with exponential backoff, silently
drops errors after timeout, and resets a package-level counter that
another goroutine reads without a lock. The call chain is correct.

The operational reality is invisible.

---

## The Donkey's Idea

This is where the Turkish proverb earns its place: The graph tool
is the "*karpuz kabugu*" (*the watermelon rind placed into the
agent's mind*). 

Before the tool existed, the agent had no choice but to read deeply. 
With the tool available, a new idea appears: *why read 500 lines of code when 
I can query the call graph?*

The agent isn't lazy. **It's rational**. 

Graph queries are faster, more reliable, and produce verifiably correct output. 
The agent is optimizing. It's satisficing (*finding answers that are good
enough*), instead of maximizing (*finding everything there is to
know*).

**Satisficing produces watermelon rinds**.

---

## The Two-Pass Compiler

Session 3 taught us that you can't fix shallow analysis by adding
more tools after the fact. The enrichment pass added verified graph
data (*blast radius numbers, registration sites, execution flow
confirmation*) but it couldn't recover the intimate code knowledge
that Session 1 had produced through sheer necessity.

You can't enrich your way out of a depth deficit.

So we redesigned. Instead of one skill with optional tools, we
built a **two-pass compiler for architecture understanding**:

**Pass 1: Semantic parsing.** The `/ctx-architecture` skill
deliberately has no access to graph query tools. The agent must read
code, build mental models, and produce architecture artifacts
through human-style comprehension. Constraint is the feature.

**Pass 2: Static analysis.** The `/ctx-architecture-enrich`
skill takes Pass 1 output as input and runs comprehensive
verification through code intelligence: blast radius analysis,
registration site discovery, execution flow tracing, domain
clustering comparison. It extends and verifies, but it doesn't
replace.

The key insight: **these must be separate skills with separate tool
permissions.** If you give the agent graph tools during Pass 1, it
will use them. The "*karpuz kabugu*" will be in its mind. The only way
to prevent satisficing is to remove the option.

---

## The Principle

We call this **constraint-as-feature**: deliberately withholding
capabilities to force deeper engagement.

It sounds paradoxical. You built sophisticated code intelligence
tools and then... forbid the agent from using them? During the most
important phase?

Yes. Because the tools don't make the agent smarter. They make it
faster. And faster, in architecture analysis, is the enemy of deep.

What's actually happening is subtler: tools reduce the agent's
search space. A graph query collapses thousands of possible
observations into one precise answer. That's efficient for known
questions. But architecture understanding depends on *unknown
unknowns*: and you only find those by wandering through code
with nothing to shortcut the journey.

The constraint forces the agent into a mode of operation that
produces better output than any amount of tooling can achieve. The
limitation *is* the capability.

---

## When Does This Apply?

Not always. The watermelon-rind antipattern is specific to
**exploratory analysis**: tasks where the value comes from
discovering unknowns, not from answering known questions.

Graph tools are excellent for:

* **Verification**: "Does X actually call Y?" (binary question,
  precise answer)
* **Impact analysis**: "What breaks if I change Z?" (bounded scope,
  enumerable results)
* **Navigation**: "Where is this interface implemented?" (lookup,
  not analysis)

Graph tools produce watermelon rinds when:

* **The goal is understanding**, not answering
* **The unknowns are unknown**: you don't know what to ask
* **Depth matters more than breadth**: operational details,
  edge cases, implicit coupling

The two-pass approach preserves both: deep reading first, tool
verification second.

---

## Takeaway

The two-pass approach is the slowest way to analyze a codebase. It
is also the only way that produces both depth and accuracy. We
accept the cost because architecture analysis is not a speed game:
it is a coverage game.

**Esegin aklina karpuz kabugu sokma!**

(*don't put the watermelon rind to a donkey's mind*)

If the agent never struggles, it never discovers. And if it never
discovers, you are not doing architecture; you are doing
autocomplete.

---

*This post is part of the [ctx field notes][blog] series,
documenting what we learn building persistent context
infrastructure for AI coding sessions.*

[blog]: index.md

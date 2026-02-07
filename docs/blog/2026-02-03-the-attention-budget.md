---
title: "The Attention Budget: Why Your AI Forgets What You Just Told It"
date: 2026-02-03
author: Jose Alekhinne
topics:
  - attention mechanics
  - context engineering
  - progressive disclosure
  - ctx primitives
  - token budgets
---

# The Attention Budget

![ctx](../images/ctx-banner.png)

## Why Your AI Forgets What You Just Told It

*Jose Alekhinne / 2026-02-03*

!!! question "Ever wondered why AI gets worse the longer you talk?"
    You paste a 2000-line file, explain the bug in detail, provide three
    examples...

    ...and the AI still suggests a fix that ignores half of what you said.

This isn't a bug. It is **physics**.

Understanding that single fact shaped every design decision behind `ctx`.

## The Finite Resource Nobody Talks About

Here's something that took me too long to internalize: **context is not free**.

Every token you send to an AI model consumes a finite resource I call the
*attention budget*. 

The model doesn't just read tokens; it forms relationships
between them: For `n` tokens, that's roughly `n^2` relationships.
Double the context, and the computation quadruples.

But the more important constraint isn't cost: It's **attention density**.

!!! info "Attention Density"
    **Attention density** is how much focus each token receives relative to all
    other tokens in the context window.

As context grows, attention density drops: Each token gets a smaller slice
of the model's focus. Nothing is ignored; but everything becomes blurrier.

Think of it like a **flashlight**: In a small room, it illuminates everything
clearly. In a warehouse, it becomes a dim glow that barely reaches the corners.

This is why `ctx agent` has an explicit `--budget` flag:

```bash
ctx agent --budget 4000 # Force prioritization
ctx agent --budget 8000 # More context, lower attention density
```

The budget isn't just about cost. It's about **preserving signal**.

## The Middle Gets Lost

This one surprised me.

Research shows that transformer-based models tend to attend more strongly to
the **beginning** and **end** of a context window than to its middle (*a 
phenomenon often called "lost in the middle"*). 

**Positional anchors matter, and the middle has fewer of them**.

In practice, this means that information placed "*somewhere in the middle*"
is statistically less salient, even if it's important.

`ctx` orders context files by **logical progression**—what the agent needs to
know before it can understand the next thing:

1. `CONSTITUTION.md`: Constraints before action
2. `TASKS.md`: Focus before patterns
3. `CONVENTIONS.md`: How to write before where to write
4. `ARCHITECTURE.md`: Structure before history
5. `DECISIONS.md`: Past choices before gotchas
6. `LEARNINGS.md`: Lessons before terminology
7. `GLOSSARY.md`: Reference material
8. `DRIFT.md`: Staleness indicators
9. `AGENT_PLAYBOOK.md`: Meta instructions last

This ordering is about logical dependencies, not attention engineering.
But it happens to be **attention-friendly** too:

The files that matter most—**CONSTITUTION**, **TASKS**, **CONVENTIONS**—land
at the **beginning** of the context window, where attention is strongest.

Reference material like **GLOSSARY** and **DRIFT** sit in the middle, where
lower salience is acceptable.

And **AGENT_PLAYBOOK**—the operating manual for the context system itself—sits
at the **end**, also outside the "lost in the middle" zone. The agent reads
*what* to work with before learning *how* the system works.

This is `ctx`'s first primitive: **hierarchical importance**.
Not all context is equal.

## `ctx` Primitives

`ctx` is built on four primitives that directly address the **attention
budget** problem.

### Primitive 1: Separation of Concerns

Instead of a single mega-document, `ctx` uses **separate files for separate
purposes**:

| File              | Purpose               | Load When                 |
|-------------------|-----------------------|---------------------------|
| CONSTITUTION.md   | Inviolable rules      | Always                    |
| TASKS.md          | Current work          | Session start             |
| CONVENTIONS.md    | How to write code     | Before coding             |
| ARCHITECTURE.md   | System structure      | Before making changes     |
| DECISIONS.md      | Architectural choices | When questioning approach |
| LEARNINGS.md      | Gotchas               | When stuck                |
| GLOSSARY.md       | Domain terminology    | When clarifying terms     |
| DRIFT.md          | Staleness indicators  | During maintenance        |
| AGENT_PLAYBOOK.md | Operating manual      | Session start             |
| sessions/         | Deep history          | On demand                 |
| journal/          | Session journal       | On demand                 |

This isn't just "*organization*": It is **progressive disclosure**.

Load only what's relevant to the task at hand. Preserve attention density.

### Primitive 2: Explicit Budgets

The `--budget` flag forces a choice:

```bash
ctx agent --budget 4000
```

Here is a sample allocation:

```
Constitution: ~200 tokens (never truncated)
Tasks: ~500 tokens (current phase)
Conventions: ~800 tokens (key patterns)
Recent decisions: ~400 tokens (last 3)
…budget exhausted, stop loading
```

The constraint is the feature: It enforces ruthless prioritization.

### Primitive 3: Indexes Over Full Content

`DECISIONS.md` and `LEARNINGS.md` both include index sections:

```markdown
<!-- INDEX:START -->
| Date       | Decision                            |
|------------|-------------------------------------|
| 2026-01-15 | Use PostgreSQL for primary database |
| 2026-01-20 | Adopt Cobra for CLI framework       |
<!-- INDEX:END -->
```

An AI agent can scan ~50 tokens of index and decide which 
200-token entries are worth loading.

This is **just-in-time context**.

References are cheaper than full text.

### Primitive 4: Filesystem as Navigation

`ctx` uses the filesystem itself as a context structure:

```
.context/
├── CONSTITUTION.md
├── TASKS.md
├── sessions/
│   ├── 2026-01-15-*.md
│   └── 2026-01-20-*.md
└── archive/
    └── tasks-2026-01.md
```

The AI doesn't need every session loaded;
it needs to know **where to look**.

```bash
ls .context/sessions/
cat .context/sessions/2026-01-20-auth-discussion.md
```

File names, timestamps, and directories encode relevance.

**Navigation is cheaper than loading**.

## Progressive Disclosure in Practice

The naive approach to context is dumping everything upfront:

> "Here's my entire codebase, all my documentation, every decision I've ever
> made—now help me fix this typo."

This is an **antipattern**.

!!! warning "Antipattern: Context Hoarding"
    Dumping everything "*just in case*" will silently destroy the **attention 
    density**.

`ctx` takes the opposite approach:

```bash
ctx status                      # Quick overview (~100 tokens)
ctx agent --budget 4000         # Typical session
cat .context/sessions/...       # Deep dive when needed
```

| Command                   | Tokens | Use Case      |
|---------------------------|--------|---------------|
| `ctx status`              | ~100   | Human glance  |
| `ctx agent --budget 4000` | 4000   | Normal work   |
| `ctx agent --budget 8000` | 8000   | Complex tasks |
| Full session read         | 10000+ | Investigation |

Summaries first. Details on demand.

## Quality Over Quantity

Here's the counterintuitive part: **more context can make AI worse**.

Extra tokens add noise, not clarity:

* Hallucinated connections increase.
* Signal per token drops.

The goal isn't maximum context. It's **maximum signal per token**.

This principle drives several `ctx` features:

| Design Choice    | Rationale                 |
|------------------|---------------------------|
| Separate files   | Load only what's relevant |
| Explicit budgets | Enforce prioritization    |
| Index sections   | Cheap scanning            |
| Task archiving   | Keep active context clean |
| `ctx compact`    | Periodic noise reduction  |

Completed work isn't deleted: It is moved somewhere cold.

## Designing for Degradation

Here is the uncomfortable truth:

**Context will degrade.**

Long sessions stretch attention thin. Important details fade.

The real question isn't how to prevent degradation, 
but how to **design for it**.

`ctx`'s answer is **persistence**:

**Persist early. Persist often.**

The `AGENT_PLAYBOOK` asks:

> "If this session ended right now, would the next one know what happened?"

Capture learnings as they occur:

```bash
ctx add learning "JWT tokens require explicit cache invalidation" \
  --context "Debugging auth failures" \
  --lesson "Token refresh doesn't clear old tokens" \
  --application "Always invalidate cache on refresh"
```

**Structure beats prose**: Bullet points survive compression.

Headings remain scannable. Tables pack density.

And above all: **single source of truth**.

Reference decisions; don't duplicate them.

## The `ctx` Philosophy

!!! info "Context as Infrastructure"
    `ctx` is not a prompt: It is **infrastructure**.

    `ctx` creates versioned files that persist across time and sessions.

The attention budget is fixed. You can't expand it.
But you can **spend it wisely**:

1. Hierarchical importance
2. Progressive disclosure
3. Explicit budgets
4. Indexes over full content
5. Filesystem as structure

This is why `ctx` exists: **not** to cram more context into AI sessions,
**but** to curate the *right* context for each moment.

## The Mental Model

I now approach every AI interaction with one question:

> "Given a fixed attention budget, what's the highest-signal thing I can load?"

Not "*how do I explain everything*," but "*what's the minimum that matters*."

That shift (*from abundance to curation*) is the difference between
frustrating sessions and **productive** ones.

----

**Spend your tokens wisely**.

Your AI will thank you.

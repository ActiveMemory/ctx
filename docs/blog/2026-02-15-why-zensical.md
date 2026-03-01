---
title: "Version Numbers Are Lagging Indicators"
date: 2026-02-15
author: Jose Alekhinne
reviewed_and_finalized: true
topics:
  - tooling decisions
  - static site generators
  - infrastructure thinking
  - journal system
  - context engineering
---

# Version Numbers Are Lagging Indicators

![ctx](../images/ctx-banner.png)

## Why ctx's Journal Site Runs on a v0.0.21 Tool

*Jose Alekhinne / 2026-02-15*

!!! question "Would You Ship Production Infrastructure on a v0.0.21 Dependency?"
    Most engineers wouldn't. Version numbers signal maturity. Pre-1.0
    means unstable API, missing features, risk.

    But version numbers tell you where a project **has been**.
    They say nothing about where it's **going**.

I just bet `ctx`'s entire journal site on a tool that hasn't
hit `v0.1.0`. 

Here's why I'd do it again.

## The Problem

When [v0.2.0][v020] shipped the journal system, the pipeline was clear:

* **Export** sessions to Markdown; 
* **Enrich** them with YAML frontmatter; 
* And render them into something **browsable**. 

The first two steps were solved; the third needed a tool.

[v020]: 2026-02-01-ctx-v0.2.0-the-archaeology-release.md

The journal entries are standard Markdown with YAML frontmatter, tables,
and fenced code blocks. That is the entire format: 

* No JSX;
* No shortcodes;
* No custom templating. 

Just Markdown rendered well.

The requirements are modest:

* Read a configuration file (*such as `mkdocs.yml`*);
* Render Markdown with extensions (*admonitions, tabs, tables*);
* Search;
* Handle 100+ files without choking on incremental rebuilds;
* Look good out of the box;
* **Not lock me in**.

The obvious candidates were as follows:

| Tool                  | Language | Strengths                                          | Pain Points                                                           |
|-----------------------|----------|----------------------------------------------------|-----------------------------------------------------------------------|
| **Hugo**              | Go       | Blazing fast, mature                               | Templating is painful; Go templates fight you on anything non-trivial |
| **Astro**             | JS/TS    | Modern, flexible                                   | JS ecosystem overhead; overkill for a docs site                       |
| **MkDocs + Material** | Python   | Beautiful defaults, massive community (22k+ stars) | Slow incremental rebuilds on large sites; limited extensibility model |
| **Zensical**          | Python   | Built to fix MkDocs' limits; 4-5x faster rebuilds  | v0.0.21; module system not yet shipped                                |

The instinct was Hugo. Same language as `ctx`. Fast. Well-established.

But **instinct is not analysis**. I picked the one with the lowest version
number.

---

## The Evaluation

Here is what I actually evaluated, in order:

### 1. The Team

[Zensical][zensical] is built by [squidfunk](https://github.com/squidfunk):
The same person behind [Material for MkDocs][material], the most popular
MkDocs theme with 22,000+ stars. It powers documentation sites for projects
across every language and framework.

[zensical]: https://github.com/zensical/zensical
[material]: https://github.com/squidfunk/mkdocs-material

* This is not someone learning how to build static site generators.
* This is someone who spent **years** understanding exactly where MkDocs
  breaks and decided to fix it from the ground up.

They did not build zensical because MkDocs was bad: 
They built it because MkDocs hit a ceiling:

* **Incremental rebuilds**: 4-5x faster during serve. When you have
  hundreds of journal entries and you edit one, the difference between
  "rebuild everything" and "rebuild this page" is the difference
  between a usable workflow and a frustrating one.

* **Large site performance**: Specifically designed for tens of
  thousands of pages. The journal grows with every session. A tool
  that slows down as content accumulates is a tool you will eventually
  replace.

**A proven team starting fresh is more predictable than an
unproven team at v3.0.**

### 2. The Architecture

Zensical is investing in a Rust-based Markdown parser with CommonMark
support. That signals something about the team's priorities:

**Performance foundations first; features second**.

`ctx`'s journal will grow: 

* Every exported session adds files.
* Every enrichment pass adds metadata. 

Choosing a tool that gets *slower* as you add content means choosing to
migrate later.

Choosing one **built for scale** means the decision holds.

### 3. The Migration Path

Zensical reads `mkdocs.yml` natively. If it doesn't work out,
I can move back to MkDocs + Material with zero content changes:

* The Markdown is standard; 
* The frontmatter is standard; 
* The configuration is compatible.

This is the [infrastructure pattern][irc] again: The same way `ZNC`
decouples presence from the client, `zensical` decouples rendering from
the generator: 

* The Markdown is yours. 
* The frontmatter is standard YAML. 
* The configuration is MkDocs-compatible.

You are not locked into anything except your own content.

[irc]: 2026-02-14-irc-as-context.md

**No lock-in** is **not** a feature: **It's a design philosophy**: 

It's the same reason `ctx` uses plain Markdown files in `.context/`
instead of a database: **the format should outlive the tool**.

!!! tip "Lock-in Is the Real Risk, Not Version Numbers"
    A mature tool with a proprietary format is riskier than a
    young tool with a standard one. Version numbers measure
    time invested. Portability measures respect for the user.

### 4. The Dependency Tree

Here is what `pip install zensical` actually pulls in:

* click
* Markdown
* Pygments
* pymdown-extensions
* PyYAML

Only **five** dependencies. All well-known. No framework bloat. No
bundler. No transpiler. No `node_modules` black hole.

3k GitHub stars at `v0.0.21` is a strong early traction for a `pre-1.0`
project. 

The dependency tree is thin: No bloat.

### 5. The Fit

This is the same principle behind the [attention budget][attention]:
**do not overfit the tool to hypothetical requirements**. The right
amount of capability is the minimum needed for the current task.

[attention]: 2026-02-03-the-attention-budget.md

Hugo is a powerful static site generator. It is also a powerful
templating engine, a powerful asset pipeline, and a powerful taxonomy
system. For rendering Markdown journals, that power is overhead:

It is the complexity you pay for but never use.

`ctx`'s journal files are standard Markdown with YAML frontmatter,
tables, and fenced code blocks. That is exactly the sweet spot
Zensical inherits from Material for MkDocs:

* **No** custom plugins needed;
* **No** special syntax; 
* **No** templating gymnastics.

The requirements match the capabilities: **Not** the capabilities
that are promised, **but** the ones that exist today, at `v0.0.21`.

---

## The Caveat

It would be dishonest not to mention what's missing.

The module system for third-party extensions opens in early 2026.

If `ctx` ever needs custom plugins (*for example, auto-linking session IDs,
rendering special journal metadata, etc.*) that infrastructure isn't
there yet.

The installation experience is rough: 

We discovered this firsthand: `pip install zensical` often fails on MacOS 
(*system Python stubs, Homebrew's PEP 668 restrictions*). The answer is
[pipx](https://pipx.pypa.io/), which creates an isolated environment
with the correct Python version automatically. 

That kind of friction is typical for young Python tooling, and it is documented
in the [Getting Started guide](../home/getting-started.md#journal-site).

And `3,000` stars at `v0.0.21` is **strong** early traction, but it's
still **early**: The community is small. When something breaks, you're reading 
source  code, *not* documentation.

**These are real costs. I chose to pay them because the
alternative costs are higher.**

For example:

* Hugo's templating pain would cost me time on every site change.
* Astro's JS ecosystem would add complexity I don't need. 
* MkDocs would work today but hit scaling walls tomorrow. 

Zensical's costs are **front-loaded** and **shrinking**. 

The others compound.

---

## The Evaluation Framework

For anyone facing a similar choice, here is the framework that
emerged:

| Signal                | What It Tells You                      | Weight |
|-----------------------|----------------------------------------|--------|
| **Team track record** | Whether the architecture will be sound | High   |
| **Migration path**    | Whether you can leave if wrong         | High   |
| **Current fit**       | Whether it solves *your* problem today | High   |
| **Dependency tree**   | How much complexity you're inheriting  | Medium |
| **Version number**    | How long the project has existed       | Low    |
| **Star count**        | Community interest (not quality)       | Low    |
| **Feature list**      | What's possible (not what you need)    | Low    |

The bottom three are the metrics most engineers optimize for.

The top four are the ones that predict whether you'll still be
happy with the choice in a year.

!!! warning "Features You Don't Need Are Not Free"
    Every feature in a dependency is code you inherit but don't
    control. 

    A tool with 200 features where you use 5 means
    195 features worth of surface area for bugs, breaking changes,
    and security issues that have nothing to do with your use case.

    **Fit is the inverse of feature count.**

---

## The Broader Pattern

This is part of a theme I keep encountering in this project:

**Leading indicators beat lagging indicators.**

| Domain                  | Lagging Indicator              | Leading Indicator                          |
|-------------------------|--------------------------------|--------------------------------------------|
| **Tooling**             | Version number, star count     | Team track record, architecture            |
| **Code quality**        | Test coverage percentage       | Whether tests catch real bugs              |
| **Context persistence** | Number of files in `.context/` | Whether the AI makes fewer mistakes        |
| **Skills**              | Number of skills created       | Whether each skill fires at the right time |
| **Consolidation**       | Lines of code refactored       | Whether drift stops accumulating           |

Version numbers, star counts, coverage percentages, file counts...

...these are all measures of *effort expended*. 

They say **nothing** about *value delivered*.

The question is never "*how mature is this tool?*" 

The question is "**does this tool's trajectory intersect with my needs?**"

Zensical's trajectory: 

* A **proven team** fixing known problems, 
* in a **proven architecture*, 
* with a **standard format**,
* and **no lock-in**.

`ctx`'s needs: 

Tender standard Markdown into a browsable site, **at scale**, 
**without complexity**.

The intersection is **clean**; the version number is **noise**.

This is the same kind of decision that shows up throughout `ctx`:

* [Skills that fight the platform][fight] taught that the best
  integration extends existing behavior, not replaces it.
* [You can't import expertise][import] taught that tools should
  grow from your project's actual needs, not from feature checklists.
* [Context as infrastructure][infra] argues that the format should
  outlive the tool; and, `zensical` honors that principle by reading
  standard Markdown and standard MkDocs configuration.

[fight]: 2026-02-04-skills-that-fight-the-platform.md
[import]: 2026-02-05-you-cant-import-expertise.md
[infra]: 2026-02-17-context-as-infrastructure.md

---

!!! quote "**If You Remember One Thing From This Post...**"
    **Version numbers measure where a project has been.**

    **The team and the architecture tell you where it's going.**

    A `v0.0.21` tool built by the right team on the right
    foundations is a safer bet than a `v5.0` tool that doesn't
    fit your problem.

    **Bet on trajectories, not timestamps.**

---

*This post started as an evaluation note in `ideas/` and a separate
decision log. The analysis held up. The two merged into one.
The meta continues.*

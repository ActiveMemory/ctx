---
title: "Why Zensical"
date: 2026-02-15
author: Jose Alekhinne
topics:
  - tooling
  - static site generators
  - journal system
  - infrastructure decisions
  - context engineering
---

# Why Zensical

![ctx](../images/ctx-banner.png)

## Choosing the Right Renderer for the Wrong Reasons Is Still Wrong

*Jose Alekhinne / February 15, 2026*

!!! question "How do you pick a static site generator for a Go project?"
    You pick the one that respects what you already have.

When [v0.2.0][v020] shipped the journal system, the pipeline was clear:
export sessions to Markdown, enrich them with YAML frontmatter, and
render them into something browsable. The first two steps were solved.
The third needed a tool.

[v020]: 2026-02-01-ctx-v0.2.0-the-archaeology-release.md

The obvious candidates were:

| Generator     | Language | Strength                 | Problem for ctx                     |
|---------------|----------|--------------------------|-------------------------------------|
| Hugo          | Go       | Blazing fast             | Templating is painful               |
| Astro         | JS       | Modern, flexible         | JS ecosystem overkill for docs      |
| MkDocs        | Python   | Great defaults           | Slow on large sites                 |
| Zensical      | Python   | MkDocs DNA, built for scale | Pre-1.0                          |

The instinct was Hugo. Same language as `ctx`. Fast. Well-established.

But instinct is not analysis.

---

## What ctx Actually Needs

The journal entries are standard Markdown with YAML frontmatter, tables,
and fenced code blocks. That is the entire format. No JSX. No shortcodes.
No custom templating. Just Markdown rendered well.

This is the same principle behind the [attention budget][attention]:
**do not overfit the tool to hypothetical requirements**. The right
amount of capability is the minimum needed for the current task.

[attention]: 2026-02-03-the-attention-budget.md

Hugo is a powerful static site generator. It is also a powerful
templating engine, a powerful asset pipeline, and a powerful taxonomy
system. For rendering Markdown journals, that power is overhead:
complexity you pay for but never use.

What `ctx` needs is:

* Read `mkdocs.yml`
* Render Markdown with extensions (admonitions, tabs, tables)
* Search
* Fast incremental rebuilds during `--serve`

That is exactly what [zensical][zensical] does.

[zensical]: https://github.com/zensical/zensical

---

## The squidfunk Lineage

Zensical comes from the same team that built
[Material for MkDocs][material] -- the most popular MkDocs theme with
22k+ stars. These people know docs-as-code deeply. They did not build
zensical because MkDocs was bad. They built it because MkDocs hit a
ceiling:

[material]: https://github.com/squidfunk/mkdocs-material

* **Incremental rebuilds**: 4-5x faster during serve. When you have
  hundreds of journal entries and you edit one, the difference between
  "rebuild everything" and "rebuild this page" is the difference
  between a usable workflow and a frustrating one.

* **Large site performance**: Specifically designed for tens of
  thousands of pages. The journal grows with every session. A tool
  that slows down as content accumulates is a tool you will eventually
  replace.

* **Migration-friendly**: Reads `mkdocs.yml` natively. If `ctx` ever
  needed to pivot back to MkDocs or forward to something else, there
  is no lock-in. The configuration is portable.

This is the [infrastructure pattern][irc] again. The same way ZNC
decouples presence from the client, zensical decouples rendering from
the generator. The Markdown is yours. The frontmatter is standard
YAML. The configuration is MkDocs-compatible. You are not locked into
anything except your own content.

[irc]: 2026-02-14-irc-as-context.md

---

## The Caveat: Pre-1.0

Zensical is `v0.0.21`. Genuinely early.

The module system for third-party extensions opens early 2026. If `ctx`
ever needs custom plugins -- auto-linking to session IDs, special
journal metadata rendering -- they would need to wait for that or work
around it.

The install UX is also rough. We discovered this firsthand:
`pip install zensical` often fails on macOS (system Python stubs,
Homebrew's PEP 668 restrictions). The answer is
[pipx](https://pipx.pypa.io/), which creates an isolated environment
with the correct Python version automatically. That friction is typical
for young Python tooling, and it is documented in the
[Getting Started guide](../getting-started.md#journal-site).

But for `ctx`'s use case -- rendering journal Markdown into a browsable
site -- it is almost ideal. The journals are exactly the sweet spot
zensical inherits from Material for MkDocs: standard Markdown, YAML
frontmatter, fenced code, tables. No custom plugins needed.

---

## A Rust-Shaped Future

Zensical is investing in a Rust-based Markdown parser with CommonMark
support. That signals something about the team's priorities:
**performance foundations first, features second**.

3k GitHub stars at `v0.0.21` is strong early traction for a pre-1.0
project. The dependency tree is thin: click, Markdown, Pygments,
pymdown-extensions, PyYAML. No bloat.

A pre-1.0 tool with a proven team, thin dependencies, and a Rust
performance roadmap is a better bet than a mature tool that does not
fit the use case.

---

## The Decision

The alternatives were real. Hugo is fast but its templating is painful
for pure documentation. Astro is modern but brings the entire JS
ecosystem for a problem that does not need it. MkDocs works but
struggles at scale -- and the journal will only grow.

For a Go project that generates Markdown journals, zensical hits
the right balance: Python-based but thin, great defaults, and the
team behind it has a decade of docs-as-code experience.

This is the same kind of decision that shows up throughout `ctx`:

* [Skills that fight the platform][fight] taught that the best
  integration extends existing behavior, not replaces it.
* [You can't import expertise][import] taught that tools should
  grow from your project's actual needs, not from feature checklists.

[fight]: 2026-02-04-skills-that-fight-the-platform.md
[import]: 2026-02-05-you-cant-import-expertise.md

Zensical does not try to be everything. It tries to render Markdown
well, fast, and at scale. That is what `ctx` needs. Nothing more.

---

!!! quote "If you remember one thing from this post..."
    Pick tools that respect what you already have.

    The journal entries are Markdown. The configuration is MkDocs-compatible.
    The content is portable.

    The best infrastructure decision is the one that does not
    force you to change your content to fit the tool.

---

*This post started as an evaluation note in `ideas/`. The analysis
held up. The note became a post. The meta continues.*

---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "We Broke the 3:1 Rule"
date: 2026-03-23
author: Jose Alekhinne
reviewed_and_finalized: true
topics:
  - consolidation
  - technical debt
  - development workflow
  - convention drift
  - field notes
---

# We Broke the 3:1 Rule

![ctx](../images/ctx-banner.png)

**The best time to consolidate was after every third session.
The second best time is now.**

*Jose Alekhinne / March 23, 2026*

The rule was simple: three feature sessions, then one
consolidation session. 

[The Architecture Release][arch-post] shows the result: 
This post shows the cost.

[arch-post]: 2026-03-23-ctx-v0.8.0-the-architecture-release.md

## The Rule We Wrote

In [The 3:1 Ratio][ratio-post], I documented a rhythm that worked
during `ctx`'s first month: **three feature sessions, then one
consolidation session**. The evidence was clear. The rule was simple.

The math checked out.

And then **we ignored it for five weeks**.

[ratio-post]: 2026-02-17-the-3-1-ratio.md
[broken-window]: https://blog.codinghorror.com/the-broken-window-theory/

---

## What Happened

After `v0.6.0` shipped on February 16, the feature pipeline was
irresistible. The MCP server spec was ready. The memory bridge
design was done. Webhook notifications had been deferred twice.
The VS Code extension needed 15 new commands. The `sysinfo` package
was overdue...

Each feature was important. Each feature was "*just one more
session.*" Each feature pushed the consolidation session one day
further out.

The git history tells the story in two numbers:

| Phase             | Dates          | Commits | Duration |
|-------------------|----------------|--------:|----------|
| Feature run       | Feb 16 - Mar 5 |     198 | 17 days  |
| Consolidation run | Mar 5 - Mar 23 |     181 | 18 days  |

198 feature commits before a single consolidation commit. If the
3:1 rule says consolidate every 4th session, we consolidated after
the **66th**.

!!! danger "The Actual Ratio"
    The ratio wasn't 3:1. It was **1:1**. 

    We spent as much time cleaning up as we did building. 

    The consolidation run took 18 days:
    **longer than the feature run itself**.

---

## What Compounded

The [3:1 post][ratio-post] warned about compounding. Here is what
compounding actually looked like at scale.

### The String Problem

By March 5, there were 879 user-facing strings scattered across
1,500 Go files. Not because anyone decided to put them there.
Because each feature session added 10-15 strings, and nobody
stopped to ask "*should these be in YAML?*"

Finding them all took longer than externalizing them. The
archaeology was the cost, not the migration.

### The Taxonomy Problem

24 CLI packages had accumulated their own conventions. Some put
cobra wiring in `cmd.go`. Some put it in `root.go`. Some mixed
business logic with command registration. Some had helpers at the
bottom of `run.go`. Some had separate `util.go` files.

At peak drift, adding a feature meant first figuring out which
of three competing patterns this package was using.

Restructuring one package into `cmd/root/ + core/` took 15
minutes. Restructuring 24 of them took **days**, because each one
had slightly different conventions to untangle. 

If we had restructured every 4th package as it was built, the taxonomy
would have emerged naturally.

### The Type Problem

Cross-cutting types like `SessionInfo`, `ExportParams`, and
`ParserResult` were defined in whichever package first needed
them. By March 5, the same types were imported through 3-4
layers of indirection, causing import cycles that required
`internal/entity` to break.

The entity package extracted 30+ types from 12 packages. Each
extraction risked breaking imports in packages we hadn't touched
in weeks.

### The Error Problem

Per-package `err.go` files had grown into a [broken-window
pattern][broken-window]:

An agent sees `err.go` in a package, adds another error
constructor. By March 5, there were error constructors scattered
across 22 packages with no central inventory. The consolidation
into `internal/err/` domain files required tracing every error
through every caller.

### The Output Problem

Output functions (`cmd.Println`, `fmt.Fprintf`) were mixed into
business logic. When we decided output belongs in `write/`
packages, we had to extract functions from every CLI package.
The Phase WC baseline commit (`4ec5999`) marks the starting
point of this migration. 181 commits later, it was done.

---

## The Compound Interest Math

The 3:1 rule assumes consolidation sessions of roughly equal
size to feature sessions. Here is what happens when you skip:

| Consolidation cadence | Feature sessions | Consolidation sessions | Total |
|-----------------------|:----------------:|:----------------------:|:-----:|
| Every 4th (3:1)       |        48        |           16           |  64   |
| Every 10th            |        48        |           ~8           |  ~56  |
| Never (what we did)   |   198 commits    |      181 commits       |  379  |

!!! warning "The Takeaway"
    You don't save consolidation work by skipping it: 

    **You increase its cost**.

Skipping consolidation doesn't save time: **It borrows it**. 

The interest rate is **nonlinear**: The longer you wait, the more each
individual fix costs, because fixes interact with other unfixed
drift.

Renaming a constant in week 2 touches 3 files. Renaming it in
week 6 touches 15, because five features built on the original
name.

---

## What Consolidation Actually Looked Like

The 18-day consolidation run wasn't one sweep. It was a sequence
of targeted campaigns, each revealing the next:

**Week 1 (Mar 5-11)**: Error consolidation and `write/` migration.
Move output functions out of `core/`. Split monolithic `errors.go`
into 22 domain files. Remove `fatih/color`. This exposed the scope
of the string problem.

**Week 2 (Mar 12-18)**: String externalization. Create
`commands.yaml`, `flags.yaml`, split `text.yaml` into 6 domain
files. Add 879 `DescKey`/`TextDescKey` constants. Build exhaustive
test. Normalize all import aliases to camelCase. This exposed the
taxonomy problem.

**Week 3 (Mar 19-23)**: Taxonomy enforcement. Singularize command
directories. Add `doc.go` to all 75 packages. Standardize import
aliases project-wide. Fix `lint-drift` false positives. This was
the "polish" phase, except it took 5 days because the
inconsistencies had compounded across 461 packages.

Each week's work would have been a single session if done
incrementally.

---

## Lessons (Again)

The [3:1 post][ratio-post] listed the symptoms of drift. This
post adds the consequences of ignoring them:

**Consolidation is not optional; it is deferred or paid**: We
didn't avoid 16 consolidation sessions by skipping them. We
compressed them into 18 days of uninterrupted cleanup. The work
was the same; the experience was worse.

**Feature velocity creates an illusion of progress**: 198 commits
felt productive. But the codebase on March 5 was harder to modify
than the codebase on February 16, despite having more features.

!!! tip "Speed Without Structure"
    Speed without structure is negative progress.

**Agents amplify both building and debt**: The same AI that can
restructure 24 packages in a day can also create 24 slightly
different conventions in a day. The 3:1 rule matters more with
AI-assisted development, not less.

**The consolidation baseline is the most important commit to
record**: We tracked ours in `TASKS.md` (`4ec5999`). Without that
marker, knowing where to start the cleanup would have been its
own archaeological expedition.

---

## The Updated Rule

The 3:1 ratio still works. We just didn't follow it. The updated
practice:

1. **After every 3rd feature session, schedule consolidation.**
   Not "*when it feels right.*" Not "*when things get bad.*" After
   the 3rd session.

2. **Record the baseline commit.** When you start a consolidation
   phase, write down the commit hash. It marks where the debt
   starts.

3. **Run `make audit` before feature work.** If it doesn't pass,
   you are already in debt. Consolidate before building.

4. **Treat consolidation as a feature.** It gets a branch. It
   gets commits. It gets a blog post. It is not overhead; it is
   the work that makes the next three features possible.

!!! quote "The Rule"
    The 3:1 ratio is not aspirational: **It is structural**.

    Ignore consolidation, and the system will schedule it for you.

---

## The Arc

This is the eighth post in the ctx blog series:

1. [The Attention Budget](2026-02-03-the-attention-budget.md):
   why context windows are a scarce resource
2. [Before Context Windows, We Had Bouncers](2026-02-14-irc-as-context.md):
   the IRC lineage of context engineering
3. [Context as Infrastructure](2026-02-17-context-as-infrastructure.md):
   treating context as persistent files, not ephemeral prompts
4. [When a System Starts Explaining Itself](2026-02-17-when-a-system-starts-explaining-itself.md):
   the journal as a first-class artifact
5. [The Homework Problem](2026-02-25-the-homework-problem.md):
   what happens when AI writes code but humans own the outcome
6. [Agent Memory Is Infrastructure](2026-03-04-agent-memory-is-infrastructure.md):
   L2 memory vs L3 project knowledge
7. [The Architecture Release](2026-03-23-ctx-v0.8.0-the-architecture-release.md):
   what v0.8.0 looks like from the inside
8. **We Broke the 3:1 Rule** (this post):
   what happens when you don't consolidate

*See also: [The 3:1 Ratio](2026-02-17-the-3-1-ratio.md):
the original observation. This post is the empirical follow-up,
five weeks and 379 commits later.*

---

Key commits marking the consolidation arc:

| Commit     | Milestone                                              |
|------------|--------------------------------------------------------|
| `4ec5999`  | Phase WC baseline (consolidation starts)               |
| `ff6cf19e` | All CLI packages restructured into `cmd/ + core/`      |
| `d295e49c` | All command descriptions externalized to YAML          |
| `3a0bae86` | Error package split into 22 domain files               |
| `0fcbd11c` | `fatih/color` removed; 2 dependencies remain           |
| `5b32e435` | `doc.go` added to all 75 packages                      |
| `a82af4bc` | Import aliases standardized project-wide               |
| `692f86cd` | `lint-drift` false positives fixed; `make audit` green |

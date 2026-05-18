---
title: "Build a Knowledge Base"
icon: lucide/library
---

![ctx](../images/ctx-banner.png)

## The Problem

You are doing knowledge-shaped work (vendor-spec analysis, a
research project, a post-incident review, domain modeling) and
the standard five context files (`TASKS.md`, `DECISIONS.md`,
`LEARNINGS.md`, `CONVENTIONS.md`, `CONSTITUTION.md`) don't fit.
Because those files are tuned for *code-development context*, not for
*evidence-tracked knowledge* with confidence bands,
contradictions, and external citations.

You need a place where:

- Every claim is pinned to a source you can re-verify.
- Topics grow into folders as they earn their depth.
- Two passes against the same source don't silently disagree.
- The next session knows *what's incomplete*, not just *what's done*.

That's what the **editorial pipeline** is for.

!!! tip "Prefer Skills to Raw Commands"
    The pipeline is driven by skills (`/ctx-kb-ingest`,
    `/ctx-kb-ask`, etc.). The CLI form (`ctx kb ingest`, etc.)
    exists for scripting and for non-Claude environments; the
    skill is the natural surface.

## TL;DR

```bash
git init && ctx init                    # lays down the kb + ingest tree
ctx kb topic new "Cursor Hooks"         # scaffold a topic folder
/ctx-kb-ingest ./docs/cursor-hooks.md "cursor hooks" # editorial pass
/ctx-kb-ask "does the kb say hooks fire async?"      # grounded Q&A
/ctx-wrap-up                            # ceremony; delegates to /ctx-handover
                                        # for the per-session handover
```

## Commands and Skills Used

| Tool                       | Type    | Purpose                                                |
|----------------------------|---------|--------------------------------------------------------|
| `ctx init`                 | Command | Scaffold `.context/kb/`, `.context/ingest/`, etc.      |
| `ctx kb topic new <name>`  | Command | Sole writer of topic-page scaffolds (folder shape)     |
| `ctx kb note "<text>"`     | Command | Lightweight capture into `.context/ingest/findings.md` |
| `ctx kb reindex`           | Command | Refresh the `CTX:KB:TOPICS` managed block              |
| `ctx handover write`       | Command | Per-session handover with closeout fold                |
| `/ctx-kb-ingest`           | Skill   | Mode-aware editorial pass (topic-page/triage/evidence) |
| `/ctx-kb-ask`              | Skill   | Q&A grounded in the kb                                 |
| `/ctx-kb-site-review`      | Skill   | Mechanical structural audit                            |
| `/ctx-kb-ground`           | Skill   | Read-only freshness audit over the kb's tracked sources |
| `/ctx-kb-note`             | Skill   | Capture a finding for the next ingest pass             |
| `/ctx-wrap-up`             | Skill   | End-of-session ceremony; delegates to the handover step |

## Step 0: Initialize and Declare Scope

```bash
git init && ctx init
```

`ctx init` lays down the editorial scaffolding alongside the
standard context files:

```
.context/
├── kb/
│   ├── index.md
│   └── topics/.gitkeep
├── ingest/
│   ├── KB-RULES.md             # editorial constitution
│   ├── 00-GROUND.md
│   ├── 30-INGEST.md
│   ├── 40-ASK.md
│   ├── 50-SITE_REVIEW.md
│   ├── OPERATOR.md
│   ├── PROMPT.md               # hand-fallback router
│   ├── closeouts/.gitkeep
│   └── schemas/
│       └── *.md                # 10 schema templates
└── handovers/.gitkeep
```

**Open `.context/kb/index.md` and replace the placeholder `## Scope`
paragraph with a one-paragraph statement of what this kb covers
and what it does not.** `/ctx-kb-ingest` refuses to run against
an undeclared kb; scope is the precondition.

!!! warning "Git is required"
    `ctx init` now refuses to run without `.git/`. The
    editorial pipeline's provenance (closeout `sha`/`branch`,
    evidence-index in-repo SHA pins) depends on it. Run
    `git init` first if the project does not already have one.

## Step 1: Scaffold a Topic

Topic pages live in folders, not flat files:

```bash
ctx kb topic new "Cursor Hooks"
```

This creates `.context/kb/topics/cursor-hooks/index.md` from the
embedded template. The slug is computed by lowercasing + kebab-
casing; vendor-namespaced shapes like `cursor/hooks` are
preserved so you can grow into nested topology
(`topics/cursor/hooks/`, `topics/cursor/skills/`,
`topics/cursor/rules/`) without breaking citations.

**`ctx kb topic new` is the sole writer of topic-page scaffolds.**
Skills invoke this command rather than synthesize a scaffold by
hand; the embedded template is the single source of truth.

## Step 2: Run an Editorial Pass

```text
/ctx-kb-ingest ./inputs/2026-04-12-call.md "cursor hooks"
```

The skill begins with a **pass-mode declaration**:

> **Pass-mode:** `topic-page`
> **Reason:** the user supplied one primary source and the intended topic is clear.
> **Definition of done:** create or extend `kb/topics/cursor-hooks/index.md`, 
> cite EV rows, run `ctx kb site build`, record cold-reader orientation.

Then it:

1. **Resolves sources** (paths / URLs / MCP resources) and updates
   the **source-coverage ledger** at
   `.context/kb/source-coverage.md` (a state machine across all
   sources the kb has touched).
2. **Scans for adjacent incomplete topics** in the ledger and
   surfaces them so the new page acknowledges sibling gaps.
3. **Synthesizes prose** section by section into the topic page,
   minting `EV-###` rows in `evidence-index.md` for every cited
   claim.
4. **Sets the Confidence floor** (the page never claims more
   certainty than its weakest cited band).
5. **Writes a closeout** under
   `.context/ingest/closeouts/<TS>-ingest-closeout.md` with
   frontmatter, the cold-reader orientation rubric, and a
   ledger-state advance per source.

Three pass modes:

- **`topic-page`** (default): write or extend a topic page.
- **`triage`**: admit / skip sources against scope; no `EV-###` minted.
- **`evidence-only`**: mint `EV-###` rows tagged `evidence-only`;
  do not touch a topic page (explicit-request-only escape hatch).

**Mid-pass mode-switching is forbidden**: the skill commits to
one mode and aborts cleanly if the work no longer fits.

## Step 3: Q&A Grounded in the KB

```text
/ctx-kb-ask "does the kb say hooks fire async?"
```

`/ctx-kb-ask` reads the kb's prose, cites `EV-###` rows, and
**refuses to web-jump**. If the kb cannot answer, it opens a
`Q-###` row in `outstanding-questions.md` and reports the gap,
which a future `/ctx-kb-ingest` pass can close.

## Step 4: Audit + Re-Ground

```text
/ctx-kb-site-review        # mechanical structural audit
/ctx-kb-ground             # refresh sources listed in grounding-sources.md
```

`site-review` coerces malformed Confidence-band capitalization,
flags malformed closeout frontmatter, and **refuses to make
judgment calls that require evidence** (those go through ingest).

`ground` reads `.context/ingest/grounding-sources.md` — the kb's
persistent watch list — and walks each declared source (URL,
in-tree path, or MCP resource) to check whether it has drifted
since the kb last cited it. The pass is **read-only on the kb's
prose and evidence**: it annotates the source-coverage ledger's
`Residue` / `Next action` cells and writes a ground closeout, but
does NOT re-extract claims, mint `EV-###` rows, or touch topic
pages. Drifted or new-to-kb sources are flagged for a follow-up
`/ctx-kb-ingest`. Use ground for "are the docs still current?"
hygiene; use `/ctx-kb-ingest` to actually absorb new material.

## Step 5: Browse the KB Locally

`.context/kb/` is a tree of Markdown files: topic pages live
under `topics/<slug>/index.md` and cross-cutting artifacts
(`glossary.md`, `evidence-index.md`,
`outstanding-questions.md`, `domain-decisions.md`,
`contradictions.md`, `timeline.md`, `source-map.md`,
`source-coverage.md`, `relationship-map.md`) sit alongside
them. Drop a minimal `zensical.toml` into `.context/kb/` and
hand it to [`ctx serve`](../cli/serve.md):

```bash
ctx serve .context/kb/
```

The KB renders the same way the docs site you are reading
right now does. Use the in-place evidence-index links to jump
from a topic page to its `EV-###` rows and back. The site
build is read-only: no skill or CLI writes through it.

## Step 6: Wrap Up with a Handover

Run `/ctx-wrap-up` at session end; it owns the ceremony and
delegates to the handover step (`/ctx-handover`) as its final
action:

```text
/ctx-wrap-up "Cursor Hooks deep dive"
```

The handover artifact lands at
`.context/handovers/<TS>-<slug>.md` (timestamped so concurrent
agent runs never overwrite). It **folds postdated closeouts**
into a `## Folded closeouts` section and **archives the source
closeout files** under `.context/archive/closeouts/`. The next
session's `/ctx-remember` reads the latest handover and folds
any closeouts whose `generated-at` postdates it.

The legitimate direct-invocation cases for `/ctx-handover`
are `--no-fold` for a mid-session checkpoint, or recovery
when a prior session ended before its wrap-up step. For the
underlying CLI, see
[`ctx handover write`](../cli/handover.md#ctx-handover-write-title).

## How It Ladders Together

```
sources you supply
       │
       ▼
/ctx-kb-ingest (mode-declared, source-coverage advanced)
       │
       ├──▶ topic-page  ──▶ .context/kb/topics/<slug>/index.md
       ├──▶ evidence    ──▶ .context/kb/evidence-index.md (EV-###)
       ├──▶ side rails  ──▶ glossary.md, contradictions.md,
       │                    outstanding-questions.md, timeline.md,
       │                    source-map.md, relationship-map.md
       └──▶ closeout    ──▶ .context/ingest/closeouts/<TS>-...md
                               │
                               ▼
                       (next session)
                               │
                               ▼
                  /ctx-wrap-up → /ctx-handover folds
                  → .context/handovers/<TS>-<slug>.md
                  + archives source closeouts under
                  .context/archive/closeouts/
                               │
                               ▼
                  /ctx-remember reads handover + postdated
                  unfolded closeouts as the recall surface
```

## What the Editorial Pipeline Is NOT

- **Not a substitute for `DECISIONS.md`.** Project-level
  architectural decisions stay in `.context/DECISIONS.md`. The
  kb's `domain-decisions.md` is a *kb-scoped* artifact (different
  schema, different write authority, different lifecycle).
- **Not a substitute for `LEARNINGS.md`.** Learnings have author
  intent; kb claims have evidence backing. They're different
  truth bases; do not cross-feed.
- **Not for casual notes.** Use `/ctx-kb-note` or `ctx kb note
  "<text>"` to park a finding for the next ingest pass.

## Reference

- Editorial constitution: `.context/ingest/KB-RULES.md` (laid
  down by `ctx init`)
- Skills reference:
  [`/ctx-kb-ingest`](../reference/skills.md#ctx-kb-ingest),
  [`/ctx-kb-ask`](../reference/skills.md#ctx-kb-ask),
  [`/ctx-kb-site-review`](../reference/skills.md#ctx-kb-site-review),
  [`/ctx-kb-ground`](../reference/skills.md#ctx-kb-ground),
  [`/ctx-kb-note`](../reference/skills.md#ctx-kb-note),
  [`/ctx-handover`](../reference/skills.md#ctx-handover)
- Related recipes:
  [Typical KB Session](typical-kb-session.md),
  [Recover an Aborted Session](recover-aborted-session.md)

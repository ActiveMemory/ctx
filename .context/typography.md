# Typography and Document Shape

Conventions for agents and contributors editing this repo. Loaded by the
agent alongside `CONVENTIONS.md`; the linter scripts under `hack/` enforce
the hard rules on every edit so the agent should not need to be told
twice.

## Headings

### Title Case for Every Heading

All Markdown headings, at every depth, use Title Case. Minor words
(`a`, `an`, `the`, `and`, `or`, `but`, `of`, `for`, `in`, `on`, `to`,
`with`) stay lowercase except when they open the heading.

Enforced by [`hack/title-case-headings.py`](https://github.com/ActiveMemory/ctx/blob/main/hack/title-case-headings.py).

Yes:

```markdown
## Build a Knowledge Base
### When to Use a Skill
#### Folding Closeouts Into the Handover
```

No:

```markdown
## Build a knowledge base
### When to use a skill
#### Folding closeouts into the handover
```

### Code Spans in Headings Stay as Written

Backticked identifiers inside a heading keep their literal casing.
Title Case applies to the prose, not the code.

```markdown
## `ctx init` and the Bootstrap Path
### Why `internal/write/handover` Is the Sole Writer
```

## Inline Code

### `ctx` Is Always Monotype

The tool name `ctx` is wrapped in backticks every time it appears in
prose, unless it is part of a heading where it is the first word and a
backtick would distort the section anchor. The same rule applies to
sibling tool names (`gitnexus`, `gemini-search`) and to command names,
flags, paths, file names, package paths, and configuration keys.

Yes:

```markdown
Run `ctx status` to inspect the working tree. The `.context/` directory
lives at the project root. Configure `ctx.handover.subdir` in `.ctxrc`.
```

No:

```markdown
Run ctx status to inspect the working tree. The .context/ directory
lives at the project root.
```

### Exception for Prose Where the Brand Is the Subject

Marketing taglines that treat the project name as a brand may drop the
backticks. Reserve this for landing pages and the manifesto; every other
surface uses the monotype form.

## Punctuation

### No Em-Dashes, No En-Dashes

Prose uses ASCII only. Hyphens (`-`) connect compound modifiers; colons,
parentheses, or sentence breaks carry the load that an em-dash would.

Enforced by [`hack/detect-ai-typography.sh`](https://github.com/ActiveMemory/ctx/blob/main/hack/detect-ai-typography.sh).
The script flags U+2013 (`–`) and U+2014 (`—`) on sight; both are
heuristic signals of AI prose that was not human-reviewed.

Yes:

```markdown
The handover is the sole authoritative recall artifact: the next session
reads it before anything else.
```

No:

```markdown
The handover is the sole authoritative recall artifact — the next session
reads it before anything else.
```

### No Smart Quotes

Use straight ASCII quotes (`"` and `'`). Curly quotes (U+2018, U+2019,
U+201C, U+201D) are flagged by the detector and auto-fixed by
[`hack/fix-smart-quotes.sh`](https://github.com/ActiveMemory/ctx/blob/main/hack/fix-smart-quotes.sh).

### Space-Padded Double Hyphens Are Not a Dash

The detector flags `" -- "` as an AI-prose tell. CLI flag examples
(`--summary`, `git rebase --no-edit`) are unaffected because they have
no surrounding spaces. Table separators (`| -- |`) are also fine.

### No Quad Backticks

Triple backticks (` ``` `) are the project maximum for code fences. AI
often emits ` ```` ` (quad) which the Zensical renderer does not
support. The detector flags these.

## Code Fences

Every fence carries a language tag (`bash`, `go`, `markdown`, `yaml`,
`text`, `json`). Untagged fences render without syntax highlighting and
should not appear in the docs.

Open with `` ```bash `` (tagged) rather than bare `` ``` `` (untagged):

```bash
ctx status
```

Triple backticks are the project maximum. Quad backticks (`` ```` ``)
are an AI artifact and do not render in Zensical. When you need to show
a fenced block inside another, use a four-space indent for the outer
demonstration rather than nesting fences.

## Document Header (Docs Site)

### Frontmatter Shape

Every `docs/` Markdown file opens with YAML frontmatter that carries the
SPDX banner as YAML comments, then the `title:` and `icon:` fields:

```markdown
---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Title in Title Case
icon: lucide/<icon-name>
---
```

The `title:` value uses Title Case. The `icon:` value selects from the
Lucide icon set (the Material for MkDocs default), referenced as
`lucide/<name>`. Browse available icons at
[lucide.dev/icons](https://lucide.dev/icons).

### Banner Image as the First Body Line

The first non-frontmatter line is the project banner image, with a
relative path from the file's location:

```markdown
![ctx](../images/ctx-banner.png)
```

Top-level pages (`docs/index.md`) use `images/ctx-banner.png`; pages one
directory deep use `../images/ctx-banner.png`. The banner anchors every
page visually and gives the rendered output a consistent opening.

### Optional `# H1` Page Title

Material for MkDocs renders the `title:` frontmatter as the page title.
A leading `# H1` is optional and used when the in-body title differs
from the navigation title (longer, more descriptive, or with backticks).

## Document Structure (Recipes)

Recipe pages under `docs/recipes/` follow a shared narrative arc that
makes them comparable side by side. Not every page hits every section,
but the order is fixed:

```markdown
## The Problem
## TL;DR
## The Solution            (or: ## Setup, ## Process)
## When to Use
## When NOT to Use
## Quality Checklist       (optional)
```

The `## The Problem` section is the strongest pattern: 34 of 48 recipes
open this way. Reach for it before inventing a custom opening.

## Admonitions

Use Material for MkDocs admonition syntax. Five variants are in active
use across the corpus:

```markdown
!!! note
    Neutral context the reader should not miss.

!!! tip
    Optional improvement or shortcut.

!!! warning
    Behavior that may surprise; recoverable.

!!! danger
    Destructive or irreversible operation.

!!! info
    Background detail; safe to skip.
```

Admonition titles use Title Case when given:

```markdown
!!! tip "Prefer Skills Over Raw Commands"
    ...
```

## Enforcement Summary

| Rule | Enforced By |
|------|-------------|
| Title Case headings | [`hack/title-case-headings.py`](https://github.com/ActiveMemory/ctx/blob/main/hack/title-case-headings.py) |
| No em-dashes / en-dashes / smart quotes / quad backticks / padded `--` | [`hack/detect-ai-typography.sh`](https://github.com/ActiveMemory/ctx/blob/main/hack/detect-ai-typography.sh) |
| Auto-fix curly quotes | [`hack/fix-smart-quotes.sh`](https://github.com/ActiveMemory/ctx/blob/main/hack/fix-smart-quotes.sh) |

Run the detector before submitting docs changes:

```bash
./hack/detect-ai-typography.sh docs
```

It exits non-zero on any hit and prints file + line number so an editor
jump-to-line takes you straight to the offender.

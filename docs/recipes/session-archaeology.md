---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "Browsing and Enriching Past Sessions"
icon: lucide/scroll-text
---

![ctx](../images/ctx-banner.png)

## Problem

After weeks of AI-assisted development you have dozens of sessions scattered
across JSONL files in `~/.claude/projects/`. Finding that session where you
debugged the Redis connection pool -- or remembering what you decided about
the caching strategy three Tuesdays ago -- means grepping raw JSON. There is
no table of contents, no search, no summaries.

This recipe shows how to turn that raw session history into a browsable,
searchable, and enriched journal site you can navigate in your browser.

## Commands and Skills Used

| Tool                     | Type    | Purpose                                      |
|--------------------------|---------|----------------------------------------------|
| `ctx recall list`        | Command | List parsed sessions with metadata           |
| `ctx recall show`        | Command | Inspect a specific session in detail         |
| `ctx recall export`      | Command | Export sessions to editable journal Markdown |
| `ctx journal site`       | Command | Generate a static site from journal entries  |
| `ctx serve`              | Command | Serve the journal site locally               |
| `/ctx-recall`            | Skill   | Browse sessions inside your AI assistant     |
| `/ctx-journal-normalize` | Skill   | Fix rendering issues in exported Markdown    |
| `/ctx-journal-enrich`    | Skill   | Add frontmatter metadata and summaries       |

## The Workflow

The session journal follows a four-stage pipeline. Each stage is idempotent
-- safe to re-run -- and each skips entries that have already been processed.

```
export --> normalize --> enrich --> rebuild
```

| Stage         | Tool                       | What it does                                | Skips if                     |
|---------------|----------------------------|---------------------------------------------|------------------------------|
| **Export**    | `ctx recall export --all`  | Converts session JSONL to Markdown          | `--skip-existing` flag       |
| **Normalize** | `/ctx-journal-normalize`   | Fixes nested fences and metadata formatting | `<!-- normalized -->` marker |
| **Enrich**    | `/ctx-journal-enrich`      | Adds frontmatter, summaries, topic tags     | Frontmatter already present  |
| **Rebuild**   | `ctx journal site --build` | Generates browsable static HTML             | --                           |

### Step 1: List Your Sessions

Start by seeing what sessions exist for the current project:

```bash
ctx recall list
```

Sample output:

```
Sessions (newest first)
=======================

  Slug                           Project   Date         Duration  Turns  Tokens
  gleaming-wobbling-sutherland   ctx       2026-02-07   1h 23m    47     82,341
  twinkly-stirring-kettle        ctx       2026-02-06   0h 45m    22     38,102
  bright-dancing-hopper          ctx       2026-02-05   2h 10m    63     124,500
  quiet-flowing-dijkstra         ctx       2026-02-04   0h 18m    11     15,230
  ...
```

Filter by project or tool if you work across multiple codebases:

```bash
ctx recall list --project ctx --limit 10
ctx recall list --tool claude-code
ctx recall list --all-projects
```

### Step 2: Inspect a Specific Session

Before exporting everything, you can inspect a single session to see its
metadata and conversation summary:

```bash
ctx recall show --latest
```

Or look up a specific session by its slug, partial ID, or UUID:

```bash
ctx recall show gleaming-wobbling-sutherland
ctx recall show twinkly
ctx recall show abc123
```

Add `--full` to see the complete message content instead of the summary view:

```bash
ctx recall show --latest --full
```

This is useful for quickly checking what happened in a session before deciding
whether to export and enrich it.

### Step 3: Export Sessions to Journal

Export converts raw session data into editable Markdown files in
`.context/journal/`:

```bash
# Export all sessions from the current project
ctx recall export --all

# Export a single session
ctx recall export gleaming-wobbling-sutherland

# Include sessions from all projects
ctx recall export --all --all-projects
```

Each exported file contains session metadata (date, time, duration, model,
project, git branch), a tool usage summary, and the full conversation
transcript.

**Re-exporting is safe.** By default, re-running `ctx recall export --all`
regenerates conversation content while preserving any YAML frontmatter you
or the enrichment skill have added. Use `--skip-existing` to leave exported
files completely untouched, or `--force` for a full overwrite (frontmatter
will be lost).

### Step 4: Normalize Rendering

Raw exported sessions often have rendering problems: nested code fences that
break Markdown parsers, malformed metadata tables, or broken list formatting.
The normalize skill fixes these issues in the source files before site
generation.

Inside your AI assistant:

```
/ctx-journal-normalize
```

The skill backs up `.context/journal/` before modifying anything and marks
each processed file with a `<!-- normalized: YYYY-MM-DD -->` comment so
subsequent runs skip already-normalized entries.

**Run normalize before enrich.** The enrichment skill reads conversation
content to extract metadata, and clean Markdown produces better results.

### Step 5: Enrich with Metadata

Raw exports have timestamps and transcripts but lack the semantic metadata
that makes sessions searchable: topics, technology tags, outcome status,
and summaries. The enrich skill adds this structured frontmatter.

Inside your AI assistant:

```
/ctx-journal-enrich twinkly-stirring-kettle
```

You can match by slug, partial slug, date, or partial UUID:

```
/ctx-journal-enrich twinkly
/ctx-journal-enrich 2026-02-06
/ctx-journal-enrich 76fe2ab9
```

The skill analyzes the conversation and proposes frontmatter:

```yaml
---
title: "Implement Redis caching middleware"
date: 2026-02-06
type: feature
outcome: completed
topics:
  - caching
  - api-performance
technologies:
  - go
  - redis
libraries:
  - go-redis/redis
key_files:
  - internal/cache/redis.go
  - internal/api/middleware/cache.go
---
```

It also generates a summary and extracts decisions, learnings, and tasks
mentioned during the session. You review a diff and confirm before it writes.

Enrichment is intentionally interactive to ensure accuracy. To find sessions
that still need enrichment:

```bash
grep -L "^---$" .context/journal/*.md | head -10
```

### Step 6: Generate and Serve the Site

With exported, normalized, and enriched journal files, generate the static
site:

```bash
# Generate site structure only
ctx journal site

# Generate and build static HTML
ctx journal site --build

# Generate, build, and serve locally
ctx journal site --serve
```

Then open [http://localhost:8000](http://localhost:8000) to browse.

The site includes a date-sorted index, individual session pages with full
conversations, search (press `/`), dark mode, and enriched titles in the
navigation when frontmatter exists.

You can also serve an existing site without regenerating using `ctx serve`.
The site generator requires [zensical](https://pypi.org/project/zensical/)
(`pip install zensical`).

## Putting It Together

The complete pipeline from raw sessions to browsable site:

```bash
# Terminal: export and generate
ctx recall export --all
ctx journal site --serve
```

```
# AI assistant: normalize and enrich
/ctx-journal-normalize
/ctx-journal-enrich twinkly-stirring-kettle
/ctx-journal-enrich gleaming-wobbling-sutherland
```

```bash
# Terminal: rebuild with enrichments
ctx journal site --serve
```

If your project includes `Makefile.ctx` (deployed by `ctx init`), use
`make journal` to combine the export and rebuild stages, then normalize
and enrich inside Claude Code, then `make journal` again to pick up
the enrichments.

## Tips

**Start with `/ctx-recall` inside your AI assistant.** If you just want to
quickly check what happened in a recent session without leaving your editor,
the `/ctx-recall` skill lets you browse interactively without exporting.

**Large sessions are split automatically.** Sessions with 200+ messages get
split into multiple parts (`session-abc123.md`, `session-abc123-p2.md`,
`session-abc123-p3.md`) with navigation links between them. The site
generator handles this transparently.

**Suggestion sessions are separated.** Claude Code generates short
"suggestion" sessions for auto-complete. These appear under a separate
"Suggestions" section in the site index so they don't clutter your main
session list.

**Enrich your best sessions first.** You don't need to enrich every session.
Focus on productive sessions where you made decisions, learned something, or
completed significant work. Exploration and debugging sessions can stay
unenriched -- they're still searchable by content.

**Journal files should be gitignored.** Session transcripts contain sensitive
data: file contents, commands, error messages with stack traces, and
potentially API keys. Add `.context/journal/` and `.context/journal-site/`
to your `.gitignore`.

## See Also

- [The Complete Session](session-lifecycle.md) -- where session saving fits in the daily workflow
- [Turning Activity into Content](publishing.md) -- generating blog posts from session history
- [Session Journal](../session-journal.md) -- full documentation of the journal system
- [CLI Reference: ctx recall](../cli-reference.md#ctx-recall) -- all recall subcommands and flags
- [CLI Reference: ctx journal](../cli-reference.md#ctx-journal) -- site generation options
- [CLI Reference: ctx serve](../cli-reference.md#ctx-serve) -- local serving options
- [Context Files](../context-files.md) -- the `.context/` directory structure

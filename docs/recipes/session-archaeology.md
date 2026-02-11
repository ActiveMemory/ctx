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

## The Problem

After weeks of AI-assisted development you have dozens of sessions scattered
across JSONL files in `~/.claude/projects/`. Finding the session where you
debugged the Redis connection pool, or remembering what you decided about the
caching strategy three Tuesdays ago, often means grepping raw JSON.

There is no table of contents, no search, and no summaries.

This recipe shows how to turn that raw session history into a **browsable**,
**searchable**, and **enriched **journal site** you can navigate 
**in your browser**.

## Commands and Skills Used

| Tool                      | Type    | Purpose                                      |
|---------------------------|---------|----------------------------------------------|
| `ctx recall list`         | Command | List parsed sessions with metadata           |
| `ctx recall show`         | Command | Inspect a specific session in detail         |
| `ctx recall export`       | Command | Export sessions to editable journal Markdown |
| `ctx journal site`        | Command | Generate a static site from journal entries  |
| `ctx serve`               | Command | Serve the journal site locally               |
| `/ctx-recall`             | Skill   | Browse sessions inside your AI assistant     |
| `/ctx-journal-normalize`  | Skill   | Fix rendering issues in exported Markdown    |
| `/ctx-journal-enrich`     | Skill   | Add frontmatter metadata to a single entry   |
| `/ctx-journal-enrich-all` | Skill   | Batch-enrich all unenriched entries          |

## The Workflow

The session journal follows a four-stage pipeline.

Each stage is idempotent and safe to re-run.
Each stage skips entries that have already been processed.

```text
export -> normalize -> enrich -> rebuild
```

| Stage     | Tool                       | What it does                                | Skips if                     | Where        |
|-----------|----------------------------|---------------------------------------------|------------------------------|--------------|
| Export    | `ctx recall export --all`  | Converts session JSONL to Markdown          | `--skip-existing` flag       | CLI or agent |
| Normalize | `/ctx-journal-normalize`   | Fixes nested fences and metadata formatting | `<!-- normalized -->` marker | Agent only   |
| Enrich    | `/ctx-journal-enrich-all`  | Adds frontmatter, summaries, topic tags     | Frontmatter already present  | Agent only   |
| Rebuild   | `ctx journal site --build` | Generates browsable static HTML             | N/A                          | CLI only     |

!!! tip "Where to run Each Stage"
    Export (*Steps 1 to 3*) works equally well from the terminal or inside your
    AI assistant via `/ctx-recall`. The CLI is fine here: the agent adds no
    special intelligence, it just runs the same command.
    
    Normalize and enrich (Steps 4 to 5) require the agent: they need it to read,
    analyze, and edit Markdown.
    
    Rebuild and serve (Step 6) is a terminal operation that starts a
    long-running server.

### Step 1: List Your Sessions

Start by seeing what sessions exist for the current project:

```bash
ctx recall list
```

Sample output:

```text
Sessions (newest first)
=======================

  Slug                           Project   Date         Duration  Turns  Tokens
  gleaming-wobbling-sutherland   ctx       2026-02-07   1h 23m    47     82,341
  twinkly-stirring-kettle        ctx       2026-02-06   0h 45m    22     38,102
  bright-dancing-hopper          ctx       2026-02-05   2h 10m    63     124,500
  quiet-flowing-dijkstra         ctx       2026-02-04   0h 18m    11     15,230
  ...
```

!!! tip "Slugs Look Cryptic?"
    These auto-generated slugs (`gleaming-wobbling-sutherland`) are hard to
    recognize later.

    Use `/ctx-journal-enrich` to add human-readable titles, topic tags, and
    summaries to exported journal entries, making them easier to find.

Filter by project or tool if you work across multiple codebases:

```bash
ctx recall list --project ctx --limit 10
ctx recall list --tool claude-code
ctx recall list --all-projects
```

### Step 2: Inspect a Specific Session

Before exporting everything, inspect a single session to see its metadata and
conversation summary:

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

This is useful for checking what happened before deciding whether to export and
enrich it.

### Step 3: Export Sessions to the Journal

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

Each exported file contains session metadata (*date, time, duration, model,
project, git branch*), a tool usage summary, and the full conversation transcript.

Re-exporting is safe. By default, re-running `ctx recall export --all`
regenerates conversation content while preserving any YAML frontmatter you
or the enrichment skill have added.

You can also use `--skip-existing` to leave exported files completely untouched.

!!! warning "--force Overwrites Journal Files"
    If you want to overwrite existing files, use `--force`.

    This triggers a **full overwrite** and frontmatter will be lost.
    
    **Back up your journal before using this flag**.

### Step 4: Normalize Rendering

Raw exported sessions often have rendering problems: nested code fences that
break Markdown parsers, malformed metadata tables, or broken list formatting.

The `/ctx-journal-normalize` skill fixes these issues in the source files
before site generation.

Inside your AI assistant:

```text
/ctx-journal-normalize
```

The skill backs up `.context/journal/` before modifying anything and marks each
processed file with a `<!-- normalized: YYYY-MM-DD -->` comment so subsequent
runs skip already-normalized entries.

Run normalize before enrich. The enrichment skill reads conversation content to
extract metadata, and clean Markdown produces better results.

### Step 5: Enrich with Metadata

Raw exports have timestamps and transcripts but lack the semantic metadata that
makes sessions searchable: topics, technology tags, outcome status, and
summaries. The `/ctx-journal-enrich*` skills add this structured frontmatter.

Batch enrichment (recommended):

```text
/ctx-journal-enrich-all
```

The skill finds all unenriched entries, filters out noise (*suggestion sessions,
very short sessions, multipart continuations*), and processes each one by
extracting **titles**, **topics**, **technologies**, and 
**summaries** from the conversation.

It shows you a grouped summary before applying changes so you can scan quickly
rather than reviewing one by one.

For large backlogs (*20+ entries*), the skill can spawn subagents to process
entries in parallel.

Single-entry enrichment:

```text
/ctx-journal-enrich twinkly
/ctx-journal-enrich 2026-02-06
```

Each enriched entry gets YAML frontmatter like this:

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

The skill also generates a summary and can extract **decisions**, 
**learnings**, and **tasks** mentioned during the session.

### Step 6: Generate and Serve the Site

With exported, normalized, and enriched journal files, generate the static site:

```bash
# Generate site structure only
ctx journal site

# Generate and build static HTML
ctx journal site --build

# Generate, build, and serve locally
ctx journal site --serve
```

Then open `http://localhost:8000` to browse.

The site includes a date-sorted index, individual session pages with full
conversations, search (press `/`), dark mode, and enriched titles in the
navigation when frontmatter exists.

You can also serve an existing site without regenerating using `ctx serve`.

The site generator requires `zensical` (`pip install zensical`).

## Where the Agent Adds Value

Export, list, and show are mechanical. The agent runs the same CLI commands you
would, so you can stay in your terminal for those.

The agent earns its keep in **normalize** and **enrich**. 

These require reading conversation content, understanding what happened, 
and producing structured metadata. **That is agent work, not CLI work**.

You can also ask your agent to browse sessions conversationally instead of
remembering flags:

```text
What did we work on last week?
Show me the session about Redis.
Export everything to the journal.
```

This is convenient but not required: `ctx recall list` gives you the same
inventory.

Where the agent genuinely helps is chaining the pipeline:

```text
You:   What happened last Tuesday?
Agent: Last Tuesday you worked on two sessions:
       - bright-dancing-hopper (2h 10m): refactored the middleware
         pipeline and added Redis caching
       - quiet-flowing-dijkstra (18m): quick fix for a nil pointer
         in the config loader
       Want me to export and enrich them?
You:   Yes, do it.
Agent: Exports both, normalizes, enriches, then proposes frontmatter.
```

The value is staying in one context while the agent runs export -> normalize ->
enrich without you manually switching tools.

## Putting It All Together

A typical pipeline from raw sessions to a browsable site:

```bash
# Terminal: export and generate
ctx recall export --all
ctx journal site --serve
```

```text
# AI assistant: normalize and enrich
/ctx-journal-normalize
/ctx-journal-enrich-all
```

```bash
# Terminal: rebuild with enrichments
ctx journal site --serve
```

If your project includes `Makefile.ctx` (deployed by `ctx init`), use
`make journal` to combine export and rebuild stages. Then normalize and enrich
inside Claude Code, then `make journal` again to pick up enrichments.

## Tips

* Start with `/ctx-recall` inside your AI assistant. If you want to quickly check
what happened in a recent session without leaving your editor, `/ctx-recall`
lets you browse interactively without exporting.
* Large sessions may be split automatically. Sessions with 200+ messages can be
split into multiple parts (`session-abc123.md`, `session-abc123-p2.md`,
`session-abc123-p3.md`) with navigation links between them. The site generator
can handle this.
* Suggestion sessions can be separated. Claude Code can generate short
suggestion sessions for autocomplete. These may appear under a separate section
in the site index, so they do not clutter your main session list.
* Your agent is a good session browser. You do not need to remember slugs, dates,
or flags. Ask "*what did we do yesterday?*" or "*find the session about Redis*" 
and it can map the question to recall commands.

!!! warning "Journal files are sensitive"
    Journal files **MUST** be `.gitignore`d.
    
    Session transcripts can contain sensitive data such as file contents,
    commands, error messages with stack traces, and potentially API keys.
    
    Add `.context/journal/` and `.context/journal-site/` to your `.gitignore`.

## Next Up

**[Running an Unattended AI Agent](autonomous-loops.md)**:
Set up an AI agent that works through tasks overnight without you at the
keyboard.

## See Also

* [The Complete Session](session-lifecycle.md): where session saving fits in the daily workflow
* [Turning Activity into Content](publishing.md): generating blog posts from session history
* [Session Journal](../session-journal.md): full documentation of the journal system
* [CLI Reference: ctx recall](../cli-reference.md#ctx-recall): all recall subcommands and flags
* [CLI Reference: ctx journal](../cli-reference.md#ctx-journal): site generation options
* [CLI Reference: ctx serve](../cli-reference.md#ctx-serve): local serving options
* [Context Files](../context-files.md): the `.context/` directory structure

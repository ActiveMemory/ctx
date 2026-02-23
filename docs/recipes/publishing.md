---
title: Turning Activity into Content
icon: lucide/pen-line
---

![ctx](../images/ctx-banner.png)

## Problem

Your `.context/` directory is full of decisions, learnings, and session history.

Your `git log` tells the story of a project evolving.

But none of this is visible to anyone outside your terminal.

You want to turn this raw activity into:

- a browsable journal site
- blog posts
- changelog posts

!!! tip "TL;DR"
    ```bash
    ctx recall export --all             # 1. export sessions to markdown
    ```
    ```text
    /ctx-journal-normalize              # 2. fix rendering issues
    /ctx-journal-enrich-all             # 3. add metadata and tags
    ```
    ```bash
    ctx journal site --serve            # 4. build and serve the journal
    ```
    ```text
    /ctx-blog about the caching layer   # 5. draft a blog post
    /ctx-blog-changelog v0.1.0 "v0.2"   # 6. write a changelog post
    ```

    Read on for details on each stage.

## Commands and Skills Used

| Tool                      | Type     | Purpose                                             |
|---------------------------|----------|-----------------------------------------------------|
| `ctx recall export`       | Command  | Export session JSONL to editable markdown           |
| `ctx journal site`        | Command  | Generate a static site from journal entries         |
| `ctx journal obsidian`    | Command  | Generate an Obsidian vault from journal entries     |
| `ctx serve`               | Command  | Serve the journal site locally                      |
| `make journal`            | Makefile | Shortcut for export + site rebuild                  |
| `/ctx-journal-normalize`  | Skill    | Fix markdown rendering in exported entries          |
| `/ctx-journal-enrich-all` | Skill    | Batch-enrich all unenriched entries (recommended)   |
| `/ctx-journal-enrich`     | Skill    | Add metadata, summaries, and tags to one entry      |
| `/ctx-blog`               | Skill    | Draft a blog post from recent project activity      |
| `/ctx-blog-changelog`     | Skill    | Write a themed post from a commit range             |

## The Workflow

### Step 1: Export Sessions to Markdown

Raw session data lives as JSONL files in Claude Code's internal storage. The
first step is converting these into readable, editable markdown.

```bash
# Export all sessions from the current project
ctx recall export --all

# Export from all projects (if you work across multiple repos)
ctx recall export --all --all-projects

# Export a single session by ID or slug
ctx recall export abc123
ctx recall export gleaming-wobbling-sutherland
````

Exported files land in `.context/journal/` as individual Markdown files with
session metadata and the full conversation transcript.

`--all` is safe by default — only new sessions are exported. Existing files
are skipped. Use `--regenerate` to re-export existing files (YAML frontmatter
is preserved). Use `--force -y` to overwrite everything.

### Step 2: Normalize Exported Entries

Raw exports can have rendering issues: nested code fences that break syntax
highlighting, metadata blocks that render poorly, and malformed lists. The
`/ctx-journal-normalize` skill fixes these in the source files before site
generation.

```text
/ctx-journal-normalize
```

The skill:

1. Backs up `.context/journal/` before modifying anything
2. Converts `**Key**: value` metadata blocks into collapsible HTML tables
3. Fixes fence nesting so code blocks render with proper highlighting
4. Marks processed files with `<!-- normalized -->` so reruns skip them

Run normalize before enrich. Clean Markdown produces better metadata extraction
in the next step.

### Step 3: Enrich Entries with Metadata

Raw entries have timestamps and conversations but lack the structured metadata
that makes a journal searchable. Use `/ctx-journal-enrich-all` to process your
entire backlog at once:

```text
/ctx-journal-enrich-all
```

The skill finds all unenriched entries, filters out noise (*suggestion sessions,
very short sessions, multipart continuations*), and processes each one by
extracting titles, topics, technologies, and summaries from the conversation.

For large backlogs (*20+ entries*), it can spawn subagents to process entries in
parallel.

To enrich a single entry instead:

```text
/ctx-journal-enrich twinkly-stirring-kettle
/ctx-journal-enrich 2026-01-24
```

After enrichment, an entry gains YAML frontmatter:

```yaml
---
title: "Implement Redis caching for API endpoints"
date: 2026-01-24
type: feature
outcome: completed
topics:
  - caching
  - api-performance
technologies:
  - go
  - redis
key_files:
  - internal/api/middleware/cache.go
  - internal/cache/redis.go
---
```

This metadata powers better navigation in the journal site: 

* titles replace slugs, 
* summaries appear in the index, 
* and search covers topics and technologies.

### Step 4: Generate the Journal Site

With entries exported, normalized, and enriched, generate the static site:

```bash
# Generate site files
ctx journal site

# Generate and build static HTML
ctx journal site --build

# Generate and serve locally (opens at http://localhost:8000)
ctx journal site --serve

# Custom output directory
ctx journal site --output ~/my-journal
```

The site is generated in `.context/journal-site/` by default. It uses
[zensical](https://pypi.org/project/zensical/) for static site generation
(`pipx install zensical`).

Or use the Makefile shortcut that combines export and rebuild:

```bash
make journal
```

This runs `ctx recall export --all` followed by `ctx journal site --build`, then
reminds you to normalize and enrich before rebuilding. To serve the built site,
use `make journal-serve` or `ctx serve`.

### Alternative: Export to Obsidian Vault

If you use [Obsidian](https://obsidian.md/) for knowledge management, generate
a vault instead of (or alongside) the static site:

```bash
ctx journal obsidian
ctx journal obsidian --output ~/vaults/ctx-journal
```

This produces an Obsidian-ready directory with wikilinks, MOC (Map of Content)
pages for topics/files/types, and a "Related Sessions" footer on each entry for
graph connectivity. Open the output directory in Obsidian as a vault.

The vault uses the same enriched source entries as the static site. Both outputs
can coexist — the static site goes to `.context/journal-site/`, the vault to
`.context/journal-obsidian/`.

### Step 5: Draft Blog Posts from Activity

When your project reaches a milestone worth sharing, use `/ctx-blog` to draft a
post from recent activity. The skill gathers context from multiple sources:
`git log`, `DECISIONS.md`, `LEARNINGS.md`, completed tasks, and journal entries.

```text
/ctx-blog about the caching layer we just built
/ctx-blog last week's refactoring work
/ctx-blog lessons learned from the migration
```

The skill gathers recent commits, decisions, and learnings; identifies a
narrative arc; drafts an outline for approval; writes the full post; and saves
it to `docs/blog/YYYY-MM-DD-slug.md`.

Posts are written in first person with code snippets, commit references, and an
honest discussion of what went wrong.

### Step 6: Write Changelog Posts from Commit Ranges

For release notes or "what changed" posts, `/ctx-blog-changelog` takes a
starting commit and a theme, then analyzes everything that changed:

```text
/ctx-blog-changelog 040ce99 "building the journal system"
/ctx-blog-changelog HEAD~30 "what's new in v0.2.0"
/ctx-blog-changelog v0.1.0 "the road to v0.2.0"
```

The skill diffs the commit range, identifies the most-changed files, and
constructs a narrative organized by theme rather than chronology, including a
key commits table and before/after comparisons.

## The Conversational Approach

You do not need to remember any commands. When the Agent Playbook is active,
your agent can suggest content at natural moments:

> "We just shipped the caching layer and closed 3 tasks. Want me to draft a blog post about it?"

> "Your journal has 6 new entries since the last rebuild. Want me to normalize, enrich, and regenerate the site?"

You can also drive it with natural language:

```text
"write about what we did this week"
"turn today's session into a blog post"
"make a changelog post covering everything since the last release"
"enrich the last few journal entries"
```

The agent has full visibility into your `.context/` state (tasks completed,
decisions recorded, learnings captured), so its suggestions are grounded in what
actually happened.

## Putting It Together

The full pipeline from raw transcripts to published content:

```bash
# 1. Export all sessions
ctx recall export --all

# 2. In Claude Code: normalize rendering
/ctx-journal-normalize

# 3. In Claude Code: enrich all entries with metadata
/ctx-journal-enrich-all

# 4. Build and serve the journal site
make journal
make journal-serve

# 4b. Or generate an Obsidian vault
ctx journal obsidian

# 5. In Claude Code: draft a blog post
/ctx-blog about the features we shipped this week

# 6. In Claude Code: write a changelog post
/ctx-blog-changelog v0.1.0 "what's new in v0.2.0"
```

The journal pipeline is idempotent at every stage. You can rerun `ctx recall
export --all` without losing enrichment. You can rerun `/ctx-journal-normalize`
and it skips already-normalized files. You can rebuild the site as many times
as you want.

## Tips

* Export regularly. Run `ctx recall export --all` after each session to keep
  your journal current. Only new sessions are exported — existing files are
  skipped by default.
* Normalize before enriching. The enrichment skill reads conversation content to
  extract metadata. Clean Markdown produces better results than raw exports with
  broken fences.
* Use batch enrichment. `/ctx-journal-enrich-all` filters noise (suggestion
  sessions, trivial sessions, multipart continuations) so you do not have to
  decide what is worth enriching.
* Keep journal files in `.gitignore`. Session journals can contain sensitive
  data: file contents, commands, internal discussions, and error messages with
  stack traces. Add `.context/journal/` and `.context/journal-site/` to
  `.gitignore`.
* Use `/ctx-blog` for narrative posts and `/ctx-blog-changelog` for release
  posts. One finds a story in recent activity, the other explains a commit
  range by theme.
* Let the agent remind you. A proactive agent can suggest content generation
  after milestones: shipping a feature, closing tasks, or finishing a long
  debugging session.
* Edit the drafts. These skills produce drafts, not final posts. Review the
  narrative, add your perspective, and remove anything that does not serve the
  reader.

## Next Up

**[Running an Unattended AI Agent →](autonomous-loops.md)**: Set up
an AI agent that works through tasks overnight without you at the
keyboard.

## See Also

* [Session Journal](../reference/session-journal.md): journal system, enrichment schema, context monitor
* [CLI Reference: ctx recall](../reference/cli-reference.md#ctx-recall): export, list, show session history
* [CLI Reference: ctx journal site](../reference/cli-reference.md#ctx-journal-site): static site generation
* [CLI Reference: ctx journal obsidian](../reference/cli-reference.md#ctx-journal-obsidian): Obsidian vault export
* [CLI Reference: ctx serve](../reference/cli-reference.md#ctx-serve): local site serving
* [Browsing and Enriching Past Sessions](session-archaeology.md): journal browsing workflow
* [The Complete Session](session-lifecycle.md): capturing context during a session

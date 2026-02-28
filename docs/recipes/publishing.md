---
title: Turning Activity into Content
icon: lucide/pen-line
---

![ctx](../images/ctx-banner.png)

## The Problem

Your `.context/` directory is full of decisions, learnings, and session history.

Your `git log` tells the story of a project evolving.

But none of this is visible to anyone outside your terminal.

You want to turn this raw activity into:

* a browsable journal site
* blog posts
* changelog posts

## TL;DR

```bash
ctx recall export --all             # 1. export sessions to markdown

/ctx-journal-enrich-all             # 2. add metadata and tags

ctx journal site --serve            # 3. build and serve the journal

/ctx-blog about the caching layer   # 4. draft a blog post
/ctx-blog-changelog v0.1.0 "v0.2"   # 5. write a changelog post
```

Read on for details on each stage.

## Commands and Skills Used

| Tool                      | Type     | Purpose                                             |
|---------------------------|----------|-----------------------------------------------------|
| `ctx recall export`       | Command  | Export session JSONL to editable markdown           |
| `ctx journal site`        | Command  | Generate a static site from journal entries         |
| `ctx journal obsidian`    | Command  | Generate an Obsidian vault from journal entries     |
| `ctx serve`               | Command  | Serve any zensical directory (default: journal)     |
| `ctx site feed`           | Command  | Generate Atom feed from finalized blog posts        |
| `make journal`            | Makefile | Shortcut for export + site rebuild                  |
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
is preserved). Use `--regenerate --keep-frontmatter=false -y` to regenerate
everything including frontmatter.

### Step 2: Enrich Entries with Metadata

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

### Step 3: Generate the Journal Site

With entries exported and enriched, generate the static site:

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
reminds you to enrich before rebuilding. To serve the built site, use
`make journal-serve` or `ctx serve` (serve-only, no regeneration).

### Alternative: Export to Obsidian Vault

If you use [Obsidian](https://obsidian.md/) for knowledge management, generate
a vault instead of (*or alongside*) the static site:

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

### Step 4: Draft Blog Posts from Activity

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

!!! info "The Output is `zensical`-Flavored Markdown"
    The blog skills produce Markdown tuned for a
    [zensical](https://pypi.org/project/zensical/) site: `topics:`
    frontmatter (zensical's tag field), a `docs/blog/` output path,
    and a banner image reference. 

    The content is still standard Markdown and can be adapted to other 
    static site generators, but the defaults assume a `zensical` 
    project structure.

### Step 5: Write Changelog Posts from Commit Ranges

For release notes or "*what changed*" posts, `/ctx-blog-changelog` takes a
starting commit and a theme, then analyzes everything that changed:

```text
/ctx-blog-changelog 040ce99 "building the journal system"
/ctx-blog-changelog HEAD~30 "what's new in v0.2.0"
/ctx-blog-changelog v0.1.0 "the road to v0.2.0"
```

The skill diffs the commit range, identifies the most-changed files, and
constructs a narrative organized by theme rather than chronology, including a
key commits table and before/after comparisons.

### Step 6: Generate the Blog Feed

After publishing blog posts, generate the Atom feed so readers and
automation can discover new content:

```bash
ctx site feed
```

This scans `docs/blog/` for finalized posts (`reviewed_and_finalized: true`),
extracts title, date, author, topics, and summary, and writes a valid
Atom 1.0 feed to `site/feed.xml`. The feed is also generated automatically
as part of `make site`.

The feed is available at [ctx.ist/feed.xml](https://ctx.ist/feed.xml).

## The Conversational Approach

You can also drive your publishing anytime with **natural language**:

```text
"write about what we did this week"
"turn today's session into a blog post"
"make a changelog post covering everything since the last release"
"enrich the last few journal entries"
```

The agent has full visibility into your `.context/` state (*tasks completed,
decisions recorded, learnings captured*), so its suggestions are grounded in what
actually happened.

## Putting It All Together

The full pipeline from raw transcripts to published content:

```bash
# 1. Export all sessions
ctx recall export --all

# 2. In Claude Code: enrich all entries with metadata
/ctx-journal-enrich-all

# 3. Build and serve the journal site
make journal
make journal-serve

# 3b. Or generate an Obsidian vault
ctx journal obsidian

# 4. In Claude Code: draft a blog post
/ctx-blog about the features we shipped this week

# 5. In Claude Code: write a changelog post
/ctx-blog-changelog v0.1.0 "what's new in v0.2.0"
```

The journal pipeline is idempotent at every stage. You can rerun `ctx recall
export --all` without losing enrichment. You can rebuild the site as many times
as you want.

## Tips

* Export regularly. Run `ctx recall export --all` after each session to keep
  your journal current. Only new sessions are exported — existing files are
  skipped by default.
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
* Edit the drafts. These skills produce drafts, not final posts. Review the
  narrative, add your perspective, and remove anything that does not serve the
  reader.

## Next Up

**[Running an Unattended AI Agent →](autonomous-loops.md)**: Set up
an AI agent that works through tasks overnight without you at the
keyboard.

## See Also

* [Session Journal](../reference/session-journal.md): journal system, enrichment schema
* [CLI Reference: ctx recall](../cli/recall.md#ctx-recall): export, list, show session history
* [CLI Reference: ctx journal site](../cli/recall.md#ctx-journal-site): static site generation
* [CLI Reference: ctx journal obsidian](../cli/recall.md#ctx-journal-obsidian): Obsidian vault export
* [CLI Reference: ctx serve](../cli/recall.md#ctx-serve): serve-only (no regeneration)
* [Browsing and Enriching Past Sessions](session-archaeology.md): journal browsing workflow
* [The Complete Session](session-lifecycle.md): capturing context during a session

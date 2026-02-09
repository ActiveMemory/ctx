---
title: Turning Activity into Content
icon: lucide/pen-line
---

![ctx](../images/ctx-banner.png)

## Problem

Your `.context/` directory is full of decisions, learnings, and session
history. Your git log tells the story of a project evolving. But none of
this is visible to anyone outside your terminal. You want to turn this
raw activity into a browsable journal site, blog posts, and changelog
posts.

## Commands and Skills Used

| Tool                     | Type     | Purpose                                        |
|--------------------------|----------|------------------------------------------------|
| `ctx recall export`      | Command  | Export session JSONL to editable markdown      |
| `ctx journal site`       | Command  | Generate a static site from journal entries    |
| `ctx serve`              | Command  | Serve the journal site locally                 |
| `make journal`           | Makefile | Shortcut for export + site rebuild             |
| `/ctx-journal-normalize` | Skill    | Fix markdown rendering in exported entries     |
| `/ctx-journal-enrich`    | Skill    | Add metadata, summaries, and tags to entries   |
| `/ctx-blog`              | Skill    | Draft a blog post from recent project activity |
| `/ctx-blog-changelog`    | Skill    | Write a themed post from a commit range        |

## The Workflow

### Step 1: Export Sessions to Markdown

Raw session data lives as JSONL files in Claude Code's internal storage.
The first step is converting these into readable, editable markdown.

```bash
# Export all sessions from the current project
ctx recall export --all

# Export from all projects (if you work across multiple repos)
ctx recall export --all --all-projects

# Export a single session by ID or slug
ctx recall export abc123
ctx recall export gleaming-wobbling-sutherland
```

Exported files land in `.context/journal/` as individual markdown files
with session metadata and the full conversation transcript. Re-exporting
preserves any YAML frontmatter added by enrichment. Use `--skip-existing`
to leave existing files untouched, or `--force` to overwrite everything.

### Step 2: Normalize Exported Entries

Raw exports can have rendering issues: nested code fences that break
syntax highlighting, metadata blocks that render as raw bold text, and
malformed lists. The `/ctx-journal-normalize` skill fixes these in the
source files before site generation.

```text
/ctx-journal-normalize
```

The skill:

1. Backs up `.context/journal/` before modifying anything
2. Converts `**Key**: value` metadata blocks into collapsible HTML tables
3. Fixes fence nesting so code blocks render with proper highlighting
4. Marks processed files with `<!-- normalized -->` so re-runs skip them

Run normalize before enrich. Clean markdown produces better metadata
extraction in the next step.

### Step 3: Enrich Entries with Metadata

Raw entries have timestamps and conversations but lack the structured
metadata that makes a journal searchable. The `/ctx-journal-enrich`
skill analyzes each conversation and adds semantic frontmatter.

```text
/ctx-journal-enrich twinkly-stirring-kettle
/ctx-journal-enrich 2026-01-24
/ctx-journal-enrich 76fe2ab9
```

The skill reads the conversation, proposes metadata, and asks for
confirmation before writing. After enrichment, an entry gains YAML
frontmatter:

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

This metadata powers better navigation in the journal site: titles
replace slugs, summaries appear in the index, and search covers topics
and technologies.

### Step 4: Generate the Journal Site

With entries exported, normalized, and enriched, generate the static
site:

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
[zensical](https://pypi.org/project/zensical/) for static site
generation (`pip install zensical`).

Or use the Makefile shortcut that combines export and rebuild:

```bash
make journal
```

This runs `ctx recall export --all` followed by `ctx journal site
--build`, then reminds you to normalize and enrich before rebuilding.
To serve the built site: `make journal-serve` or `ctx serve`.

### Step 5: Draft Blog Posts from Activity

When your project reaches a milestone worth sharing, use `/ctx-blog` to
draft a post from your recent activity. The skill gathers context from
multiple sources: git log, DECISIONS.md, LEARNINGS.md, completed tasks,
and journal entries.

```text
/ctx-blog about the caching layer we just built
/ctx-blog last week's refactoring work
/ctx-blog lessons learned from the migration
```

The skill gathers recent commits, decisions, and learnings; identifies a
narrative arc; drafts an outline for approval; writes the full post; and
saves it to `docs/blog/YYYY-MM-DD-slug.md`. Posts are written in first
person with actual code snippets, commit references, and honest
discussion of what went wrong.

### Step 6: Write Changelog Posts from Commit Ranges

For release notes or "what changed" posts, `/ctx-blog-changelog` takes a
starting commit and a theme, then analyzes everything that changed:

```text
/ctx-blog-changelog 040ce99 "building the journal system"
/ctx-blog-changelog HEAD~30 "what's new in v0.2.0"
/ctx-blog-changelog v0.1.0 "the road to v0.2.0"
```

The skill diffs the commit range, identifies the most-changed files,
and constructs a narrative organized by theme rather than chronology,
including a key commits table and before/after comparisons.

## The Conversational Approach

You do not need to remember any of the commands above. When the Agent
Playbook is active, your AI agent tracks what you are working on and
proactively suggests content at natural moments:

> "We just shipped the caching layer and closed 3 tasks. Want me to
> draft a blog post about it?"

> "Your journal has 6 new entries since last rebuild. Want me to
> normalize, enrich, and regenerate the site?"

You can also drive it with natural language instead of skills:

```text
"write about what we did this week"
"turn today's session into a blog post"
"make a changelog post covering everything since the last release"
"enrich the last few journal entries"
```

The agent has full visibility into your `.context/` state — tasks
completed, decisions recorded, learnings captured — so its suggestions
are grounded in what actually happened, not guesswork.

## Putting It Together

The full pipeline from raw transcripts to published content:

```bash
# 1. Export all sessions
ctx recall export --all

# 2. In Claude Code: normalize rendering
/ctx-journal-normalize

# 3. In Claude Code: enrich entries with metadata
/ctx-journal-enrich twinkly-stirring-kettle
/ctx-journal-enrich gleaming-wobbling-sutherland

# 4. Build and serve the journal site
make journal
make journal-serve

# 5. In Claude Code: draft a blog post
/ctx-blog about the features we shipped this week

# 6. In Claude Code: write a changelog post
/ctx-blog-changelog v0.1.0 "what's new in v0.2.0"
```

The journal pipeline is idempotent at every stage. You can re-run
`ctx recall export --all` without losing enrichment. You can re-run
`/ctx-journal-normalize` and it skips already-normalized files. You
can rebuild the site as many times as you want.

## Tips

- **Export regularly.** Run `ctx recall export --all --skip-existing`
  after each session to keep your journal current without re-processing
  old entries.

- **Normalize before enriching.** The enrichment skill reads the
  conversation content to extract metadata. Clean markdown with proper
  formatting produces significantly better results than raw exports
  with broken fences.

- **Enrich selectively.** Not every session needs enrichment.
  Short suggestion sessions and trivial debugging sessions can be
  left as-is. Focus enrichment on sessions where meaningful work
  happened.

- **Keep journal files gitignored.** Session journals contain sensitive
  data: file contents, commands, API keys, internal discussions, and
  error messages with stack traces. The `.context/journal/` and
  `.context/journal-site/` directories must be in `.gitignore`.

- **Use `/ctx-blog` for narrative posts, `/ctx-blog-changelog` for
  release posts.** The blog skill looks at recent activity and finds a
  story. The changelog skill takes a commit range and a theme. They
  complement each other: one for "what I learned" posts, the other
  for "what changed" posts.

- **Let the agent remind you.** You do not need to remember to run
  `/ctx-blog` or `/ctx-journal-enrich`. A proactive agent will suggest
  content generation after productive milestones — shipping a feature,
  closing a batch of tasks, or finishing a long debugging session. The
  best content gets written while the context is fresh.

- **Edit the drafts.** Both blog skills produce drafts, not final
  posts. Review the narrative, add your personal perspective, and
  remove anything that does not serve the reader.

## See Also

- [Session Journal](../session-journal.md): Full documentation of the
  journal system, enrichment schema, and context monitor
- [CLI Reference: ctx recall](../cli-reference.md#ctx-recall):
  Export, list, and show session history
- [CLI Reference: ctx journal](../cli-reference.md#ctx-journal):
  Site generation commands
- [CLI Reference: ctx serve](../cli-reference.md#ctx-serve): Local
  site serving
- [Browsing and Enriching Past Sessions](session-archaeology.md):
  Recipe focused on the journal browsing workflow
- [The Complete Session](session-lifecycle.md): How to capture
  context during a session so there is material to publish later

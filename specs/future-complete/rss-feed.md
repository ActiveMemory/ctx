# RSS Feed for ctx.ist

**Status**: Proposed
**Date**: 2026-02-24

## Motivation

RSS in the ctx universe is not a "user feature" — it is infrastructure.

The question is not "do people still use RSS readers?" but rather:

> Is this a durable, pull-based, machine-readable event log of canonical state?

For ctx, RSS is:

* a replication protocol
* a zero-auth public API
* an append-only timeline
* automation glue

That is **perfectly on brand**.

## Why It Matters Even If Few Humans Subscribe

The real consumers are not casual readers. They are:

### A) Power Users

The exact people who run IRC, self-host things, and automate their
environment. Those people **still use RSS heavily** — and they are
our early adopters.

### B) Future Automation

RSS gives us:

* blog → event stream
* no scraping
* no GitHub API dependency
* no HTML parsing
* no auth

Agents, scripts, CLIs, and dashboards can do:

```
curl feed.xml → act
```

That is extremely ctx-native.

### C) Email Without Email Platforms

RSS → email relay gives us a zero-cost mailing list backend without
ever touching an email address.

## Durability Argument

Think in decades, not trends.

RSS is:

* stable for 20+ years
* trivial to generate
* trivial to consume
* platform-independent

Twitter, Substack, etc. come and go. RSS is still here.

That matches ctx's **durability / archaeology** themes almost
poetically.

## Generation Cost Is ~Zero

We already have:

* structured content (blog posts with frontmatter)
* predictable paths (`docs/blog/YYYY-MM-DD-slug.md`)
* publish events (`git push`)

Generating `feed.xml` is a minimal addition to the build pipeline.

## Workflows It Unlocks

Without building anything else today, RSS makes these possible
**without redesign**:

* `ctx doctor` → warn if blog feed changed since last run
* show latest release narrative inside the CLI
* personal dashboards
* local knowledge mirrors
* cross-site federation later

## Cultural Signaling

Having RSS on a site like ctx.ist signals:

> This is infrastructure, not a content funnel.

To the exact kind of people ctx attracts, that is a strong positive
signal. It says: no tracking, no lock-in, no platform dependency.

## Implementation Strategy

Since zensical does not support RSS natively, generate it at site
build time from the blog index. Keep it static — not on demand per
request.

Pipeline:

```
content → zensical build → feed.xml (generated post-build)
```

A `git push` that includes `./site` automatically deploys it. Zero
effort. One extra file. No runtime logic. No maintenance.

### Approach: Go Subcommand

Add a `ctx site feed` (or similar) subcommand that:

1. Reads `docs/blog/*.md` (or the built `site/blog/` index)
2. Parses frontmatter for title, date, summary, topics
3. Emits a valid Atom 1.0 feed to `site/feed.xml`
4. Runs as a post-build step in the site generation pipeline

Go's `encoding/xml` makes Atom generation trivial (~50-80 lines).

### Minimum Viable Feed

Just:

* title
* URL
* date
* summary (from frontmatter or first paragraph)

That is enough for readers, automation, and email bridges.
Perfection is unnecessary.

### Integration Points

* `Makefile` / `make site`: add feed generation as a post-build step
* `ctx journal site --build`: optionally generate a journal feed too
* `hack/release.sh`: feed regeneration on release

## Resolved Questions

1. **One feed or two?** Blog feed only (`site/feed.xml`). Journal
   content is sensitive — keep it local. Anything worth surfacing
   from other docs can be canonicalized as a blog post linking them.
2. **Feed format**: Atom 1.0. Better spec (required dates, proper
   content types, XML namespace).
3. **Where in the CLI?** `ctx site feed`. The `site` parent command
   leaves room for future subcommands (`ctx site build`, `ctx site
   serve`, etc.).

## Command

```
ctx site feed [flags]
```

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--out` | `-o` | string | `site/feed.xml` | Output path |
| `--base-url` | | string | `https://ctx.ist` | Base URL for entry links |

No `--dry-run` — the command is idempotent (same input → same output).
Re-run to regenerate.

### Output

```
$ ctx site feed

Generated site/feed.xml (18 entries)

Skipped:
  2026-02-25-the-homework-problem.md — not finalized

Warnings:
  2026-02-09-defense-in-depth.md — no summary paragraph found
```

Three buckets: **included** (count), **skipped** (with reason),
**warnings** (included but degraded). Exit 0 always — warnings
don't block generation, they inform the user so they can fix and
re-run.

## Frontmatter Contract

Source: `docs/blog/YYYY-MM-DD-slug.md`

### Required for inclusion

| Field | Type | Feed mapping |
|-------|------|-------------|
| `title` | string | `<title>` |
| `date` | YYYY-MM-DD | `<updated>`, sort order |

Posts missing `title` or `date` are **skipped** with a warning.

### Optional, used if present

| Field | Type | Feed mapping | Fallback |
|-------|------|-------------|----------|
| `author` | string | `<author><name>` | Site-level default ("Context contributors") |
| `topics` | list | `<category term="...">` | Omitted |

### Summary extraction

No frontmatter field. Auto-extracted from the first paragraph after
the `# Heading` line. If no paragraph found, `<summary>` is omitted
and a warning is reported.

### Draft gate

Posts with `reviewed_and_finalized: false` (or absent) are **skipped**.
This is the natural editorial gate — no new field needed. Only
finalized posts appear in the feed.

## Behavior

### Happy path

1. Scan `docs/blog/*.md` for files matching `YYYY-MM-DD-*.md`.
2. Parse YAML frontmatter from each file.
3. Skip drafts (`reviewed_and_finalized` is not `true`).
4. Skip posts missing `title` or `date` (report as warning).
5. Extract summary from first paragraph after `#` heading.
6. Sort entries by date descending (newest first).
7. Generate Atom 1.0 XML.
8. Write to output path.
9. Print report (included, skipped, warnings).

### Edge cases

| Case | Behavior |
|------|----------|
| No blog posts found | Write valid empty feed, print "0 entries" |
| Post without `date` | Skip, report warning |
| Post without `title` | Skip, report warning |
| Post without summary paragraph | Include, omit `<summary>`, report warning |
| `reviewed_and_finalized` absent | Treat as draft, skip |
| `reviewed_and_finalized: false` | Skip |
| Malformed frontmatter | Skip, report warning with parse error |
| Output directory doesn't exist | Create it |
| Duplicate dates | Both included, sorted by filename as tiebreak |
| Non-blog .md files in blog dir | Ignored (filename must match `YYYY-MM-DD-*.md`) |

### Error handling

| Error | Message | Recovery |
|-------|---------|----------|
| `docs/blog/` doesn't exist | "No blog directory found at docs/blog/" | Exit 1 |
| Can't write output file | "Cannot write feed: {err}" | Exit 1 |
| Individual post parse error | Warning in report | Skip post, continue |

## Atom Feed Structure

```xml
<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>ctx blog</title>
  <link href="https://ctx.ist/blog/" />
  <link rel="self" href="https://ctx.ist/feed.xml" />
  <id>https://ctx.ist/feed.xml</id>
  <updated>2026-02-25T00:00:00Z</updated>

  <entry>
    <title>The Dog Ate My Homework</title>
    <link href="https://ctx.ist/blog/2026-02-25-the-homework-problem/" />
    <id>https://ctx.ist/blog/2026-02-25-the-homework-problem/</id>
    <updated>2026-02-25T00:00:00Z</updated>
    <summary>Teaching AI agents to read before they write...</summary>
    <author><name>Jose Alekhinne</name></author>
    <category term="hooks" />
    <category term="agent behavior" />
  </entry>
</feed>
```

Entry URLs derived from filename: `YYYY-MM-DD-slug.md` →
`{base-url}/blog/YYYY-MM-DD-slug/`.

## Implementation

### New files

```
internal/cli/site/
  site.go          — Cmd() parent command ("ctx site")
  feed.go          — feedCmd(), runFeed(), scanBlogPosts(), generateAtom()
  feed_test.go     — tests with constructed blog entries
  atom.go          — Atom XML types (Feed, Entry, Author, Link, Category)
```

### Modified files

| File | Change |
|------|--------|
| `internal/bootstrap/bootstrap.go` | Register `site.Cmd` |

### Key types

```go
// internal/cli/site/atom.go

type AtomFeed struct {
    XMLName xml.Name    `xml:"feed"`
    NS      string      `xml:"xmlns,attr"`
    Title   string      `xml:"title"`
    Links   []AtomLink  `xml:"link"`
    ID      string      `xml:"id"`
    Updated string      `xml:"updated"`
    Entries []AtomEntry `xml:"entry"`
}

type AtomEntry struct {
    Title      string         `xml:"title"`
    Links      []AtomLink     `xml:"link"`
    ID         string         `xml:"id"`
    Updated    string         `xml:"updated"`
    Summary    string         `xml:"summary,omitempty"`
    Author     *AtomAuthor    `xml:"author,omitempty"`
    Categories []AtomCategory `xml:"category,omitempty"`
}

type AtomLink struct {
    Href string `xml:"href,attr"`
    Rel  string `xml:"rel,attr,omitempty"`
}

type AtomAuthor struct {
    Name string `xml:"name"`
}

type AtomCategory struct {
    Term string `xml:"term,attr"`
}
```

```go
// internal/cli/site/feed.go

type blogPost struct {
    filename string
    title    string
    date     string
    author   string
    topics   []string
    summary  string
}

type feedReport struct {
    included int
    skipped  []string // "filename — reason"
    warnings []string // "filename — reason"
}
```

### Helpers to reuse

- `config.RegExFrontmatter` or manual `---` split — frontmatter parsing
- `filepath.Glob` — file scanning
- `encoding/xml` — Atom generation
- `sort.Slice` — date sorting

## Testing

| Test | Scenario |
|------|----------|
| `TestFeed_Basic` | 3 finalized posts → valid Atom XML with 3 entries |
| `TestFeed_SkipsDrafts` | Mix of finalized and draft → only finalized in feed |
| `TestFeed_MissingTitle` | Post without title → skipped, warning reported |
| `TestFeed_MissingDate` | Post without date → skipped, warning reported |
| `TestFeed_NoSummary` | Post without paragraph after heading → included, no `<summary>`, warning |
| `TestFeed_EmptyBlog` | No posts → valid empty feed, 0 entries |
| `TestFeed_SortOrder` | Posts in random order → sorted by date descending |
| `TestFeed_MalformedFrontmatter` | Broken YAML → skipped, warning with parse error |
| `TestFeed_Idempotent` | Run twice → identical output |
| `TestFeed_Categories` | Post with topics → `<category>` elements |
| `TestFeed_CustomBaseURL` | `--base-url` flag → URLs use custom base |
| `TestFeed_FilenameFilter` | Non-matching filenames in blog dir → ignored |

## Non-Goals

- **Journal feed** — journal content is sensitive. Blog only.
- **Full HTML content in feed** — summary only. Readers follow the link.
- **Pagination** — all posts in one feed. At 20-50 posts this is fine.
- **Incremental generation** — always regenerate from scratch. Idempotent.
- **Feed validation service** — use `xmllint` or online validators manually.
- **Auto-discovery `<link>` tag in HTML** — would require modifying zensical
  templates. Add later if needed.

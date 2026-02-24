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

## Open Questions

1. **One feed or two?** Blog feed (`/feed.xml`) is the obvious one.
   A journal feed (`/journal/feed.xml`) could be useful for personal
   use but journal content is sensitive — probably keep it local only.
2. **Feed format**: Atom 1.0 vs RSS 2.0. Atom is the better spec
   (required dates, proper content types, XML namespace). Prefer Atom.
3. **Where in the CLI?** `ctx site feed`, `ctx feed`, or folded into
   `ctx journal site --feed`? TBD based on where it fits cleanest.

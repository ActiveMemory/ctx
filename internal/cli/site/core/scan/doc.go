//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package scan reads **blog-post directories** and turns
// them into the typed records the feed generator and
// blog-changelog skill need: title, slug, date, summary,
// canonical URL.
//
// # Public Surface
//
//   - **[BlogPosts](dir)** — walks `dir` (typically
//     `docs/blog/`), parses every `*.md`'s
//     frontmatter via [internal/parse], and returns
//     a slice of [BlogPost] sorted by `date`
//     descending.
//   - **[ParsePost](path)** — reads one file and
//     returns a [BlogPost]. Used when the caller
//     needs a single post by path.
//   - **[ExtractSummary](body)** — extracts a feed-
//     friendly summary from the post body: the first
//     paragraph, or the explicit `<!-- summary -->`-
//     marked block when present.
//
// # Frontmatter Schema
//
// A blog-post frontmatter must declare:
//
//   - `title` — the post's display title.
//   - `date`  — the publication date (`YYYY-MM-DD`).
//
// Optional: `topics`, `summary`, `author`. Posts
// missing required fields are skipped with a
// warning.
//
// # Concurrency
//
// Filesystem-bound. Concurrent calls each pay the
// full read cost; no module-level cache.
package scan

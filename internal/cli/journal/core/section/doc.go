//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package section builds **topic-based index pages** for the
// journal site: the page that lists every entry tagged with
// `#auth`, the page that lists every entry tagged with
// `#hooks`, and so on. It also assembles the section content
// (collated month/topic groupings) the site renderer drops
// into the navigation tree.
//
// The package is one of three site-rendering helpers — the
// other two are [moc] (Maps of Content for both the site and
// Obsidian) and [generate] (top-level page templates).
//
// # The Surface
//
//   - **[BuildTopicIndex](entries, threshold)** — buckets
//     entries by topic (frontmatter `topics:` list).
//     Topics that appear in fewer than `threshold` entries
//     are folded into a tail "other" bucket so the index
//     stays readable as the journal grows. Returns a
//     [TopicIndex] keyed by canonical topic slug.
//   - **[GenerateTopicsIndex](idx)** — renders the
//     topics-overview page: every topic name + entry
//     count, sorted by popularity descending. Output is
//     a `string` ready to be written to
//     `site/topics/index.md`.
//   - **[GenerateTopicPage](topic, entries)** — renders a
//     single topic's entry list (date + title + slug
//     link). Used per-topic by the site builder to
//     produce `site/topics/<slug>.md`.
//   - **[WriteFormatted](sb, entries)** — appends a
//     formatted entry list into a `*strings.Builder`.
//   - **[WriteMonths](sb, entries)** — appends entries
//     grouped by year-month with month sub-headings.
//
// # Popularity Threshold
//
// The threshold is configurable via the site builder
// invocation; the default is "show topics with 3+
// entries individually, fold the rest into 'other'".
// This is a tunable balance: too low and the topics
// page is dominated by one-off tags; too high and
// long-tail tags disappear entirely.
//
// # Concurrency
//
// All functions are pure data transformations over
// `[]Entry` / topic maps. Concurrent callers never race.
package section

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package generate builds the **top-level pages** of the
// journal site from parsed entries — the README zensical
// reads at build time, the chronological index page, and the
// AI-generated summary insertion that decorates pages already
// produced upstream.
//
// The package is the third leg of the site-building tripod
// alongside [section] (topic indexes) and [moc] (Maps of
// Content). Together they cover everything the journal site
// renders.
//
// # The Surface
//
//   - **[SiteReadme](opts)** — produces the
//     `site/README.md` zensical reads at build time.
//     Embeds the zensical configuration block (theme,
//     navigation, search settings) and the site-wide
//     description. Idempotent: a call with identical
//     `opts` produces byte-identical output.
//   - **[Index](entries)** — produces the chronological
//     index page: entries grouped by month, newest at
//     the top. Output is markdown ready to land at
//     `site/index.md` (or `site/journal/index.md`,
//     depending on layout).
//   - **[InjectedSummary](existing, summary)** — splices
//     an AI-generated summary into existing page
//     content **at a stable insertion point** (a
//     marker comment) so re-running site generation
//     does not duplicate the summary or push other
//     content around. The marker pattern matches what
//     `/ctx-blog` and `/ctx-blog-changelog` skills emit.
//
// # Idempotency Contract
//
// All three generators are idempotent under the same
// inputs. This is what makes `ctx journal site` safe to
// re-run during a CI build: identical entries → identical
// output → no spurious git diffs.
//
// # Concurrency
//
// All functions are pure data transformations. Concurrent
// callers never race.
//
// # Related Packages
//
//   - [internal/cli/journal/cmd/site]   — invokes
//     [SiteReadme] / [Index] when assembling the
//     zensical site.
//   - [internal/cli/journal/core/section] /
//     [internal/cli/journal/core/moc]   — produce the
//     topic and MOC pages that complement the
//     top-level pages this package generates.
//   - [internal/cli/journal/core/normalize]   — runs
//     before generation to sanitize each entry's
//     markdown.
//   - [internal/entity]                 — [Entry],
//     [GroupedIndex] domain types.
package generate

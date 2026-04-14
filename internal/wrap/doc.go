//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package wrap soft-wraps long lines in markdown files to a
// target column width (default 80) without breaking
// markdown semantics — preserving fenced code blocks, tables,
// frontmatter, and list continuation indentation.
//
// The package is what backs `ctx fmt` for journal and context
// files; it is also called by the journal-import pipeline
// before writing enriched entries so reviewers see the same
// shape on disk that they would see in a code review.
//
// # The Three Public Functions
//
//   - **[Soft](line, width)** — wraps a **single** line at
//     word boundaries and returns the resulting `[]string`.
//     Preserves the leading indent of the original line on
//     each continuation. Never breaks inside a word.
//   - **[Content](text, width)** — wraps every line in a
//     **journal entry**. Recognizes YAML frontmatter and
//     skips it (frontmatter values may not be wrapped),
//     skips lines inside fenced code blocks, and skips
//     table rows.
//   - **[ContextFile](text, width)** — same intent as
//     [Content] but tuned for `.context/*.md` files: aware
//     of the markdown list continuation convention
//     (2-space indent for follow-on lines under a
//     bullet) so wrapped continuations look like the
//     original input. Used by `ctx fmt` and the post-add
//     formatter.
//
// # What Stays Unwrapped
//
// The wrap functions deliberately leave several constructs
// alone:
//
//   - **YAML frontmatter** — keys and scalar values must
//     stay on one line.
//   - **Fenced code blocks** (` ``` ` / ` ~~~ `) — code
//     is wrapped by the language, not the markdown
//     renderer.
//   - **Table rows** (lines that match the markdown table
//     pattern) — rewrapping would break column alignment.
//   - **Heading lines** — wrapping a heading mid-phrase
//     would change semantics in many renderers.
//   - **Lines that have no whitespace inside the body
//     beyond the column limit** (e.g. a single long URL)
//     — better to overflow than to break the link.
//
// # Concurrency
//
// All functions are pure: input string → output string.
// Concurrent callers never race.
//
// # Related Packages
//
//   - [internal/cli/fmt]                  — the `ctx fmt`
//     CLI surface.
//   - [internal/cli/journal/core/normalize] — invokes
//     [Content] when normalizing imported journal entries.
//   - [internal/config/wrap]              — column-width
//     constant.
package wrap

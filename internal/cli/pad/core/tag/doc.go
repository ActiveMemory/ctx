//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tag extracts and matches **`#word` tags** in
// scratchpad entries. Tags are convention-based: any
// `#word` token in entry text counts as a tag. The
// package owns the predicates that the `ctx pad` CLI
// uses for filtering and the `ctx pad tags` subcommand
// uses to list every tag in the scratchpad.
//
// # Public Surface
//
//   - **[Extract](text)** — returns every `#word`
//     occurrence in `text` as a `[]string`.
//     De-duplicates and lower-cases.
//   - **[Has](text, tag)** — predicate: does `text`
//     contain `#tag`?
//   - **[Match](entry, query)** — true when `entry`
//     matches the query (single tag).
//   - **[MatchAll](entry, queries)** — true when
//     `entry` matches every tag in `queries` (AND
//     semantics).
//   - **[ScanText](text, fn)** — visitor: invoke
//     `fn` for every tag in `text`.
//
// # Tag Syntax
//
// `#word` where `word` is `[a-z0-9_-]+`. Anchored
// to a word boundary (so `class#1` does not produce
// a `1` tag). Comparison is case-insensitive.
//
// # Concurrency
//
// All functions are pure. Concurrent callers never
// race.
//
// # Related Packages
//
//   - [internal/cli/pad]            — chief consumer
//     for filtering.
//   - [internal/cli/pad/cmd/tags]   — the dedicated
//     tags subcommand.
//   - [internal/cli/pad/core/parse] — supplies the
//     entry stream this package scans.
package tag

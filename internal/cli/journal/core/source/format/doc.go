//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package format provides the **fine-grained formatting
// primitives** used to render a parsed AI session into
// human-readable markdown — part navigation, duration
// strings, tool-call summaries, file references.
//
// The package sits one level below
// [internal/cli/journal/core/source]: this package answers
// "how do I render *this fragment*", `source` answers
// "how do I render the whole session".
//
// # Public Surface
//
//   - **[PartNavigation](currentPart, totalParts, slug)**
//     — generates Previous / Next links for multipart
//     sessions (sessions long enough to be split across
//     several files). Returns markdown ready to splice
//     into the per-part frontmatter.
//   - **[Duration](d)** — formats a `time.Duration` as
//     "23m 14s" / "2h 5m" / "3 days" depending on
//     magnitude. Empty when zero.
//   - **[ToolUse](tu)** — one-line summary of a tool
//     call: tool name, key argument(s), success/error.
//     Used in the per-turn header and in the
//     compressed view.
//   - **[ToolResult](tr)** — one-line summary of the
//     tool's output, truncated to a configurable
//     preview length.
//
// # Local Time vs UTC
//
// Date headers use **local time** so the user sees
// timestamps in their own timezone. UTC is reserved for
// stored timestamps (frontmatter fields) where
// timezone-stable comparisons matter.
//
// # Concurrency
//
// All functions are pure. Concurrent callers never race.
//
// # Related Packages
//
//   - [internal/cli/journal/core/source]   — parent
//     package that orchestrates the per-fragment
//     renderers exposed here.
//   - [internal/cli/journal/core/source/frontmatter] —
//     sister sub-package for YAML frontmatter
//     assembly.
//   - [internal/format]                    — general
//     time/number formatters this package builds on.
//   - [internal/entity]                    — [Session],
//     [Message], [ToolUse], [ToolResult] domain types.
package format

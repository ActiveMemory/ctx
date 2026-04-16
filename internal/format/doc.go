//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package format converts typed Go values into the
// **human-readable display strings** ctx prints in CLI
// output, hook nudges, and journal headers: relative time
// ("3 hours ago"), durations ("23m 14s"), truncated previews,
// and grouped numbers ("1,234,567").
//
// The package is the small, well-tested layer below every
// renderer; centralizing the formatters keeps presentation
// consistent across the CLI and prevents subtle drift like
// "3h ago" in one place vs "3 hours ago" in another.
//
// # Public Surface
//
//   - **[TimeAgo](t)**: relative time vs `now`:
//     "just now", "5 minutes ago", "3 hours ago",
//     "yesterday", "3 days ago", "Mar 12". The
//     break-points and phrasing match what most CLIs
//     have converged on.
//   - **[Duration](d)**: formats a `time.Duration` as
//     "23m 14s" / "2h 5m" / "3d 4h" depending on
//     magnitude. Drops the smaller unit when the
//     larger is ≥ 10 (so "12h 0m" → "12h").
//   - **[DurationAgo](d)**: convenience; takes a
//     duration and renders the [TimeAgo] form for "now
//     minus d".
//   - **[TruncateFirstLine](text, n)**: returns the
//     first line of `text`, truncated to `n` runes
//     (rune-aware, not byte-aware) with an ellipsis
//     when truncation occurs.
//   - **[Number](n)**: thousands-grouped integer
//     formatting ("1,234,567"). Uses comma regardless
//     of locale (ctx is English-only at present).
//
// # Design Choices
//
//   - **Rune-aware truncation**: byte truncation
//     would split multi-byte characters and produce
//     mojibake. [TruncateFirstLine] counts runes.
//   - **Stable break-points**: relative-time
//     phrasing is deterministic per input, so a
//     re-render after a small clock advance does not
//     produce noisy diffs in journal output.
//   - **No localization**: single-locale today;
//     when localization arrives, the per-locale
//     phrase tables will plug in here without
//     changing call sites.
//
// # Concurrency
//
// All functions are pure. Concurrent callers never
// race.
package format

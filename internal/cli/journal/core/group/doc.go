//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package group aggregates journal entries for index
// generation.
//
// The journal index organizes entries by time and by
// topic. This package provides the grouping logic that
// the index renderer consumes.
//
// # Temporal Grouping
//
// [ByMonth] partitions a slice of journal entries by
// their YYYY-MM date prefix. It returns a map keyed by
// month string together with a slice of month strings
// in first-seen order. The cmd layer iterates the
// ordered slice to render month headings while looking
// up entries from the map.
//
// # Topic Aggregation
//
// [GroupedIndex] builds frequency-ranked groups from
// arbitrary keys extracted via a caller-supplied
// function. For every entry the extractor may return
// one or more keys (e.g. tags, key files). The function
// counts how many entries share each key, marks groups
// that meet the popularity threshold (two or more
// sessions) as popular, and sorts the result by count
// descending then alphabetically. The cmd layer uses
// the popularity flag to split output into a "popular"
// section and a long-tail section.
//
// # Data Flow
//
// The cmd/journal layer calls [ByMonth] and
// [GroupedIndex] after loading parsed journal entries
// from disk. Results flow into template rendering in
// the write/journal package.
package group

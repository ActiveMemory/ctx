//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parse holds the small **string-to-typed-value**
// converters that more than one ctx package needs: dates,
// section ranges, frontmatter splits, system-reminder
// stripping, and word-set helpers. Each function is a thin,
// well-tested wrapper around standard-library or
// canonical-format primitives so callers do not have to
// duplicate the same edge-case handling.
//
// # Functions
//
//   - **[Date](s)**: parses `YYYY-MM-DD` into a
//     `time.Time` at midnight UTC. Empty input returns the
//     zero time with no error so callers can branch on
//     `.IsZero()` instead of comparing strings.
//   - **[SplitFrontmatter](data)**: splits a `---`-fenced
//     YAML frontmatter from the markdown body and returns
//     the two byte slices plus a parse error. Used by the
//     skill, steering, and journal-entry parsers.
//   - **[StripSystemReminders](text)**: Claude Code injects
//     `<system-reminder>` tags into tool results that the
//     user did not write. This function strips them so the
//     journal pipeline records what the user actually said.
//   - **[FixCodeFenceSpacing](text)**: users often type
//     `text: ```code` without proper line spacing around
//     the fence; this function normalizes the spacing so
//     the renderer treats it as a code block.
//   - **[WordSet](words)**: builds a `map[string]struct{}`
//     from a slice for O(1) membership; used by the
//     classifier and several lint helpers.
//
// # Why a Shared Package
//
// Every one of these conversions sat in two or three
// places before it was hoisted here. The package's
// existence is enforced by the audit suite: a duplicate
// implementation in another package fails CI.
//
// # Concurrency
//
// All functions are pure and stateless. Concurrent
// callers never race.
package parse

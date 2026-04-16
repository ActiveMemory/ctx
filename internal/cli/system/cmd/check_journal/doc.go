//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_journal implements the
// **`ctx system check-journal`** hidden hook, which
// reminds the user about unprocessed journal entries
// at session start.
//
// # What It Does
//
// The hook checks two sources of journal backlog:
//
//  1. **Unimported sessions** -- counts Claude Code
//     session files (.jsonl) in the projects directory
//     that are newer than the most recent journal entry,
//     indicating sessions that have not been imported
//     into the journal yet.
//  2. **Unenriched entries** -- counts journal entries
//     (.md) that lack frontmatter enrichment (tags,
//     summary, metadata).
//
// When either count is nonzero, the hook emits a nudge
// box showing the counts and suggesting the appropriate
// journal processing commands. The nudge variant adapts
// to which backlog types are present (both, unimported
// only, or unenriched only).
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// On backlog: a formatted nudge box with counts and
// processing hints. On no backlog or throttled: no
// output.
//
// # Throttling
//
// The hook is throttled to fire at most once per day
// using a marker file in the state directory.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble], counts unimported
// and unenriched entries through [core/journal], loads
// the appropriate message template via
// [core/message.Load], and emits the nudge through
// [write/setup.Nudge] and [core/nudge].
package check_journal

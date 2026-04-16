//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package render formats change detection results for
// terminal output and hook injection.
//
// # List
//
// [List] renders the full CLI output for ctx changes.
// It builds a Markdown document with three sections:
//
//   - A reference point label showing the time anchor.
//   - A context changes section listing each modified
//     .context/ file with its modification timestamp.
//   - A code changes section showing commit count,
//     latest commit message, affected directories, and
//     contributing authors.
//
// When no changes are found in either category, List
// emits a "no changes" message instead.
//
// # ChangesForHook
//
// [ChangesForHook] renders the same data in a compact
// single-line format suitable for hook relay injection.
// It concatenates context file names and commit counts
// into a brief summary prefixed with a standard label.
// Returns an empty string when there are no changes,
// allowing the hook layer to skip injection entirely.
//
// # Commit Count Formatting
//
// The unexported commitCount helper formats an integer
// commit count with correct singular/plural text using
// localized templates from the embedded assets. It
// returns "1 commit" or "N commits" as appropriate.
//
// # Data Flow
//
// The cmd/ layer calls List for terminal output or
// ChangesForHook for hook relay mode. Both functions
// receive the reference label, context changes, and
// code summary produced by the scan subpackage.
package render

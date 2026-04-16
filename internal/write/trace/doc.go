//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trace provides terminal output for the
// context trace commands (ctx trace, ctx trace tag,
// ctx trace enable/disable).
//
// The trace system attaches context references to
// git commits and renders commit history with
// resolved context annotations. Output functions
// format commits, references, and hook status.
//
// # Commit Display
//
// [CommitHeader] prints the hash, subject, and date
// for a single commit. [CommitContext] prints the
// "Context:" label before resolved references.
// [CommitNoContext] prints when a commit has no
// attached context.
//
// # File Trace
//
// [FileEntry] prints a single-line trace entry with
// hash, date, subject, and formatted ref summary.
// [LastEntry] prints a compact entry for the
// last-N listing mode.
//
// # Reference Resolution
//
// [Resolved] prints a single resolved context
// reference with its type label, raw value, and
// optional title and detail. It formats the output
// differently depending on whether the reference
// was found and has metadata.
//
// # Tagging and Hooks
//
// [Tagged] confirms a commit was annotated with a
// context note. [HooksEnabled] and [HooksDisabled]
// report trace hook installation and removal.
// [Trailer] prints a collected context trailer
// line when non-empty.
package trace

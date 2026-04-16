//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reindex provides the ctx reindex convenience
// command.
//
// It regenerates quick-reference indices for both
// DECISIONS.md and LEARNINGS.md in a single invocation,
// replacing the need to run ctx decision reindex and
// ctx learning reindex separately. This is the
// recommended way to update indices after adding,
// editing, or removing entries.
//
// # What It Does
//
// The reindex command scans each file for numbered
// entries, extracts their one-line summaries, and
// rebuilds the index table at the top of the file.
// Entry numbering and summary text are synchronized
// with the full entries below.
//
// # Subpackages
//
//	cmd/root: cobra command definition that invokes
//	  both decision and learning reindex in sequence
package reindex

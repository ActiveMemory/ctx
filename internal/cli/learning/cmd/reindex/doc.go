//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reindex implements the "ctx learning reindex"
// command.
//
// # Overview
//
// The reindex command regenerates the index table at the
// top of LEARNINGS.md. It parses all timestamped entry
// headers in the file, sorts them by date, rebuilds the
// index section, and writes the result back in place.
//
// This is useful after manual edits to LEARNINGS.md that
// may have left the index out of sync with the actual
// entries, or after bulk imports that appended entries
// without updating the index.
//
// # Flags
//
// This command accepts no flags.
//
// # Behavior
//
// [Cmd] builds a simple cobra.Command with no flags.
// [Run] resolves the LEARNINGS.md path from the context
// directory and delegates to the shared index.Reindex
// function, which handles parsing, sorting, and
// rewriting the file.
//
// The reindex operation is idempotent: running it
// multiple times on an already-indexed file produces
// the same output.
//
// # Output
//
// Prints a confirmation message to stdout indicating
// that the index was regenerated, along with the number
// of entries found.
package reindex

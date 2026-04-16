//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reindex implements the "ctx decision reindex"
// subcommand that regenerates the index table at the
// top of DECISIONS.md.
//
// # What It Does
//
// Parses all decision entry headers in DECISIONS.md,
// sorts them, and rebuilds the quick-reference index
// table that appears at the top of the file. This is
// useful after manual edits that leave the index out
// of sync with the entries below.
//
// # Arguments
//
// None required. The command operates on the
// DECISIONS.md file in the active context directory.
//
// # Flags
//
// None.
//
// # Output
//
// Prints a confirmation message indicating how many
// entries were indexed. The file is updated in place.
//
// # Delegation
//
// [Cmd] builds the cobra.Command and delegates
// directly to [Run]. Run resolves the file path via
// [rc.ContextDir], then calls [index.Reindex] with
// the decision-specific update function to parse
// headers and regenerate the table.
package reindex

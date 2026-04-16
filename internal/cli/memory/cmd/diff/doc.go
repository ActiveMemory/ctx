//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package diff implements the "ctx memory diff" command.
//
// # Overview
//
// The diff command computes and displays a line-based
// diff between the mirror copy of MEMORY.md (stored in
// .context/memory/) and the current source MEMORY.md.
// This lets the user see what has changed since the last
// sync without performing the sync itself.
//
// # Flags
//
// This command accepts no flags.
//
// # Behavior
//
// [Cmd] builds a simple cobra.Command with no flags.
// [Run] resolves the project root from the context
// directory, discovers the source MEMORY.md path, and
// calls mem.Diff to compute the line-based difference
// between the mirror and the source.
//
// If the source file cannot be discovered (no MEMORY.md
// exists in the project), the command returns a "not
// found" error.
//
// # Output
//
// When differences exist, prints the unified diff to
// stdout. When the mirror and source are identical,
// prints a "no changes" message. The diff format uses
// standard addition/removal markers for easy scanning.
package diff

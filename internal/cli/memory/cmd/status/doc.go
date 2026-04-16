//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package status implements the "ctx memory status"
// command.
//
// # Overview
//
// The status command prints a dashboard of the memory
// bridge state, showing the source MEMORY.md location,
// mirror path, last sync timestamp, line counts for
// both source and mirror, drift detection, and archive
// count.
//
// # Flags
//
// This command accepts no flags.
//
// # Behavior
//
// [Cmd] builds a simple cobra.Command with no flags.
// [Run] resolves the project root, discovers the source
// MEMORY.md, loads the sync state, and prints each
// status section:
//
//  1. Bridge header with source and mirror paths.
//  2. Last sync timestamp with relative duration.
//  3. Source and mirror line counts.
//  4. Drift indicator (whether source differs from
//     mirror).
//  5. Archive count from the memory archive directory.
//
// # Exit Codes
//
//	0    No drift detected; source and mirror match.
//	2    Drift detected; source has changed since the
//	     last sync. This exit code enables scripted
//	     checks in CI or automation.
//
// If the source MEMORY.md cannot be discovered, the
// command prints a "not active" message and returns
// an error.
//
// # Output
//
// Prints a structured status report to stdout with
// labeled fields for each metric. The drift line
// uses a visual indicator for quick scanning.
package status

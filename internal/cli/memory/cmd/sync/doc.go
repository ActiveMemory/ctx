//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync implements the "ctx memory sync" command.
//
// # Overview
//
// The sync command discovers the source MEMORY.md file,
// mirrors it into .context/memory/, archives the
// previous mirror copy, and updates the sync state
// timestamp. This establishes a snapshot that other
// memory commands (diff, import, status) use as their
// baseline.
//
// # Flags
//
//	--dry-run    Report what would happen without
//	             writing any files. Shows source and
//	             mirror paths and whether drift exists.
//
// # Behavior
//
// [Cmd] builds the cobra.Command and registers the
// --dry-run flag. [Run] performs these steps:
//
//  1. Resolves the project root and discovers the
//     source MEMORY.md path.
//  2. In dry-run mode, reports the paths and drift
//     status, then exits.
//  3. In normal mode, calls memory.Sync which copies
//     the source to the mirror directory and archives
//     the previous mirror.
//  4. Loads the sync state, marks it as synced with
//     the current timestamp, and saves it.
//
// If the source MEMORY.md cannot be discovered, the
// command prints a warning and returns a "not found"
// error.
//
// # Output
//
// In dry-run mode, prints the source path, mirror path,
// and drift status. In normal mode, prints the source
// name, mirror path, archive filename, and line counts
// for both source and mirror.
package sync

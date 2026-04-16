//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package snapshot implements the "ctx task snapshot"
// cobra subcommand.
//
// This command creates a point-in-time copy of
// TASKS.md without modifying the original. Snapshots
// are stored in .context/archive/ with timestamped
// filenames, making it easy to compare task state
// across sessions.
//
// # Usage
//
//	ctx task snapshot [name]
//
// # Arguments
//
// An optional positional argument:
//
//   - name: a label for the snapshot. Defaults to
//     "snapshot" when omitted. The name is sanitized
//     to produce a safe filename.
//
// # Behavior
//
// The command:
//
//   - Reads the current TASKS.md content.
//   - Ensures the .context/archive/ directory exists,
//     creating it if necessary.
//   - Generates a filename using the pattern
//     "<name>-<timestamp>.md".
//   - Prepends a header containing the snapshot name
//     and an RFC 3339 timestamp to the content.
//   - Writes the combined content to the archive
//     directory.
//
// Unlike "ctx task archive", this command never
// removes tasks from TASKS.md. It is purely a
// read-only backup operation on the task list.
//
// # Output
//
// Prints the path to the saved snapshot file.
//
// # Delegation
//
// Path resolution uses task/core/path. Filename
// sanitization uses the sanitize package. Content
// formatting and output use write/archive.
package snapshot

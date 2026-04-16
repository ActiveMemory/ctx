//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core contains business logic for the compact
// command, which moves completed tasks from TASKS.md
// into a Completed section and optionally archives them.
//
// This package delegates its work to the task
// subpackage. It does not export functions directly.
//
// # Task Compaction (task/)
//
// The task subpackage provides [task.CompactTasks],
// which orchestrates the full compaction pipeline:
//
//  1. Calls tidy.CompactContext to scan TASKS.md for
//     checked items ("- [x]") outside the Completed
//     section. Only top-level tasks where all nested
//     subtasks are also complete are considered.
//  2. Reports each moved and skipped task via the
//     write/compact output helpers, truncating long
//     descriptions for display.
//  3. Writes the updated TASKS.md content to disk.
//  4. When the --archive flag is set, writes completed
//     task blocks to a dated file in .context/archive/
//     using tidy.WriteArchive.
//
// CompactTasks returns the number of tasks moved and
// an error if the file write fails. A zero count with
// nil error means no completed tasks were found.
//
// # Data Flow
//
// The cmd/ layer loads the context, calls
// task.CompactTasks with the cobra command, loaded
// context, and archive flag. The write/ layer handles
// user-facing output messages. The tidy package owns
// the pure-logic compaction algorithm.
package core

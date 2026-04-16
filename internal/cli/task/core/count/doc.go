//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package count counts pending top-level tasks in
// TASKS.md. Subtasks (indented checkbox lines) are
// excluded from the count so the result reflects only
// the primary work items.
//
// # Counting Logic
//
// [Pending] iterates lines from TASKS.md and counts
// those that match the task regex, are unchecked (pending),
// and are not subtasks. A line is a subtask when it has
// leading whitespace before the checkbox marker.
//
// The function takes pre-split lines rather than a file
// path, so callers can reuse lines they have already
// read for other purposes (e.g., the archive package
// reads TASKS.md once and passes the lines to both
// block parsing and pending counting).
//
// # Usage
//
// The archive package calls Pending to report how many
// tasks remain after archival. The task cmd/ layer uses
// it for the "ctx task count" subcommand.
package count

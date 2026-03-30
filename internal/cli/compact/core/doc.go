//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core contains the pure-logic compaction algorithm for
// the compact command.
//
// [CompactTasks] takes the current TASKS.md content, identifies
// completed top-level tasks (with all children complete), moves
// them to the Completed section, and returns a [CompactResult]
// with no I/O side effects — callers own file writes.
package core

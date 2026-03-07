//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

// TaskStats holds counts of completed and pending tasks.
//
// Used by SeparateTasks to report how many tasks were processed during
// an archive operation.
//
// Fields:
//   - Completed: Number of tasks marked with [x]
//   - Pending: Number of tasks marked with [ ]
type TaskStats struct {
	Completed int
	Pending   int
}

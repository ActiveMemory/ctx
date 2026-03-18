//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// ImportResult tracks per-type counts for memory import operations.
//
// Fields:
//   - Conventions: convention entries imported
//   - Decisions: decision entries imported
//   - Learnings: learning entries imported
//   - Tasks: task entries imported
//   - Skipped: entries skipped (unclassified)
//   - Dupes: duplicate entries skipped
type ImportResult struct {
	Conventions int
	Decisions   int
	Learnings   int
	Tasks       int
	Skipped     int
	Dupes       int
}

// Total returns the number of entries actually imported (excludes
// skips and duplicates).
//
// Returns:
//   - int: sum of conventions, decisions, learnings, and tasks
func (r ImportResult) Total() int {
	return r.Conventions + r.Decisions + r.Learnings + r.Tasks
}

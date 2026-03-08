//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

// ImportResult tracks counts per target for import reporting.
//
// Fields:
//   - Conventions: number of convention entries imported.
//   - Decisions: number of decision entries imported.
//   - Learnings: number of learning entries imported.
//   - Tasks: number of task entries imported.
//   - Skipped: number of entries skipped (unclassified).
//   - Dupes: number of duplicate entries skipped.
type ImportResult struct {
	Conventions int
	Decisions   int
	Learnings   int
	Tasks       int
	Skipped     int
	Dupes       int
}

// Total returns the number of entries actually imported (excludes skips
// and duplicates).
//
// Returns:
//   - int: sum of conventions, decisions, learnings, and tasks.
func (r ImportResult) Total() int {
	return r.Conventions + r.Decisions + r.Learnings + r.Tasks
}

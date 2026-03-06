//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

// importResult tracks counts per target for import reporting.
type importResult struct {
	conventions int
	decisions   int
	learnings   int
	tasks       int
	skipped     int
	dupes       int
}

// total returns the number of entries actually imported (excludes skips
// and duplicates).
func (r importResult) total() int {
	return r.conventions + r.decisions + r.learnings + r.tasks
}

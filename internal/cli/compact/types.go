//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import "time"

// TaskBlock represents a task and its nested content.
//
// Fields:
//   - Lines: All lines in the block (parent and children)
//   - StartIndex: Index of first line in original content
//   - EndIndex: Index of last line (exclusive)
//   - IsCompleted: The parent task is checked
//   - IsArchivable: Completed and no unchecked children
//   - DoneTime: When the task was marked done (from #done: timestamp), nil if not present
type TaskBlock struct {
	Lines        []string
	StartIndex   int
	EndIndex     int
	IsCompleted  bool
	IsArchivable bool
	DoneTime     *time.Time
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

// Pending holds the index and content of a Pending top-level task.
//
// Fields:
//   - Index: Zero-based position in the task list
//   - Content: Full task line text
type Pending struct {
	Index   int
	Content string
}

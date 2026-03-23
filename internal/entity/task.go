//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

import (
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/task"
)

// TaskBlock represents a task and its nested content.
//
// Fields:
//   - Lines: All lines in the block (parent and children)
//   - StartIndex: Index of first line in original content
//   - EndIndex: Index of last line (exclusive)
//   - IsCompleted: The parent task is checked
//   - IsArchivable: Completed and no unchecked children
//   - DoneTime: When the task was marked done (from #done: timestamp),
//     nil if not present
type TaskBlock struct {
	Lines        []string
	StartIndex   int
	EndIndex     int
	IsCompleted  bool
	IsArchivable bool
	DoneTime     *time.Time
}

// BlockContent returns the full content of a block as a single string.
//
// Returns:
//   - string: All lines joined with newlines
func (b *TaskBlock) BlockContent() string {
	return strings.Join(b.Lines, token.NewlineLF)
}

// ParentTaskText extracts just the task text from the parent line.
//
// Returns:
//   - string: Task text without the checkbox prefix, empty if no lines
func (b *TaskBlock) ParentTaskText() string {
	if len(b.Lines) == 0 {
		return ""
	}
	match := regex.Task.FindStringSubmatch(b.Lines[0])
	if match != nil {
		return task.Content(match)
	}
	return ""
}

// OlderThan checks if the task was completed more than the specified days ago.
//
// Parameters:
//   - days: Number of days threshold
//
// Returns:
//   - bool: True if DoneTime is set and older than days ago, false otherwise
func (b *TaskBlock) OlderThan(days int) bool {
	if b.DoneTime == nil {
		return false
	}
	threshold := time.Now().AddDate(0, 0, -days)
	return b.DoneTime.Before(threshold)
}

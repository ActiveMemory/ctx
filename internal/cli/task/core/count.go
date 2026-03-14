//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/task"
)

// CountPendingTasks counts top-level unchecked tasks in the lines.
//
// Parameters:
//   - lines: Lines from TASKS.md to scan
//
// Returns:
//   - int: Number of top-level unchecked tasks
func CountPendingTasks(lines []string) int {
	count := 0
	for _, line := range lines {
		match := regex.Task.FindStringSubmatch(line)
		if match != nil && task.Pending(match) && !task.SubTask(match) {
			count++
		}
	}
	return count
}

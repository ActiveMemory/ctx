//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

import (
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/task"
)

// countPendingTasks counts top-level unchecked tasks in the lines.
//
// Parameters:
//   - lines: Lines from TASKS.md to scan
//
// Returns:
//   - int: Number of top-level unchecked tasks
func countPendingTasks(lines []string) int {
	count := 0
	for _, line := range lines {
		match := config.RegExTask.FindStringSubmatch(line)
		if match != nil && task.Pending(match) && !task.SubTask(match) {
			count++
		}
	}
	return count
}

//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

import (
	"bufio"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/task"
)

// separateTasks parses TASKS.md and separates completed from pending tasks.
//
// The function scans TASKS.md line by line, identifying task items by their
// checkbox markers ([x] for completed, [ ] for pending). It preserves phase
// headers (### Phase ...) in the archived content for traceability.
//
// Subtasks (indented task items) follow their parent task:
//   - Subtasks of completed tasks are archived with the parent
//   - Subtasks of pending tasks remain with the parent
//
// Parameters:
//   - content: Full content of TASKS.md as a string
//
// Returns:
//   - remaining: Content with only pending tasks (to write back to TASKS.md)
//   - archived: Content with completed tasks and their phase headers
//   - stats: Counts of completed and pending tasks processed
func separateTasks(content string) (string, string, taskStats) {
	var remaining strings.Builder
	var archived strings.Builder
	var stats taskStats
	nl := config.NewlineLF

	// Track the current phase header
	var currentPhase string
	var phaseHasArchivedTasks bool
	var phaseArchiveBuffer strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(content))
	var inCompletedTask bool

	for scanner.Scan() {
		line := scanner.Text()

		// Check for phase headers
		if config.RegExPhase.MatchString(line) {
			// Flush previous phase's archived tasks
			if phaseHasArchivedTasks {
				archived.WriteString(currentPhase + nl)
				archived.WriteString(phaseArchiveBuffer.String())
				archived.WriteString(nl)
			}

			currentPhase = line
			phaseHasArchivedTasks = false
			phaseArchiveBuffer.Reset()
			remaining.WriteString(line + nl)
			inCompletedTask = false
			continue
		}

		// Check if the line is a task item
		match := config.RegExTask.FindStringSubmatch(line)
		if match != nil {
			if task.IsSubTask(match) {
				// Handle subtasks - follow their parent
				if inCompletedTask {
					phaseArchiveBuffer.WriteString(line + nl)
				} else {
					remaining.WriteString(line + nl)
				}
				continue
			}

			// Top-level task
			if task.Completed(match) {
				stats.completed++
				phaseHasArchivedTasks = true
				phaseArchiveBuffer.WriteString(line + nl)
				inCompletedTask = true
			} else {
				stats.pending++
				remaining.WriteString(line + nl)
				inCompletedTask = false
			}
			continue
		}

		// Non-task lines go to the remaining
		remaining.WriteString(line + nl)
		inCompletedTask = false
	}

	// Flush final phase's archived tasks
	if phaseHasArchivedTasks {
		archived.WriteString(currentPhase + nl)
		archived.WriteString(phaseArchiveBuffer.String())
	}

	return remaining.String(), archived.String(), stats
}

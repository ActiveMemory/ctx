//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package complete

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/spf13/cobra"

	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/task"
	"github.com/ActiveMemory/ctx/internal/write"
)

// CompleteTask finds a task in TASKS.md by number or text match and marks
// it complete by changing "- [ ]" to "- [x]".
//
// Parameters:
//   - query: Task number (e.g. "1") or search text to match
//   - contextDir: Path to .context/ directory; if empty, uses rc.ContextDir()
//
// Returns:
//   - string: The text of the completed task
//   - error: Non-nil if the task is not found, multiple matches, or file
//     operations fail
func CompleteTask(query, contextDir string) (string, error) {
	if contextDir == "" {
		contextDir = rc.ContextDir()
	}

	filePath := filepath.Join(contextDir, ctx.Task)

	// Check if the file exists
	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		return "", ctxerr.TaskFileNotFound()
	}

	// Read existing content
	content, readErr := os.ReadFile(filepath.Clean(filePath))
	if readErr != nil {
		return "", ctxerr.TaskFileRead(readErr)
	}

	// Parse tasks and find matching one
	lines := strings.Split(string(content), token.NewlineLF)

	var taskNumber int
	isNumber := false
	if num, parseErr := strconv.Atoi(query); parseErr == nil {
		taskNumber = num
		isNumber = true
	}

	currentTaskNum := 0
	matchedLine := -1
	matchedTask := ""

	for i, line := range lines {
		match := regex.Task.FindStringSubmatch(line)
		if match != nil && task.Pending(match) {
			currentTaskNum++
			taskText := task.Content(match)

			// Match by number
			if isNumber && currentTaskNum == taskNumber {
				matchedLine = i
				matchedTask = taskText
				break
			}

			// Match by text (case-insensitive partial match)
			if !isNumber && strings.Contains(
				strings.ToLower(taskText), strings.ToLower(query),
			) {
				if matchedLine != -1 {
					return "", ctxerr.TaskMultipleMatches(query)
				}
				matchedLine = i
				matchedTask = taskText
			}
		}
	}

	if matchedLine == -1 {
		return "", ctxerr.TaskNotFound(query)
	}

	// Mark the task as complete
	lines[matchedLine] = regex.Task.ReplaceAllString(
		lines[matchedLine], regex.TaskCompleteReplace,
	)

	// Write back
	newContent := strings.Join(lines, token.NewlineLF)
	if writeErr := os.WriteFile(filePath, []byte(newContent), fs.PermFile); writeErr != nil {
		return "", ctxerr.TaskFileWrite(writeErr)
	}

	return matchedTask, nil
}

// Run executes the complete command logic.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Command arguments (first arg is the query)
//
// Returns:
//   - error: Non-nil on task match or write failure
func Run(cmd *cobra.Command, args []string) error {
	matchedTask, completeErr := CompleteTask(args[0], "")
	if completeErr != nil {
		return completeErr
	}

	write.InfoCompletedTask(cmd, matchedTask)

	return nil
}

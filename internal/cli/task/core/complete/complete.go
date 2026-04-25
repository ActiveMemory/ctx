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
	errTask "github.com/ActiveMemory/ctx/internal/err/task"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/task"
)

// Complete finds a task in TASKS.md by number or text match and marks
// it complete by changing "- [ ]" to "- [x]".
//
// Parameters:
//   - query: Task number (e.g. "1") or search text to match
//   - contextDir: Path to .context/ directory; if empty, uses rc.ContextDir()
//
// Returns:
//   - string: The text of the completed task
//   - int: The 1-based task number that was matched
//   - error: Non-nil if the task is not found, multiple matches, or file
//     operations fail
func Complete(query, contextDir string) (string, int, error) {
	if contextDir == "" {
		declared, ctxErr := rc.ContextDir()
		if ctxErr != nil {
			return "", 0, ctxErr
		}
		contextDir = declared
	}

	filePath := filepath.Join(contextDir, ctx.Task)

	// Check if the file exists
	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		return "", 0, errTask.FileNotFound()
	}

	// Read existing content
	content, readErr := io.SafeReadUserFile(filepath.Clean(filePath))
	if readErr != nil {
		return "", 0, errTask.FileRead(readErr)
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
	matchedNum := 0

	for i, line := range lines {
		match := regex.Task.FindStringSubmatch(line)
		if match != nil && task.Pending(match) {
			currentTaskNum++
			taskText := task.Content(match)

			// Match by number
			if isNumber && currentTaskNum == taskNumber {
				matchedLine = i
				matchedTask = taskText
				matchedNum = currentTaskNum
				break
			}

			// Match by text (case-insensitive partial match)
			if !isNumber && strings.Contains(
				strings.ToLower(taskText), strings.ToLower(query),
			) {
				if matchedLine != -1 {
					return "", 0, errTask.MultipleMatches(query)
				}
				matchedLine = i
				matchedTask = taskText
				matchedNum = currentTaskNum
			}
		}
	}

	if matchedLine == -1 {
		return "", 0, errTask.NotFound(query)
	}

	// Mark the task as complete
	lines[matchedLine] = regex.Task.ReplaceAllString(
		lines[matchedLine], regex.TaskCompleteReplace,
	)

	// Write back
	newContent := strings.Join(lines, token.NewlineLF)
	if writeErr := io.SafeWriteFile(
		filePath, []byte(newContent), fs.PermFile,
	); writeErr != nil {
		return "", 0, errTask.FileWrite(writeErr)
	}

	return matchedTask, matchedNum, nil
}

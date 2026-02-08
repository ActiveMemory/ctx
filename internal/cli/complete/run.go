//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package complete

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/task"
)

// runComplete executes the complete command logic.
//
// Finds a task in TASKS.md by number or text match and marks it complete
// by changing "- [ ]" to "- [x]".
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - args: Command arguments; args[0] is the task number or search text
//
// Returns:
//   - error: Non-nil if the task is not found, multiple matches, or file
//     operations fail
func runComplete(cmd *cobra.Command, args []string) error {
	query := args[0]

	filePath := filepath.Join(rc.ContextDir(), config.FileTask)

	// Check if the file exists
	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		return fmt.Errorf("TASKS.md not found. Run 'ctx init' first")
	}

	// Read existing content
	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return fmt.Errorf("failed to read TASKS.md: %w", err)
	}

	// Parse tasks and find matching one
	lines := strings.Split(string(content), config.NewlineLF)

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
		match := config.RegExTask.FindStringSubmatch(line)
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
					// Multiple matches - be more specific
					return fmt.Errorf(
						"multiple tasks match %q. Be more specific or use task number",
						query,
					)
				}
				matchedLine = i
				matchedTask = taskText
			}
		}
	}

	if matchedLine == -1 {
		if isNumber {
			return fmt.Errorf(
				"task #%d not found. Use 'ctx status' to see tasks", taskNumber,
			)
		}
		return fmt.Errorf(
			"no task matching %q found. Use 'ctx status' to see tasks", query,
		)
	}

	// Mark the task as complete
	lines[matchedLine] = config.RegExTask.ReplaceAllString(
		lines[matchedLine], "$1- [x] $3",
	)

	// Write back
	newContent := strings.Join(lines, config.NewlineLF)
	if writeErr := os.WriteFile(filePath, []byte(newContent), config.PermFile); writeErr != nil {
		return fmt.Errorf("failed to write TASKS.md: %w", writeErr)
	}

	green := color.New(color.FgGreen).SprintFunc()
	cmd.Println(fmt.Sprintf("%s Completed: %s", green("✓"), matchedTask))

	return nil
}

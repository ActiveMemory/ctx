//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/compact"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/task"
)

// applyFixes attempts to auto-fix issues in the drift report.
//
// Currently, supports fixing:
//   - staleness: Archives completed tasks from TASKS.md
//   - missing_file: Creates missing required files from templates
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - ctx: Loaded context
//   - report: Drift report containing issues to fix
//
// Returns:
//   - *fixResult: Summary of fixes applied
func applyFixes(
	cmd *cobra.Command, ctx *context.Context, report *drift.Report,
) *fixResult {
	result := &fixResult{}
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Process warnings (staleness, missing_file, dead_path)
	for _, issue := range report.Warnings {
		switch issue.Type {
		case drift.IssueStaleness:
			if fixErr := fixStaleness(cmd, ctx); fixErr != nil {
				result.errors = append(result.errors,
					fmt.Sprintf("staleness: %v", fixErr))
			} else {
				cmd.Println(
					fmt.Sprintf(
						"%s Fixed staleness in %s (archived completed tasks)",
						green("✓"), issue.File),
				)
				result.fixed++
			}

		case drift.IssueMissing:
			if fixErr := fixMissingFile(issue.File); fixErr != nil {
				result.errors = append(result.errors,
					fmt.Sprintf("missing %s: %v", issue.File, fixErr))
			} else {
				cmd.Println(
					fmt.Sprintf("%s Created missing file: %s", green("✓"), issue.File),
				)
				result.fixed++
			}

		case drift.IssueDeadPath:
			cmd.Println(fmt.Sprintf("%s Cannot auto-fix dead path in %s:%d (%s)",
				yellow("○"), issue.File, issue.Line, issue.Path))
			result.skipped++

		case drift.IssueStaleAge:
			cmd.Println(fmt.Sprintf("%s Cannot auto-fix file age: %s",
				yellow("○"), issue.File))
			result.skipped++
		}
	}

	// Process violations (potential_secret) - never auto-fix
	for _, issue := range report.Violations {
		if issue.Type == drift.IssueSecret {
			cmd.Println(fmt.Sprintf("%s Cannot auto-fix potential secret: %s",
				yellow("○"), issue.File))
			result.skipped++
		}
	}

	return result
}

// fixStaleness archives completed tasks from TASKS.md.
//
// Moves completed tasks to .context/archive/tasks-YYYY-MM-DD.md and removes
// them from the Completed section in TASKS.md.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - ctx: Loaded context containing the files
//
// Returns:
//   - error: Non-nil if file operations fail
func fixStaleness(cmd *cobra.Command, ctx *context.Context) error {
	tasksFile := ctx.File(config.FileTask)

	if tasksFile == nil {
		return errTasksNotFound()
	}

	nl := config.NewlineLF
	content := string(tasksFile.Content)
	lines := strings.Split(content, nl)

	// Find completed tasks in the Completed section
	var completedTasks []string
	var newLines []string
	inCompletedSection := false

	for _, line := range lines {
		// Track if we're in the Completed section
		if strings.HasPrefix(line, config.HeadingCompleted) {
			inCompletedSection = true
			newLines = append(newLines, line)
			continue
		}
		if strings.HasPrefix(
			line, config.HeadingLevelTwoStart,
		) && inCompletedSection {
			inCompletedSection = false
		}

		// Collect completed tasks from the Completed section for archiving
		match := config.RegExTask.FindStringSubmatch(line)
		if inCompletedSection && match != nil && task.Completed(match) {
			completedTasks = append(completedTasks, task.Content(match))
			continue // Remove from the file
		}

		newLines = append(newLines, line)
	}

	if len(completedTasks) == 0 {
		return errNoCompletedTasks()
	}

	// Build archive content
	var archiveContent string
	for _, t := range completedTasks {
		archiveContent += config.PrefixTaskDone + " " + t + nl
	}

	archiveFile, writeErr := compact.WriteArchive("tasks", config.HeadingArchivedTasks, archiveContent)
	if writeErr != nil {
		return writeErr
	}

	// Write updated TASKS.md
	newContent := strings.Join(newLines, nl)
	if writeErr := os.WriteFile(
		tasksFile.Path, []byte(newContent), config.PermFile,
	); writeErr != nil {
		return errFileWrite(tasksFile.Path, writeErr)
	}

	cmd.Println(fmt.Sprintf("  Archived %d completed tasks to %s",
		len(completedTasks), archiveFile))

	return nil
}

// fixMissingFile creates a missing required context file from template.
//
// Parameters:
//   - filename: Name of the file to create (e.g., "CONSTITUTION.md")
//
// Returns:
//   - error: Non-nil if the template is not found or file write fails
func fixMissingFile(filename string) error {
	content, err := assets.Template(filename)
	if err != nil {
		return errNoTemplate(filename, err)
	}

	targetPath := filepath.Join(rc.ContextDir(), filename)

	// Ensure .context/ directory exists
	if mkErr := os.MkdirAll(rc.ContextDir(), config.PermExec); mkErr != nil {
		return errMkdir(rc.ContextDir(), mkErr)
	}

	if writeErr := os.WriteFile(
		targetPath, content, config.PermFile,
	); writeErr != nil {
		return errFileWrite(targetPath, writeErr)
	}

	return nil
}

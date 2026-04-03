//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package fix

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/read/template"
	"github.com/ActiveMemory/ctx/internal/config/archive"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/entity"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/err/prompt"
	errTask "github.com/ActiveMemory/ctx/internal/err/task"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/task"
	"github.com/ActiveMemory/ctx/internal/tidy"
	writeDrift "github.com/ActiveMemory/ctx/internal/write/drift"
)

// Result tracks fixes applied during drift fix.
//
// Fields:
//   - Fixed: Number of issues successfully fixed
//   - Skipped: Number of issues skipped (not auto-fixable)
//   - Errors: Error messages from failed fix attempts
type Result struct {
	Fixed   int
	Skipped int
	Errors  []string
}

// Apply attempts to auto-fix issues in the drift report.
//
// Currently supports fixing:
//   - staleness: Archives completed tasks from TASKS.md
//   - missing_file: Creates missing required files
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - ctx: Loaded context
//   - report: Drift report containing issues to fix
//
// Returns:
//   - *Result: Summary of fixes applied
func Apply(
	cmd *cobra.Command, ctx *entity.Context,
	report *drift.Report,
) *Result {
	result := &Result{}

	for _, issue := range report.Warnings {
		switch issue.Type {
		case drift.IssueStaleness:
			if fixErr := Staleness(cmd, ctx); fixErr != nil {
				result.Errors = append(result.Errors,
					fmt.Sprintf(
						desc.Text(text.DescKeyDriftFixStalenessErr),
						fixErr,
					))
			} else {
				writeDrift.FixStaleness(cmd, issue.File)
				result.Fixed++
			}

		case drift.IssueMissing:
			if fixErr := MissingFile(issue.File); fixErr != nil {
				result.Errors = append(result.Errors,
					fmt.Sprintf(
						desc.Text(text.DescKeyDriftFixMissingErr),
						issue.File, fixErr,
					))
			} else {
				writeDrift.FixMissing(cmd, issue.File)
				result.Fixed++
			}

		case drift.IssueDeadPath:
			writeDrift.SkipDeadPath(
				cmd, issue.File, issue.Line, issue.Path,
			)
			result.Skipped++

		case drift.IssueStaleAge:
			writeDrift.SkipStaleAge(cmd, issue.File)
			result.Skipped++
		}
	}

	for _, issue := range report.Violations {
		if issue.Type == drift.IssueSecret {
			writeDrift.SkipSensitiveFile(cmd, issue.File)
			result.Skipped++
		}
	}

	return result
}

// Staleness archives completed tasks from TASKS.md.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - ctx: Loaded context containing the files
//
// Returns:
//   - error: Non-nil if file operations fail
func Staleness(
	cmd *cobra.Command, ctx *entity.Context,
) error {
	tasksFile := ctx.File(cfgCtx.Task)

	if tasksFile == nil {
		return errTask.FileNotFound()
	}

	nl := token.NewlineLF
	content := string(tasksFile.Content)
	lines := strings.Split(content, nl)

	var completedTasks []string
	var newLines []string
	inCompletedSection := false

	for _, line := range lines {
		if strings.HasPrefix(
			line,
			desc.Text(text.DescKeyHeadingCompleted),
		) {
			inCompletedSection = true
			newLines = append(newLines, line)
			continue
		}
		if strings.HasPrefix(
			line, token.HeadingLevelTwoStart,
		) && inCompletedSection {
			inCompletedSection = false
		}

		match := regex.Task.FindStringSubmatch(line)
		if inCompletedSection && match != nil &&
			task.Completed(match) {
			completedTasks = append(
				completedTasks, task.Content(match),
			)
			continue
		}

		newLines = append(newLines, line)
	}

	if len(completedTasks) == 0 {
		return errTask.NoneCompleted()
	}

	var archiveContent string
	for _, t := range completedTasks {
		archiveContent += marker.PrefixTaskDone +
			" " + t + nl
	}

	archiveFile, writeErr := tidy.WriteArchive(
		archive.ScopeTasks,
		desc.Text(text.DescKeyHeadingArchivedTasks),
		archiveContent,
	)
	if writeErr != nil {
		return writeErr
	}

	newContent := strings.Join(newLines, nl)
	if writeErr := ctxIo.SafeWriteFile(
		tasksFile.Path, []byte(newContent), fs.PermFile,
	); writeErr != nil {
		return errTask.FileWrite(writeErr)
	}

	writeDrift.Archived(
		cmd, len(completedTasks), archiveFile,
	)

	return nil
}

// MissingFile creates a missing required context file from
// the template.
//
// Parameters:
//   - filename: Name of the file to create
//
// Returns:
//   - error: Non-nil if template not found or write fails
func MissingFile(filename string) error {
	content, err := template.Template(filename)
	if err != nil {
		return prompt.NoTemplate(filename, err)
	}

	targetPath := filepath.Join(rc.ContextDir(), filename)

	if mkErr := ctxIo.SafeMkdirAll(
		rc.ContextDir(), fs.PermExec,
	); mkErr != nil {
		return errFs.Mkdir(rc.ContextDir(), mkErr)
	}

	if writeErr := ctxIo.SafeWriteFile(
		targetPath, content, fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(targetPath, writeErr)
	}

	return nil
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// InfoMovingTask reports a completed task being moved.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: Truncated task description
func InfoMovingTask(cmd *cobra.Command, taskText string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteMovingTask), taskText))
}

// InfoSkippingTask reports a task skipped due to incomplete children.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: Truncated task description
func InfoSkippingTask(cmd *cobra.Command, taskText string) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(
			desc.Text(text.DescKeyTaskArchiveSkipping), taskText,
		),
	)
}

// InfoArchivedTasks reports the number of tasks archived.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of tasks archived
//   - archiveFile: Path to the archive file
//   - days: Age threshold in days
func InfoArchivedTasks(
	cmd *cobra.Command, count int, archiveFile string, days int,
) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(
			desc.Text(text.DescKeyTaskArchiveSuccessWithAge),
			count, archiveFile, days,
		),
	)
}

// ReportHeading prints the compact report heading, separator, and blank line.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func ReportHeading(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyCompactHeading))
	cmd.Println(desc.Text(text.DescKeyCompactSeparator))
	cmd.Println()
}

// TaskError prints a task processing error line.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - err: The error encountered during task processing
func TaskError(cmd *cobra.Command, err error) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyCompactTaskError), err))
}

// SectionsRemoved prints the count of empty sections removed from a file.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - removed: Number of sections removed
//   - fileName: Name of the file that was cleaned
func SectionsRemoved(cmd *cobra.Command, removed int, fileName string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyCompactSectionsRemoved),
		removed, fileName))
}

// ReportClean prints the message shown when no changes were needed.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func ReportClean(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyCompactClean))
}

// ReportSummary prints the final summary with total changes count.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - changes: Total number of changes made
func ReportSummary(cmd *cobra.Command, changes int) {
	if cmd == nil {
		return
	}
	cmd.Println()
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyCompactSummary), changes))
}

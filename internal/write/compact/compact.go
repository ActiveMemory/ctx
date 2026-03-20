//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// InfoMovingTask reports a completed task being moved.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: Truncated task description
func InfoMovingTask(cmd *cobra.Command, taskText string) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteMovingTask), taskText))
}

// InfoSkippingTask reports a task skipped due to incomplete children.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: Truncated task description
func InfoSkippingTask(cmd *cobra.Command, taskText string) {
	cmd.Println(
		fmt.Sprintf(
			desc.TextDesc(text.DescKeyTaskArchiveSkipping), taskText,
		),
	)
}

// InfoArchivedTasks reports the number of tasks archived.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of tasks archived
//   - archiveFile: Path to the archive file
//   - days: Age threshold in days
func InfoArchivedTasks(
	cmd *cobra.Command, count int, archiveFile string, days int,
) {
	cmd.Println(
		fmt.Sprintf(
			desc.TextDesc(text.DescKeyTaskArchiveSuccessWithAge),
			count, archiveFile, days,
		),
	)
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Skipping prints a notice that a task block was skipped due to
// incomplete children.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - taskText: the parent task description.
func Skipping(cmd *cobra.Command, taskText string) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyTaskArchiveSkipping), taskText,
		),
	)
}

// SkipIncomplete prints a summary when no tasks could be archived
// due to incomplete children.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - skippedCount: number of skipped task blocks.
func SkipIncomplete(cmd *cobra.Command, skippedCount int) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyTaskArchiveSkipIncomplete),
			skippedCount,
		),
	)
}

// NoCompleted prints a message when there are no completed tasks
// to archive.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func NoCompleted(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(assets.TextDesc(assets.TextDescKeyTaskArchiveNoCompleted))
}

// DryRun prints the dry-run preview for task archiving.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - archivableCount: number of tasks that would be archived.
//   - pendingCount: number of pending tasks remaining.
//   - preview: the archived content preview string.
//   - separator: the separator string for framing the preview.
func DryRun(
	cmd *cobra.Command,
	archivableCount, pendingCount int, preview, separator string,
) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		assets.TextDesc(assets.TextDescKeyTaskArchiveDryRunBlock),
		archivableCount, pendingCount, separator, preview, separator,
	))
}

// Success prints the result of a successful task archive operation.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - archivedCount: number of tasks archived.
//   - archiveFilePath: path to the created archive file.
//   - pendingCount: the number of pending tasks remaining.
func Success(
	cmd *cobra.Command,
	archivedCount int, archiveFilePath string, pendingCount int,
) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyTaskArchiveSuccess),
			archivedCount, archiveFilePath,
		),
	)
	cmd.Println(
		fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyTaskArchivePendingRemain),
			pendingCount,
		),
	)
}

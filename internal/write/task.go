//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// ArchiveSkipping prints a notice that a task block was skipped due to
// incomplete children.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - taskText: the parent task description.
func ArchiveSkipping(cmd *cobra.Command, taskText string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyTaskArchiveSkipping), taskText))
}

// ArchiveSkipIncomplete prints a summary when no tasks could be archived
// due to incomplete children.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - skippedCount: number of skipped task blocks.
func ArchiveSkipIncomplete(cmd *cobra.Command, skippedCount int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyTaskArchiveSkipIncomplete), skippedCount))
}

// ArchiveNoCompleted prints a message when there are no completed tasks
// to archive.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func ArchiveNoCompleted(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(assets.TextDesc(assets.TextDescKeyTaskArchiveNoCompleted))
}

// ArchiveDryRun prints the dry-run preview for task archiving.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - archivableCount: number of tasks that would be archived.
//   - pendingCount: number of pending tasks remaining.
//   - preview: the archived content preview string.
//   - separator: the separator string for framing the preview.
func ArchiveDryRun(cmd *cobra.Command, archivableCount, pendingCount int, preview, separator string) {
	if cmd == nil {
		return
	}
	cmd.Println(assets.TextDesc(assets.TextDescKeyTaskArchiveDryRunHeader))
	cmd.Println()
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyTaskArchiveDryRunSummary), archivableCount, pendingCount))
	cmd.Println()
	cmd.Println(assets.TextDesc(assets.TextDescKeyTaskArchiveContentPreview))
	cmd.Println(separator)
	cmd.Print(preview)
	cmd.Println(separator)
}

// ArchiveSuccess prints the result of a successful task archive operation.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - archivedCount: number of tasks archived.
//   - archiveFilePath: path to the created archive file.
//   - pendingCount: number of pending tasks remaining.
func ArchiveSuccess(cmd *cobra.Command, archivedCount int, archiveFilePath string, pendingCount int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyTaskArchiveSuccess), archivedCount, archiveFilePath))
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyTaskArchivePendingRemain), pendingCount))
}

// SnapshotSaved prints the result of a successful task snapshot.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - snapshotPath: path to the created snapshot file.
func SnapshotSaved(cmd *cobra.Command, snapshotPath string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyTaskSnapshotSaved), snapshotPath))
}

// SnapshotContent builds the snapshot file content with header and body.
//
// Parameters:
//   - name: snapshot name.
//   - created: RFC3339 formatted creation timestamp.
//   - separator: the separator string.
//   - nl: newline string.
//   - body: the original TASKS.md content.
//
// Returns:
//   - string: formatted snapshot content.
func SnapshotContent(name, created, separator, nl, body string) string {
	return fmt.Sprintf(
		assets.TextDesc(assets.TextDescKeyTaskSnapshotHeaderFormat)+
			nl+nl+
			assets.TextDesc(assets.TextDescKeyTaskSnapshotCreatedFormat)+
			nl+nl+separator+nl+nl+"%s",
		name, created, body,
	)
}

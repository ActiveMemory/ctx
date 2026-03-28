//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cfgArchive "github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/tidy"
	"github.com/ActiveMemory/ctx/internal/write/compact"
)

// CompactTasks moves completed tasks to the "Completed" section in TASKS.md.
//
// Scans TASKS.md for checked items ("- [x]") outside the Completed section,
// including their nested content (indented lines below the task).
// This only moves tasks where all nested subtasks are also complete.
// Optionally archives them to .context/archive/.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - ctx: Loaded context containing the files
//   - archive: If true, write completed tasks to a dated archive file
//
// Returns:
//   - int: Number of tasks moved
//   - error: Non-nil if file write fails
func CompactTasks(
	cmd *cobra.Command, ctx *entity.Context, archive bool,
) (int, error) {
	result := tidy.CompactContext(ctx)

	// Report what happened.
	for _, taskText := range result.TasksMoved {
		compact.InfoMovingTask(cmd, tidy.TruncateString(taskText, token.TruncateLen))
	}
	for _, taskText := range result.TasksSkipped {
		compact.InfoSkippingTask(cmd, tidy.TruncateString(taskText, token.TruncateLen))
	}

	if len(result.TasksMoved) == 0 {
		return 0, nil
	}

	// Write TASKS.md.
	if result.TasksFileUpdate != nil {
		if writeErr := os.WriteFile(
			result.TasksFileUpdate.Path,
			result.TasksFileUpdate.Content,
			fs.PermFile,
		); writeErr != nil {
			return 0, writeErr
		}
	}

	// Archive if requested.
	if archive && len(result.ArchivableBlocks) > 0 {
		archiveDays := rc.ArchiveAfterDays()
		var blocksToArchive []entity.TaskBlock
		for _, block := range result.ArchivableBlocks {
			if block.OlderThan(archiveDays) {
				blocksToArchive = append(blocksToArchive, block)
			}
		}

		if len(blocksToArchive) > 0 {
			nl := token.NewlineLF
			var archiveContent string
			for _, block := range blocksToArchive {
				archiveContent += block.BlockContent() + nl + nl
			}
			if archiveFile, archiveErr := tidy.WriteArchive(
				cfgArchive.ScopeTasks,
				desc.Text(text.DescKeyHeadingArchivedTasks),
				archiveContent,
			); archiveErr == nil {
				compact.InfoArchivedTasks(
					cmd, len(blocksToArchive), archiveFile, archiveDays)
			}
		}
	}

	return len(result.TasksMoved), nil
}

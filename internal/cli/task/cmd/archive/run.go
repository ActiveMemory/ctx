//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import (
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	compactcore "github.com/ActiveMemory/ctx/internal/cli/compact/core"
	"github.com/ActiveMemory/ctx/internal/cli/task/core"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
)

// runArchive executes the archive subcommand logic.
//
// Moves completed tasks (marked with [x]) from TASKS.md to a timestamped
// archive file, including all nested content (subtasks, metadata). Tasks
// with incomplete children are skipped to avoid orphaning pending work.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dryRun: If true, preview changes without modifying files
//
// Returns:
//   - error: Non-nil if TASKS.md doesn't exist or file operations fail
func runArchive(cmd *cobra.Command, dryRun bool) error {
	tasksPath := core.TasksFilePath()
	nl := token.NewlineLF

	// Check if TASKS.md exists
	if _, statErr := os.Stat(tasksPath); os.IsNotExist(statErr) {
		return ctxerr.TaskFileNotFound()
	}

	// Read TASKS.md
	content, readErr := io.SafeReadUserFile(tasksPath)
	if readErr != nil {
		return ctxerr.TaskFileRead(readErr)
	}

	lines := strings.Split(string(content), nl)

	// Parse task blocks using block-based parsing
	blocks := compactcore.ParseTaskBlocks(lines)

	// Filter to only archivable blocks (completed with no incomplete children)
	var archivableBlocks []compactcore.TaskBlock
	var skippedCount int
	for _, block := range blocks {
		if block.IsArchivable {
			archivableBlocks = append(archivableBlocks, block)
		} else {
			skippedCount++
			write.ArchiveSkipping(cmd, block.ParentTaskText())
		}
	}

	// Count pending tasks
	pendingCount := core.CountPendingTasks(lines)

	if len(archivableBlocks) == 0 {
		if skippedCount > 0 {
			write.ArchiveSkipIncomplete(cmd, skippedCount)
		} else {
			write.ArchiveNoCompleted(cmd)
		}
		return nil
	}

	// Build archived content
	var archivedContent strings.Builder
	for _, block := range archivableBlocks {
		archivedContent.WriteString(block.BlockContent())
		archivedContent.WriteString(nl)
	}

	if dryRun {
		write.ArchiveDryRun(cmd, len(archivableBlocks), pendingCount,
			archivedContent.String(), token.Separator)
		return nil
	}

	// Write to archive
	archiveFilePath, writeErr := compactcore.WriteArchive(archive.ArchiveScopeTasks, assets.HeadingArchivedTasks, archivedContent.String())
	if writeErr != nil {
		return writeErr
	}

	// Remove archived blocks from lines and write back
	newLines := compactcore.RemoveBlocksFromLines(lines, archivableBlocks)
	newContent := strings.Join(newLines, nl)

	if updateErr := os.WriteFile(
		tasksPath, []byte(newContent), fs.PermFile,
	); updateErr != nil {
		return ctxerr.TaskFileWrite(updateErr)
	}

	write.ArchiveSuccess(cmd, len(archivableBlocks), archiveFilePath, pendingCount)

	return nil
}

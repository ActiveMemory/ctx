//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	compactcore "github.com/ActiveMemory/ctx/internal/cli/compact/core"
	"github.com/ActiveMemory/ctx/internal/cli/task/core"
	"github.com/ActiveMemory/ctx/internal/config"
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
	nl := config.NewlineLF

	// Check if TASKS.md exists
	if _, statErr := os.Stat(tasksPath); os.IsNotExist(statErr) {
		return fmt.Errorf("no TASKS.md found")
	}

	// Read TASKS.md
	content, readErr := os.ReadFile(filepath.Clean(tasksPath))
	if readErr != nil {
		return fmt.Errorf("failed to read TASKS.md: %w", readErr)
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
			cmd.Println(fmt.Sprintf(
				"! Skipping (has incomplete children): %s",
				block.ParentTaskText(),
			))
		}
	}

	// Count pending tasks
	pendingCount := core.CountPendingTasks(lines)

	if len(archivableBlocks) == 0 {
		if skippedCount > 0 {
			cmd.Println(fmt.Sprintf(
				"No tasks to archive (%d skipped due to incomplete children).",
				skippedCount,
			))
		} else {
			cmd.Println("No completed tasks to archive.")
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
		cmd.Println("Dry run - no files modified")
		cmd.Println()
		cmd.Println(fmt.Sprintf(
			"Would archive %d completed tasks (keeping %d pending)",
			len(archivableBlocks), pendingCount,
		))
		cmd.Println()
		cmd.Println("Archived content preview:")
		cmd.Println(config.Separator)
		cmd.Print(archivedContent.String())
		cmd.Println(config.Separator)
		return nil
	}

	// Write to archive
	archiveFilePath, writeErr := compactcore.WriteArchive("tasks", config.HeadingArchivedTasks, archivedContent.String())
	if writeErr != nil {
		return writeErr
	}

	// Remove archived blocks from lines and write back
	newLines := compactcore.RemoveBlocksFromLines(lines, archivableBlocks)
	newContent := strings.Join(newLines, nl)

	if updateErr := os.WriteFile(
		tasksPath, []byte(newContent), config.PermFile,
	); updateErr != nil {
		return fmt.Errorf("failed to update TASKS.md: %w", updateErr)
	}

	cmd.Println(fmt.Sprintf(
		"✓ Archived %d completed tasks to %s",
		len(archivableBlocks),
		archiveFilePath,
	))
	cmd.Println(fmt.Sprintf("  %d pending tasks remain in TASKS.md", pendingCount))

	return nil
}

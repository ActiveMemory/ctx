//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/compact"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// runTasksSnapshot executes the snapshot subcommand logic.
//
// Creates a point-in-time copy of TASKS.md in the archive directory.
// The snapshot includes a header with the name and timestamp.
//
// Parameters:
//   - cmd: Cobra command (unused, for interface compliance)
//   - args: Optional snapshot name as first argument
//
// Returns:
//   - error: Non-nil if TASKS.md doesn't exist or file operations fail
func runTasksSnapshot(cmd *cobra.Command, args []string) error {
	green := color.New(color.FgGreen).SprintFunc()
	tasksPath := tasksFilePath()
	archivePath := archiveDirPath()

	// Check if TASKS.md exists
	if _, statErr := os.Stat(tasksPath); os.IsNotExist(statErr) {
		return fmt.Errorf("no TASKS.md found")
	}

	// Read TASKS.md
	content, readErr := os.ReadFile(filepath.Clean(tasksPath))
	if readErr != nil {
		return fmt.Errorf("failed to read TASKS.md: %w", readErr)
	}

	// Ensure the archive directory exists
	if mkdirErr := os.MkdirAll(archivePath, config.PermExec); mkdirErr != nil {
		return fmt.Errorf("failed to create archive directory: %w", mkdirErr)
	}

	// Generate snapshot filename
	now := time.Now()
	name := "snapshot"
	if len(args) > 0 {
		name = validation.SanitizeFilename(args[0])
	}
	snapshotFilename := fmt.Sprintf(
		"tasks-%s-%s.md", name, now.Format("2006-01-02-1504"),
	)
	snapshotPath := filepath.Join(archivePath, snapshotFilename)

	// Add snapshot header
	nl := config.NewlineLF
	snapshotContent := fmt.Sprintf(
		"# TASKS.md Snapshot — %s"+
			nl+nl+
			"Created: %s"+nl+nl+config.Separator+nl+nl+"%s",
		name, now.Format(time.RFC3339), string(content),
	)

	// Write snapshot
	if writeErr := os.WriteFile(
		snapshotPath, []byte(snapshotContent), config.PermFile,
	); writeErr != nil {
		return fmt.Errorf("failed to write snapshot: %w", writeErr)
	}

	cmd.Println(fmt.Sprintf("%s Snapshot saved to %s", green("✓"), snapshotPath))

	return nil
}

// runTaskArchive executes the archive subcommand logic.
//
// Moves completed tasks (marked with [x]) from TASKS.md to a timestamped
// archive file, including all nested content (subtasks, metadata). Tasks
// with incomplete children are skipped to avoid orphaning pending work.
//
// Parameters:
//   - cmd: Cobra command (unused, for interface compliance)
//   - dryRun: If true, preview changes without modifying files
//
// Returns:
//   - error: Non-nil if TASKS.md doesn't exist or file operations fail
func runTaskArchive(cmd *cobra.Command, dryRun bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	tasksPath := tasksFilePath()
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
	blocks := compact.ParseTaskBlocks(lines)

	// Filter to only archivable blocks (completed with no incomplete children)
	var archivableBlocks []compact.TaskBlock
	var skippedCount int
	for _, block := range blocks {
		if block.IsArchivable {
			archivableBlocks = append(archivableBlocks, block)
		} else {
			skippedCount++
			cmd.Println(fmt.Sprintf(
				"%s Skipping (has incomplete children): %s",
				yellow("!"), block.ParentTaskText(),
			))
		}
	}

	// Count pending tasks
	pendingCount := countPendingTasks(lines)

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
		cmd.Println(yellow("Dry run - no files modified"))
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
	archiveFilePath, writeErr := compact.WriteArchive("tasks", config.HeadingArchivedTasks, archivedContent.String())
	if writeErr != nil {
		return writeErr
	}

	// Remove archived blocks from lines and write back
	newLines := compact.RemoveBlocksFromLines(lines, archivableBlocks)
	newContent := strings.Join(newLines, nl)

	if updateErr := os.WriteFile(
		tasksPath, []byte(newContent), config.PermFile,
	); updateErr != nil {
		return fmt.Errorf("failed to update TASKS.md: %w", updateErr)
	}

	cmd.Println(fmt.Sprintf(
		"%s Archived %d completed tasks to %s",
		green("✓"),
		len(archivableBlocks),
		archiveFilePath,
	))
	cmd.Println(fmt.Sprintf("  %d pending tasks remain in TASKS.md", pendingCount))

	return nil
}

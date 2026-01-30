//   /    Context:                     https://ctx.ist
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
	if _, err := os.Stat(tasksPath); os.IsNotExist(err) {
		return fmt.Errorf("no TASKS.md found")
	}

	// Read TASKS.md
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		return fmt.Errorf("failed to read TASKS.md: %w", err)
	}

	// Ensure the archive directory exists
	if err := os.MkdirAll(archivePath, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
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
	snapshotContent := fmt.Sprintf(
		"# TASKS.md Snapshot — %s\n\nCreated: %s\n\n---\n\n%s",
		name, now.Format(time.RFC3339), string(content),
	)

	// Write snapshot
	if err := os.WriteFile(
		snapshotPath, []byte(snapshotContent), 0644,
	); err != nil {
		return fmt.Errorf("failed to write snapshot: %w", err)
	}

	fmt.Printf("%s Snapshot saved to %s\n", green("✓"), snapshotPath)

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
	archiveDir := archiveDirPath()

	// Check if TASKS.md exists
	if _, err := os.Stat(tasksPath); os.IsNotExist(err) {
		return fmt.Errorf("no TASKS.md found")
	}

	// Read TASKS.md
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		return fmt.Errorf("failed to read TASKS.md: %w", err)
	}

	lines := strings.Split(string(content), "\n")

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
			fmt.Printf(
				"%s Skipping (has incomplete children): %s\n",
				yellow("!"), block.ParentTaskText(),
			)
		}
	}

	// Count pending tasks
	pendingCount := countPendingTasks(lines)

	if len(archivableBlocks) == 0 {
		if skippedCount > 0 {
			fmt.Printf(
				"No tasks to archive (%d skipped due to incomplete children).\n",
				skippedCount,
			)
		} else {
			fmt.Println("No completed tasks to archive.")
		}
		return nil
	}

	// Build archived content
	var archivedContent strings.Builder
	for _, block := range archivableBlocks {
		archivedContent.WriteString(block.BlockContent())
		archivedContent.WriteString("\n")
	}

	if dryRun {
		fmt.Println(yellow("Dry run - no files modified"))
		fmt.Println()
		fmt.Printf(
			"Would archive %d completed tasks (keeping %d pending)\n",
			len(archivableBlocks), pendingCount,
		)
		fmt.Println()
		fmt.Println("Archived content preview:")
		fmt.Println("---")
		fmt.Print(archivedContent.String())
		fmt.Println("---")
		return nil
	}

	// Ensure the archive directory exists
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	// Generate archive filename
	now := time.Now()
	archiveFilename := fmt.Sprintf("tasks-%s.md", now.Format("2006-01-02"))
	archiveFilePath := filepath.Join(archiveDir, archiveFilename)

	// Check if the archive file already exists for today - append if so
	var finalArchiveContent string
	if existingContent, err := os.ReadFile(archiveFilePath); err == nil {
		finalArchiveContent = string(existingContent) + "\n" + archivedContent.String()
	} else {
		finalArchiveContent = fmt.Sprintf(
			"# Task Archive — %s\n\nArchived from TASKS.md\n\n%s",
			now.Format("2006-01-02"),
			archivedContent.String(),
		)
	}

	// Write the archive file
	if err := os.WriteFile(
		archiveFilePath, []byte(finalArchiveContent), 0644,
	); err != nil {
		return fmt.Errorf("failed to write archive: %w", err)
	}

	// Remove archived blocks from lines and write back
	newLines := compact.RemoveBlocksFromLines(lines, archivableBlocks)
	newContent := strings.Join(newLines, "\n")

	if err := os.WriteFile(tasksPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to update TASKS.md: %w", err)
	}

	fmt.Printf(
		"%s Archived %d completed tasks to %s\n",
		green("✓"),
		len(archivableBlocks),
		archiveFilePath,
	)
	fmt.Printf("  %d pending tasks remain in TASKS.md\n", pendingCount)

	return nil
}

// countPendingTasks counts top-level unchecked tasks in the lines.
//
// Parameters:
//   - lines: Lines from TASKS.md to scan
//
// Returns:
//   - int: Number of top-level unchecked tasks
func countPendingTasks(lines []string) int {
	count := 0
	pattern := compact.UncheckedTaskPattern()
	for _, line := range lines {
		if pattern.MatchString(line) && compact.GetIndentLevel(line) == 0 {
			count++
		}
	}
	return count
}

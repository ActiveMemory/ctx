//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2025-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/session"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	tasksFileName  = ".context/TASKS.md"
	archiveDirName = ".context/archive"
)

// TasksCmd returns the tasks command with subcommands.
func TasksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "Manage task archival and snapshots",
		Long: `Manage task archival and snapshots.

Tasks can be archived to move completed items out of TASKS.md while
preserving them for historical reference. Snapshots create point-in-time
copies without modifying the original.

Subcommands:
  archive   Move completed tasks to timestamped archive file
  snapshot  Create point-in-time snapshot of TASKS.md`,
	}

	cmd.AddCommand(tasksArchiveCmd())
	cmd.AddCommand(tasksSnapshotCmd())

	return cmd
}

var archiveDryRun bool

// tasksArchiveCmd returns the tasks archive subcommand.
func tasksArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "Move completed tasks to timestamped archive file",
		Long: `Move completed tasks from TASKS.md to an archive file.

Archive files are stored in .context/archive/ with timestamped names:
  .context/archive/tasks-YYYY-MM-DD.md

The archive preserves Phase structure for traceability. Completed tasks
(marked with [x]) are moved; pending tasks ([ ]) remain in TASKS.md.

Use --dry-run to preview changes without modifying files.`,
		RunE: runTasksArchive,
	}

	cmd.Flags().BoolVar(&archiveDryRun, "dry-run", false, "Preview changes without modifying files")

	return cmd
}

func runTasksArchive(cmd *cobra.Command, args []string) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Check if TASKS.md exists
	if _, err := os.Stat(tasksFileName); os.IsNotExist(err) {
		return fmt.Errorf("no %s found", tasksFileName)
	}

	// Read TASKS.md
	content, err := os.ReadFile(tasksFileName)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", tasksFileName, err)
	}

	// Parse and separate completed vs pending tasks
	remaining, archived, stats := separateTasks(string(content))

	if stats.completed == 0 {
		fmt.Println("No completed tasks to archive.")
		return nil
	}

	if archiveDryRun {
		fmt.Println(yellow("Dry run - no files modified"))
		fmt.Println()
		fmt.Printf("Would archive %d completed tasks (keeping %d pending)\n", stats.completed, stats.pending)
		fmt.Println()
		fmt.Println("Archived content preview:")
		fmt.Println("---")
		fmt.Println(archived)
		fmt.Println("---")
		return nil
	}

	// Ensure archive directory exists
	if err := os.MkdirAll(archiveDirName, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	// Generate archive filename
	now := time.Now()
	archiveFilename := fmt.Sprintf("tasks-%s.md", now.Format("2006-01-02"))
	archivePath := filepath.Join(archiveDirName, archiveFilename)

	// Check if archive file already exists for today - append if so
	var archiveContent string
	if existingContent, err := os.ReadFile(archivePath); err == nil {
		archiveContent = string(existingContent) + "\n" + archived
	} else {
		archiveContent = fmt.Sprintf("# Task Archive — %s\n\nArchived from TASKS.md\n\n%s", now.Format("2006-01-02"), archived)
	}

	// Write archive file
	if err := os.WriteFile(archivePath, []byte(archiveContent), 0644); err != nil {
		return fmt.Errorf("failed to write archive: %w", err)
	}

	// Write updated TASKS.md
	if err := os.WriteFile(tasksFileName, []byte(remaining), 0644); err != nil {
		return fmt.Errorf("failed to update %s: %w", tasksFileName, err)
	}

	fmt.Printf("%s Archived %d completed tasks to %s\n", green("✓"), stats.completed, archivePath)
	fmt.Printf("  %d pending tasks remain in TASKS.md\n", stats.pending)

	return nil
}

type taskStats struct {
	completed int
	pending   int
}

// separateTasks parses TASKS.md and separates completed from pending tasks.
// Returns: remaining content (pending), archived content (completed), stats
func separateTasks(content string) (string, string, taskStats) {
	var remaining strings.Builder
	var archived strings.Builder
	var stats taskStats

	// Track current phase header
	var currentPhase string
	var phaseHasArchivedTasks bool
	var phaseArchiveBuffer strings.Builder

	completedPattern := regexp.MustCompile(`^\s*-\s*\[x\]`)
	pendingPattern := regexp.MustCompile(`^\s*-\s*\[\s*\]`)
	phasePattern := regexp.MustCompile(`^###\s+Phase`)
	subTaskPattern := regexp.MustCompile(`^\s{2,}-\s*\[`)

	scanner := bufio.NewScanner(strings.NewReader(content))
	var inCompletedTask bool

	for scanner.Scan() {
		line := scanner.Text()

		// Check for phase headers
		if phasePattern.MatchString(line) {
			// Flush previous phase's archived tasks
			if phaseHasArchivedTasks {
				archived.WriteString(currentPhase + "\n")
				archived.WriteString(phaseArchiveBuffer.String())
				archived.WriteString("\n")
			}

			currentPhase = line
			phaseHasArchivedTasks = false
			phaseArchiveBuffer.Reset()
			remaining.WriteString(line + "\n")
			inCompletedTask = false
			continue
		}

		// Check for completed tasks
		if completedPattern.MatchString(line) {
			stats.completed++
			phaseHasArchivedTasks = true
			phaseArchiveBuffer.WriteString(line + "\n")
			inCompletedTask = true
			continue
		}

		// Check for pending tasks
		if pendingPattern.MatchString(line) {
			stats.pending++
			remaining.WriteString(line + "\n")
			inCompletedTask = false
			continue
		}

		// Handle subtasks (indented task items)
		if subTaskPattern.MatchString(line) {
			if inCompletedTask {
				// Subtask of a completed task - archive it
				phaseArchiveBuffer.WriteString(line + "\n")
			} else {
				// Subtask of a pending task - keep it
				remaining.WriteString(line + "\n")
			}
			continue
		}

		// Non-task lines go to remaining
		remaining.WriteString(line + "\n")
		inCompletedTask = false
	}

	// Flush final phase's archived tasks
	if phaseHasArchivedTasks {
		archived.WriteString(currentPhase + "\n")
		archived.WriteString(phaseArchiveBuffer.String())
	}

	return remaining.String(), archived.String(), stats
}

// tasksSnapshotCmd returns the tasks snapshot subcommand.
func tasksSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot [name]",
		Short: "Create point-in-time snapshot of TASKS.md",
		Long: `Create a point-in-time snapshot of TASKS.md without modifying the original.

Snapshots are stored in .context/archive/ with timestamped names:
  .context/archive/tasks-snapshot-YYYY-MM-DD-HHMM.md

Unlike archive, snapshot copies the entire file as-is.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runTasksSnapshot,
	}

	return cmd
}

func runTasksSnapshot(cmd *cobra.Command, args []string) error {
	green := color.New(color.FgGreen).SprintFunc()

	// Check if TASKS.md exists
	if _, err := os.Stat(tasksFileName); os.IsNotExist(err) {
		return fmt.Errorf("no %s found", tasksFileName)
	}

	// Read TASKS.md
	content, err := os.ReadFile(tasksFileName)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", tasksFileName, err)
	}

	// Ensure archive directory exists
	if err := os.MkdirAll(archiveDirName, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	// Generate snapshot filename
	now := time.Now()
	name := "snapshot"
	if len(args) > 0 {
		name = session.sanitizeFilename(args[0])
	}
	snapshotFilename := fmt.Sprintf("tasks-%s-%s.md", name, now.Format("2006-01-02-1504"))
	snapshotPath := filepath.Join(archiveDirName, snapshotFilename)

	// Add snapshot header
	snapshotContent := fmt.Sprintf("# TASKS.md Snapshot — %s\n\nCreated: %s\n\n---\n\n%s",
		name, now.Format(time.RFC3339), string(content))

	// Write snapshot
	if err := os.WriteFile(snapshotPath, []byte(snapshotContent), 0644); err != nil {
		return fmt.Errorf("failed to write snapshot: %w", err)
	}

	fmt.Printf("%s Snapshot saved to %s\n", green("✓"), snapshotPath)

	return nil
}

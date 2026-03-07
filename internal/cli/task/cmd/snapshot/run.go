//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/task/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// runSnapshot executes the snapshot subcommand logic.
//
// Creates a point-in-time copy of TASKS.md in the archive directory.
// The snapshot includes a header with the name and timestamp.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Optional snapshot name as first argument
//
// Returns:
//   - error: Non-nil if TASKS.md doesn't exist or file operations fail
func runSnapshot(cmd *cobra.Command, args []string) error {
	green := color.New(color.FgGreen).SprintFunc()
	tasksPath := core.TasksFilePath()
	archivePath := core.ArchiveDirPath()

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

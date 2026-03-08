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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/task/core"
	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/validation"
	"github.com/ActiveMemory/ctx/internal/write"
)

// Run executes the snapshot subcommand logic.
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
func Run(cmd *cobra.Command, args []string) error {
	tasksPath := core.TasksFilePath()
	archivePath := core.ArchiveDirPath()

	// Check if TASKS.md exists
	if _, statErr := os.Stat(tasksPath); os.IsNotExist(statErr) {
		return ctxerr.TaskFileNotFound()
	}

	// Read TASKS.md
	content, readErr := os.ReadFile(filepath.Clean(tasksPath))
	if readErr != nil {
		return ctxerr.TaskFileRead(readErr)
	}

	// Ensure the archive directory exists
	if mkdirErr := os.MkdirAll(archivePath, config.PermExec); mkdirErr != nil {
		return ctxerr.CreateArchiveDir(mkdirErr)
	}

	// Generate snapshot filename
	now := time.Now()
	name := config.DefaultSnapshotName
	if len(args) > 0 {
		name = validation.SanitizeFilename(args[0])
	}
	snapshotFilename := fmt.Sprintf(
		config.SnapshotFilenameFormat, name, now.Format(config.SnapshotTimeFormat),
	)
	snapshotPath := filepath.Join(archivePath, snapshotFilename)

	// Build snapshot content
	nl := config.NewlineLF
	snapshotContent := write.SnapshotContent(
		name, now.Format(time.RFC3339), config.Separator, nl, string(content),
	)

	// Write snapshot
	if writeErr := os.WriteFile(
		snapshotPath, []byte(snapshotContent), config.PermFile,
	); writeErr != nil {
		return ctxerr.SnapshotWrite(writeErr)
	}

	write.SnapshotSaved(cmd, snapshotPath)

	return nil
}

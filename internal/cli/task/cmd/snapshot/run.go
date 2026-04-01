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

	"github.com/ActiveMemory/ctx/internal/cli/task/core/path"
	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/err/backup"
	errTask "github.com/ActiveMemory/ctx/internal/err/task"
	"github.com/ActiveMemory/ctx/internal/sanitize"
	writeArchive "github.com/ActiveMemory/ctx/internal/write/archive"
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
	tasksPath := path.FilePath()
	archivePath := path.ArchiveDir()

	// Check if TASKS.md exists
	if _, statErr := os.Stat(tasksPath); os.IsNotExist(statErr) {
		return errTask.FileNotFound()
	}

	// Read TASKS.md
	content, readErr := os.ReadFile(filepath.Clean(tasksPath))
	if readErr != nil {
		return errTask.FileRead(readErr)
	}

	// Ensure the archive directory exists
	if mkdirErr := os.MkdirAll(archivePath, fs.PermExec); mkdirErr != nil {
		return backup.CreateArchiveDir(mkdirErr)
	}

	// Generate snapshot filename
	now := time.Now()
	name := archive.DefaultSnapshotName
	if len(args) > 0 {
		name = sanitize.Filename(args[0])
	}
	snapshotFilename := fmt.Sprintf(
		archive.SnapshotFilenameFormat, name, now.Format(archive.SnapshotTimeFormat),
	)
	snapshotPath := filepath.Join(archivePath, snapshotFilename)

	// Build snapshot content
	nl := token.NewlineLF
	snapshotContent := writeArchive.SnapshotContent(
		name, now.Format(time.RFC3339), token.Separator, nl, string(content),
	)

	// Write snapshot
	if writeErr := os.WriteFile( //nolint:gosec // path built from rc.ContextDir + archive dir
		snapshotPath, []byte(snapshotContent), fs.PermFile,
	); writeErr != nil {
		return errTask.SnapshotWrite(writeErr)
	}

	writeArchive.SnapshotSaved(cmd, snapshotPath)

	return nil
}

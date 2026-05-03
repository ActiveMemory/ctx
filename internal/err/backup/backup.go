//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backup

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Create wraps a failure to create a backup (.bak) file during
// ctx init --reset.
//
// Parameters:
//   - name: backup filename that could not be created
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create backup <name>: <cause>"
func Create(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupCreateBackup), name, cause,
	)
}

// CreateArchiveDir wraps a failure to create the archive directory
// under .context/archive/ during task archival.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create archive directory: <cause>"
func CreateArchiveDir(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupCreateArchiveDir),
		cause)
}

// WriteArchive wraps a failure to write an archive file during
// task archival.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to write archive: <cause>"
func WriteArchive(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupWriteArchive),
		cause,
	)
}

// ContextDirNotFound returns an error when the context directory
// does not exist.
//
// Parameters:
//   - dir: the missing context directory path.
//
// Returns:
//   - error: "context directory not found: <dir>: run 'ctx init'"
func ContextDirNotFound(dir string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupContextDirNotFound), dir,
	)
}

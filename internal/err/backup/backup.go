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

// Create wraps a failure to create a backup file.
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

// CreateGeneric wraps a generic backup creation failure.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create backup: <cause>"
func CreateGeneric(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupCreateBackupGeneric),
		cause,
	)
}

// Create wraps an archive creation failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "create archive file: <cause>"
func CreateArchive(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupCreateArchive),
		cause,
	)
}

// CreateArchiveDir wraps a failure to create the archive directory.
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

// WriteArchive wraps a failure to write an archive file.
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

// SMBConfig wraps an SMB configuration parse failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "parse SMB config: <cause>"
func SMBConfig(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupBackupSMBConfig),
		cause,
	)
}

// Project wraps a project backup failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "project backup: <cause>"
func Project(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupBackupProject),
		cause,
	)
}

// Global wraps a global backup failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "global backup: <cause>"
func Global(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupBackupGlobal), cause,
	)
}

// InvalidScope returns an error for an unrecognized backup scope value.
//
// Parameters:
//   - scope: the invalid scope string
//
// Returns:
//   - error: "invalid scope '<scope>': must be project, global, or all"
func InvalidScope(scope string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupInvalidBackupScope), scope,
	)
}

// SourceNotFound returns an error when a backup source path is missing.
//
// Parameters:
//   - path: the missing source path
//
// Returns:
//   - error: "source not found: <path>"
func SourceNotFound(path string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupSourceNotFound), path,
	)
}

// InvalidSMBURL returns an error for a malformed SMB URL.
//
// Parameters:
//   - url: the invalid SMB URL
//
// Returns:
//   - error: "invalid SMB URL: <url>"
func InvalidSMBURL(url string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupInvalidSMBURL), url,
	)
}

// SMBMissingShare returns an error when an SMB URL has no share name.
//
// Parameters:
//   - url: the SMB URL missing a share name
//
// Returns:
//   - error: "SMB URL missing share name: <url>"
func SMBMissingShare(url string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupSMBMissingShare), url,
	)
}

// MountFailed wraps a failure to mount an SMB share.
//
// Parameters:
//   - source: the SMB source URL
//   - cause: the underlying mount error
//
// Returns:
//   - error: "failed to mount <source>: <cause>"
func MountFailed(source string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupMountFailed), source, cause,
	)
}

// WriteSMB wraps a failure to write to an SMB share.
//
// Parameters:
//   - cause: the underlying write error
//
// Returns:
//   - error: "write to SMB: <cause>"
func WriteSMB(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackupWriteSMB), cause,
	)
}

// ContextDirNotFound returns an error when the context directory does not exist.
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

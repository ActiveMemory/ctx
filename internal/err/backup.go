//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import "fmt"

// CreateBackup wraps a failure to create a backup file.
//
// Parameters:
//   - name: backup filename that could not be created
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create backup <name>: <cause>"
func CreateBackup(name string, cause error) error {
	return fmt.Errorf("failed to create backup %s: %w", name, cause)
}

// CreateBackupGeneric wraps a generic backup creation failure.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create backup: <cause>"
func CreateBackupGeneric(cause error) error {
	return fmt.Errorf("failed to create backup: %w", cause)
}

// CreateArchive wraps an archive creation failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "create archive file: <cause>"
func CreateArchive(cause error) error {
	return fmt.Errorf("create archive file: %w", cause)
}

// CreateArchiveDir wraps a failure to create the archive directory.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create archive directory: <cause>"
func CreateArchiveDir(cause error) error {
	return fmt.Errorf("failed to create archive directory: %w", cause)
}

// WriteArchive wraps a failure to write an archive file.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to write archive: <cause>"
func WriteArchive(cause error) error {
	return fmt.Errorf("failed to write archive: %w", cause)
}

// BackupSMBConfig wraps an SMB configuration parse failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "parse SMB config: <cause>"
func BackupSMBConfig(cause error) error {
	return fmt.Errorf("parse SMB config: %w", cause)
}

// BackupProject wraps a project backup failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "project backup: <cause>"
func BackupProject(cause error) error {
	return fmt.Errorf("project backup: %w", cause)
}

// BackupGlobal wraps a global backup failure.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "global backup: <cause>"
func BackupGlobal(cause error) error {
	return fmt.Errorf("global backup: %w", cause)
}

// InvalidBackupScope returns an error for an unrecognized backup scope value.
//
// Parameters:
//   - scope: the invalid scope string
//
// Returns:
//   - error: "invalid scope '<scope>': must be project, global, or all"
func InvalidBackupScope(scope string) error {
	return fmt.Errorf("invalid scope %q: must be project, global, or all", scope)
}

// SourceNotFound returns an error when a backup source path is missing.
//
// Parameters:
//   - path: the missing source path
//
// Returns:
//   - error: "source not found: <path>"
func SourceNotFound(path string) error {
	return fmt.Errorf("source not found: %s", path)
}

// ContextDirNotFound returns an error when the context directory does not exist.
//
// Parameters:
//   - dir: the missing context directory path.
//
// Returns:
//   - error: "context directory not found: <dir> — run 'ctx init'"
func ContextDirNotFound(dir string) error {
	return fmt.Errorf("context directory not found: %s — run 'ctx init'", dir)
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"errors"
	"fmt"
)

// Mkdir wraps a directory creation failure.
//
// Parameters:
//   - desc: human description of the directory (e.g. "journal directory").
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "failed to create <desc>: <cause>"
func Mkdir(desc string, cause error) error {
	return fmt.Errorf("failed to create %s: %w", desc, cause)
}

// ReadDir wraps a directory read failure.
//
// Parameters:
//   - desc: human description of the directory (e.g. "journal directory").
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "read <desc>: <cause>"
func ReadDir(desc string, cause error) error {
	return fmt.Errorf("read %s: %w", desc, cause)
}

// DirNotFound returns an error when a directory does not exist.
//
// Parameters:
//   - dir: the missing directory path.
//
// Returns:
//   - error: "directory not found: <dir>"
func DirNotFound(dir string) error {
	return fmt.Errorf("directory not found: %s", dir)
}

// FileWrite wraps a file write failure.
//
// Parameters:
//   - path: file path that could not be written.
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "failed to write <path>: <cause>"
func FileWrite(path string, cause error) error {
	return fmt.Errorf("failed to write %s: %w", path, cause)
}

// FileRead wraps a file read failure with path context.
//
// Parameters:
//   - path: file path that could not be read.
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "failed to read <path>: <cause>"
func FileRead(path string, cause error) error {
	return fmt.Errorf("failed to read %s: %w", path, cause)
}

// FileAmend wraps a failure to amend an existing file.
//
// Parameters:
//   - path: file path that could not be amended
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to amend <path>: <cause>"
func FileAmend(path string, cause error) error {
	return fmt.Errorf("failed to amend %s: %w", path, cause)
}

// FileUpdate wraps a failure to update a file.
//
// Parameters:
//   - path: file path that could not be updated
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to update <path>: <cause>"
func FileUpdate(path string, cause error) error {
	return fmt.Errorf("failed to update %s: %w", path, cause)
}

// WriteFileFailed wraps a file write failure.
//
// Parameters:
//   - cause: the underlying write error.
//
// Returns:
//   - error: "write file: <cause>"
func WriteFileFailed(cause error) error {
	return fmt.Errorf("write file: %w", cause)
}

// WriteMerged wraps a failure to write a merged file.
//
// Parameters:
//   - path: file path that could not be written
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to write merged <path>: <cause>"
func WriteMerged(path string, cause error) error {
	return fmt.Errorf("failed to write merged %s: %w", path, cause)
}

// OpenFile wraps a file open failure.
//
// Parameters:
//   - path: the file path.
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "open <path>: <cause>"
func OpenFile(path string, cause error) error {
	return fmt.Errorf("open %s: %w", path, cause)
}

// StatPath wraps a stat failure.
//
// Parameters:
//   - path: the path that failed.
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "stat <path>: <cause>"
func StatPath(path string, cause error) error {
	return fmt.Errorf("stat %s: %w", path, cause)
}

// NotDirectory returns an error when a path is not a directory.
//
// Parameters:
//   - path: the path.
//
// Returns:
//   - error: "<path> is not a directory"
func NotDirectory(path string) error {
	return fmt.Errorf("%s is not a directory", path)
}

// ReadDirectory wraps a directory read failure.
//
// Parameters:
//   - path: the directory path.
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "read directory <path>: <cause>"
func ReadDirectory(path string, cause error) error {
	return fmt.Errorf("read directory %s: %w", path, cause)
}

// CreateDir wraps a directory creation failure.
//
// Parameters:
//   - dir: the directory path that could not be created
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to create directory <dir>: <cause>"
func CreateDir(dir string, cause error) error {
	return fmt.Errorf("failed to create directory %s: %w", dir, cause)
}

// BoundaryViolation wraps a boundary validation error with a hint
// to use --allow-outside-cwd.
//
// Parameters:
//   - cause: the underlying validation error
//
// Returns:
//   - error: "<cause>\nUse --allow-outside-cwd to override this check"
func BoundaryViolation(cause error) error {
	return fmt.Errorf("%w\nUse --allow-outside-cwd to override this check", cause)
}

// ReadFile wraps a file read failure.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "read file: <cause>"
func ReadFile(cause error) error {
	return fmt.Errorf("read file: %w", cause)
}

// ReadInput wraps a failure to read user input.
//
// Parameters:
//   - cause: the underlying error from the read operation.
//
// Returns:
//   - error: "failed to read input: <cause>"
func ReadInput(cause error) error {
	return fmt.Errorf("failed to read input: %w", cause)
}

// ReadInputStream wraps a failure to read from the input stream.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "error reading input: <cause>"
func ReadInputStream(cause error) error {
	return fmt.Errorf("error reading input: %w", cause)
}

// NoInput returns an error for missing stdin input.
//
// Returns:
//   - error: "no input received"
func NoInput() error {
	return errors.New("no input received")
}

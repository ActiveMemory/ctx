//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import "fmt"

// ErrTasksNotFound returns an error when TASKS.md is not in the context.
//
// Returns:
//   - error: Descriptive error
func ErrTasksNotFound() error {
	return fmt.Errorf("TASKS.md not found")
}

// ErrNoCompletedTasks returns an error when there are no completed tasks to archive.
//
// Returns:
//   - error: Descriptive error
func ErrNoCompletedTasks() error {
	return fmt.Errorf("no completed tasks to archive")
}

// ErrMkdir wraps a directory creation failure.
//
// Parameters:
//   - path: Directory path that failed
//   - err: Underlying error
//
// Returns:
//   - error: Wrapped error with path context
func ErrMkdir(path string, err error) error {
	return fmt.Errorf("failed to create %s: %w", path, err)
}

// ErrFileWrite wraps a file write failure.
//
// Parameters:
//   - path: File path that failed
//   - err: Underlying error
//
// Returns:
//   - error: Wrapped error with path context
func ErrFileWrite(path string, err error) error {
	return fmt.Errorf("failed to write %s: %w", path, err)
}

// ErrNoTemplate returns an error when no template is available for a file.
//
// Parameters:
//   - filename: Name of the file without a template
//   - err: Underlying error
//
// Returns:
//   - error: Wrapped error with filename context
func ErrNoTemplate(filename string, err error) error {
	return fmt.Errorf("no template available for %s: %w", filename, err)
}

// ErrViolationsFound returns an error when drift violations are detected.
//
// Returns:
//   - error: Descriptive error
func ErrViolationsFound() error {
	return fmt.Errorf("drift detection found violations")
}

// ErrNoContext returns an error when .context/ directory is not found.
//
// Returns:
//   - error: Descriptive error
func ErrNoContext() error {
	return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
}

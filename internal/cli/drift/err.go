//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import "fmt"

// errTasksNotFound returns an error when TASKS.md is not in the context.
func errTasksNotFound() error {
	return fmt.Errorf("TASKS.md not found")
}

// errNoCompletedTasks returns an error when there are no completed tasks to archive.
func errNoCompletedTasks() error {
	return fmt.Errorf("no completed tasks to archive")
}

// errMkdir wraps a directory creation failure.
func errMkdir(path string, err error) error {
	return fmt.Errorf("failed to create %s: %w", path, err)
}

// errFileWrite wraps a file write failure.
func errFileWrite(path string, err error) error {
	return fmt.Errorf("failed to write %s: %w", path, err)
}

// errNoTemplate returns an error when no template is available for a file.
func errNoTemplate(filename string, err error) error {
	return fmt.Errorf("no template available for %s: %w", filename, err)
}

// errViolationsFound returns an error when drift violations are detected.
func errViolationsFound() error {
	return fmt.Errorf("drift detection found violations")
}

// errNoContext returns an error when .context/ directory is not found.
func errNoContext() error {
	return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
}

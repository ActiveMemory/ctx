//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package fmt

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// NoContextDir returns an error when the context directory is missing.
//
// Returns:
//   - error: Suggests running ctx init
func NoContextDir() error {
	return errors.New(desc.Text(text.DescKeyErrFmtNoContextDir))
}

// FileRead wraps a file read failure for a context file.
//
// Parameters:
//   - name: Context filename (e.g., "TASKS.md")
//   - cause: Underlying read error
//
// Returns:
//   - error: "failed to format <name>: <cause>"
func FileRead(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrFmtFileRead), name, cause,
	)
}

// FileWrite wraps a file write failure for a context file.
//
// Parameters:
//   - name: Context filename (e.g., "TASKS.md")
//   - cause: Underlying write error
//
// Returns:
//   - error: "failed to format <name>: <cause>"
func FileWrite(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrFmtFileWrite), name, cause,
	)
}

// NoFiles returns an error when no context files are found.
//
// Parameters:
//   - dir: Context directory path
//
// Returns:
//   - error: "no context files found in <dir>"
func NoFiles(dir string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrFmtNoFiles), dir,
	)
}

// NeedsFormatting returns an error for check mode when files
// need formatting.
//
// Returns:
//   - error: "files need formatting"
func NeedsFormatting() error {
	return errors.New(
		desc.Text(text.DescKeyErrFmtNeedsFormatting),
	)
}

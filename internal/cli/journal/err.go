//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/config"
)

// errNoJournalDir returns an error when the journal directory does not exist.
//
// Parameters:
//   - path: Absolute path to the missing journal directory
//
// Returns:
//   - error: Includes a hint to run 'ctx recall export --all'
func errNoJournalDir(path string) error {
	return fmt.Errorf(
		"no journal directory found at %s"+config.NewlineLF+
			"Run 'ctx recall export --all' first", path)
}

// errScanJournal wraps a journal scanning failure.
//
// Parameters:
//   - err: Underlying scan error
//
// Returns:
//   - error: Wrapped with the context message
func errScanJournal(err error) error {
	return fmt.Errorf("failed to scan journal: %w", err)
}

// errNoEntries returns an error when the journal directory has no entries.
//
// Parameters:
//   - path: Path to the empty journal directory
//
// Returns:
//   - error: Includes a hint to run 'ctx recall export --all'
func errNoEntries(path string) error {
	return fmt.Errorf(
		"no journal entries found in %s"+config.NewlineLF+
			"Run 'ctx recall export --all' first", path)
}

// errMkdir wraps a directory creation failure.
//
// Parameters:
//   - path: Path that could not be created
//   - err: Underlying OS error
//
// Returns:
//   - error: Wrapped with the context message
func errMkdir(path string, err error) error {
	return fmt.Errorf("failed to create %s: %w", path, err)
}

// errFileWrite wraps a file write failure.
//
// Parameters:
//   - path: Path that could not be written
//   - err: Underlying OS error
//
// Returns:
//   - error: Wrapped with the context message
func errFileWrite(path string, err error) error {
	return fmt.Errorf("failed to write %s: %w", path, err)
}

// warnFileErr prints a non-fatal file operation warning to stderr.
//
// Parameters:
//   - cmd: Cobra command (or any type with PrintErrln)
//   - path: Path of the file that caused the warning
//   - err: Underlying error
func warnFileErr(
	cmd interface{ PrintErrln(...any) }, path string, err error,
) {
	cmd.PrintErrln(fmt.Sprintf("  ! %s: %v", path, err))
}

// errZensicalNotFound returns an error when zensical is not installed.
//
// Returns:
//   - error: Includes installation instructions
func errZensicalNotFound() error {
	return fmt.Errorf("zensical not found. Install with: pipx install zensical (requires Python >= 3.10)")
}

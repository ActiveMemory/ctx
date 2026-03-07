//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import "fmt"

// ErrSettingsNotFound returns an error when settings.local.json is missing.
//
// Returns:
//   - error: Descriptive error for missing settings file
func ErrSettingsNotFound() error {
	return fmt.Errorf("no .claude/settings.local.json found")
}

// ErrGoldenNotFound returns an error when settings.golden.json is missing.
//
// Returns:
//   - error: Descriptive error advising the user to run snapshot first
func ErrGoldenNotFound() error {
	return fmt.Errorf("no .claude/settings.golden.json found — run 'ctx permissions snapshot' first")
}

// ErrReadFile wraps a file read failure.
//
// Parameters:
//   - path: File path that failed to read
//   - err: Underlying read error
//
// Returns:
//   - error: Wrapped error with file path context
func ErrReadFile(path string, err error) error {
	return fmt.Errorf("failed to read %s: %w", path, err)
}

// ErrWriteFile wraps a file write failure.
//
// Parameters:
//   - path: File path that failed to write
//   - err: Underlying write error
//
// Returns:
//   - error: Wrapped error with file path context
func ErrWriteFile(path string, err error) error {
	return fmt.Errorf("failed to write %s: %w", path, err)
}

// ErrParseSettings wraps a JSON parse failure for a settings file.
//
// Parameters:
//   - path: File path that failed to parse
//   - err: Underlying parse error
//
// Returns:
//   - error: Wrapped error with file path context
func ErrParseSettings(path string, err error) error {
	return fmt.Errorf("failed to parse %s: %w", path, err)
}

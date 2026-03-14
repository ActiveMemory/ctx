//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SafeReadFile resolves filename within baseDir, verifies the result
// stays within the base directory boundary, and reads the file content.
//
// Unlike [SafeReadUserFile], this function enforces containment: the
// resolved path must remain under baseDir. Use it when the path is
// constructed from a trusted base and a filename component.
//
// Parameters:
//   - baseDir: trusted root directory
//   - filename: file name (or relative path) to join and validate
//
// Returns:
//   - []byte: file content
//   - error: non-nil if resolution fails, path escapes baseDir, or read fails
func SafeReadFile(baseDir, filename string) ([]byte, error) {
	absBase, absErr := filepath.Abs(baseDir)
	if absErr != nil {
		return nil, fmt.Errorf("resolve base: %w", absErr)
	}

	safe := filepath.Join(absBase, filepath.Base(filename))

	if !strings.HasPrefix(safe, absBase+string(os.PathSeparator)) {
		return nil, fmt.Errorf("path escapes base directory: %s", filename)
	}

	data, readErr := os.ReadFile(safe) //nolint:gosec // validated by boundary check above
	if readErr != nil {
		return nil, readErr
	}

	return data, nil
}

// SafeOpenUserFile opens a file for reading after cleaning the path
// and rejecting system directory prefixes.
//
// Parameters:
//   - path: file path to open
//
// Returns:
//   - *os.File: open file handle (caller must close)
//   - error: non-nil on validation or open failure
func SafeOpenUserFile(path string) (*os.File, error) {
	clean, validateErr := cleanAndValidate(path)
	if validateErr != nil {
		return nil, validateErr
	}
	return os.Open(clean) //nolint:gosec // validated by cleanAndValidate
}

// SafeReadUserFile reads a file after cleaning the path and rejecting
// system directory prefixes.
//
// Parameters:
//   - path: file path to read
//
// Returns:
//   - []byte: file content
//   - error: non-nil on validation or read failure
func SafeReadUserFile(path string) ([]byte, error) {
	clean, validateErr := cleanAndValidate(path)
	if validateErr != nil {
		return nil, validateErr
	}
	return os.ReadFile(clean) //nolint:gosec // validated by cleanAndValidate
}

// SafeWriteFile writes data to a file after cleaning the path and
// rejecting system directory prefixes.
//
// Parameters:
//   - path: file path to write
//   - data: content to write
//   - perm: file permission bits
//
// Returns:
//   - error: non-nil on validation or write failure
func SafeWriteFile(path string, data []byte, perm os.FileMode) error {
	clean, validateErr := cleanAndValidate(path)
	if validateErr != nil {
		return validateErr
	}
	return os.WriteFile(clean, data, perm) //nolint:gosec // validated by cleanAndValidate
}

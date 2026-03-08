//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidateBoundary checks that dir resolves to a path within the current
// working directory. Returns an error if the resolved path escapes the
// project root.
func ValidateBoundary(dir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("validate boundary: %w", err)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("validate boundary: %w", err)
	}

	// Resolve symlinks in both paths so traversal via symlinked parents
	// is caught.
	resolvedCwd, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		return fmt.Errorf("validate boundary: %w", err)
	}

	resolvedDir, err := filepath.EvalSymlinks(absDir)
	if err != nil {
		// If the target doesn't exist yet (e.g. before init), fall back
		// to the absolute path for the prefix check.
		resolvedDir = filepath.Clean(absDir)
	}

	// Ensure the resolved dir is equal to or nested under the project root.
	// Append os.PathSeparator to avoid "/foo/bar" matching "/foo/b".
	root := resolvedCwd + string(os.PathSeparator)
	if resolvedDir != resolvedCwd && !strings.HasPrefix(resolvedDir, root) {
		return fmt.Errorf("context directory %q resolves outside project root %q", dir, resolvedCwd)
	}

	return nil
}

// SafeReadFile resolves filename within baseDir, verifies the result stays
// within the base directory boundary, and reads the file content.
//
// Use this instead of raw os.ReadFile when the path is constructed from
// a base directory and a filename component, to prove containment
// statically and avoid per-site nolint directives.
//
// Parameters:
//   - baseDir: Trusted root directory
//   - filename: File name (or relative path) to join and validate
//
// Returns:
//   - []byte: File content
//   - error: Non-nil if resolution fails, path escapes baseDir, or read fails
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

// OpenUserFile opens a file at a user-provided path for reading.
//
// Use this instead of raw os.Open when the path comes directly from
// user input. Centralises the gosec suppression.
//
// Parameters:
//   - path: user-provided file path.
//
// Returns:
//   - *os.File: open file handle (caller must close).
//   - error: non-nil on open failure.
func OpenUserFile(path string) (*os.File, error) {
	clean := filepath.Clean(path)
	return os.Open(clean) //nolint:gosec // user-provided path is intentional
}

// ReadUserFile reads a file at a user-provided path.
//
// Use this instead of raw os.ReadFile when the path comes directly from
// user input (CLI argument, flag, or interactive prompt). Centralises
// the gosec suppression so call sites stay clean.
//
// Parameters:
//   - path: user-provided file path.
//
// Returns:
//   - []byte: file content.
//   - error: non-nil on read failure.
func ReadUserFile(path string) ([]byte, error) {
	clean := filepath.Clean(path)
	return os.ReadFile(clean) //nolint:gosec // user-provided path is intentional
}

// WriteFile writes data to a cleaned file path.
//
// This centralises the gosec suppression for WriteFile calls where the
// path is constructed internally but flagged by the linter.
//
// Parameters:
//   - path: file path (will be cleaned).
//   - data: content to write.
//   - perm: file permission bits.
//
// Returns:
//   - error: non-nil on write failure.
func WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filepath.Clean(path), data, perm) //nolint:gosec // path is internally constructed
}

// CheckSymlinks checks whether dir itself or any of its immediate children
// are symlinks. Returns an error describing the first symlink found.
func CheckSymlinks(dir string) error {
	// Check the directory itself.
	info, err := os.Lstat(dir)
	if err != nil {
		// Non-existent dir is not our concern — let the caller handle it.
		return nil
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("context directory %q is a symlink", dir)
	}

	// Check immediate children.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		child := filepath.Join(dir, entry.Name())
		ci, err := os.Lstat(child)
		if err != nil {
			continue
		}
		if ci.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("context file %q is a symlink", child)
		}
	}

	return nil
}

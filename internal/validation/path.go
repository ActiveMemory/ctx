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
	"runtime"
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
	// On Windows, use case-insensitive comparison since NTFS paths are
	// case-insensitive but EvalSymlinks normalizes casing only for the
	// existing cwd, not the non-existent target — creating a mismatch.
	root := resolvedCwd + string(os.PathSeparator)
	if runtime.GOOS == "windows" {
		if !strings.EqualFold(resolvedDir, resolvedCwd) && !strings.HasPrefix(strings.ToLower(resolvedDir), strings.ToLower(root)) {
			return fmt.Errorf("context directory %q resolves outside project root %q", dir, resolvedCwd)
		}
	} else {
		if resolvedDir != resolvedCwd && !strings.HasPrefix(resolvedDir, root) {
			return fmt.Errorf("context directory %q resolves outside project root %q", dir, resolvedCwd)
		}
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

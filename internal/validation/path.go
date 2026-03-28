//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validation

import (
	"os"
	"path/filepath"
	"strings"

	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
)

// ValidateBoundary checks that dir resolves to a path within the current
// working directory. Returns an error if the resolved path escapes the
// project root.
//
// Parameters:
//   - dir: Directory path to validate
//
// Returns:
//   - error: Non-nil if the path escapes the project root
func ValidateBoundary(dir string) error {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return errFs.BoundaryViolation(cwdErr)
	}

	absDir, absErr := filepath.Abs(dir)
	if absErr != nil {
		return errFs.BoundaryViolation(absErr)
	}

	// Resolve symlinks in both paths so traversal via symlinked parents
	// is caught.
	resolvedCwd, resolveErr := filepath.EvalSymlinks(cwd)
	if resolveErr != nil {
		return errFs.BoundaryViolation(resolveErr)
	}

	resolvedDir, dirResolveErr := filepath.EvalSymlinks(absDir)
	if dirResolveErr != nil {
		// If the target doesn't exist yet (e.g. before init), fall back
		// to the absolute path for the prefix check.
		resolvedDir = filepath.Clean(absDir)
	}

	// Ensure the resolved dir is equal to or nested under the project root.
	// Append os.PathSeparator to avoid "/foo/bar" matching "/foo/b".
	root := resolvedCwd + string(os.PathSeparator)
	if resolvedDir != resolvedCwd && !strings.HasPrefix(resolvedDir, root) {
		return errCtx.OutsideRoot(dir, resolvedCwd)
	}

	return nil
}

// CheckSymlinks checks whether dir itself or any of its immediate children
// are symlinks. Returns an error describing the first symlink found.
//
// Parameters:
//   - dir: Directory path to check for symlinks
//
// Returns:
//   - error: Non-nil if a symlink is found in the directory or its children
func CheckSymlinks(dir string) error {
	// Check the directory itself.
	info, lstatErr := os.Lstat(dir)
	if lstatErr != nil {
		// Non-existent dir is not our concern: let the caller handle it.
		return nil
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return errCtx.DirSymlink(dir)
	}

	// Check immediate children.
	entries, readDirErr := os.ReadDir(dir)
	if readDirErr != nil {
		return nil
	}

	for _, entry := range entries {
		child := filepath.Join(dir, entry.Name())
		ci, childLstatErr := os.Lstat(child)
		if childLstatErr != nil {
			continue
		}
		if ci.Mode()&os.ModeSymlink != 0 {
			return errCtx.FileSymlink(child)
		}
	}

	return nil
}

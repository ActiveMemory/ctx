//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
)

// osWindows is the GOOS value for Windows, extracted to satisfy goconst.
const osWindows = "windows"

// Boundary checks that dir resolves to a path within the current
// working directory. Returns an error if the resolved path escapes the
// project root.
//
// Parameters:
//   - dir: Directory path to validate
//
// Returns:
//   - error: Non-nil if the path escapes the project root
func Boundary(dir string) error {
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

	// On Windows, path comparisons must be case-insensitive because
	// filepath.EvalSymlinks resolves to actual disk casing while
	// os.Getwd preserves the casing from the caller (e.g. VS Code
	// passes a lowercase drive letter via fsPath).
	equal := func(a, b string) bool { return a == b }
	hasPrefix := strings.HasPrefix
	if runtime.GOOS == osWindows {
		equal = strings.EqualFold
		hasPrefix = func(s, prefix string) bool {
			return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
		}
	}

	// Ensure the resolved dir is equal to or nested under the project root.
	// Append os.PathSeparator to avoid "/foo/bar" matching "/foo/b".
	// On Windows, use case-insensitive comparison since NTFS paths are
	// case-insensitive but EvalSymlinks normalizes casing only for the
	// existing cwd, not the non-existent target — creating a mismatch.
	root := resolvedCwd + string(os.PathSeparator)
	if !equal(resolvedDir, resolvedCwd) && !hasPrefix(resolvedDir, root) {
		return errCtx.OutsideRoot(dir, resolvedCwd)
	}

	return nil
}

// Symlinks checks whether dir itself or any of its immediate children
// are symlinks. Returns an error describing the first symlink found.
//
// Parameters:
//   - dir: Directory path to check for symlinks
//
// Returns:
//   - error: Non-nil if a symlink is found in the directory or its children
func Symlinks(dir string) error {
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

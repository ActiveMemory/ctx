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

	errctx "github.com/ActiveMemory/ctx/internal/err/context"
	fserr "github.com/ActiveMemory/ctx/internal/err/fs"
)

// ValidateBoundary checks that dir resolves to a path within the current
// working directory. Returns an error if the resolved path escapes the
// project root.
func ValidateBoundary(dir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fserr.BoundaryViolation(err)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fserr.BoundaryViolation(err)
	}

	// Resolve symlinks in both paths so traversal via symlinked parents
	// is caught.
	resolvedCwd, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		return fserr.BoundaryViolation(err)
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
		return errctx.OutsideRoot(dir, resolvedCwd)
	}

	return nil
}

// CheckSymlinks checks whether dir itself or any of its immediate children
// are symlinks. Returns an error describing the first symlink found.
func CheckSymlinks(dir string) error {
	// Check the directory itself.
	info, err := os.Lstat(dir)
	if err != nil {
		// Non-existent dir is not our concern: let the caller handle it.
		return nil
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return errctx.DirSymlink(dir)
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
			return errctx.FileSymlink(child)
		}
	}

	return nil
}

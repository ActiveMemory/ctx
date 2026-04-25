//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"os"
	"path/filepath"

	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
)

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

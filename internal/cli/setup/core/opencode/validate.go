//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"os"

	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
)

// validateManagedTarget rejects symlinks and non-regular files before
// OpenCode-managed files are read or refreshed in place. Missing files are
// allowed and reported as absent.
//
// Parameters:
//   - targetFile: file path to validate
//
// Returns:
//   - bool: true when the path exists and passed validation.
//   - error: non-nil when stat fails or the path is not a regular file.
func validateManagedTarget(targetFile string) (bool, error) {
	fi, lstatErr := os.Lstat(targetFile)
	if lstatErr != nil {
		if os.IsNotExist(lstatErr) {
			return false, nil
		}
		return false, errFs.StatPath(targetFile, lstatErr)
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		return false, errFs.FileRead(targetFile, os.ErrInvalid)
	}

	if !fi.Mode().IsRegular() {
		return false, errFs.FileRead(targetFile, os.ErrInvalid)
	}

	return true, nil
}

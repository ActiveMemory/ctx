//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	errTrigger "github.com/ActiveMemory/ctx/internal/err/trigger"
)

// parentDir is the relative parent directory component used
// in boundary checks.
const parentDir = ".."

// ValidatePath checks that a hook script path:
//  1. Is not a symlink
//  2. Resolves within the hooksDir boundary
//  3. Has the executable permission bit set
//
// Returns a descriptive error if any check fails.
//
// Parameters:
//   - hooksDir: the root hooks directory (e.g. .context/hooks)
//   - hookPath: the path to the hook script to validate
//
// Returns:
//   - error: non-nil if the path is a symlink, escapes the boundary,
//     or lacks the executable bit
func ValidatePath(hooksDir, hookPath string) error {
	// 1. Symlink check via os.Lstat (does not follow symlinks).
	fi, lstatErr := os.Lstat(hookPath)
	if lstatErr != nil {
		return errTrigger.StatPath(hookPath, lstatErr)
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		return errTrigger.Symlink(hookPath)
	}

	// 2. Boundary check — hookPath must resolve within hooksDir.
	absHooksDir, absHooksDirErr := filepath.Abs(hooksDir)
	if absHooksDirErr != nil {
		return errTrigger.ResolveHooksDir(hooksDir, absHooksDirErr)
	}

	absHookPath, absHookPathErr := filepath.Abs(hookPath)
	if absHookPathErr != nil {
		return errTrigger.ResolvePath(hookPath, absHookPathErr)
	}

	rel, relErr := filepath.Rel(absHooksDir, absHookPath)
	if relErr != nil {
		return errTrigger.Boundary(hookPath, hooksDir)
	}

	sep := string(filepath.Separator)
	if rel == parentDir || len(rel) >= 3 && rel[:3] == parentDir+sep {
		return errTrigger.Boundary(hookPath, hooksDir)
	}

	// 3. Executable permission bit check.
	if fi.Mode().Perm()&fs.ExecBitMask == 0 {
		return errTrigger.NotExecutable(hookPath)
	}

	return nil
}

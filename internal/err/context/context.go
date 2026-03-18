//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// NotFoundError is returned when the context directory does not exist.
type NotFoundError struct {
	Dir string
}

// Error implements the error interface for NotFoundError.
//
// Returns:
//   - string: Error message including the missing directory path
func (e *NotFoundError) Error() string {
	return assets.TextDesc(assets.TextDescKeyErrContextDirNotFound) + e.Dir
}

// NotFound returns a NotFoundError for the given directory.
//
// Parameters:
//   - dir: path to the missing context directory
//
// Returns:
//   - *NotFoundError: typed error for errors.As matching
func NotFound(dir string) *NotFoundError {
	return &NotFoundError{Dir: dir}
}

// OutsideRoot returns an error when .context/ resolves outside the
// project root.
//
// Parameters:
//   - dir: the context directory path
//   - root: the project root path
//
// Returns:
//   - error: "context directory <dir> resolves outside project root <root>"
func OutsideRoot(dir, root string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidateContextOutsideRoot), dir, root,
	)
}

// DirSymlink returns an error when .context/ is a symlink.
//
// Parameters:
//   - dir: the context directory path
//
// Returns:
//   - error: "context directory <dir> is a symlink"
func DirSymlink(dir string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidateContextDirSymlink), dir,
	)
}

// FileSymlink returns an error when a file inside .context/ is a
// symlink.
//
// Parameters:
//   - file: the symlinked file path
//
// Returns:
//   - error: "context file <file> is a symlink"
func FileSymlink(file string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidateContextFileSymlink), file,
	)
}

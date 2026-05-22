//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgRc "github.com/ActiveMemory/ctx/internal/config/rc"
	"github.com/ActiveMemory/ctx/internal/entity"
)

const (
	// ErrNoCtxHere is the sentinel returned by rc.ContextDir
	// when $PWD/.context/ does not exist as a directory. Under the
	// cwd-anchored resolution model
	// (spec: specs/cwd-anchored-context.md), this is the canonical
	// "this process is not at a ctx project root" signal. Callers
	// that can legitimately proceed without a project (init,
	// bootstrap diagnostics) check with errors.Is; everyone else
	// should propagate. Wrap via [NoCtxHere] for user-facing
	// messages so the offending cwd is shown.
	ErrNoCtxHere = entity.Sentinel(
		text.DescKeyErrContextNoContextHereMsg,
	)

	// ErrContextDirNotADirectory is the sentinel returned when
	// $PWD/.context exists but is not a directory (typically a
	// regular file). Symlinks pointing at directories pass.
	ErrContextDirNotADirectory = entity.Sentinel(
		text.DescKeyErrContextDirNotADirectoryMsg,
	)

	// ErrContextDirStat is the sentinel returned when [os.Stat]
	// on $PWD/.context fails for a reason other than not-exist
	// (permission denied, I/O error). Wrap via [StatFailed] to
	// attach the underlying cause.
	ErrContextDirStat = entity.Sentinel(
		text.DescKeyErrContextDirStatMsg,
	)

	// ErrNotInitialized is the sentinel returned when $PWD/.context
	// exists as a directory but the project lacks the required
	// context files (i.e., `ctx init` has not run there). Distinct
	// from [ErrNoCtxHere] (no .context at all) and from
	// [ErrContextDirNotADirectory] (path exists but is the wrong
	// type). Wrap via [NotInitialized] for user-facing messages so
	// the offending path is shown.
	ErrNotInitialized = entity.Sentinel(
		text.DescKeyErrContextNotInitializedMsg,
	)
)

// NoCtxHere wraps [ErrNoCtxHere] with the offending cwd so
// the user sees where ctx was looking.
//
// Parameters:
//   - cwd: the working directory that lacked .context/
//
// Returns:
//   - error: wrapping [ErrNoCtxHere] for [errors.Is] matches
func NoCtxHere(cwd string) error {
	return fmt.Errorf(cfgRc.FmtWrapColon,
		ErrNoCtxHere,
		fmt.Sprintf(desc.Text(text.DescKeyErrContextNoContextHere), cwd),
	)
}

// NotADir wraps [ErrContextDirNotADirectory] with the offending
// path so the user sees what was rejected.
//
// Parameters:
//   - path: absolute path that exists but is not a directory
//
// Returns:
//   - error: wrapping [ErrContextDirNotADirectory] for [errors.Is]
func NotADir(path string) error {
	return fmt.Errorf(cfgRc.FmtWrapColon,
		ErrContextDirNotADirectory,
		fmt.Sprintf(desc.Text(text.DescKeyErrContextDirNotADirectory), path),
	)
}

// StatFailed wraps [ErrContextDirStat] with the path and the
// underlying [os.Stat] failure.
//
// Parameters:
//   - path: absolute path that failed to stat
//   - cause: the underlying stat error
//
// Returns:
//   - error: wrapping both [ErrContextDirStat] and the underlying
//     cause; supports [errors.Is] for either
func StatFailed(path string, cause error) error {
	return fmt.Errorf(cfgRc.FmtWrapColon,
		ErrContextDirStat,
		fmt.Errorf(desc.Text(text.DescKeyErrContextDirStat), path, cause),
	)
}

// NotInitialized wraps [ErrNotInitialized] with the offending
// directory so the user sees which project is not initialized.
//
// Parameters:
//   - path: absolute path to the (existing but uninitialized) context dir
//
// Returns:
//   - error: wrapping [ErrNotInitialized] for [errors.Is] matches
func NotInitialized(path string) error {
	return fmt.Errorf(cfgRc.FmtWrapColon,
		ErrNotInitialized,
		fmt.Sprintf(desc.Text(text.DescKeyErrContextNotInitialized), path),
	)
}

// NotFoundError is returned by the context loader when the loaded
// directory has no readable files. Distinct from a bare
// [ErrNoCtxHere]: the directory exists but its contents are
// empty or unreadable. Callers using `errors.AsType[*NotFoundError]`
// can still recover the missing path; callers using
// `errors.Is(err, ErrNoCtxHere)` will also match via the
// type's [NotFoundError.Is] method.
type NotFoundError struct {
	Dir string
}

// Error implements the error interface for NotFoundError.
//
// Returns:
//   - string: Error message including the missing directory path
func (e *NotFoundError) Error() string {
	return desc.Text(text.DescKeyErrContextDirNotFound) + e.Dir
}

// Is reports whether target matches the no-context-here sentinel.
// Lets callers using `errors.Is(err, ErrNoCtxHere)` match
// instances of [NotFoundError] without rewriting them.
//
// Parameters:
//   - target: error to compare against
//
// Returns:
//   - bool: true when target is [ErrNoCtxHere]
func (e *NotFoundError) Is(target error) bool {
	return target == ErrNoCtxHere
}

// NotFound returns a NotFoundError for the given directory.
//
// Parameters:
//   - path: path to the missing context directory
//
// Returns:
//   - *NotFoundError: typed error for errors.AsType matching
func NotFound(path string) *NotFoundError {
	return &NotFoundError{Dir: path}
}

// DirSymlink returns an error when .context/ is a symlink.
//
// Parameters:
//   - path: the context directory path
//
// Returns:
//   - error: "context directory <path> is a symlink"
func DirSymlink(path string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrValidateContextDirSymlink), path,
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
		desc.Text(text.DescKeyErrValidateContextFileSymlink), file,
	)
}

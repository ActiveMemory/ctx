//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgRc "github.com/ActiveMemory/ctx/internal/config/rc"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// ErrDirNotDeclared is the sentinel returned by rc.ContextDir when
// CTX_DIR is unset or empty. Callers that can legitimately proceed
// without a declared context directory (init, activate, deactivate,
// bootstrap) check with errors.Is; everyone else should propagate
// the error or call rc.RequireContextDir for a user-facing message
// (see NotDeclared below).
//
// The message lives in config/rc (not resolved through desc.Text)
// because sentinel values are initialized at package load time,
// before the embedded YAML lookup is populated. Callers that print
// this to users should wrap it via NotDeclared; the sentinel itself
// is for errors.Is comparisons, not for display.
var ErrDirNotDeclared = errors.New(cfgRc.ErrMsgDirNotDeclared)

// ErrRelativeNotAllowed is the sentinel returned when CTX_DIR is
// declared as a relative path. Absolute-only is a hardline: a
// relative CTX_DIR would resolve differently in every cwd, exactly
// the silent cwd-dependency this resolver is meant to eliminate.
//
// Wrap via [RelativeNotAllowed] for user-facing messages so the
// offending value is shown.
var ErrRelativeNotAllowed = errors.New(cfgRc.ErrMsgRelativeNotAllowed)

// ErrNonCanonicalBasename is the sentinel returned when CTX_DIR's
// basename is not the canonical [cfgDir.Context]. It catches the
// common footgun `export CTX_DIR=$(pwd)` (project root instead of
// the `.context` subdirectory) on first use rather than letting init
// deposit canonical files into the project root.
//
// Wrap via [NonCanonicalBasename] for user-facing messages.
var ErrNonCanonicalBasename = errors.New(cfgRc.ErrMsgNonCanonicalBasename)

// ErrContextDirNotFound is the sentinel returned by
// rc.RequireContextDir when CTX_DIR is shape-valid but the directory
// does not exist on disk. Distinct from [ErrDirNotDeclared], which
// fires before any filesystem check.
//
// Construct via [Missing]; the legacy [NotFoundError] type also
// carries this sentinel through its [NotFoundError.Is] method, so
// callers using either pattern can compare with [errors.Is].
var ErrContextDirNotFound = errors.New(cfgRc.ErrMsgContextDirNotFound)

// ErrContextDirNotADirectory is the sentinel returned when CTX_DIR
// points at an existing path that is not a directory (typically a
// regular file). Symlinks pointing at directories pass.
var ErrContextDirNotADirectory = errors.New(cfgRc.ErrMsgContextDirNotADirectory)

// ErrContextDirStat is the sentinel returned when [os.Stat] on
// CTX_DIR fails for a reason other than not-exist (permission
// denied, I/O error). Wrap via [StatFailed] to attach the
// underlying cause.
var ErrContextDirStat = errors.New(cfgRc.ErrMsgContextDirStat)

// RelativeNotAllowed wraps [ErrRelativeNotAllowed] with the
// offending value so the user sees what they declared.
//
// Parameters:
//   - raw: the rejected CTX_DIR value
//
// Returns:
//   - error: wrapping [ErrRelativeNotAllowed] for [errors.Is] matches
func RelativeNotAllowed(raw string) error {
	return fmt.Errorf(cfgRc.FmtWrapColon,
		ErrRelativeNotAllowed,
		fmt.Sprintf(desc.Text(text.DescKeyErrContextRelativeNotAllowed), raw),
	)
}

// NonCanonicalBasename wraps [ErrNonCanonicalBasename] with the
// offending basename so the user sees how their declaration deviated
// from the canonical `.context`.
//
// Parameters:
//   - base: the rejected basename (e.g., "tmp", "myctx")
//
// Returns:
//   - error: wrapping [ErrNonCanonicalBasename] for [errors.Is] matches
func NonCanonicalBasename(base string) error {
	return fmt.Errorf(cfgRc.FmtWrapColon,
		ErrNonCanonicalBasename,
		fmt.Sprintf(
			desc.Text(text.DescKeyErrContextNonCanonicalBasename),
			cfgDir.Context, base,
		),
	)
}

// Missing wraps [ErrContextDirNotFound] with the missing path so
// the user sees which directory was expected.
//
// Parameters:
//   - path: absolute path that does not exist
//
// Returns:
//   - error: wrapping [ErrContextDirNotFound] for [errors.Is] matches
func Missing(path string) error {
	return fmt.Errorf(cfgRc.FmtWrapBare,
		ErrContextDirNotFound,
		path,
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

// NotFoundError is returned when the context directory does not exist.
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

// Is reports whether target matches the not-found sentinel. Lets
// callers using errors.Is(err, ErrContextDirNotFound) match instances
// of [NotFoundError] without rewriting them.
//
// Parameters:
//   - target: error to compare against
//
// Returns:
//   - bool: true when target is the not-found sentinel
func (e *NotFoundError) Is(target error) bool {
	return target == ErrContextDirNotFound
}

// NotFound returns a NotFoundError for the given directory.
//
// Parameters:
//   - path: path to the missing context directory
//
// Returns:
//   - *NotFoundError: typed error for errors.As matching
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

// NotDeclared returns the standard "no context directory specified"
// error used by rc.RequireContextDir when CTX_DIR has not been
// declared.
//
// The returned message is tailored by how many .context/ candidates
// are visible from the caller's CWD, so users get a next-step hint
// specific to their situation:
//
//   - zero candidates:  suggest `ctx init`.
//   - one candidate:    name it as the likely target and suggest
//     `eval "$(ctx activate)"`.
//   - many candidates:  list all of them and refer the user to
//     `ctx activate` from a more specific cwd.
//
// The scan that produces candidates is read-only (rc.ScanCandidates)
// and never binds anything; resolution itself stays explicit.
//
// Parameters:
//   - candidates: absolute paths of every visible .context/
//     directory, ordered innermost-first. Empty/nil when none.
//
// Returns:
//   - error: a multi-line, actionable message ready to be returned
//     from a Cobra Run function.
func NotDeclared(candidates []string) error {
	switch len(candidates) {
	case 0:
		return errors.New(desc.Text(text.DescKeyErrContextNotDeclaredZero))
	case 1:
		return fmt.Errorf(
			desc.Text(text.DescKeyErrContextNotDeclaredOne),
			candidates[0],
		)
	default:
		var b strings.Builder
		for _, p := range candidates {
			b.WriteString(token.Indent2)
			b.WriteString(p)
			b.WriteString(token.NewlineLF)
		}
		return fmt.Errorf(
			desc.Text(text.DescKeyErrContextNotDeclaredMany),
			strings.TrimRight(b.String(), token.NewlineLF),
		)
	}
}

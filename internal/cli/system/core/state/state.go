//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/validate"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Dir returns the project-scoped runtime state directory
// (`<context dir>/state/`). Ensures the directory exists on each call;
// MkdirAll is a no-op when the directory is already present.
//
// **Always returns an error when the path is empty.** Specifically,
// when CTX_DIR is not declared, Dir returns
// ("", [errCtx.ErrDirNotDeclared]) so callers that gate on
// `dirErr != nil` are uniformly safe. Defensive callers that need
// to special-case the legitimate-absence path can match with
// `errors.Is(dirErr, errCtx.ErrDirNotDeclared)`.
//
// The contract was tightened from the earlier ("", nil) form because
// that form silently invited `filepath.Join("", rel)` traps:
// callers that only checked `dirErr != nil` would join to a
// CWD-relative path and write to the wrong location. Returning an
// explicit error makes the empty-path case unrepresentable in a
// "looks fine" branch.
//
// Returns:
//   - string: Absolute path to the state directory; always non-empty
//     when the error is nil.
//   - error: [errCtx.ErrDirNotDeclared] when CTX_DIR is unset,
//     resolver errors otherwise, mkdir failures otherwise.
func Dir() (string, error) {
	if dirOverride != "" {
		return dirOverride, nil
	}
	ctxDir, err := rc.ContextDir()
	if err != nil {
		// Propagate every resolver error (including
		// ErrDirNotDeclared) so callers can match on it via
		// errors.Is when they need to special-case the absence.
		return "", err
	}
	d := filepath.Join(ctxDir, dir.State)
	if mkdirErr := ctxIo.SafeMkdirAll(d, fs.PermRestrictedDir); mkdirErr != nil {
		return "", mkdirErr
	}
	return d, nil
}

// dirOverride allows tests to redirect Dir() to a temp directory.
var dirOverride string

// SetDirForTest overrides Dir() for testing. Pass an empty string
// to restore the default behavior. Only call from tests.
//
// Parameters:
//   - d: Directory path to use, or empty string to restore default
func SetDirForTest(d string) {
	dirOverride = d
}

// Initialized reports whether the context directory has been properly set up
// via "ctx init". Hooks should no-op when this returns false to avoid
// creating a partial state (e.g., logs/) before initialization.
//
// Returns (false, nil) when the context directory is not declared: there
// is no directory to inspect, which is a legitimate "not initialized"
// answer. Any other resolver failure is propagated so callers can
// distinguish "properly not initialized" from "we could not tell" and
// surface the failure instead of letting hooks silently stop firing.
//
// Returns:
//   - bool: True if the context directory is initialized
//   - error: non-nil on resolver failure (other than not-declared)
func Initialized() (bool, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		if errors.Is(err, errCtx.ErrDirNotDeclared) {
			return false, nil
		}
		return false, err
	}
	return ctxContext.Initialized(ctxDir), nil
}

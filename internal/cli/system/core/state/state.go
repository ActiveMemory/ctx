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
// (`<context dir>/state/`). Creates the directory on demand only
// when the project is initialized; MkdirAll is a no-op when the
// directory is already present.
//
// **Always returns an error when the path is empty.** Three
// legitimate-absence sentinels can fire and should be treated by
// hook callers as silent no-ops, by interactive callers as
// user-facing errors:
//
//   - [errCtx.ErrDirNotDeclared]: CTX_DIR is unset.
//   - [errCtx.ErrNotInitialized]: CTX_DIR is set, but the project
//     lacks the required context files (`ctx init` has not run).
//
// The Initialized gate inside Dir was added to close a leak: a
// caller that reached Dir before checking [Initialized] would
// mkdir `.context/state/` (mode 0750) into a non-ctx project. That
// happens in practice when Cursor imports the ctx Claude plugin's
// hooks and fires them in every workspace it opens; see
// specs/state-dir-no-mkdir-when-uninitialized.md. Gating mkdir on
// Initialized makes the invariant ("no .context/state/ in
// uninitialized projects") structural rather than convention.
//
// The empty-path-on-error contract was tightened from the earlier
// ("", nil) form because that form silently invited
// `filepath.Join("", rel)` traps: callers that only checked
// `dirErr != nil` would join to a CWD-relative path and write to
// the wrong location. Returning an explicit error makes the
// empty-path case unrepresentable in a "looks fine" branch.
//
// Returns:
//   - string: Absolute path to the state directory; always non-empty
//     when the error is nil.
//   - error: [errCtx.ErrDirNotDeclared] when CTX_DIR is unset;
//     [errCtx.ErrNotInitialized] (wrapped via [errCtx.NotInitialized]
//     for the user-facing path) when the project is not initialized;
//     resolver errors otherwise; mkdir failures otherwise.
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
	if !ctxContext.Initialized(ctxDir) {
		// Refuse to mkdir state/ in a project that has not been
		// initialized. Wrap the sentinel so interactive callers
		// surface a path-bearing message; hook callers that gate
		// on dirErr != nil keep working as before.
		return "", errCtx.NotInitialized(ctxDir)
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

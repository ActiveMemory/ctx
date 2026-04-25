//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"errors"
	"os"

	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
)

// RequireContextDir returns the declared context directory after
// validating both its declaration shape (via [ContextDir]) and that
// the path actually exists on disk as a directory.
//
// This is the **operating-command boundary**: every non-exempt
// command calls it at the start of its Run function (or via
// [PersistentPreRunE]). Diagnostic and exempt callers (init,
// activate, bootstrap, hooks like check-anchor-drift) must use
// [ContextDir] directly so they observe declared state without
// erroring on broken state.
//
// Convention: operating callers use this; only diagnostic / exempt
// callers may use raw [ContextDir]. Without that rule, operating
// callers would receive shape-valid but non-existent paths and
// surface confusing downstream errors instead of the friendly
// tailored not-found message.
//
// Rejection conditions:
//
//  1. CTX_DIR truly unset ([errCtx.ErrDirNotDeclared]) is rewrapped
//     as [errCtx.NotDeclared] tailored to how many .context/
//     candidates are visible from CWD. The user said "I haven't
//     told you anything yet"; the message offers a next step.
//  2. CTX_DIR set to a relative or non-canonical-basename value
//     ([errCtx.ErrRelativeNotAllowed] / [errCtx.ErrNonCanonicalBasename])
//     is propagated unchanged. The user told us a specific value;
//     the diagnostic should name what's wrong with that value
//     ("must be absolute, got '...'", "basename must be '.context',
//     got 'tmp'") rather than pretend nothing was declared.
//  3. Path does not exist: [errCtx.ErrContextDirNotFound] (wrapped
//     via [errCtx.Missing]).
//  4. Stat failed for a reason other than not-exist (permission
//     denied, I/O error): [errCtx.ErrContextDirStat] (wrapped via
//     [errCtx.StatFailed]).
//  5. Path exists but is not a directory:
//     [errCtx.ErrContextDirNotADirectory].
//
// Exempt commands (ctx init, ctx activate, ctx deactivate,
// ctx version, ctx help, ctx system bootstrap) must not call this
// helper; they handle the unset case themselves, either by creating
// the directory (init), walking to emit shell integration (activate),
// or reporting resolution state for diagnostics (bootstrap).
//
// Returns:
//   - string: absolute path to the declared context directory.
//   - error: non-nil with a multi-line actionable message when the
//     context directory has not been declared, does not exist, or
//     does not name a directory; the error is already formatted
//     for direct return from a Cobra Run function.
func RequireContextDir() (string, error) {
	path, err := ContextDir()
	if err != nil {
		// Discriminate by error kind: only truly-unset gets the
		// tailored multi-line "no context directory specified"
		// message with candidate hints. Relative-path and
		// non-canonical-basename errors are propagated with their
		// precise "what's wrong with the value you gave us"
		// message; collapsing them into the unset form would tell
		// the user "you didn't declare it" when they did declare
		// it (just to the wrong shape): exactly the silent /
		// confusing diagnostic the spec was meant to eliminate.
		if errors.Is(err, errCtx.ErrDirNotDeclared) {
			cwd, _ := os.Getwd()
			return "", errCtx.NotDeclared(ScanCandidates(cwd))
		}
		return "", err
	}
	info, statErr := os.Stat(path)
	if statErr != nil {
		if errors.Is(statErr, os.ErrNotExist) {
			return "", errCtx.Missing(path)
		}
		return "", errCtx.StatFailed(path, statErr)
	}
	if !info.IsDir() {
		return "", errCtx.NotADir(path)
	}
	return path, nil
}

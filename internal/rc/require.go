//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

// RequireContextDir returns the project's context directory.
//
// Under the cwd-anchored resolution model
// (spec: specs/cwd-anchored-context.md), this is a thin wrapper
// around [ContextDir]: both perform a single [os.Stat] of
// `$PWD/.context/` and either return its absolute path or one of
// the typed `errCtx` errors. The wrapper exists so operating
// commands can keep calling the longer name and so the call site
// reads as a contract gate rather than a getter.
//
// Convention: operating commands call this from PersistentPreRunE
// (or their Run function); exempt commands (init, bootstrap
// diagnostics, version, help) may call [ContextDir] directly and
// decide how to handle [errCtx.ErrNoCtxHere].
//
// Returns:
//   - string: absolute path to `$PWD/.context` when present.
//   - error: see [ContextDir] for the full failure taxonomy.
func RequireContextDir() (string, error) {
	return ContextDir()
}

//	/    ctx:                         https://ctx.ist
//
// ,'`./    do you remember?
//
//	`.,'\
//	  \    Copyright 2026-present Context contributors.
//	                SPDX-License-Identifier: Apache-2.0

package env

// Environment variable names.
const (
	// Home is the environment variable for the user's home directory.
	Home = "HOME"
	// Shell is the environment variable that names the user's login
	// shell (e.g. "/bin/bash"). Read by `ctx activate` /
	// `ctx deactivate` to auto-detect the emitter dialect.
	Shell = "SHELL"
	// CtxDir is the environment variable that declares the context
	// directory. Single-source-anchor model:
	// specs/single-source-context-anchor.md.
	CtxDir = "CTX_DIR"
	// CtxDirInherited is the diagnostic-only sibling of CtxDir set by
	// the check-anchor-drift hook line so the hook can compare the
	// parent shell's pre-injection CTX_DIR against the
	// CLAUDE_PROJECT_DIR-anchored CTX_DIR. Not read by the resolver
	// or any operating command; consumed only by
	// `ctx system check-anchor-drift`.
	CtxDirInherited = "CTX_DIR_INHERITED"
	// CtxTokenBudget is the environment variable for overriding
	// the token budget.
	//nolint:gosec // G101: env var name, not a credential
	CtxTokenBudget = "CTX_TOKEN_BUDGET"
	// SessionID is the environment variable for the active AI session ID.
	// Used by ctx trace for context linking.
	SessionID = "CTX_SESSION_ID"
	// SkipPathCheck is the environment variable that skips the PATH
	// validation during init. Set to True in tests.
	SkipPathCheck = "CTX_SKIP_PATH_CHECK"
)

// Environment toggle values.
const (
	// True is the canonical truthy value for environment variable toggles.
	True = "1"
)

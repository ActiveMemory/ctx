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
	// User is the environment variable for the current username.
	User = "USER"
	// CtxBudget is the environment variable for overriding
	// the token budget.
	CtxBudget = "CTX_TOKEN_BUDGET"
	// SessionID is the environment variable for the active AI session ID.
	// Used by ctx trace for context linking.
	SessionID = "CTX_SESSION_ID"
	// SkipPathCheck is the environment variable that skips the PATH
	// validation during init. Set to True in tests.
	SkipPathCheck = "CTX_SKIP_PATH_CHECK"
	// HubAdmin is the environment variable holding the hub
	// admin token, used as a fallback when --token is not passed
	// to admin-gated commands like `ctx hub revoke`.
	HubAdmin = "CTX_HUB_ADMIN_TOKEN"
)

// Environment toggle values.
const (
	// True is the canonical truthy value for environment variable toggles.
	True = "1"
)

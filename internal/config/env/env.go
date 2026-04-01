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
	// CtxDir is the environment variable for overriding the context directory.
	CtxDir = "CTX_DIR"
	// CtxTokenBudget is the environment variable for overriding
	// the token budget.
	//nolint:gosec // G101: env var name, not a credential
	CtxTokenBudget = "CTX_TOKEN_BUDGET"
	// BackupSMBURL is the environment variable for the SMB share URL.
	BackupSMBURL = "CTX_BACKUP_SMB_URL"
	// BackupSMBSubdir is the environment variable for the SMB share subdirectory.
	BackupSMBSubdir = "CTX_BACKUP_SMB_SUBDIR"
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

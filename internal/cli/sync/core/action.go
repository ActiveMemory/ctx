//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import "github.com/ActiveMemory/ctx/internal/entity"

// DetectSyncActions scans the codebase and returns suggested sync actions.
//
// Runs multiple checks to identify discrepancies between the codebase and
// context documentation:
//   - New directories not documented in ARCHITECTURE.md
//   - Package manager files without dependency documentation
//   - Config files not mentioned in CONVENTIONS.md
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []Action: List of suggested actions to reconcile context with codebase
func DetectSyncActions(ctx *entity.Context) []Action {
	var actions []Action

	// Check for new top-level directories not mentioned in ARCHITECTURE.md
	actions = append(actions, CheckNewDirectories(ctx)...)

	// Check for package manager files
	actions = append(actions, CheckPackageFiles(ctx)...)

	// Check for common config files that might need documenting
	actions = append(actions, CheckConfigFiles(ctx)...)

	return actions
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package action

import (
	"github.com/ActiveMemory/ctx/internal/cli/sync/core/validate"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// Detect scans the codebase and returns suggested sync actions.
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
//   - error: non-nil when a check cannot confirm its answer (e.g. the
//     project root directory cannot be read); callers surface this
//     rather than printing a confident empty suggestion list
func Detect(ctx *entity.Context) ([]validate.Action, error) {
	var actions []validate.Action

	// Check for new top-level directories not mentioned in ARCHITECTURE.md
	newDirs, newDirsErr := validate.CheckNewDirectories(ctx)
	if newDirsErr != nil {
		return nil, newDirsErr
	}
	actions = append(actions, newDirs...)

	// Check for package manager files
	actions = append(actions, validate.CheckPackageFiles(ctx)...)

	// Check for common config files that might need documenting
	actions = append(actions, validate.CheckConfigFiles(ctx)...)

	return actions, nil
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import (
	"github.com/spf13/cobra"

	coreArchive "github.com/ActiveMemory/ctx/internal/cli/task/core/archive"
	"github.com/ActiveMemory/ctx/internal/config/token"
	writeArchive "github.com/ActiveMemory/ctx/internal/write/archive"
)

// Run executes the archive subcommand logic.
//
// Moves completed tasks (marked with [x]) from TASKS.md to a timestamped
// archive file, including all nested content (subtasks, metadata). Tasks
// with incomplete children are skipped to avoid orphaning pending work.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dryRun: If true, preview changes without modifying files
//
// Returns:
//   - error: Non-nil if TASKS.md doesn't exist or file operations fail
func Run(cmd *cobra.Command, dryRun bool) error {
	r, planErr := coreArchive.Plan()
	if planErr != nil {
		return planErr
	}

	for _, name := range r.SkippedNames {
		writeArchive.Skipping(cmd, name)
	}

	if len(r.Archivable) == 0 {
		if len(r.SkippedNames) > 0 {
			writeArchive.SkipIncomplete(cmd, len(r.SkippedNames))
		} else {
			writeArchive.NoCompleted(cmd)
		}
		return nil
	}

	if dryRun {
		writeArchive.DryRun(cmd, len(r.Archivable), r.PendingCount,
			r.Content, token.Separator)
		return nil
	}

	archivePath, execErr := coreArchive.Execute(r)
	if execErr != nil {
		return execErr
	}

	writeArchive.Success(cmd, len(r.Archivable), archivePath, r.PendingCount)
	return nil
}

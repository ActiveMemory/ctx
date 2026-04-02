//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/compact/core"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/context/load"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/tidy"
	writeCompact "github.com/ActiveMemory/ctx/internal/write/compact"
)

// Run executes the compact command logic.
//
// Loads context, processes TASKS.md for completed tasks, and removes
// empty sections from all context files.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - archive: If true, archive old completed tasks to .context/archive/
//
// Returns:
//   - error: Non-nil if context loading fails or .context/ is not found
func Run(cmd *cobra.Command, archive bool) error {
	ctx, err := load.Do("")
	if err != nil {
		if _, ok := errors.AsType[*errCtx.NotFoundError](err); ok {
			return errInit.ContextNotInitialized()
		}
		return err
	}

	// Enable archiving if configured in .ctxrc
	if rc.AutoArchive() {
		archive = true
	}

	writeCompact.ReportHeading(cmd)

	changes := 0

	// Process TASKS.md
	tasksChanges, compactErr := core.CompactTasks(cmd, ctx, archive)
	if compactErr != nil {
		writeCompact.TaskError(cmd, compactErr)
	} else {
		changes += tasksChanges
	}

	// Reload context to pick up TASKS.md changes, then clean sections.
	ctx, err = load.Do("")
	if err == nil {
		result := tidy.CompactContext(ctx)
		for i, sc := range result.SectionsCleaned {
			if writeErr := ctxIo.SafeWriteFile(
				result.SectionFileUpdates[i].Path,
				result.SectionFileUpdates[i].Content,
				fs.PermFile,
			); writeErr == nil {
				writeCompact.SectionsRemoved(cmd, sc.Removed, sc.FileName)
				changes += sc.Removed
			}
		}
	}

	if changes == 0 {
		writeCompact.ReportClean(cmd)
	} else {
		writeCompact.ReportSummary(cmd, changes)
	}

	return nil
}

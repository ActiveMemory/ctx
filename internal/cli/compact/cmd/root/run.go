//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"
	"fmt"
	"os"

	"github.com/ActiveMemory/ctx/internal/context/load"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/compact/core"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errctx "github.com/ActiveMemory/ctx/internal/err/context"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/tidy"
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
		var notFoundError *errctx.NotFoundError
		if errors.As(err, &notFoundError) {
			return ctxerr.ContextNotInitialized()
		}
		return err
	}

	// Enable archiving if configured in .ctxrc
	if rc.AutoArchive() {
		archive = true
	}

	cmd.Println("Compact Analysis")
	cmd.Println("================")
	cmd.Println()

	changes := 0

	// Process TASKS.md
	tasksChanges, compactErr := core.CompactTasks(cmd, ctx, archive)
	if compactErr != nil {
		cmd.Println(fmt.Sprintf("⚠ Error processing TASKS.md: %v", compactErr))
	} else {
		changes += tasksChanges
	}

	// Reload context to pick up TASKS.md changes, then clean sections.
	ctx, err = load.Do("")
	if err == nil {
		result := tidy.CompactContext(ctx)
		for i, sc := range result.SectionsCleaned {
			if writeErr := os.WriteFile(
				result.SectionFileUpdates[i].Path,
				result.SectionFileUpdates[i].Content,
				fs.PermFile,
			); writeErr == nil {
				cmd.Println(
					fmt.Sprintf("✓ Removed %d empty sections from %s",
						sc.Removed, sc.FileName),
				)
				changes += sc.Removed
			}
		}
	}

	if changes == 0 {
		cmd.Println("✓ Nothing to compact — context is already clean")
	} else {
		cmd.Println()
		cmd.Println(fmt.Sprintf("✓ Compacted %d items", changes))
	}

	return nil
}

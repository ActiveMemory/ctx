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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/compact/core"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/rc"
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
	ctx, err := context.Load("")
	if err != nil {
		var notFoundError *context.NotFoundError
		if errors.As(err, &notFoundError) {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
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
	tasksChanges, err := core.CompactTasks(cmd, ctx, archive)
	if err != nil {
		cmd.Println(fmt.Sprintf("⚠ Error processing TASKS.md: %v", err))
	} else {
		changes += tasksChanges
	}

	// Process other files for empty sections
	for _, f := range ctx.Files {
		if f.Name == ctxCfg.Task {
			continue
		}
		cleaned, count := core.RemoveEmptySections(string(f.Content))
		if count > 0 {
			if err := os.WriteFile(f.Path, []byte(cleaned), fs.PermFile); err == nil {
				cmd.Println(
					fmt.Sprintf("✓ Removed %d empty sections from %s", count, f.Name),
				)
				changes += count
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

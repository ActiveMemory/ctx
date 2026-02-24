//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// runCompact executes the compact command logic.
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
func runCompact(cmd *cobra.Command, archive bool) error {
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

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	cmd.Println(cyan("Compact Analysis"))
	cmd.Println(cyan("================"))
	cmd.Println()

	changes := 0

	// Process TASKS.md
	tasksChanges, err := compactTasks(cmd, ctx, archive)
	if err != nil {
		cmd.Println(fmt.Sprintf("%s Error processing TASKS.md: %v", yellow("⚠"), err))
	} else {
		changes += tasksChanges
	}

	// Process other files for empty sections
	for _, f := range ctx.Files {
		if f.Name == config.FileTask {
			continue
		}
		cleaned, count := removeEmptySections(string(f.Content))
		if count > 0 {
			if err := os.WriteFile(f.Path, []byte(cleaned), config.PermFile); err == nil {
				cmd.Println(
					fmt.Sprintf("%s Removed %d empty sections from %s", green("✓"), count, f.Name),
				)
				changes += count
			}
		}
	}

	if changes == 0 {
		cmd.Println(fmt.Sprintf("%s Nothing to compact — context is already clean", green("✓")))
	} else {
		cmd.Println()
		cmd.Println(fmt.Sprintf("%s Compacted %d items", green("✓"), changes))
	}

	return nil
}

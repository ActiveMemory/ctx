//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prune

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// PruneDryRunLine prints a single dry-run prune candidate.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: file name being considered for pruning.
//   - age: human-readable age string.
func PruneDryRunLine(cmd *cobra.Command, name, age string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyPruneDryRunLine), name, age))
}

// PruneErrorLine prints an error encountered while removing a file.
//
// Parameters:
//   - cmd: Cobra command for error output. Nil is a no-op.
//   - name: file name that failed to remove.
//   - err: the removal error.
func PruneErrorLine(cmd *cobra.Command, name string, err error) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyPruneErrorLine), name, err))
}

// PruneSummary prints the prune results summary.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - dryRun: whether this was a dry-run invocation.
//   - pruned: number of files pruned (or would be pruned).
//   - skipped: number of files skipped (too recent).
//   - preserved: number of global files preserved.
func PruneSummary(cmd *cobra.Command, dryRun bool, pruned, skipped, preserved int) {
	if cmd == nil {
		return
	}
	if dryRun {
		cmd.Println()
		cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyPruneDryRunSummary),
			pruned, skipped, preserved))
	} else {
		cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyPruneSummary),
			pruned, skipped, preserved))
	}
}

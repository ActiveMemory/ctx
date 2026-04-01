//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run regenerates the index for both DECISIONS.md and LEARNINGS.md.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - args: Command arguments (unused)
//
// Returns:
//   - error: Non-nil if either file read/write fails
func Run(cmd *cobra.Command, _ []string) error {
	w := cmd.OutOrStdout()
	ctxDir := rc.ContextDir()

	decisionsPath := filepath.Join(ctxDir, ctx.Decision)
	decisionsErr := index.Reindex(
		w,
		decisionsPath,
		ctx.Decision,
		index.UpdateDecisions,
		entry.Decisions,
	)
	if decisionsErr != nil {
		return decisionsErr
	}

	learningsPath := filepath.Join(ctxDir, ctx.Learning)
	return index.Reindex(
		w,
		learningsPath,
		ctx.Learning,
		index.UpdateLearnings,
		entry.Learnings,
	)
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// runReindex regenerates the index for both DECISIONS.md and LEARNINGS.md.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - args: Command arguments (unused)
//
// Returns:
//   - error: Non-nil if either file read/write fails
func runReindex(cmd *cobra.Command, _ []string) error {
	w := cmd.OutOrStdout()
	ctxDir := rc.ContextDir()

	decisionsPath := filepath.Join(ctxDir, config.FileDecision)
	decisionsErr := index.ReindexFile(
		w,
		decisionsPath,
		config.FileDecision,
		index.UpdateDecisions,
		config.EntryPlural[config.EntryDecision],
	)
	if decisionsErr != nil {
		return decisionsErr
	}

	learningsPath := filepath.Join(ctxDir, config.FileLearning)
	return index.ReindexFile(
		w,
		learningsPath,
		config.FileLearning,
		index.UpdateLearnings,
		config.EntryPlural[config.EntryLearning],
	)
}

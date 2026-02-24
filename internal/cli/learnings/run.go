//	/    ctx:                         https://ctx.ist
//
// ,'`./    do you remember?
//
//	`.,'\
//	  \    Copyright 2026-present Context contributors.
//	                SPDX-License-Identifier: Apache-2.0

package learnings

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/spf13/cobra"
)

// runReindex regenerates the LEARNINGS.md index.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - args: Command arguments (unused)
//
// Returns:
//   - error: Non-nil if file read/write fails
func runReindex(cmd *cobra.Command, _ []string) error {
	filePath := filepath.Join(rc.ContextDir(), config.FileLearning)
	return index.ReindexFile(
		cmd.OutOrStdout(),
		filePath,
		config.FileLearning,
		index.UpdateLearnings,
		config.EntryPlural[config.EntryLearning],
	)
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	cfgEntry "github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run regenerates the LEARNINGS.md index.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - args: Command arguments (unused)
//
// Returns:
//   - error: Non-nil if the file read/write fails
func Run(cmd *cobra.Command, _ []string) error {
	ctxDir, err := rc.RequireContextDir()
	if err != nil {
		cmd.SilenceUsage = true
		return err
	}
	filePath := filepath.Join(ctxDir, ctx.Learning)
	return index.Reindex(
		cmd.OutOrStdout(),
		filePath,
		ctx.Learning,
		index.UpdateLearnings,
		cfgEntry.Learning,
	)
}

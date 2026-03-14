//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run regenerates the DECISIONS.md index.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - args: Command arguments (unused)
//
// Returns:
//   - error: Non-nil if the file read/write fails
func Run(cmd *cobra.Command, _ []string) error {
	filePath := filepath.Join(rc.ContextDir(), ctx.Decision)
	return index.ReindexFile(
		cmd.OutOrStdout(),
		filePath,
		ctx.Decision,
		index.UpdateDecisions,
		"decisions",
	)
}

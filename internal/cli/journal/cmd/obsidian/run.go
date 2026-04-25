//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package obsidian

import (
	"path/filepath"

	"github.com/spf13/cobra"

	coreObsidian "github.com/ActiveMemory/ctx/internal/cli/journal/core/obsidian"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/obsidian"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run generates an Obsidian vault from journal entries.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - output: Output directory for the vault; when empty, defaults
//     to <context dir>/<obsidian.DirName>
//
// Returns:
//   - error: Non-nil if generation fails
func Run(cmd *cobra.Command, output string) error {
	ctxDir, err := rc.RequireContextDir()
	if err != nil {
		cmd.SilenceUsage = true
		return err
	}
	if output == "" {
		output = filepath.Join(ctxDir, obsidian.DirName)
	}
	return coreObsidian.BuildVault(
		cmd, filepath.Join(ctxDir, dir.Journal), output,
	)
}

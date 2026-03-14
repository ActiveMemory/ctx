//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx compact" command for cleaning up context files.
//
// The command moves completed tasks to a "Completed (Recent)" section,
// optionally archives old content, and removes empty sections from all
// context files.
//
// Flags:
//   - --archive: Create .context/archive/ for old completed tasks
//
// Returns:
//   - *cobra.Command: Configured compact command with flags registered
func Cmd() *cobra.Command {
	var archive bool

	short, long := assets.CommandDesc(assets.CmdDescKeyCompact)

	cmd := &cobra.Command{
		Use:   "compact",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, archive)
		},
	}

	cmd.Flags().BoolVar(
		&archive,
		"archive",
		false,
		assets.FlagDesc(assets.FlagDescKeyCompactArchive),
	)

	return cmd
}

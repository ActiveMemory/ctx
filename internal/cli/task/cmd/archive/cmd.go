//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the tasks archive subcommand.
//
// The archive command moves completed tasks (marked with [x]) from TASKS.md
// to a timestamped archive file in .context/archive/. Pending tasks ([ ])
// remain in TASKS.md.
//
// Flags:
//   - --dry-run: Preview changes without modifying files
//
// Returns:
//   - *cobra.Command: Configured archive subcommand
func Cmd() *cobra.Command {
	var dryRun bool

	short, long := assets.CommandDesc(assets.CmdDescKeyTaskArchive)

	cmd := &cobra.Command{
		Use:   "archive",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runArchive(cmd, dryRun)
		},
	}

	cmd.Flags().BoolVar(
		&dryRun,
		"dry-run",
		false,
		assets.FlagDesc(assets.FlagDescKeyTaskArchiveDryRun),
	)

	return cmd
}

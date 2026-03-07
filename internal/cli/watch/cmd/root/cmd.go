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

// Cmd returns the watch command.
//
// Flags:
//   - --log: Log file to watch (default: stdin)
//   - --dry-run: Show updates without applying
//
// Returns:
//   - *cobra.Command: Configured watch command with flags registered
func Cmd() *cobra.Command {
	var (
		logPath string
		dryRun  bool
	)

	short, long := assets.CommandDesc(assets.CmdDescKeyWatch)

	cmd := &cobra.Command{
		Use:   "watch",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, logPath, dryRun)
		},
	}

	cmd.Flags().StringVar(
		&logPath, "log", "", assets.FlagDesc(assets.FlagDescKeyWatchLog),
	)
	cmd.Flags().BoolVar(
		&dryRun, "dry-run", false, assets.FlagDesc(assets.FlagDescKeyWatchDryRun),
	)

	return cmd
}

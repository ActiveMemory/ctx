//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
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

	short, long := desc.CommandDesc(cmd.DescKeyWatch)

	cmd := &cobra.Command{
		Use:   cmd.DescKeyWatch,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, logPath, dryRun)
		},
	}

	cmd.Flags().StringVar(
		&logPath, "log", "", desc.FlagDesc(flag.DescKeyWatchLog),
	)
	cmd.Flags().BoolVar(
		&dryRun, "dry-run", false, desc.FlagDesc(flag.DescKeyWatchDryRun),
	)

	return cmd
}

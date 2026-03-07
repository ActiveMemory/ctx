//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	watchroot "github.com/ActiveMemory/ctx/internal/cli/watch/cmd/root"
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

	short, long := assets.CommandDesc("watch")

	cmd := &cobra.Command{
		Use:   "watch",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return watchroot.Run(cmd, logPath, dryRun)
		},
	}

	cmd.Flags().StringVar(
		&logPath, "log", "", "Log file to watch (default: stdin)",
	)
	cmd.Flags().BoolVar(
		&dryRun, "dry-run", false, "Show updates without applying",
	)

	return cmd
}

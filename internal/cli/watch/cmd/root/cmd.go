//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
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

	short, long := desc.Command(cmd.DescKeyWatch)

	c := &cobra.Command{
		Use:     cmd.UseWatch,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyWatch),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, logPath, dryRun)
		},
	}

	flagbind.StringFlag(c, &logPath,
		cFlag.Log, flag.DescKeyWatchLog,
	)
	flagbind.BoolFlag(c, &dryRun,
		cFlag.DryRun, flag.DescKeyWatchDryRun,
	)

	return c
}

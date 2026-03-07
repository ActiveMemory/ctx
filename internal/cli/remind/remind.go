//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package remind

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/remind/cmd/add"
	"github.com/ActiveMemory/ctx/internal/cli/remind/cmd/dismiss"
	"github.com/ActiveMemory/ctx/internal/cli/remind/cmd/list"
)

// Cmd returns the remind command with subcommands.
//
// When invoked with arguments and no subcommand, it adds a reminder.
// When invoked with no arguments, it lists all reminders.
//
// Returns:
//   - *cobra.Command: Configured remind command with subcommands
func Cmd() *cobra.Command {
	var afterFlag string

	short, long := assets.CommandDesc("remind")

	cmd := &cobra.Command{
		Use:   "remind [TEXT]",
		Short: short,
		Long:  long,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return add.Run(cmd, args[0], afterFlag)
			}
			return list.Run(cmd)
		},
	}

	cmd.Flags().StringVarP(&afterFlag, "after", "a", "", assets.FlagDesc("remind.after"))

	cmd.AddCommand(add.Cmd())
	cmd.AddCommand(list.Cmd())
	cmd.AddCommand(dismiss.Cmd())

	return cmd
}

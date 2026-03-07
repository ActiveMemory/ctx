//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package remind

import (
	"github.com/spf13/cobra"

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

	cmd := &cobra.Command{
		Use:   "remind [TEXT]",
		Short: "Session-scoped reminders",
		Long: `Manage session-scoped reminders stored in .context/reminders.json.

Reminders surface verbatim at session start and repeat every session until
dismissed. Use --after to gate a reminder until a specific date.

When invoked with a text argument, adds a reminder (equivalent to "remind add").
When invoked with no arguments, lists all reminders.

Subcommands:
  add      Add a reminder (default action)
  list     Show all pending reminders
  dismiss  Dismiss one or all reminders`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return add.RunAdd(cmd, args[0], afterFlag)
			}
			return list.RunList(cmd)
		},
	}

	cmd.Flags().StringVarP(&afterFlag, "after", "a", "", "Don't surface until this date (YYYY-MM-DD)")

	cmd.AddCommand(add.Cmd())
	cmd.AddCommand(list.Cmd())
	cmd.AddCommand(dismiss.Cmd())

	return cmd
}

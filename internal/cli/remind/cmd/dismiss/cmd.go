//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	errReminder "github.com/ActiveMemory/ctx/internal/err/reminder"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the remind dismiss subcommand.
//
// Returns:
//   - *cobra.Command: Configured dismiss subcommand
func Cmd() *cobra.Command {
	var allFlag bool

	short, _ := desc.Command(cmd.DescKeyRemindDismiss)

	c := &cobra.Command{
		Use:     cmd.UseRemindDismiss,
		Aliases: []string{cmd.UseRemindDismissAlias},
		Short:   short,
		Example: desc.Example(cmd.DescKeyRemindDismiss),
		RunE: func(cmd *cobra.Command, args []string) error {
			if allFlag {
				return Run(cmd, "", allFlag)
			}
			if len(args) == 0 {
				return errReminder.IDRequired()
			}
			return Run(cmd, args[0], allFlag)
		},
	}

	flagbind.BoolFlag(c, &allFlag,
		cFlag.All, flag.DescKeyRemindDismissAll,
	)

	return c
}

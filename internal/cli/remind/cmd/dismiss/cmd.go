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
	ctxErr "github.com/ActiveMemory/ctx/internal/err/reminder"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if allFlag {
				return RunDismissAll(cmd)
			}
			if len(args) == 0 {
				return ctxErr.IDRequired()
			}
			return RunDismiss(cmd, args[0])
		},
	}

	c.Flags().BoolVar(&allFlag, cFlag.All, false,
		desc.Flag(flag.DescKeyRemindDismissAll),
	)

	return c
}

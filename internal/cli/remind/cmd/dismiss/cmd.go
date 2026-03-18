//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

import (
	"github.com/ActiveMemory/ctx/internal/config/embed"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/reminder"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the remind dismiss subcommand.
//
// Returns:
//   - *cobra.Command: Configured dismiss subcommand
func Cmd() *cobra.Command {
	var allFlag bool

	short, _ := assets.CommandDesc(embed.CmdDescKeyRemindDismiss)

	cmd := &cobra.Command{
		Use:     "dismiss [ID]",
		Aliases: []string{"rm"},
		Short:   short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if allFlag {
				return RunDismissAll(cmd)
			}
			if len(args) == 0 {
				return ctxerr.IDRequired()
			}
			return RunDismiss(cmd, args[0])
		},
	}

	cmd.Flags().BoolVar(&allFlag, "all", false,
		assets.FlagDesc(embed.FlagDescKeyRemindDismissAll),
	)

	return cmd
}

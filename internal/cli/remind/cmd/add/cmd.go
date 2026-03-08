//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the remind add subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func Cmd() *cobra.Command {
	var afterFlag string

	short, _ := assets.CommandDesc(assets.CmdDescKeyRemindAdd)

	cmd := &cobra.Command{
		Use:   "add TEXT",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], afterFlag)
		},
	}

	cmd.Flags().StringVarP(&afterFlag, "after", "a", "",
		assets.FlagDesc(assets.FlagDescKeyRemindAddAfter),
	)

	return cmd
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the remind add subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func Cmd() *cobra.Command {
	var afterFlag string

	short, _ := desc.Command(cmd.DescKeyRemindAdd)

	cmd := &cobra.Command{
		Use:   "add TEXT",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], afterFlag)
		},
	}

	cmd.Flags().StringVarP(&afterFlag, "after", "a", "",
		desc.Flag(flag.DescKeyRemindAddAfter),
	)

	return cmd
}

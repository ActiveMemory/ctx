//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the remind add subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func Cmd() *cobra.Command {
	var afterFlag string

	short, _ := desc.Command(cmd.DescKeyRemindAdd)

	c := &cobra.Command{
		Use:   cmd.UseRemindAdd,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], afterFlag)
		},
	}

	c.Flags().StringVarP(&afterFlag, cflag.After, cflag.ShortAfter, "",
		desc.Flag(flag.DescKeyRemindAddAfter),
	)

	return c
}

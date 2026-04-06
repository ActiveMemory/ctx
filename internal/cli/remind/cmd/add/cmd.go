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
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the remind add subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func Cmd() *cobra.Command {
	var afterFlag string

	short, _ := desc.Command(cmd.DescKeyRemindAdd)

	c := &cobra.Command{
		Use:     cmd.UseRemindAdd,
		Short:   short,
		Example: desc.Example(cmd.DescKeyRemindAdd),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], afterFlag)
		},
	}

	flagbind.StringFlagP(c, &afterFlag,
		cFlag.After, cFlag.ShortAfter,
		flag.DescKeyRemindAddAfter,
	)

	return c
}

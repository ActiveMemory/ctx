//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unlock

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx recall unlock" subcommand.
//
// Removes lock protection from journal entries, allowing export
// --regenerate to overwrite them again.
//
// Returns:
//   - *cobra.Command: Command for unlocking journal entries
func Cmd() *cobra.Command {
	var all bool

	short, long := desc.CommandDesc(cmd.DescKeyRecallUnlock)

	cmd := &cobra.Command{
		Use:   "unlock <pattern>",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, all)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false,
		desc.FlagDesc(flag.DescKeyRecallUnlockAll),
	)

	return cmd
}

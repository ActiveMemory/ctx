//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unlock

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
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

	short, long := assets.CommandDesc(assets.CmdDescKeyRecallUnlock)

	cmd := &cobra.Command{
		Use:   "unlock <pattern>",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, all)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false,
		assets.FlagDesc(assets.FlagDescKeyRecallUnlockAll),
	)

	return cmd
}

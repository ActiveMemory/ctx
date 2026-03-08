//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rm

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the prompt rm subcommand.
//
// Returns:
//   - *cobra.Command: Configured rm subcommand
func Cmd() *cobra.Command {
	short, _ := assets.CommandDesc(assets.CmdDescKeyPromptRm)

	return &cobra.Command{
		Use:   "rm NAME",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0])
		},
	}
}

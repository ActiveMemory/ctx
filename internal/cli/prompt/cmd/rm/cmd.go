//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rm

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the prompt rm subcommand.
//
// Returns:
//   - *cobra.Command: Configured rm subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeyPromptRm)

	return &cobra.Command{
		Use:   cmd.UsePromptRm,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0])
		},
	}
}

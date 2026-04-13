//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reset

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx message reset" subcommand.
//
// Returns:
//   - *cobra.Command: Configured reset subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeyMessageReset)

	return &cobra.Command{
		Use:     cmd.UseMessageReset,
		Short:   short,
		Example: desc.Example(cmd.DescKeyMessageReset),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], args[1])
		},
	}
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system message show" subcommand.
//
// Returns:
//   - *cobra.Command: Configured show subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeySystemMessageShow)

	return &cobra.Command{
		Use:     cmd.UseSystemMessageShow,
		Short:   short,
		Example: desc.Example(cmd.DescKeySystemMessageShow),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], args[1])
		},
	}
}

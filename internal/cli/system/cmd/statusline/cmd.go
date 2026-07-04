//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package statusline

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system statusline" subcommand.
//
// Returns:
//   - *cobra.Command: Configured statusline subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemStatusline)

	return &cobra.Command{
		Use:     cmd.UseSystemStatusline,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySystemStatusline),
		Hidden:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}

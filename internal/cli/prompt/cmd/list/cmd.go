//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"
)

// Cmd returns the prompt list subcommand.
//
// Returns:
//   - *cobra.Command: Configured list subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeyPromptList)

	return &cobra.Command{
		Use:   cmd.UsePromptList,
		Short: short,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}

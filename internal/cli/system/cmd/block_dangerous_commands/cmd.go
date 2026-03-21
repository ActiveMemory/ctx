//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package block_dangerous_commands

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system block-dangerous-commands" subcommand.
//
// Returns:
//   - *cobra.Command: Configured block-dangerous-commands subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemBlockDangerousCommands)

	return &cobra.Command{
		Use:    "block-dangerous-commands",
		Short:  short,
		Long:   long,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}

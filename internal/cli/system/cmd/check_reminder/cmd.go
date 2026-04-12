//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_reminder

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system check-reminder" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-reminder subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeySystemCheckReminder)

	return &cobra.Command{
		Use:     cmd.UseSystemCheckReminder,
		Short:   short,
		Example: desc.Example(cmd.DescKeySystemCheckReminder),
		Hidden:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}

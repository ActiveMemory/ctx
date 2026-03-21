//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_reminders

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system check-reminders" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-reminders subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeySystemCheckReminders)

	return &cobra.Command{
		Use:    cmd.UseSystemCheckReminders,
		Short:  short,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}

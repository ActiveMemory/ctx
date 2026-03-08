//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_reminders

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system check-reminders" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-reminders subcommand
func Cmd() *cobra.Command {
	short, _ := assets.CommandDesc(assets.CmdDescKeySystemCheckReminders)

	return &cobra.Command{
		Use:    "check-reminders",
		Short:  short,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}

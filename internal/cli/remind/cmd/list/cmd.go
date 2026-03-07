//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/spf13/cobra"
)

// Cmd returns the remind list subcommand.
//
// Returns:
//   - *cobra.Command: Configured list subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Show all pending reminders",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return RunList(cmd)
		},
	}
}

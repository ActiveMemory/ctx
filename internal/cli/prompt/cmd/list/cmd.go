//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/spf13/cobra"
)

// Cmd returns the prompt list subcommand.
//
// Returns:
//   - *cobra.Command: Configured list subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available prompt templates",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return RunList(cmd)
		},
	}
}

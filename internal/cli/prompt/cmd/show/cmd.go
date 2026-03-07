//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"github.com/spf13/cobra"
)

// Cmd returns the prompt show subcommand.
//
// Returns:
//   - *cobra.Command: Configured show subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show NAME",
		Short: "Print a prompt template to stdout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShow(cmd, args[0])
		},
	}
}

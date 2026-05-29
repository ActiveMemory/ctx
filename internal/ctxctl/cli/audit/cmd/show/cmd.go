//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"github.com/spf13/cobra"
)

// Cmd returns the audit show subcommand.
//
// Parameters:
//   - s: English user-facing text supplied by ctxctl
//
// Returns:
//   - *cobra.Command: Configured show subcommand
func Cmd(s Strings) *cobra.Command {
	return &cobra.Command{
		Use:     s.Use,
		Short:   s.Short,
		Example: s.Example,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0])
		},
	}
}

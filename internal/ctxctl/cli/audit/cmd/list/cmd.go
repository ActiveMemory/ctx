//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/spf13/cobra"
)

// Cmd returns the audit list subcommand.
//
// Parameters:
//   - s: English user-facing text supplied by ctxctl
//
// Returns:
//   - *cobra.Command: Configured list subcommand
func Cmd(s Strings) *cobra.Command {
	return &cobra.Command{
		Use:     s.Use,
		Short:   s.Short,
		Example: s.Example,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, s.None, s.ListItem)
		},
	}
}

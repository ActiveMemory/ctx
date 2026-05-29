//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkaudit

import (
	"os"

	"github.com/spf13/cobra"
)

// Cmd returns the "ctxctl audit-relay" subcommand.
//
// Hidden by default: the hook is invoked from
// .claude/settings.local.json (or analogous integration
// config), not by humans.
//
// Parameters:
//   - s: English user-facing text supplied by ctxctl
//
// Returns:
//   - *cobra.Command: Configured audit-relay subcommand
func Cmd(s Strings) *cobra.Command {
	return &cobra.Command{
		Use:     s.Use,
		Short:   s.Short,
		Example: s.Example,
		Hidden:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin, s)
		},
	}
}

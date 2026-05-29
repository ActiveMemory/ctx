//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

import (
	"github.com/spf13/cobra"

	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	errAudit "github.com/ActiveMemory/ctx/internal/ctxctl/err/audit"
)

// Cmd returns the audit dismiss subcommand.
//
// Parameters:
//   - s: English user-facing text supplied by ctxctl
//
// Returns:
//   - *cobra.Command: Configured dismiss subcommand
func Cmd(s Strings) *cobra.Command {
	var allFlag bool

	c := &cobra.Command{
		Use:     s.Use,
		Short:   s.Short,
		Example: s.Example,
		RunE: func(cmd *cobra.Command, args []string) error {
			if allFlag {
				return Run(cmd, nil, true, s.Dismissed, s.DismissedAll)
			}
			if len(args) == 0 {
				cmd.SilenceUsage = true
				return errAudit.ErrIDRequired
			}
			return Run(cmd, args, false, s.Dismissed, s.DismissedAll)
		},
	}

	c.Flags().BoolVar(&allFlag, cFlag.All, false, s.AllFlag)

	return c
}

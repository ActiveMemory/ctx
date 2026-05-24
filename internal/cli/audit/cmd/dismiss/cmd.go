//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dismiss

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	errAudit "github.com/ActiveMemory/ctx/internal/err/audit"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the audit dismiss subcommand.
//
// Returns:
//   - *cobra.Command: Configured dismiss subcommand
func Cmd() *cobra.Command {
	var allFlag bool

	short, _ := desc.Command(cmd.DescKeyAuditDismiss)
	c := &cobra.Command{
		Use:     cmd.UseAuditDismiss,
		Short:   short,
		Example: desc.Example(cmd.DescKeyAuditDismiss),
		RunE: func(cmd *cobra.Command, args []string) error {
			if allFlag {
				return Run(cmd, nil, true)
			}
			if len(args) == 0 {
				cmd.SilenceUsage = true
				return errAudit.IDRequired()
			}
			return Run(cmd, args, false)
		},
	}

	flagbind.BoolFlag(c, &allFlag,
		cFlag.All, flag.DescKeyAuditDismissAll,
	)

	return c
}

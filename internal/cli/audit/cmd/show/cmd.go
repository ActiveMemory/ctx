//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the audit show subcommand.
//
// Returns:
//   - *cobra.Command: Configured show subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeyAuditShow)
	return &cobra.Command{
		Use:     cmd.UseAuditShow,
		Short:   short,
		Example: desc.Example(cmd.DescKeyAuditShow),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0])
		},
	}
}

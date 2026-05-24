//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/audit/cmd/dismiss"
	"github.com/ActiveMemory/ctx/internal/cli/audit/cmd/list"
	"github.com/ActiveMemory/ctx/internal/cli/audit/cmd/show"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the audit command with subcommands.
//
// When invoked with no subcommand, defaults to `list`.
//
// Returns:
//   - *cobra.Command: Configured audit command with
//     list / show / dismiss subcommands
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyAudit)
	c := &cobra.Command{
		Use:     cmd.UseAudit,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyAudit),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return list.Run(cmd)
		},
	}

	c.AddCommand(list.Cmd())
	c.AddCommand(show.Cmd())
	c.AddCommand(dismiss.Cmd())

	return c
}

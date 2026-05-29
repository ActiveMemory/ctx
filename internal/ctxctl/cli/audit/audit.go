//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/cmd/dismiss"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/cmd/list"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/cmd/show"
)

// Cmd returns the audit command with subcommands.
//
// When invoked with no subcommand, defaults to `list`.
//
// Parameters:
//   - s: English user-facing text supplied by ctxctl
//
// Returns:
//   - *cobra.Command: Configured audit command with
//     list / show / dismiss subcommands
func Cmd(s Strings) *cobra.Command {
	c := &cobra.Command{
		Use:     s.Use,
		Short:   s.Short,
		Long:    s.Long,
		Example: s.Example,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return list.Run(cmd, s.List.None, s.List.ListItem)
		},
	}

	c.AddCommand(list.Cmd(s.List))
	c.AddCommand(show.Cmd(s.Show))
	c.AddCommand(dismiss.Cmd(s.Dismiss))

	return c
}

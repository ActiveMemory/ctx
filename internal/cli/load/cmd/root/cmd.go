//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the "ctx load" command for outputting assembled context.
//
// The command loads context files from .context/ and outputs them in the
// recommended read order, suitable for providing to an AI assistant.
//
// Flags:
//   - --budget: Token budget for assembly (default 8000)
//   - --raw: Output raw file contents without headers or assembly
//
// Returns:
//   - *cobra.Command: Configured load command with flags registered
func Cmd() *cobra.Command {
	var (
		budget int
		raw    bool
	)

	short, long := assets.CommandDesc(assets.CmdDescKeyLoad)
	cmd := &cobra.Command{
		Use:   "load",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use configured budget if flag not explicitly set
			if !cmd.Flags().Changed("budget") {
				budget = rc.TokenBudget()
			}
			return Run(cmd, budget, raw)
		},
	}

	cmd.Flags().IntVar(
		&budget, "budget",
		rc.DefaultTokenBudget,
		assets.FlagDesc(assets.FlagDescKeyLoadBudget),
	)
	cmd.Flags().BoolVar(
		&raw, "raw", false, assets.FlagDesc(assets.FlagDescKeyLoadRaw),
	)

	return cmd
}

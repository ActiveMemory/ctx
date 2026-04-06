//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
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

	short, long := desc.Command(cmd.DescKeyLoad)
	c := &cobra.Command{
		Use:     cmd.UseLoad,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyLoad),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use the configured budget if the flag is not explicitly set
			if !cmd.Flags().Changed(cFlag.Budget) {
				budget = rc.TokenBudget()
			}
			return Run(cmd, budget, raw)
		},
	}

	flagbind.IntFlag(
		c, &budget,
		cFlag.Budget, rc.DefaultTokenBudget,
		flag.DescKeyLoadBudget,
	)
	flagbind.BoolFlag(
		c, &raw,
		cFlag.Raw, flag.DescKeyLoadRaw,
	)

	return c
}

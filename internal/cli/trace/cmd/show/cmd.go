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
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the trace command.
//
// Returns:
//   - *cobra.Command: Configured trace command with flags registered
func Cmd() *cobra.Command {
	var (
		last       int
		jsonOutput bool
	)

	short, long := desc.Command(cmd.DescKeyTrace)

	c := &cobra.Command{
		Use:   cmd.UseTrace,
		Short: short,
		Long:  long,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return Run(cobraCmd, args, last, jsonOutput)
		},
	}

	c.Flags().IntVarP(&last, cFlag.Last, cFlag.ShortLast, 0, "Show context for last N commits")
	c.Flags().BoolVar(&jsonOutput, cFlag.JSON, false, "Output as JSON")

	return c
}

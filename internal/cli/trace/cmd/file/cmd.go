//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package file

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the trace file subcommand.
//
// Returns:
//   - *cobra.Command: Configured trace file command with flags registered
func Cmd() *cobra.Command {
	var last int
	short, long := desc.Command(cmd.DescKeyTraceFile)
	c := &cobra.Command{
		Use:   cmd.UseTraceFile,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return Run(cobraCmd, args[0], last)
		},
	}
	c.Flags().IntVarP(&last, cFlag.Last, cFlag.ShortLast, 20, "Max commits to show")
	return c
}

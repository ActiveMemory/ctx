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
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
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

	flagbind.IntFlagP(c, &last, cFlag.Last, cFlag.ShortLast, 0, flag.DescKeyTraceLast)
	flagbind.BoolFlag(c, &jsonOutput, cFlag.JSON, flag.DescKeyTraceJSON)

	return c
}

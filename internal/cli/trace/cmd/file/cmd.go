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
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the trace file subcommand.
//
// Returns:
//   - *cobra.Command: Configured trace file command with flags registered
func Cmd() *cobra.Command {
	var last int
	short, long := desc.Command(cmd.DescKeyTraceFile)
	c := &cobra.Command{
		Use:     cmd.UseTraceFile,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyTraceFile),
		Args:    cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return Run(cobraCmd, args[0], last)
		},
	}
	flagbind.IntFlagP(
		c, &last,
		cFlag.Last, cFlag.ShortLast,
		cfgTrace.DefaultLastFile, flag.DescKeyTraceFileLast,
	)
	return c
}

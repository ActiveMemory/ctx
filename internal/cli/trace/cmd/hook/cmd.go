//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the trace hook subcommand.
//
// Returns:
//   - *cobra.Command: Configured trace hook command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyTraceHook)
	c := &cobra.Command{
		Use:   cmd.UseTraceHook,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return Run(cobraCmd, args[0])
		},
	}
	return c
}

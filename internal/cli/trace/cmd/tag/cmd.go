//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tag

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the trace tag subcommand.
//
// Returns:
//   - *cobra.Command: Configured trace tag command with flags registered
func Cmd() *cobra.Command {
	var note string
	short, long := desc.Command(cmd.DescKeyTraceTag)
	c := &cobra.Command{
		Use:   cmd.UseTraceTag,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return Run(cobraCmd, args[0], note)
		},
	}
	c.Flags().StringVar(&note, "note", "", "Context note to attach to the commit (required)")
	_ = c.MarkFlagRequired("note")
	return c
}

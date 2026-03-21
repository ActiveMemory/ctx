//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mark_wrapped_up

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system mark-wrapped-up" subcommand.
//
// Returns:
//   - *cobra.Command: Configured mark-wrapped-up subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemMarkWrappedUp)

	return &cobra.Command{
		Use:    "mark-wrapped-up",
		Short:  short,
		Long:   long,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the pad resolve subcommand.
//
// Returns:
//   - *cobra.Command: Configured resolve subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyPadResolve)
	return &cobra.Command{
		Use:     cmd.UsePadResolve,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyPadResolve),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}

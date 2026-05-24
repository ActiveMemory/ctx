//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package undo

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the pad undo subcommand.
//
// Returns:
//   - *cobra.Command: Configured undo subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyPadUndo)
	return &cobra.Command{
		Use:     cmd.UsePadUndo,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyPadUndo),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}

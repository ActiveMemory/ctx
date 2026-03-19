//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package restore

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx permission restore" subcommand.
//
// Returns:
//   - *cobra.Command: Configured restore subcommand
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeyPermissionRestore)

	return &cobra.Command{
		Use:   "restore",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}

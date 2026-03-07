//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package restore

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx permissions restore" subcommand.
//
// Returns:
//   - *cobra.Command: Configured restore subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc("permissions.restore")

	return &cobra.Command{
		Use:   "restore",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}

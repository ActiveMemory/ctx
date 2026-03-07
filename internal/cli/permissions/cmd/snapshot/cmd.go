//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package snapshot

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx permissions snapshot" subcommand.
//
// Returns:
//   - *cobra.Command: Configured snapshot subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeyPermissionsSnapshot)

	return &cobra.Command{
		Use:   "snapshot",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}

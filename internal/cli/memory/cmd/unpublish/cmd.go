//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unpublish

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the memory unpublish subcommand.
//
// Returns:
//   - *cobra.Command: command for removing published context from MEMORY.md.
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeyMemoryUnpublish)
	return &cobra.Command{
		Use:   "unpublish",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}

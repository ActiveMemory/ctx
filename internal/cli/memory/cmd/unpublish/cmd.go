//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unpublish

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the memory unpublish subcommand.
//
// Returns:
//   - *cobra.Command: command for removing published context from MEMORY.md.
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyMemoryUnpublish)
	return &cobra.Command{
		Use:   cmd.UseMemoryUnpublish,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}

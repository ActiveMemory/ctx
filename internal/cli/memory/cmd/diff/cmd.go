//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package diff

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the memory diff subcommand.
//
// Returns:
//   - *cobra.Command: command for showing memory diff.
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyMemoryDiff)
	return &cobra.Command{
		Use:   cmd.UseMemoryDiff,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}

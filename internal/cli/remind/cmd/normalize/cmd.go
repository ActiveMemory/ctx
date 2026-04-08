//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package normalize

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the remind normalize subcommand.
//
// Returns:
//   - *cobra.Command: Configured normalize subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeyRemindNormalize)
	return &cobra.Command{
		Use:   cmd.UseRemindNormalize,
		Short: short,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}

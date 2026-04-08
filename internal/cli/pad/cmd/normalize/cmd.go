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

// Cmd returns the pad normalize subcommand.
//
// Returns:
//   - *cobra.Command: Configured normalize subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeyPadNormalize)
	return &cobra.Command{
		Use:     cmd.UsePadNormalize,
		Short:   short,
		Example: desc.Example(cmd.DescKeyPadNormalize),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}

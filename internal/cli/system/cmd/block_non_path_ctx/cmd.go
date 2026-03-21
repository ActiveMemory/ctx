//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package block_non_path_ctx

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system block-non-path-ctx" subcommand.
//
// Returns:
//   - *cobra.Command: Configured block-non-path-ctx subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemBlockNonPathCtx)

	return &cobra.Command{
		Use:    cmd.UseSystemBlockNonPathCtx,
		Short:  short,
		Long:   long,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}

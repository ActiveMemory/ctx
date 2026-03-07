//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package imp

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the pad import subcommand.
//
// Returns:
//   - *cobra.Command: Configured import subcommand
func Cmd() *cobra.Command {
	var blobs bool

	short, long := assets.CommandDesc(assets.CmdDescKeyPadImp)
	cmd := &cobra.Command{
		Use:   "import FILE",
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if blobs {
				return runImportBlobs(cmd, args[0])
			}
			return runImport(cmd, args[0])
		},
	}

	cmd.Flags().BoolVar(&blobs, "blobs", false,
		assets.FlagDesc(assets.FlagDescKeyPadImpBlobs))

	return cmd
}

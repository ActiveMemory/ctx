//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package imp

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the pad import subcommand.
//
// Returns:
//   - *cobra.Command: Configured import subcommand
func Cmd() *cobra.Command {
	var blobs bool

	short, long := desc.Command(cmd.DescKeyPadImp)
	c := &cobra.Command{
		Use:   cmd.UsePadImport,
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

	c.Flags().BoolVar(&blobs, cflag.Blob, false,
		desc.Flag(flag.DescKeyPadImpBlob))

	return c
}

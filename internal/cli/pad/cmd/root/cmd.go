//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the pad import subcommand.
//
// Returns:
//   - *cobra.Command: Configured import subcommand
func Cmd() *cobra.Command {
	var blobs bool

	short, long := desc.Command(cmd.DescKeyPadImp)
	c := &cobra.Command{
		Use:     cmd.UsePadImport,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyPadImp),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], blobs)
		},
	}

	flagbind.BoolFlag(c, &blobs,
		cFlag.Blob, flag.DescKeyPadImpBlob,
	)

	return c
}

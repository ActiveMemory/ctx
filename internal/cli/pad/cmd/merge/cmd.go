//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the pad merge subcommand.
//
// Returns:
//   - *cobra.Command: Configured merge subcommand
func Cmd() *cobra.Command {
	var keyFile string
	var dryRun bool

	short, long := desc.Command(cmd.DescKeyPadMerge)
	c := &cobra.Command{
		Use:     cmd.UsePadMerge,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyPadMerge),
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, keyFile, dryRun)
		},
	}

	flagbind.StringFlagP(c, &keyFile,
		cFlag.Key, cFlag.ShortKey,
		flag.DescKeyPadMergeKey,
	)
	flagbind.BoolFlag(c, &dryRun,
		cFlag.DryRun, flag.DescKeyPadMergeDryRun,
	)

	return c
}

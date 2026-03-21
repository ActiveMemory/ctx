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
)

// Cmd returns the pad merge subcommand.
//
// Returns:
//   - *cobra.Command: Configured merge subcommand
func Cmd() *cobra.Command {
	var keyFile string
	var dryRun bool

	short, long := desc.CommandDesc(cmd.DescKeyPadMerge)
	c := &cobra.Command{
		Use:   cmd.UsePadMerge,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, keyFile, dryRun)
		},
	}

	c.Flags().StringVarP(&keyFile, cFlag.Key, cFlag.ShortKey, "",
		desc.FlagDesc(flag.DescKeyPadMergeKey))
	c.Flags().BoolVar(&dryRun, cFlag.DryRun, false,
		desc.FlagDesc(flag.DescKeyPadMergeDryRun))

	return c
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the pad export subcommand.
//
// Returns:
//   - *cobra.Command: Configured export subcommand
func Cmd() *cobra.Command {
	var force, dryRun bool

	short, long := desc.Command(cmd.DescKeyPadExport)
	c := &cobra.Command{
		Use:   cmd.UsePadExport,
		Short: short,
		Long:  long,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			return Run(cmd, dir, force, dryRun)
		},
	}

	c.Flags().BoolVarP(
		&force, cFlag.Force, cFlag.ShortForce, false,
		desc.Flag(flag.DescKeyPadExportForce),
	)
	c.Flags().BoolVar(
		&dryRun, cFlag.DryRun, false,
		desc.Flag(flag.DescKeyPadExportDryRun),
	)

	return c
}

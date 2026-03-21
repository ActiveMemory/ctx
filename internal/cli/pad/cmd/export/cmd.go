//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the pad export subcommand.
//
// Returns:
//   - *cobra.Command: Configured export subcommand
func Cmd() *cobra.Command {
	var force, dryRun bool

	short, long := desc.CommandDesc(cmd.DescKeyPadExport)
	cmd := &cobra.Command{
		Use:   "export [DIR]",
		Short: short,
		Long:  long,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			return runExport(cmd, dir, force, dryRun)
		},
	}

	cmd.Flags().BoolVarP(
		&force, "force", "f", false,
		desc.FlagDesc(flag.DescKeyPadExportForce),
	)
	cmd.Flags().BoolVar(
		&dryRun, cflag.DryRun, false,
		desc.FlagDesc(flag.DescKeyPadExportDryRun),
	)

	return cmd
}

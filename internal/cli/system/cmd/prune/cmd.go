//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prune

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system prune" subcommand.
//
// Returns:
//   - *cobra.Command: Configured prune subcommand
func Cmd() *cobra.Command {
	var days int
	var dryRun bool

	short, long := assets.CommandDesc(assets.CmdDescKeySystemPrune)

	cmd := &cobra.Command{
		Use:   "prune",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, days, dryRun)
		},
	}

	cmd.Flags().IntVar(&days, "days", 7,
		assets.FlagDesc(assets.FlagDescKeySystemPruneDays),
	)
	cmd.Flags().BoolVar(&dryRun, "dry-run", false,
		assets.FlagDesc(assets.FlagDescKeySystemPruneDryRun),
	)

	return cmd
}

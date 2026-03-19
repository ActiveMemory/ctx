//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prune

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system prune" subcommand.
//
// Returns:
//   - *cobra.Command: Configured prune subcommand
func Cmd() *cobra.Command {
	var days int
	var dryRun bool

	short, long := desc.CommandDesc(cmd.DescKeySystemPrune)

	cmd := &cobra.Command{
		Use:   "prune",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, days, dryRun)
		},
	}

	cmd.Flags().IntVar(&days, "days", 7,
		desc.FlagDesc(flag.DescKeySystemPruneDays),
	)
	cmd.Flags().BoolVar(&dryRun, "dry-run", false,
		desc.FlagDesc(flag.DescKeySystemPruneDryRun),
	)

	return cmd
}

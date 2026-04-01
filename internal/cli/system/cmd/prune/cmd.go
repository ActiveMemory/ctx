//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prune

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/runtime"
)

// Cmd returns the "ctx system prune" subcommand.
//
// Returns:
//   - *cobra.Command: Configured prune subcommand
func Cmd() *cobra.Command {
	var days int
	var dryRun bool

	short, long := desc.Command(cmd.DescKeySystemPrune)

	c := &cobra.Command{
		Use:   cmd.UseSystemPrune,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, days, dryRun)
		},
	}

	c.Flags().IntVar(&days, cFlag.Days, runtime.DefaultPruneDays,
		desc.Flag(flag.DescKeySystemPruneDays),
	)
	c.Flags().BoolVar(&dryRun, cFlag.DryRun, false,
		desc.Flag(flag.DescKeySystemPruneDryRun),
	)

	return c
}

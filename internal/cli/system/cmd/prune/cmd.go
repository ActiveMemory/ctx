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
	"github.com/ActiveMemory/ctx/internal/flagbind"
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
		Use:     cmd.UseSystemPrune,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySystemPrune),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, days, dryRun)
		},
	}

	flagbind.IntFlag(c, &days,
		cFlag.Days, runtime.DefaultPruneDays,
		flag.DescKeySystemPruneDays,
	)
	flagbind.BoolFlag(c, &dryRun,
		cFlag.DryRun, flag.DescKeySystemPruneDryRun,
	)

	return c
}

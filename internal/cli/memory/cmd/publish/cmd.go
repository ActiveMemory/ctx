//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package publish

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/memory"
)

// Cmd returns the memory publish subcommand.
//
// Returns:
//   - *cobra.Command: command for publishing curated context to MEMORY.md.
func Cmd() *cobra.Command {
	var budget int
	var dryRun bool

	short, long := desc.Command(cmd.DescKeyMemoryPublish)
	c := &cobra.Command{
		Use:   cmd.UseMemoryPublish,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, budget, dryRun)
		},
	}

	c.Flags().IntVar(&budget,
		cFlag.Budget, memory.DefaultPublishBudget,
		desc.Flag(flag.DescKeyMemoryPublishBudget),
	)
	c.Flags().BoolVar(&dryRun,
		cFlag.DryRun, false,
		desc.Flag(flag.DescKeyMemoryPublishDryRun),
	)

	return c
}

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
	"github.com/ActiveMemory/ctx/internal/flagbind"
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
		Use:     cmd.UseMemoryPublish,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyMemoryPublish),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, budget, dryRun)
		},
	}

	flagbind.IntFlag(
		c, &budget,
		cFlag.Budget, memory.DefaultPublishBudget,
		flag.DescKeyMemoryPublishBudget,
	)
	flagbind.BoolFlag(
		c, &dryRun,
		cFlag.DryRun, flag.DescKeyMemoryPublishDryRun,
	)

	return c
}

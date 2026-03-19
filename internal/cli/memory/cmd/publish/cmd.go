//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package publish

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/memory"
	"github.com/spf13/cobra"
)

// Cmd returns the memory publish subcommand.
//
// Returns:
//   - *cobra.Command: command for publishing curated context to MEMORY.md.
func Cmd() *cobra.Command {
	var budget int
	var dryRun bool

	short, long := desc.CommandDesc(cmd.DescKeyMemoryPublish)
	cmd := &cobra.Command{
		Use:   "publish",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, budget, dryRun)
		},
	}

	cmd.Flags().IntVar(&budget,
		"budget", memory.DefaultPublishBudget,
		desc.FlagDesc(flag.DescKeyMemoryPublishBudget),
	)
	cmd.Flags().BoolVar(&dryRun,
		"dry-run", false,
		desc.FlagDesc(flag.DescKeyMemoryPublishDryRun),
	)

	return cmd
}

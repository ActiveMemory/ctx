//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreSchema "github.com/ActiveMemory/ctx/internal/cli/journal/core/schema"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the journal schema check subcommand.
//
// Returns:
//   - *cobra.Command: command for scanning JSONL schema drift
func Cmd() *cobra.Command {
	var opts coreSchema.CheckOpts

	short, long := desc.Command(cmd.DescKeyJournalSchemaCheck)

	c := &cobra.Command{
		Use:     cmd.UseJournalSchemaCheck,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyJournalSchemaCheck),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, opts)
		},
	}

	flagbind.StringFlag(
		c, &opts.Dir,
		cFlag.Dir,
		flag.DescKeyJournalSchemaCheckDir,
	)
	flagbind.BoolFlag(
		c, &opts.AllProjects,
		cFlag.AllProjects,
		flag.DescKeyJournalSchemaCheckAllProjects,
	)
	flagbind.BoolFlagP(
		c, &opts.Quiet,
		cFlag.Quiet, cFlag.ShortQuiet,
		flag.DescKeyJournalSchemaCheckQuiet,
	)

	return c
}

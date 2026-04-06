//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package importer

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the journal import subcommand.
//
// Returns:
//   - *cobra.Command: Command for importing sessions to journal files
func Cmd() *cobra.Command {
	var opts entity.ImportOpts

	short, long := desc.Command(cmd.DescKeyJournalImport)

	c := &cobra.Command{
		Use:     cmd.UseJournalImport,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyJournalImport),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, opts)
		},
	}

	flagbind.BoolFlag(
		c, &opts.All,
		cFlag.All, flag.DescKeyJournalImportAll,
	)
	flagbind.BoolFlag(
		c, &opts.AllProjects,
		cFlag.AllProjects,
		flag.DescKeyJournalImportAllProjects,
	)
	flagbind.BoolFlag(
		c, &opts.Regenerate,
		cFlag.Regenerate,
		flag.DescKeyJournalImportRegenerate,
	)
	flagbind.BoolFlagDefault(
		c, &opts.KeepFrontmatter,
		cFlag.KeepFrontmatter, true,
		flag.DescKeyJournalImportKeepFrontmatter,
	)
	flagbind.BoolFlagP(
		c, &opts.Yes,
		cFlag.Yes, cFlag.ShortYes,
		flag.DescKeyJournalImportYes,
	)
	flagbind.BoolFlag(
		c, &opts.DryRun,
		cFlag.DryRun, flag.DescKeyJournalImportDryRun,
	)

	return c
}

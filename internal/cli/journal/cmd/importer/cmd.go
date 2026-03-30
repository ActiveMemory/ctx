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
)

// Cmd returns the journal import subcommand.
//
// Returns:
//   - *cobra.Command: Command for importing sessions to journal files
func Cmd() *cobra.Command {
	var opts entity.ImportOpts

	short, long := desc.Command(cmd.DescKeyJournalImport)

	c := &cobra.Command{
		Use:   cmd.UseJournalImport,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, opts)
		},
	}

	c.Flags().BoolVar(
		&opts.All, cFlag.All, false, desc.Flag(flag.DescKeyJournalImportAll),
	)
	c.Flags().BoolVar(
		&opts.AllProjects, cFlag.AllProjects, false,
		desc.Flag(flag.DescKeyJournalImportAllProjects),
	)
	c.Flags().BoolVar(
		&opts.Regenerate, cFlag.Regenerate, false,
		desc.Flag(flag.DescKeyJournalImportRegenerate),
	)
	c.Flags().BoolVar(
		&opts.KeepFrontmatter, cFlag.KeepFrontmatter, true,
		desc.Flag(flag.DescKeyJournalImportKeepFrontmatter),
	)
	c.Flags().BoolVarP(
		&opts.Yes, cFlag.Yes, cFlag.ShortYes, false,
		desc.Flag(flag.DescKeyJournalImportYes),
	)
	c.Flags().BoolVar(
		&opts.DryRun, cFlag.DryRun, false,
		desc.Flag(flag.DescKeyJournalImportDryRun),
	)

	return c
}

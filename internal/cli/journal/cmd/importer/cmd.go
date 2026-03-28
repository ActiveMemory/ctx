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

	short, long := desc.Command(cmd.DescKeyRecallImport)

	c := &cobra.Command{
		Use:   cmd.UseRecallImport,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, opts)
		},
	}

	c.Flags().BoolVar(
		&opts.All, cFlag.All, false, desc.Flag(flag.DescKeyRecallImportAll),
	)
	c.Flags().BoolVar(
		&opts.AllProjects, cFlag.AllProjects, false,
		desc.Flag(flag.DescKeyRecallImportAllProjects),
	)
	c.Flags().BoolVar(
		&opts.Regenerate, cFlag.Regenerate, false,
		desc.Flag(flag.DescKeyRecallImportRegenerate),
	)
	c.Flags().BoolVar(
		&opts.KeepFrontmatter, cFlag.KeepFrontmatter, true,
		desc.Flag(flag.DescKeyRecallImportKeepFrontmatter),
	)
	c.Flags().BoolVarP(
		&opts.Yes, cFlag.Yes, cFlag.ShortYes, false,
		desc.Flag(flag.DescKeyRecallImportYes),
	)
	c.Flags().BoolVar(
		&opts.DryRun, cFlag.DryRun, false,
		desc.Flag(flag.DescKeyRecallImportDryRun),
	)

	return c
}

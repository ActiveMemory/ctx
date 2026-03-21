//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the recall export subcommand.
//
// Returns:
//   - *cobra.Command: Command for exporting sessions to journal files
func Cmd() *cobra.Command {
	var opts core.ExportOpts

	short, long := desc.Command(cmd.DescKeyRecallExport)

	c := &cobra.Command{
		Use:   cmd.UseRecallExport,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, opts)
		},
	}

	c.Flags().BoolVar(
		&opts.All, cFlag.All, false, desc.Flag(flag.DescKeyRecallExportAll),
	)
	c.Flags().BoolVar(
		&opts.AllProjects, cFlag.AllProjects, false,
		desc.Flag(flag.DescKeyRecallExportAllProjects),
	)
	c.Flags().BoolVar(
		&opts.Regenerate, cFlag.Regenerate, false,
		desc.Flag(flag.DescKeyRecallExportRegenerate),
	)
	c.Flags().BoolVar(
		&opts.KeepFrontmatter, cFlag.KeepFrontmatter, true,
		desc.Flag(flag.DescKeyRecallExportKeepFrontmatter),
	)
	c.Flags().BoolVarP(
		&opts.Yes, cFlag.Yes, cFlag.ShortYes, false,
		desc.Flag(flag.DescKeyRecallExportYes),
	)
	c.Flags().BoolVar(
		&opts.DryRun, cFlag.DryRun, false,
		desc.Flag(flag.DescKeyRecallExportDryRun),
	)

	return c
}

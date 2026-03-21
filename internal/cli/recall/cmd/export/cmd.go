//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/recall/core"
)

// Cmd returns the recall export subcommand.
//
// Returns:
//   - *cobra.Command: Command for exporting sessions to journal files
func Cmd() *cobra.Command {
	var opts core.ExportOpts

	short, long := desc.CommandDesc(cmd.DescKeyRecallExport)

	cmd := &cobra.Command{
		Use:   "export [session-id]",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, opts)
		},
	}

	cmd.Flags().BoolVar(
		&opts.All, "all", false, desc.FlagDesc(flag.DescKeyRecallExportAll),
	)
	cmd.Flags().BoolVar(
		&opts.AllProjects, "all-projects", false,
		desc.FlagDesc(flag.DescKeyRecallExportAllProjects),
	)
	cmd.Flags().BoolVar(
		&opts.Regenerate,
		"regenerate", false,
		desc.FlagDesc(flag.DescKeyRecallExportRegenerate),
	)
	cmd.Flags().BoolVar(
		&opts.KeepFrontmatter,
		"keep-frontmatter", true,
		desc.FlagDesc(flag.DescKeyRecallExportKeepFrontmatter),
	)

	cmd.Flags().BoolVarP(
		&opts.Yes,
		"yes", "y", false,
		desc.FlagDesc(flag.DescKeyRecallExportYes),
	)
	cmd.Flags().BoolVar(
		&opts.DryRun,
		cflag.DryRun, false,
		desc.FlagDesc(flag.DescKeyRecallExportDryRun),
	)

	// Deprecated: --skip-existing is now the default behavior for --all.
	var skipExisting bool
	cmd.Flags().BoolVar(&skipExisting, "skip-existing", false, desc.FlagDesc(flag.DescKeyRecallExportSkipExisting))
	_ = cmd.Flags().MarkDeprecated("skip-existing", "this is now the default behavior for --all")

	return cmd
}

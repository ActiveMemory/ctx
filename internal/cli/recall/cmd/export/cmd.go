//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"github.com/ActiveMemory/ctx/internal/config/embed"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core"
)

// Cmd returns the recall export subcommand.
//
// Returns:
//   - *cobra.Command: Command for exporting sessions to journal files
func Cmd() *cobra.Command {
	var opts core.ExportOpts

	short, long := assets.CommandDesc(embed.CmdDescKeyRecallExport)

	cmd := &cobra.Command{
		Use:   "export [session-id]",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, opts)
		},
	}

	cmd.Flags().BoolVar(
		&opts.All, "all", false, assets.FlagDesc(embed.FlagDescKeyRecallExportAll),
	)
	cmd.Flags().BoolVar(
		&opts.AllProjects, "all-projects", false,
		assets.FlagDesc(embed.FlagDescKeyRecallExportAllProjects),
	)
	cmd.Flags().BoolVar(
		&opts.Regenerate,
		"regenerate", false,
		assets.FlagDesc(embed.FlagDescKeyRecallExportRegenerate),
	)
	cmd.Flags().BoolVar(
		&opts.KeepFrontmatter,
		"keep-frontmatter", true,
		assets.FlagDesc(embed.FlagDescKeyRecallExportKeepFrontmatter),
	)

	cmd.Flags().BoolVarP(
		&opts.Yes,
		"yes", "y", false,
		assets.FlagDesc(embed.FlagDescKeyRecallExportYes),
	)
	cmd.Flags().BoolVar(
		&opts.DryRun,
		"dry-run", false,
		assets.FlagDesc(embed.FlagDescKeyRecallExportDryRun),
	)

	// Deprecated: --skip-existing is now the default behavior for --all.
	var skipExisting bool
	cmd.Flags().BoolVar(&skipExisting, "skip-existing", false, assets.FlagDesc(embed.FlagDescKeyRecallExportSkipExisting))
	_ = cmd.Flags().MarkDeprecated("skip-existing", "this is now the default behavior for --all")

	return cmd
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
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

	short, long := assets.CommandDesc("recall.export")

	cmd := &cobra.Command{
		Use:   "export [session-id]",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExport(cmd, args, opts)
		},
	}

	cmd.Flags().BoolVar(
		&opts.All, "all", false, assets.FlagDesc("recall.export.all"),
	)
	cmd.Flags().BoolVar(
		&opts.AllProjects, "all-projects", false, assets.FlagDesc("recall.export.all-projects"),
	)
	cmd.Flags().BoolVar(
		&opts.Regenerate,
		"regenerate", false,
		assets.FlagDesc("recall.export.regenerate"),
	)
	cmd.Flags().BoolVar(
		&opts.KeepFrontmatter,
		"keep-frontmatter", true,
		assets.FlagDesc("recall.export.keep-frontmatter"),
	)

	// Deprecated: --force is replaced by --keep-frontmatter=false.
	cmd.Flags().BoolVar(
		&opts.Force,
		"force", false,
		assets.FlagDesc("recall.export.force"),
	)
	_ = cmd.Flags().MarkDeprecated("force", "use --keep-frontmatter=false instead")
	cmd.Flags().BoolVarP(
		&opts.Yes,
		"yes", "y", false,
		assets.FlagDesc("recall.export.yes"),
	)
	cmd.Flags().BoolVar(
		&opts.DryRun,
		"dry-run", false,
		assets.FlagDesc("recall.export.dry-run"),
	)

	// Deprecated: --skip-existing is now the default behavior for --all.
	var skipExisting bool
	cmd.Flags().BoolVar(&skipExisting, "skip-existing", false, assets.FlagDesc("recall.export.skip-existing"))
	_ = cmd.Flags().MarkDeprecated("skip-existing", "this is now the default behavior for --all")

	return cmd
}

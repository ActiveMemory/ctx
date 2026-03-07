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
		&opts.All, "all", false, "Export all sessions from current project",
	)
	cmd.Flags().BoolVar(
		&opts.AllProjects, "all-projects", false, "Include sessions from all projects",
	)
	cmd.Flags().BoolVar(
		&opts.Regenerate,
		"regenerate", false,
		"Re-export existing files (preserves YAML frontmatter by default)",
	)
	cmd.Flags().BoolVar(
		&opts.KeepFrontmatter,
		"keep-frontmatter", true,
		"Preserve enriched YAML frontmatter during regeneration",
	)

	// Deprecated: --force is replaced by --keep-frontmatter=false.
	cmd.Flags().BoolVar(
		&opts.Force,
		"force", false,
		"Overwrite existing files completely (discard frontmatter)",
	)
	_ = cmd.Flags().MarkDeprecated("force", "use --keep-frontmatter=false instead")
	cmd.Flags().BoolVarP(
		&opts.Yes,
		"yes", "y", false,
		"Skip confirmation prompt",
	)
	cmd.Flags().BoolVar(
		&opts.DryRun,
		"dry-run", false,
		"Show what would be exported without writing files",
	)

	// Deprecated: --skip-existing is now the default behavior for --all.
	var skipExisting bool
	cmd.Flags().BoolVar(&skipExisting, "skip-existing", false, "Skip files that already exist")
	_ = cmd.Flags().MarkDeprecated("skip-existing", "this is now the default behavior for --all")

	return cmd
}

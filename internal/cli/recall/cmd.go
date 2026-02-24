//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"github.com/spf13/cobra"
)

// maxMessagesPerPart is the maximum number of messages per exported file.
// Sessions with more messages are split into multiple parts for browser
// performance.
const maxMessagesPerPart = 200

// recallExportCmd returns the recall export subcommand.
//
// Returns:
//   - *cobra.Command: Command for exporting sessions to journal files
func recallExportCmd() *cobra.Command {
	var opts exportOpts

	cmd := &cobra.Command{
		Use:   "export [session-id]",
		Short: "Export sessions to editable journal files",
		Long: `Export AI sessions to .context/journal/ as editable Markdown files.

Exported files include session metadata, tool usage summary, and the full
conversation. You can edit these files to add notes, highlight key moments,
or clean up the transcript.

By default, only sessions from the current project are exported. Use
--all-projects to include sessions from all projects.

Safe by default: --all only exports new sessions. Existing files are
skipped. Use --regenerate to re-export existing files (preserves YAML
frontmatter by default). Use --keep-frontmatter=false to discard
enriched frontmatter during regeneration.

Locked entries (via "ctx recall lock") are always skipped, regardless
of flags.

Examples:
  ctx recall export abc123                              # Export one session
  ctx recall export --all                               # Export only new
  ctx recall export --all --dry-run                     # Preview changes
  ctx recall export --all --regenerate                  # Re-export (prompts)
  ctx recall export --all --regenerate -y               # Re-export, no prompt
  ctx recall export --all --regenerate --keep-frontmatter=false -y  # Discard frontmatter`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecallExport(cmd, args, opts)
		},
	}

	cmd.Flags().BoolVar(
		&opts.all, "all", false, "Export all sessions from current project",
	)
	cmd.Flags().BoolVar(
		&opts.allProjects, "all-projects", false, "Include sessions from all projects",
	)
	cmd.Flags().BoolVar(
		&opts.regenerate,
		"regenerate", false,
		"Re-export existing files (preserves YAML frontmatter by default)",
	)
	cmd.Flags().BoolVar(
		&opts.keepFrontmatter,
		"keep-frontmatter", true,
		"Preserve enriched YAML frontmatter during regeneration",
	)

	// Deprecated: --force is replaced by --keep-frontmatter=false.
	cmd.Flags().BoolVar(
		&opts.force,
		"force", false,
		"Overwrite existing files completely (discard frontmatter)",
	)
	_ = cmd.Flags().MarkDeprecated("force", "use --keep-frontmatter=false instead")
	cmd.Flags().BoolVarP(
		&opts.yes,
		"yes", "y", false,
		"Skip confirmation prompt",
	)
	cmd.Flags().BoolVar(
		&opts.dryRun,
		"dry-run", false,
		"Show what would be exported without writing files",
	)

	// Deprecated: --skip-existing is now the default behavior for --all.
	var skipExisting bool
	cmd.Flags().BoolVar(&skipExisting, "skip-existing", false, "Skip files that already exist")
	_ = cmd.Flags().MarkDeprecated("skip-existing", "this is now the default behavior for --all")

	return cmd
}

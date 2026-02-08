//   /    Context:                     https://ctx.ist
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
	var (
		all          bool
		allProjects  bool
		force        bool
		skipExisting bool
	)

	cmd := &cobra.Command{
		Use:   "export [session-id]",
		Short: "Export sessions to editable journal files",
		Long: `Export AI sessions to .context/journal/ as editable Markdown files.

Exported files include session metadata, tool usage summary, and the full
conversation. You can edit these files to add notes, highlight key moments,
or clean up the transcript.

By default, only sessions from the current project are exported. Use
--all-projects to include sessions from all projects.

By default, existing files are updated: YAML frontmatter from enrichment is
preserved, conversation content is regenerated. Use --skip-existing to leave
existing files untouched, or --force to overwrite completely.

Examples:
  ctx recall export abc123              # Export one session
  ctx recall export --all               # Export/update all sessions
  ctx recall export --all --skip-existing # Skip files that already exist
  ctx recall export --all --force       # Overwrite existing exports completely`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecallExport(cmd, args, all, allProjects, force, skipExisting)
		},
	}

	cmd.Flags().BoolVar(
		&all, "all", false, "Export all sessions from current project",
	)
	cmd.Flags().BoolVar(
		&allProjects, "all-projects", false, "Include sessions from all projects",
	)
	cmd.Flags().BoolVar(
		&force,
		"force", false,
		"Overwrite existing files completely (discard frontmatter)",
	)
	cmd.Flags().BoolVar(
		&skipExisting, "skip-existing", false, "Skip files that already exist",
	)

	return cmd
}

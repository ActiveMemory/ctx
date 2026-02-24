//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"github.com/spf13/cobra"
)

// Cmd returns the journal command with subcommands.
//
// The journal system provides LLM-powered analysis and synthesis of
// exported session files in .context/journal/.
//
// Returns:
//   - *cobra.Command: The journal command with subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "journal",
		Short: "Analyze and synthesize exported sessions",
		Long: `Work with exported session files in .context/journal/.

The journal system provides tools for analyzing, enriching, and
publishing your AI session history.

Subcommands:
  site      Generate a static site from journal entries
  obsidian  Generate an Obsidian vault from journal entries

Examples:
  ctx journal site                    # Generate site in .context/journal-site/
  ctx journal site --output ~/public  # Custom output directory
  ctx journal site --serve            # Generate and serve locally
  ctx journal obsidian                # Generate Obsidian vault`,
	}

	cmd.AddCommand(journalSiteCmd())
	cmd.AddCommand(journalObsidianCmd())

	return cmd
}

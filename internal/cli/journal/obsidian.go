//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// journalObsidianCmd returns the journal obsidian subcommand.
//
// Returns:
//   - *cobra.Command: Command for generating an Obsidian vault from journal
//     entries
func journalObsidianCmd() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "obsidian",
		Short: "Generate an Obsidian vault from journal entries",
		Long: `Generate an Obsidian-compatible vault from .context/journal/ entries.

Creates a vault structure with:
  - Wikilinks for internal navigation
  - MOC (Map of Content) pages for topics, files, and types
  - Related sessions footer for graph connectivity
  - Minimal .obsidian/ configuration

Examples:
  ctx journal obsidian                          # Generate in .context/journal-obsidian/
  ctx journal obsidian --output ~/vaults/ctx    # Custom output directory`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runJournalObsidian(cmd, output)
		},
	}

	defaultOutput := filepath.Join(rc.ContextDir(), config.ObsidianDirName)
	cmd.Flags().StringVarP(
		&output, "output", "o", defaultOutput, "Output directory for vault",
	)

	return cmd
}

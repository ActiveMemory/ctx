//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/rc"
)

// journalSiteCmd returns the journal site subcommand.
//
// Returns:
//   - *cobra.Command: Command for generating a static site from journal entries
func journalSiteCmd() *cobra.Command {
	var (
		output string
		serve  bool
		build  bool
	)

	cmd := &cobra.Command{
		Use:   "site",
		Short: "Generate a static site from journal entries",
		Long: `Generate a zensical-compatible static site from .context/journal/ entries.

Creates a site structure with:
  - Index page with all sessions listed by date
  - Individual pages for each journal entry
  - Navigation and search support

Requires zensical to be installed for building/serving:
  pipx install zensical

Examples:
  ctx journal site                    # Generate in .context/journal-site/
  ctx journal site --output ~/public  # Custom output directory
  ctx journal site --build            # Generate and build HTML
  ctx journal site --serve            # Generate and serve locally`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runJournalSite(cmd, output, build, serve)
		},
	}

	defaultOutput := filepath.Join(rc.ContextDir(), "journal-site")
	cmd.Flags().StringVarP(
		&output, "output", "o", defaultOutput, "Output directory for site",
	)
	cmd.Flags().BoolVar(
		&build, "build", false, "Run zensical build after generating",
	)
	cmd.Flags().BoolVar(
		&serve, "serve", false, "Run zensical serve after generating",
	)

	return cmd
}

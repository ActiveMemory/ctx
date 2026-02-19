//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package learnings

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/compact"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// archiveCmd returns the learnings archive subcommand.
//
// The archive command moves old or superseded learnings from LEARNINGS.md
// to a timestamped archive file in .context/archive/.
//
// Flags:
//   - --days/-d: Days threshold for archiving (default from .contextrc)
//   - --keep/-k: Number of recent entries to keep (default from .contextrc)
//   - --all: Archive all entries except keepRecent
//   - --dry-run: Preview changes without modifying files
//
// Returns:
//   - *cobra.Command: Configured archive subcommand
func archiveCmd() *cobra.Command {
	var (
		days   int
		keep   int
		all    bool
		dryRun bool
	)

	cmd := &cobra.Command{
		Use:   "archive",
		Short: "Archive old learnings to .context/archive/",
		Long: `Archive old or superseded learnings from LEARNINGS.md.

Entries older than --days (default 90) or marked as superseded are moved
to a dated archive file. The most recent --keep entries are always preserved.

Use --all to archive everything except the most recent --keep entries.
Use --dry-run to preview what would be archived.

Examples:
  ctx learnings archive                    # Archive old learnings
  ctx learnings archive --days 30          # Archive learnings older than 30 days
  ctx learnings archive --keep 3           # Keep only 3 most recent
  ctx learnings archive --all              # Archive all except recent
  ctx learnings archive --dry-run          # Preview without changes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := compact.ArchiveKnowledgeFile(
				cmd,
				config.FileLearning,
				"learnings",
				config.HeadingArchivedLearnings,
				index.UpdateLearnings,
				days, keep, all, dryRun,
			)
			return err
		},
	}

	cmd.Flags().IntVarP(&days, "days", "d", rc.ArchiveKnowledgeAfterDays(),
		"Archive entries older than this many days")
	cmd.Flags().IntVarP(&keep, "keep", "k", rc.ArchiveKeepRecent(),
		"Number of recent entries to always keep")
	cmd.Flags().BoolVar(&all, "all", false,
		"Archive all entries except the most recent --keep")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false,
		"Preview changes without modifying files")

	return cmd
}

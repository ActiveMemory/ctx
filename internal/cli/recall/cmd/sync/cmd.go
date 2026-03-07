//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx recall sync" subcommand.
//
// Scans journal markdowns and syncs their frontmatter lock state into
// .state.json. This is the inverse of "ctx recall lock": the frontmatter
// is treated as the source of truth, and state is updated to match.
//
// Returns:
//   - *cobra.Command: Command for syncing lock state from frontmatter
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync lock state from journal frontmatter to state file",
		Long: `Scan journal markdowns and sync their lock state to .state.json.

This is the sister command to "ctx recall lock". Instead of marking files
locked in state and updating frontmatter, it reads "locked: true" from
each file's YAML frontmatter and updates .state.json to match.

Typical workflow:
  1. Enrich journal entries (add "locked: true" to frontmatter)
  2. Run "ctx recall sync" to propagate lock state to .state.json

Files with "locked: true" in frontmatter will be marked locked in state.
Files without a "locked:" line (or with "locked: false") will have their
lock cleared if one exists in state.

Examples:
  ctx recall sync`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSync(cmd)
		},
	}

	return cmd
}

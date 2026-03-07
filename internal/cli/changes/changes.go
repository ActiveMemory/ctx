//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package changes

import (
	"github.com/spf13/cobra"

	changesroot "github.com/ActiveMemory/ctx/internal/cli/changes/cmd/root"
)

// Cmd returns the changes command.
//
// Returns:
//   - *cobra.Command: Configured changes command with flags registered
func Cmd() *cobra.Command {
	var since string

	cmd := &cobra.Command{
		Use:   "changes",
		Short: "Show what changed since last session",
		Long: `Show changes in context files and code since the last AI session.

Automatically detects the last session boundary from state markers.
Use --since to specify a custom time range (duration like "24h" or
date like "2026-03-01").

Examples:
  ctx changes                     # changes since last session
  ctx changes --since 24h         # changes in last 24 hours
  ctx changes --since 2026-03-01  # changes since specific date`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return changesroot.Run(cmd, since)
		},
	}

	cmd.Flags().StringVar(&since, "since", "", "Time reference: duration (24h) or date (2026-03-01)")

	return cmd
}

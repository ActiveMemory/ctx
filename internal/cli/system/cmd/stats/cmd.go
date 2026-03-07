//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system stats" subcommand.
//
// Returns:
//   - *cobra.Command: Configured stats subcommand
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show session token usage stats",
		Long: `Display per-session token usage statistics from stats JSONL files.

By default, shows the last 20 entries across all sessions. Use --follow
to stream new entries as they arrive (like tail -f).

Flags:
  --follow, -f   Stream new entries as they arrive
  --session, -s  Filter by session ID (prefix match)
  --last, -n     Show last N entries (default 20)
  --json, -j     Output raw JSONL`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runStats(cmd)
		},
	}

	cmd.Flags().BoolP("follow", "f", false, "Stream new entries as they arrive")
	cmd.Flags().StringP("session", "s", "", "Filter by session ID (prefix match)")
	cmd.Flags().IntP("last", "n", 20, "Show last N entries")
	cmd.Flags().BoolP("json", "j", false, "Output raw JSONL")

	return cmd
}

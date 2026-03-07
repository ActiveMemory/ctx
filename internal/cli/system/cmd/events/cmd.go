//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package events

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system events" subcommand.
//
// Returns:
//   - *cobra.Command: Configured events subcommand
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Query the local hook event log",
		Long: `Query the local event log (requires event_log: true in .ctxrc).

Reads events from .context/state/events.jsonl and outputs them in
human-readable or raw JSONL format. All filter flags use intersection
(AND) logic.

Flags:
  --hook       Filter by hook name
  --session    Filter by session ID
  --event      Filter by event type (relay, nudge)
  --last       Show last N events (default 50)
  --json       Output raw JSONL (for piping to jq)
  --all        Include rotated log file`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runEvents(cmd)
		},
	}

	cmd.Flags().StringP("hook", "k", "", assets.FlagDesc(assets.FlagDescKeySystemEventsHook))
	cmd.Flags().StringP("session", "s", "", assets.FlagDesc(assets.FlagDescKeySystemEventsSession))
	cmd.Flags().StringP("event", "e", "", assets.FlagDesc(assets.FlagDescKeySystemEventsEvent))
	cmd.Flags().IntP("last", "n", 50, assets.FlagDesc(assets.FlagDescKeySystemEventsLast))
	cmd.Flags().BoolP("json", "j", false, assets.FlagDesc(assets.FlagDescKeySystemEventsJson))
	cmd.Flags().BoolP("all", "a", false, assets.FlagDesc(assets.FlagDescKeySystemEventsAll))

	return cmd
}

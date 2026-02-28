//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// eventsCmd returns the "ctx system events" command.
//
// Raw event log query. Filters and outputs events from
// .context/state/events.jsonl without analysis.
func eventsCmd() *cobra.Command {
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

	cmd.Flags().StringP("hook", "k", "", "Filter by hook name")
	cmd.Flags().StringP("session", "s", "", "Filter by session ID")
	cmd.Flags().StringP("event", "e", "", "Filter by event type")
	cmd.Flags().IntP("last", "n", 50, "Show last N events")
	cmd.Flags().BoolP("json", "j", false, "Output raw JSONL")
	cmd.Flags().BoolP("all", "a", false, "Include rotated log file")

	return cmd
}

func runEvents(cmd *cobra.Command) error {
	hook, _ := cmd.Flags().GetString("hook")
	session, _ := cmd.Flags().GetString("session")
	event, _ := cmd.Flags().GetString("event")
	last, _ := cmd.Flags().GetInt("last")
	jsonOut, _ := cmd.Flags().GetBool("json")
	includeAll, _ := cmd.Flags().GetBool("all")

	opts := eventlog.QueryOpts{
		Hook:           hook,
		Session:        session,
		Event:          event,
		Last:           last,
		IncludeRotated: includeAll,
	}

	events, queryErr := eventlog.Query(opts)
	if queryErr != nil {
		return fmt.Errorf("reading event log: %w", queryErr)
	}

	if len(events) == 0 {
		cmd.Println("No events logged.")
		return nil
	}

	if jsonOut {
		return outputEventsJSON(cmd, events)
	}
	return outputEventsHuman(cmd, events)
}

// outputEventsJSON writes events as raw JSONL.
func outputEventsJSON(cmd *cobra.Command, events []notify.Payload) error {
	for _, e := range events {
		line, marshalErr := json.Marshal(e)
		if marshalErr != nil {
			continue
		}
		cmd.Println(string(line))
	}
	return nil
}

// outputEventsHuman writes events in aligned columns.
func outputEventsHuman(cmd *cobra.Command, events []notify.Payload) error {
	for _, e := range events {
		ts := formatEventTimestamp(e.Timestamp)
		hookName := extractHookName(e)
		msg := truncateMessage(e.Message, 60)
		cmd.Printf("%-19s  %-5s  %-24s  %s\n", ts, e.Event, hookName, msg)
	}
	return nil
}

// formatEventTimestamp converts an RFC3339 timestamp to local time display.
func formatEventTimestamp(ts string) string {
	t, parseErr := time.Parse(time.RFC3339, ts)
	if parseErr != nil {
		return ts
	}
	return t.Local().Format("2006-01-02 15:04:05")
}

// extractHookName gets the hook name from the event payload detail.
func extractHookName(e notify.Payload) string {
	if e.Detail != nil && e.Detail.Hook != "" {
		return e.Detail.Hook
	}
	// Fall back to extracting from message prefix (e.g., "qa-reminder: ...")
	if idx := strings.Index(e.Message, ":"); idx > 0 {
		return e.Message[:idx]
	}
	return "-"
}

// truncateMessage limits message length for display.
func truncateMessage(msg string, maxLen int) string {
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen-3] + "..."
}

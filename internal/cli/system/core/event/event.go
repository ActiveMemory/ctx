//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/event"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// FormatTimestamp converts an RFC3339 timestamp to local time display
// using the DateTimePreciseFormat layout.
//
// Parameters:
//   - ts: RFC3339-formatted timestamp string
//
// Returns:
//   - string: local time formatted as "2006-01-02 15:04:05", or the
//     original string on parse failure
func FormatTimestamp(ts string) string {
	t, parseErr := time.Parse(time.RFC3339, ts)
	if parseErr != nil {
		return ts
	}
	return t.Local().Format(cfgTime.DateTimePreciseFormat)
}

// ExtractHookName gets the hook name from an event payload's detail field.
// Falls back to extracting from the message prefix (e.g., "qa-reminder: ...").
//
// Parameters:
//   - e: event payload to inspect
//
// Returns:
//   - string: hook name, or EventsHookFallback if undetermined
func ExtractHookName(e notify.Payload) string {
	if e.Detail != nil && e.Detail.Hook != "" {
		return e.Detail.Hook
	}
	// Fall back to extracting from message prefix (e.g., "qa-reminder: ...")
	if idx := strings.Index(e.Message, ":"); idx > 0 {
		return e.Message[:idx]
	}
	return event.EventsHookFallback
}

// TruncateMessage limits message length for display, appending a
// truncation suffix when the message exceeds maxLen characters.
//
// Parameters:
//   - msg: message to potentially truncate
//   - maxLen: maximum allowed length including suffix
//
// Returns:
//   - string: original or truncated message
func TruncateMessage(msg string, maxLen int) string {
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen-len(event.EventsTruncationSuffix)] +
		event.EventsTruncationSuffix
}

// FormatJSON formats events as raw JSONL lines.
//
// Parameters:
//   - evts: event payloads to serialize
//
// Returns:
//   - []string: JSON lines (marshal errors are silently skipped)
func FormatJSON(evts []notify.Payload) []string {
	var lines []string
	for _, e := range evts {
		line, marshalErr := json.Marshal(e)
		if marshalErr != nil {
			continue
		}
		lines = append(lines, string(line))
	}
	return lines
}

// FormatHuman formats events in aligned columns for human reading.
//
// Parameters:
//   - evts: event payloads to format
//
// Returns:
//   - []string: formatted event lines
func FormatHuman(evts []notify.Payload) []string {
	fmtStr := desc.Text(text.DescKeyEventsHumanFormat)
	lines := make([]string, 0, len(evts))
	for _, e := range evts {
		ts := FormatTimestamp(e.Timestamp)
		hookName := ExtractHookName(e)
		msg := TruncateMessage(e.Message, event.EventsMessageMaxLen)
		lines = append(lines, fmt.Sprintf(fmtStr, ts, e.Event, hookName, msg))
	}
	return lines
}

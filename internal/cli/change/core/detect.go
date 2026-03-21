//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/event"
	"github.com/ActiveMemory/ctx/internal/config/load_gate"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/format"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// DetectReferenceTime determines the reference time for change detection.
//
// Priority:
//  1. --since flag (duration like "24h" or date like "2026-03-01")
//  2. ctx-loaded-* marker files (second most recent by mtime)
//  3. events.jsonl (last context-load-gate event)
//  4. Fallback to 24h ago
//
// Parameters:
//   - since: User-provided time reference, or empty for auto-detection
//
// Returns:
//   - time.Time: The determined reference time
//   - string: Human-readable label describing the reference point
//   - error: Non-nil if the --since value cannot be parsed
func DetectReferenceTime(since string) (time.Time, string, error) {
	if since != "" {
		return ParseSinceFlag(since)
	}

	// Try marker files.
	if t, ok := DetectFromMarkers(); ok {
		return t, format.DurationAgo(time.Since(t)), nil
	}

	// Try events.jsonl.
	if t, ok := DetectFromEvents(); ok {
		return t, format.DurationAgo(time.Since(t)), nil
	}

	// Fallback: 24h ago.
	t := time.Now().Add(-24 * time.Hour)
	return t, desc.TextDesc(text.DescKeyChangesFallbackLabel), nil
}

// ParseSinceFlag parses a duration (like "24h") or date (like "2026-03-01").
//
// Parameters:
//   - since: Time reference string to parse
//
// Returns:
//   - time.Time: Parsed time
//   - string: Human-readable label
//   - error: Non-nil if parsing fails
func ParseSinceFlag(since string) (time.Time, string, error) {
	// Try duration first.
	if d, err := time.ParseDuration(since); err == nil {
		t := time.Now().Add(-d)
		return t, format.DurationAgo(d), nil
	}

	// Try date.
	if t, err := time.Parse(cfgTime.DateFormat, since); err == nil {
		return t, desc.TextDesc(text.DescKeyChangesSincePrefix) + since, nil
	}

	// Try RFC3339.
	if t, err := time.Parse(time.RFC3339, since); err == nil {
		return t, format.DurationAgo(time.Since(t)), nil
	}

	return time.Time{}, "", os.ErrInvalid
}

// DetectFromMarkers finds the second most recent ctx-loaded-* marker file.
// The most recent is the current session's marker.
//
// Returns:
//   - time.Time: Marker file modification time
//   - bool: True if a valid marker was found
func DetectFromMarkers() (time.Time, bool) {
	stateDir := filepath.Join(rc.ContextDir(), dir.State)
	entries, err := os.ReadDir(stateDir)
	if err != nil {
		return time.Time{}, false
	}

	type markerInfo struct {
		modTime time.Time
	}

	var markers []markerInfo
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), load_gate.PrefixCtxLoaded) {
			continue
		}
		info, infoErr := e.Info()
		if infoErr != nil {
			continue
		}
		markers = append(markers, markerInfo{modTime: info.ModTime()})
	}

	if len(markers) < 2 {
		return time.Time{}, false
	}

	// Sort by modtime descending.
	sort.Slice(markers, func(i, j int) bool {
		return markers[i].modTime.After(markers[j].modTime)
	})

	// Second most recent = previous session.
	return markers[1].modTime, true
}

// DetectFromEvents scans events.jsonl in reverse for the last
// context-load-gate event.
//
// Returns:
//   - time.Time: Event timestamp
//   - bool: True if a valid event was found
func DetectFromEvents() (time.Time, bool) {
	eventsPath := filepath.Join(rc.ContextDir(), dir.State, event.FileEventLog)
	data, err := io.SafeReadUserFile(eventsPath)
	if err != nil {
		return time.Time{}, false
	}

	lines := strings.Split(strings.TrimSpace(string(data)), token.NewlineLF)
	// Scan in reverse for the last context-load-gate event.
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if !strings.Contains(line, load_gate.EventContextLoadGate) {
			continue
		}
		if t, ok := ExtractTimestamp(line); ok {
			return t, true
		}
	}

	return time.Time{}, false
}

// ExtractTimestamp extracts a timestamp from a JSON line without full unmarshal.
// Looks for "timestamp":"..." and parses as RFC3339.
//
// Parameters:
//   - jsonLine: JSON string to extract timestamp from
//
// Returns:
//   - time.Time: Parsed timestamp
//   - bool: True if extraction succeeded
func ExtractTimestamp(jsonLine string) (time.Time, bool) {
	key := load_gate.JSONKeyTimestamp
	idx := strings.Index(jsonLine, key)
	if idx < 0 {
		return time.Time{}, false
	}
	start := idx + len(key)
	end := strings.Index(jsonLine[start:], token.DoubleQuote)
	if end < 0 {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, jsonLine[start:start+end])
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package changes

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
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
// Returns the reference time, a human-readable label, and any error.
func DetectReferenceTime(since string) (time.Time, string, error) {
	if since != "" {
		return parseSinceFlag(since)
	}

	// Try marker files.
	if t, ok := detectFromMarkers(); ok {
		return t, humanAgo(time.Since(t)), nil
	}

	// Try events.jsonl.
	if t, ok := detectFromEvents(); ok {
		return t, humanAgo(time.Since(t)), nil
	}

	// Fallback: 24h ago.
	t := time.Now().Add(-24 * time.Hour)
	return t, "24 hour(s) ago (default)", nil
}

// parseSinceFlag parses a duration (like "24h") or date (like "2026-03-01").
func parseSinceFlag(since string) (time.Time, string, error) {
	// Try duration first.
	if d, err := time.ParseDuration(since); err == nil {
		t := time.Now().Add(-d)
		return t, humanAgo(d), nil
	}

	// Try date.
	if t, err := time.Parse("2006-01-02", since); err == nil {
		return t, "since " + since, nil
	}

	// Try RFC3339.
	if t, err := time.Parse(time.RFC3339, since); err == nil {
		return t, humanAgo(time.Since(t)), nil
	}

	return time.Time{}, "", os.ErrInvalid
}

// detectFromMarkers finds the second most recent ctx-loaded-* marker file.
// The most recent is the current session's marker.
func detectFromMarkers() (time.Time, bool) {
	stateDir := filepath.Join(rc.ContextDir(), config.DirState)
	entries, err := os.ReadDir(stateDir)
	if err != nil {
		return time.Time{}, false
	}

	type markerInfo struct {
		modTime time.Time
	}

	var markers []markerInfo
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "ctx-loaded-") {
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

// detectFromEvents scans events.jsonl in reverse for the last
// context-load-gate event.
func detectFromEvents() (time.Time, bool) {
	eventsPath := filepath.Join(rc.ContextDir(), config.DirState, "events.jsonl")
	data, err := os.ReadFile(eventsPath) //nolint:gosec // state dir path
	if err != nil {
		return time.Time{}, false
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	// Scan in reverse for last context-load-gate event.
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if !strings.Contains(line, "context-load-gate") {
			continue
		}
		if t, ok := extractTimestamp(line); ok {
			return t, true
		}
	}

	return time.Time{}, false
}

// extractTimestamp extracts a timestamp from a JSON line without full unmarshal.
// Looks for "timestamp":"..." and parses as RFC3339.
func extractTimestamp(jsonLine string) (time.Time, bool) {
	const key = `"timestamp":"`
	idx := strings.Index(jsonLine, key)
	if idx < 0 {
		return time.Time{}, false
	}
	start := idx + len(key)
	end := strings.Index(jsonLine[start:], `"`)
	if end < 0 {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, jsonLine[start:start+end])
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

// humanAgo returns a human-readable "ago" string from a duration.
func humanAgo(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		return pluralize(m, "minute") + " ago"
	case d < 24*time.Hour:
		h := int(d.Hours())
		return pluralize(h, "hour") + " ago"
	default:
		days := int(d.Hours() / 24)
		return pluralize(days, "day") + " ago"
	}
}

// pluralize returns "N unit" or "N units".
func pluralize(n int, unit string) string {
	if n == 1 {
		return "1 " + unit
	}
	return itoa(n) + " " + unit + "s"
}

// itoa is a minimal int-to-string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

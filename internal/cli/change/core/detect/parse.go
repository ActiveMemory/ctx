//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package detect

import (
	"os"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/load_gate"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/format"
)

// ReferenceTime determines the reference time for change detection.
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
func ReferenceTime(since string) (time.Time, string, error) {
	if since != "" {
		return ParseSinceFlag(since)
	}

	// Try marker files.
	if t, ok := FromMarkers(); ok {
		return t, format.DurationAgo(time.Since(t)), nil
	}

	// Try events.jsonl.
	if t, ok := FromEvents(); ok {
		return t, format.DurationAgo(time.Since(t)), nil
	}

	// Fallback: 24h ago.
	t := time.Now().Add(-24 * time.Hour)
	return t, desc.Text(text.DescKeyChangesFallbackLabel), nil
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
	if d, durationErr := time.ParseDuration(since); durationErr == nil {
		t := time.Now().Add(-d)
		return t, format.DurationAgo(d), nil
	}

	// Try date.
	if t, dateErr := time.Parse(cfgTime.DateFormat, since); dateErr == nil {
		return t, desc.Text(text.DescKeyChangesSincePrefix) + since, nil
	}

	// Try RFC3339.
	if t, rfcErr := time.Parse(time.RFC3339, since); rfcErr == nil {
		return t, format.DurationAgo(time.Since(t)), nil
	}

	return time.Time{}, "", os.ErrInvalid
}

// ExtractTimestamp extracts a timestamp from a JSON
// line without full unmarshal.
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
	t, parseErr := time.Parse(time.RFC3339, jsonLine[start:start+end])
	if parseErr != nil {
		return time.Time{}, false
	}
	return t, true
}

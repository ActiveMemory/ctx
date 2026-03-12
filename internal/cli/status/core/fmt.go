//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"time"

	ctxtime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/write"
)

// FormatTimeAgo returns a human-readable relative time string.
//
// Examples: "just now", "5 minutes ago", "2 hours ago", "3 days ago",
// or a formatted date for times older than a week.
//
// Parameters:
//   - t: The time to format relative to now
//
// Returns:
//   - string: Human-readable relative time
func FormatTimeAgo(t time.Time) string {
	d := time.Since(t)
	return write.FormatTimeAgo(
		d.Hours(), int(d.Minutes()), t.Format(ctxtime.OlderFormat),
	)
}

// FormatNumber returns a number with thousand separators.
//
// Parameters:
//   - n: The number to format
//
// Returns:
//   - string: Formatted number with commas
func FormatNumber(n int) string {
	return write.FormatNumber(n)
}

// FormatBytes returns a human-readable byte-size string.
//
// Parameters:
//   - b: The byte count to format
//
// Returns:
//   - string: Human-readable size with unit
func FormatBytes(b int64) string {
	return write.FormatBytes(b)
}

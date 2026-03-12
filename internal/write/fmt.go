//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/write/config"
)

// FormatTimeAgo returns a human-readable relative time duration.
//
// Examples: "just now", "5 minutes ago", "2 hours ago", "3 days ago",
// or a formatted date for times older than a week.
//
// Parameters:
//   - d: Duration since the event
//   - fallbackDate: Formatted date string for durations older than a week
//
// Returns:
//   - string: Human-readable relative time
func FormatTimeAgo(hours float64, mins int, fallbackDate string) string {
	switch {
	case hours < 1.0/60: // less than a minute
		return config.TplTimeJustNow
	case hours < 1:
		if mins == 1 {
			return config.TplTimeMinuteAgo
		}
		return fmt.Sprintf(config.TplTimeMinutesAgo, mins)
	case hours < 24:
		h := int(hours)
		if h == 1 {
			return config.TplTimeHourAgo
		}
		return fmt.Sprintf(config.TplTimeHoursAgo, h)
	case hours < 7*24:
		days := int(hours / 24)
		if days == 1 {
			return config.TplTimeDayAgo
		}
		return fmt.Sprintf(config.TplTimeDaysAgo, days)
	default:
		return fallbackDate
	}
}

// FormatNumber returns a number with thousand separators.
//
// Examples: 500 -> "500", 1500 -> "1,500", 12345 -> "12,345"
//
// Parameters:
//   - n: The number to format
//
// Returns:
//   - string: Formatted number with commas
func FormatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%d,%03d", n/1000, n%1000)
}

// FormatBytes returns a human-readable byte-size string.
//
// Uses binary units (1024-based): B, KB, MB, GB, etc.
//
// Parameters:
//   - b: The byte count to format
//
// Returns:
//   - string: Human-readable size with unit
func FormatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

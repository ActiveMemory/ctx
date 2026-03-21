//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// TimeAgo returns a human-readable relative time duration.
//
// Examples: "just now", "5 minutes ago", "2 hours ago", "3 days ago",
// or a formatted date for times older than a week.
//
// Parameters:
//   - hours: total hours since the event
//   - mins: total minutes since the event
//   - fallbackDate: formatted date string for durations older than
//     a week
//
// Returns:
//   - string: human-readable relative time
func TimeAgo(hours float64, mins int, fallbackDate string) string {
	switch {
	case hours < 1.0/60: // less than a minute
		return desc.TextDesc(text.DescKeyWriteTimeJustNow)
	case hours < 1:
		if mins == 1 {
			return desc.TextDesc(text.DescKeyWriteTimeMinuteAgo)
		}
		return fmt.Sprintf(desc.TextDesc(text.DescKeyWriteTimeMinutesAgo), mins)
	case hours < 24:
		h := int(hours)
		if h == 1 {
			return desc.TextDesc(text.DescKeyWriteTimeHourAgo)
		}
		return fmt.Sprintf(desc.TextDesc(text.DescKeyWriteTimeHoursAgo), h)
	case hours < 7*24:
		days := int(hours / 24)
		if days == 1 {
			return desc.TextDesc(text.DescKeyWriteTimeDayAgo)
		}
		return fmt.Sprintf(desc.TextDesc(text.DescKeyWriteTimeDaysAgo), days)
	default:
		return fallbackDate
	}
}

// Pluralize returns "1 unit" or "N units" for English pluralization.
//
// Parameters:
//   - n: count
//   - unit: singular form of the unit word
//
// Returns:
//   - string: pluralized string (e.g., "1 commit", "3 commits")
func Pluralize(n int, unit string) string {
	if n == 1 {
		return "1" + token.Space + unit
	}
	return strconv.Itoa(n) + token.Space + unit + "s"
}

// Duration returns a human-readable duration string without a suffix.
//
// Parameters:
//   - d: duration to format
//
// Returns:
//   - string: human-readable representation (e.g., "3 hours", "1 day", "just now")
func Duration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return desc.TextDesc(text.DescKeyTimeJustNow)
	case d < time.Hour:
		return Pluralize(int(d.Minutes()), desc.TextDesc(text.DescKeyTimeMinute))
	case d < 24*time.Hour:
		return Pluralize(int(d.Hours()), desc.TextDesc(text.DescKeyTimeHour))
	default:
		return Pluralize(int(d.Hours()/24), desc.TextDesc(text.DescKeyTimeDay))
	}
}

// DurationAgo returns a human-readable duration with an "ago" suffix.
//
// Parameters:
//   - d: duration to format
//
// Returns:
//   - string: relative time (e.g., "3 hours ago", "just now")
func DurationAgo(d time.Duration) string {
	base := Duration(d)
	if d < time.Minute {
		return base
	}
	return base + desc.TextDesc(text.DescKeyTimeAgo)
}

// TruncateFirstLine returns the first line of s, capped at max characters.
//
// Parameters:
//   - s: Input string (may be multi-line)
//   - max: Maximum length including ellipsis
//
// Returns:
//   - string: Truncated first line
func TruncateFirstLine(s string, max int) string {
	line, _, _ := strings.Cut(s, token.NewlineLF)
	if len(line) <= max {
		return line
	}
	return line[:max-len(token.Ellipsis)] + token.Ellipsis
}

// Number returns a number with thousand separators.
//
// Examples: 500 → "500", 1500 → "1,500", 12345 → "12,345"
//
// Parameters:
//   - n: the number to format
//
// Returns:
//   - string: formatted number with commas
func Number(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%d,%03d", n/1000, n%1000)
}

// Bytes returns a human-readable byte-size string.
//
// Uses binary units (1024-based): B, KB, MB, GB, etc.
//
// Parameters:
//   - b: the byte count to format
//
// Returns:
//   - string: human-readable size with unit
func Bytes(b int64) string {
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

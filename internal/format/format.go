//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

import (
	"fmt"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgFmt "github.com/ActiveMemory/ctx/internal/config/format"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
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
	case hours < 1.0/cfgTime.MinutesPerHour: // less than a minute
		return desc.Text(text.DescKeyWriteTimeJustNow)
	case hours < 1:
		if mins == 1 {
			return desc.Text(text.DescKeyWriteTimeMinuteAgo)
		}
		return fmt.Sprintf(desc.Text(text.DescKeyWriteTimeMinutesAgo), mins)
	case hours < cfgTime.HoursPerDay:
		h := int(hours)
		if h == 1 {
			return desc.Text(text.DescKeyWriteTimeHourAgo)
		}
		return fmt.Sprintf(desc.Text(text.DescKeyWriteTimeHoursAgo), h)
	case hours < 7*cfgTime.HoursPerDay:
		days := int(hours / cfgTime.HoursPerDay)
		if days == 1 {
			return desc.Text(text.DescKeyWriteTimeDayAgo)
		}
		return fmt.Sprintf(desc.Text(text.DescKeyWriteTimeDaysAgo), days)
	default:
		return fallbackDate
	}
}

// Duration returns a human-readable duration string without a suffix.
//
// Parameters:
//   - d: duration to format
//
// Returns:
//   - string: human-readable representation
//     (e.g., "3 hours", "1 day", "just now")
func Duration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return desc.Text(text.DescKeyTimeJustNow)
	case d < time.Hour:
		n := int(d.Minutes())
		if n == 1 {
			return desc.Text(text.DescKeyTimeMinuteCount)
		}
		return fmt.Sprintf(desc.Text(text.DescKeyTimeMinutesCount), n)
	case d < cfgTime.HoursPerDay*time.Hour:
		n := int(d.Hours())
		if n == 1 {
			return desc.Text(text.DescKeyTimeHourCount)
		}
		return fmt.Sprintf(desc.Text(text.DescKeyTimeHoursCount), n)
	default:
		n := int(d.Hours() / cfgTime.HoursPerDay)
		if n == 1 {
			return desc.Text(text.DescKeyTimeDayCount)
		}
		return fmt.Sprintf(desc.Text(text.DescKeyTimeDaysCount), n)
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
	return base + desc.Text(text.DescKeyTimeAgo)
}

// Today returns today's date as YYYY-MM-DD.
//
// Returns:
//   - string: Current date formatted per cfgTime.DateFormat
func Today() string {
	return time.Now().Format(cfgTime.DateFormat)
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
	if n < cfgFmt.SIThreshold {
		return fmt.Sprintf(desc.Text(text.DescKeyWriteFormatSIInteger), n)
	}
	return fmt.Sprintf(
		desc.Text(text.DescKeyWriteFormatThousands),
		n/cfgFmt.SIThreshold, n%cfgFmt.SIThreshold,
	)
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
	if b < cfgFmt.IECUnit {
		return fmt.Sprintf(desc.Text(text.DescKeyWriteFormatBytesRaw), b)
	}
	div, exp := int64(cfgFmt.IECUnit), 0
	for n := b / cfgFmt.IECUnit; n >= cfgFmt.IECUnit; n /= cfgFmt.IECUnit {
		div *= cfgFmt.IECUnit
		exp++
	}
	return fmt.Sprintf(
		desc.Text(text.DescKeyWriteFormatBytesUnit),
		float64(b)/float64(div), "KMGTPE"[exp],
	)
}

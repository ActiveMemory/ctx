//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"strconv"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// FormatDuration returns a human-readable duration string.
//
// Parameters:
//   - d: duration to format
//
// Returns:
//   - string: human-readable representation (e.g. "3 hours", "1 day")
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return assets.TextDesc(assets.TextDescKeyTimeJustNow)
	}
	if d < time.Hour {
		return pluralize(int(d.Minutes()),
			assets.TextDesc(assets.TextDescKeyTimeMinute))
	}
	h := int(d.Hours())
	if h < 24 {
		return pluralize(h,
			assets.TextDesc(assets.TextDescKeyTimeHour))
	}
	return pluralize(h/24,
		assets.TextDesc(assets.TextDescKeyTimeDay))
}

// pluralize returns "1 unit" or "N units".
func pluralize(n int, unit string) string {
	if n == 1 {
		return "1 " + unit
	}
	return strconv.Itoa(n) + " " + unit + "s"
}

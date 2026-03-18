//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"time"

	ctxtime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/format"
)

// FormatTimeAgo returns a human-readable relative time string.
//
// Parameters:
//   - t: the time to format relative to now
//
// Returns:
//   - string: human-readable relative time
func FormatTimeAgo(t time.Time) string {
	d := time.Since(t)
	return format.TimeAgo(
		d.Hours(), int(d.Minutes()), t.Format(ctxtime.OlderFormat),
	)
}

// FormatNumber returns a number with thousand separators.
//
// Parameters:
//   - n: the number to format
//
// Returns:
//   - string: formatted number with commas
func FormatNumber(n int) string {
	return format.Number(n)
}

// FormatBytes returns a human-readable byte-size string.
//
// Parameters:
//   - b: the byte count to format
//
// Returns:
//   - string: human-readable size with unit
func FormatBytes(b int64) string {
	return format.Bytes(b)
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
)

// ImportResult tracks counts per target for import reporting.
type ImportResult struct {
	Conventions int
	Decisions   int
	Learnings   int
	Tasks       int
	Skipped     int
	Dupes       int
}

// Total returns the number of entries actually imported (excludes skips
// and duplicates).
//
// Returns:
//   - int: count of imported entries.
func (r ImportResult) Total() int {
	return r.Conventions + r.Decisions + r.Learnings + r.Tasks
}

// CountFileLines counts the number of newline characters in data.
//
// Parameters:
//   - data: raw file bytes.
//
// Returns:
//   - int: number of newline characters.
func CountFileLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	count := 0
	for _, b := range data {
		if b == '\n' {
			count++
		}
	}
	return count
}

// FormatDuration returns a human-readable duration string.
//
// Parameters:
//   - d: duration to format.
//
// Returns:
//   - string: human-readable representation (e.g. "3 hours", "1 day").
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", m)
	}
	h := int(d.Hours())
	if h == 1 {
		return "1 hour"
	}
	if h < 24 {
		return fmt.Sprintf("%d hours", h)
	}
	days := h / 24
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}

// Truncate returns the first line of s, capped at max characters.
//
// Parameters:
//   - s: input string (may be multi-line).
//   - max: maximum length including ellipsis.
//
// Returns:
//   - string: truncated first line.
func Truncate(s string, max int) string {
	line := strings.SplitN(s, config.NewlineLF, 2)[0]
	if len(line) <= max {
		return line
	}
	return line[:max-3] + "..."
}

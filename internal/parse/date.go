//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

import (
	"time"

	time2 "github.com/ActiveMemory/ctx/internal/config/time"
)

// Date parses a YYYY-MM-DD string into a time.Time at midnight UTC.
// An empty string returns the zero time with no error.
//
// Parameters:
//   - s: date string in YYYY-MM-DD format. Empty is a no-op.
//
// Returns:
//   - time.Time: parsed date at midnight UTC, or zero time if s is empty.
//   - error: non-nil if the string is non-empty and malformed.
func Date(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(time2.DateFormat, s)
}

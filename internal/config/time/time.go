//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package time

import "time"

// Date and time format constants.
const (
	// DateFormat is the canonical YYYY-MM-DD date layout for time.Parse.
	DateFormat = "2006-01-02"
	// DateTimeFormat is DateFormat with hours and minutes (HH:MM).
	DateTimeFormat = "2006-01-02 15:04"
	// DateTimePreciseFormat is DateFormat with hours, minutes, and seconds.
	DateTimePreciseFormat = "2006-01-02 15:04:05"
	// Format is the hours:minutes:seconds layout for timestamps.
	Format = "15:04:05"
	// TimestampCompact is the YYYYMMDD-HHMMSS layout used in entry headers
	// and task timestamps (e.g., 2026-01-28-143022).
	TimestampCompact = "2006-01-02-150405"
)

// InclusiveUntilOffset is the duration added to an --until date to make
// it inclusive of the entire day (23:59:59).
const InclusiveUntilOffset = 24*time.Hour - time.Second

// OlderFormat is the Go time layout for dates older than a week.
// Exported because callers must format the fallback date before calling FormatTimeAgo.
const OlderFormat = "Jan 2, 2006"

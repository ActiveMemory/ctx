//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package time defines date and time format layouts,
// duration constants, and date-parsing helpers used
// throughout ctx.
//
// ctx timestamps appear in journal entries, task
// headers, commit trailers, feed entries, and CLI
// output. This package provides the canonical Go
// time.Parse layout strings so every subsystem
// formats and parses dates identically.
//
// # Format Layouts
//
//   - [DateFormat] ("2006-01-02") — the canonical
//     YYYY-MM-DD layout for dates.
//   - [DateTimeFmt] — date with hours and minutes.
//   - [DateTimePreciseFmt] — date with full HH:MM:SS.
//   - [Format] ("15:04:05") — time-only layout.
//   - [CompactTimestamp] — YYYYMMDD-HHMMSS layout
//     used in entry headers and task timestamps.
//   - [OlderFormat] ("Jan 2, 2006") — human-friendly
//     layout for dates older than a week.
//
// # Duration Constants
//
//   - [HoursPerDay] (24), [MinutesPerHour] (60),
//     [DaysPerWeek] (7) — integer constants for
//     duration arithmetic that avoids magic numbers.
//
// # Date Parsing Helpers
//
//   - [DateMinLen] — minimum string length for a
//     YYYY-MM-DD date (10 characters).
//   - [DateHyphenPos1], [DateHyphenPos2] — byte
//     positions of hyphens for fast validation.
//   - [InclusiveUntilOffset] — duration added to
//     --until flags to include the entire end day.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package time

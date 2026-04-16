//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check provides shared preamble logic and
// throttling for system hook checks.
//
// Most hooks share an identical startup sequence:
// read input, resolve the session ID, and check the
// pause state. This package extracts that common
// preamble so each hook does not duplicate it.
//
// # Preamble
//
// [Preamble] reads hook input from stdin as JSON,
// extracts the session ID (falling back to a
// configured unknown-session constant), and checks
// whether the session is currently paused via the
// nudge pause system. Returns the parsed input,
// resolved session ID, and pause flag.
//
// # Daily Throttle
//
// [DailyThrottled] checks whether a marker file was
// touched today. Hooks use this to limit certain
// expensive checks to once per day. The function
// compares the marker file's modification date to
// today's date.
//
// # Wrap-Up Recency
//
// [WrappedUpRecently] checks whether the wrap-up
// marker exists and is younger than the configured
// expiry window. When true, ceremony nudges and
// other post-session reminders should be suppressed
// because the user recently performed a wrap-up.
package check

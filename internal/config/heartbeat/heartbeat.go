//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package heartbeat

// Heartbeat state file prefixes.
const (
	// CounterPrefix is the state file prefix for per-session
	// heartbeat prompt counters.
	CounterPrefix = "heartbeat-"
	// MtimePrefix is the state file prefix for per-session
	// heartbeat context mtime tracking.
	MtimePrefix = "heartbeat-mtime-"
	// LogFile is the log filename for heartbeat events.
	LogFile = "heartbeat.log"
)

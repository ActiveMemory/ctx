//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package heartbeat

// Heartbeat state file prefixes.
const (
	// HeartbeatCounterPrefix is the state file prefix for per-session
	// heartbeat prompt counters.
	HeartbeatCounterPrefix = "heartbeat-"
	// HeartbeatMtimePrefix is the state file prefix for per-session
	// heartbeat context mtime tracking.
	HeartbeatMtimePrefix = "heartbeat-mtime-"
	// HeartbeatLogFile is the log filename for heartbeat events.
	HeartbeatLogFile = "heartbeat.log"
)

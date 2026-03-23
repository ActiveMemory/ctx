//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

// Event log constants for .context/state/ directory.
const (
	// FileEventLog is the current event log file.
	FileEventLog = "events.jsonl"
	// FileEventLogPrev is the rotated (previous) event log file.
	FileEventLogPrev = "events.1.jsonl"
	// LogMaxBytes is the size threshold for log rotation (1MB).
	LogMaxBytes = 1 << 20
	// HookLogMaxBytes is the size threshold for hook log rotation (1MB).
	HookLogMaxBytes = 1 << 20
	// RotationSuffix is the suffix appended to log files during rotation.
	RotationSuffix = ".1"
	// DefaultEventsLast is the default number of events shown by ctx system events.
	DefaultEventsLast = 50
)

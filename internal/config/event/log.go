//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

// Event log constants for .context/state/ directory.
const (
	// FileLog is the current event log file.
	FileLog = "events.jsonl"
	// FileLogPrev is the rotated (previous) event log file.
	FileLogPrev = "events.1.jsonl"
	// LogMaxBytes is the size threshold for log rotation (1MB).
	LogMaxBytes = 1 << 20
	// HookLogMaxBytes is the size threshold for hook log rotation (1MB).
	HookLogMaxBytes = 1 << 20
	// RotationSuffix is the suffix appended to log files during rotation.
	RotationSuffix = ".1"
	// DefaultLast is the default number of events shown by ctx system events.
	DefaultLast = 50
)

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package event defines constants for the ctx event
// logging subsystem: event types, log file names,
// rotation thresholds, display limits, and context-size
// event names.
//
// ctx records structured events in a JSONL log inside
// .context/state/. The constants here govern what gets
// written, how the log rotates, and how events are
// rendered in the CLI.
//
// # Event Types
//
// Session lifecycle events:
//
//   - TypeStart -- emitted when a session begins
//   - TypeEnd -- emitted when a session ends
//
// Notification events:
//
//   - TypeTest -- used for test/ping notifications
//   - TestMessage -- the payload for test events
//
// # Categories and Variables
//
//   - CategorySession -- groups session lifecycle
//     events for filtering
//   - VarCaller -- template variable key for the
//     calling editor in session events
//
// # Log Files and Rotation
//
//   - FileLog ("events.jsonl") -- the active log file
//   - FileLogPrev ("events.1.jsonl") -- the rotated
//     previous log
//   - LogMaxBytes (1 MB) -- size threshold that
//     triggers rotation
//   - HookLogMaxBytes (1 MB) -- size threshold for
//     hook-specific log rotation
//   - RotationSuffix (".1") -- suffix appended during
//     rotation
//   - DefaultLast (50) -- default event count shown
//     by ctx hook event
//
// # Display Formatting
//
//   - MessageMaxLen (60) -- max characters before an
//     event message is truncated
//   - HookFallback ("-") -- placeholder when no hook
//     name can be determined
//   - TruncationSuffix ("...") -- appended to
//     truncated messages
//
// # Context-Size Events
//
// These event names track context window behavior:
//
//   - Suppressed -- prompt was suppressed
//   - Silent -- prompt produced no action
//   - Checkpoint -- context checkpoint emitted
//   - WindowWarning -- context window nearing limit
//
// # Why Centralized
//
// Event constants are shared between the event logger,
// hook runner, CLI display code, and doctor checks.
// Centralizing them prevents string drift and import
// cycles across these subsystems.
package event

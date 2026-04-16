//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package heartbeat implements the hidden
// "ctx system heartbeat" cobra subcommand.
//
// The heartbeat hook fires on every prompt cycle and
// silently collects session telemetry. It never writes
// to stdout, so the agent is unaware of it.
//
// # Behavior
//
// On each invocation the hook:
//
//   - Increments a per-session prompt counter stored
//     in a file under the temporary state directory.
//   - Compares the latest modification time of files
//     in the context directory against the previously
//     recorded mtime to detect context changes.
//   - Reads token usage and context-window size for the
//     current session, computing a usage percentage.
//   - Sends a notification (via the heartbeat channel)
//     with prompt count, context-modified flag, and
//     optional token usage percentage.
//   - Appends an event log entry with the same data.
//   - Writes a timestamped line to the heartbeat log
//     file inside the context directory.
//
// The hook is skipped entirely when the context
// directory is not initialized or the session is
// paused.
//
// # Flags
//
// None. The command reads hook JSON from stdin.
//
// # Output
//
// No stdout output. All data flows to notification
// channels, event log, and the heartbeat log file.
//
// # Delegation
//
// Prompt counting uses system/core/counter. Mtime
// tracking uses system/core/heartbeat. Token info
// is read via system/core/session. Notifications
// go through the notify package.
package heartbeat

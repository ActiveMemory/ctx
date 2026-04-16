//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pause implements the hidden
// "ctx system pause" cobra subcommand.
//
// This plumbing command creates a pause marker file
// that suppresses all subsequent hook invocations for
// a given session. Hooks check for the marker at the
// start of their preamble and return immediately when
// it is present. Use "ctx system resume" to remove
// the marker and re-enable hooks.
//
// # Usage
//
//	ctx system pause [--session-id <id>]
//
// # Flags
//
//	--session-id   Explicit session identifier.
//	               When omitted the session ID is
//	               read from the hook JSON on stdin.
//	               Falls back to "unknown" if neither
//	               source provides a value.
//
// # Behavior
//
// The command resolves the session ID from the flag,
// stdin JSON, or a fallback, then writes a counter
// file at the pause marker path. All hooks that call
// nudge.Paused see the marker and skip execution.
//
// # Output
//
// Prints a confirmation line with the session ID
// that was paused.
//
// # Delegation
//
// Marker path resolution uses system/core/nudge.
// Counter file writing uses system/core/counter.
// Output formatting uses write/pause.
package pause

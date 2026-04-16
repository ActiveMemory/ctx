//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resume implements the hidden
// "ctx system resume" cobra subcommand.
//
// This plumbing command removes the pause marker file
// created by "ctx system pause", re-enabling all hook
// invocations for the given session.
//
// # Usage
//
//	ctx system resume [--session-id <id>]
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
// stdin JSON, or a fallback, then deletes the pause
// marker file. If the file does not exist the removal
// error is logged as a warning but the command still
// succeeds.
//
// After removal, all hooks that check nudge.Paused
// will proceed normally again.
//
// # Output
//
// Prints a confirmation line with the session ID
// that was resumed.
//
// # Delegation
//
// Marker path resolution uses system/core/nudge.
// Session ID reading uses system/core/session.
// Output formatting uses write/session.
package resume

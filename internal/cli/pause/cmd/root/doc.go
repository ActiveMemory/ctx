//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx hook pause" command
// for suppressing context nudges during a session.
//
// # Behavior
//
// The command pauses context hook nudges for the
// specified session. Once paused, hooks will not emit
// reminders, task suggestions, or other advisory
// output until the session is explicitly resumed.
//
// The session ID can be provided via the --session-id
// flag or read from stdin. When neither is provided,
// stdin is read to obtain the session identifier.
//
// Pausing is idempotent: calling pause on an already-
// paused session has no effect beyond confirming the
// state.
//
// # Flags
//
//	--session-id <id>   Session ID to pause. When
//	                     omitted, reads from stdin.
//
// # Output
//
// Prints a one-line confirmation showing the session
// ID that was paused. The command always succeeds
// (returns nil error).
//
// # Delegation
//
// Session pause state is managed by [nudge.Pause]
// in the system/core/nudge package. Session ID
// reading from stdin uses [coreSession.ReadID].
// Output is routed through [session.Paused].
package root

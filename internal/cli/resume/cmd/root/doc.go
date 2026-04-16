//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx hook resume"
// command for re-enabling context nudges after a
// pause.
//
// # Behavior
//
// The command resumes context hook nudges for the
// specified session. After resuming, hooks will
// again emit reminders, task suggestions, and other
// advisory output that was suppressed by a prior
// "ctx hook pause" invocation.
//
// The session ID can be provided via the --session-id
// flag or read from stdin. When neither is provided,
// stdin is read to obtain the session identifier.
//
// Resuming is idempotent: calling resume on a session
// that is not paused has no effect beyond confirming
// the state.
//
// # Flags
//
//	--session-id <id>   Session ID to resume. When
//	                     omitted, reads from stdin.
//
// # Output
//
// Prints a one-line confirmation showing the session
// ID that was resumed. The command always succeeds
// (returns nil error).
//
// # Delegation
//
// Session resume state is managed by [nudge.Resume]
// in the system/core/nudge package. Session ID
// reading from stdin uses [coreSession.ReadID].
// Output is routed through [session.Resumed].
package root

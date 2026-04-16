//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pause implements the "ctx hook pause" command
// that suppresses all context hooks for the current
// session.
//
// When paused, system hooks (check-context-size,
// check-persistence, check-ceremony, and others) stop
// firing, eliminating nudge messages from the AI
// session. This is useful during focused work where
// hook interruptions are unwanted or during debugging
// when hook side effects should be isolated.
//
// Pause state is stored as a per-session flag in the
// .context/state/ directory. It persists until the
// session ends or the user explicitly resumes hooks
// via ctx hook resume.
//
// # Subpackages
//
//	cmd/root -- cobra command definition and pause
//	  state management
//
// [Cmd] returns the cobra command that writes the pause
// flag to the session state directory and prints a
// confirmation message.
package pause

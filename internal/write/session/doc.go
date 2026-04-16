//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package session provides terminal output for session
// lifecycle commands.
//
// Sessions are the fundamental unit of agent
// interaction in ctx. Each session has a start event,
// optional pause/resume cycles, and an end event.
// The output functions confirm each lifecycle
// transition.
//
// # Lifecycle Events
//
// [Event] confirms a session start or end event was
// recorded, printing the event type and the calling
// editor identifier (e.g. "vscode", "claude").
//
// # Hook Control
//
// [Paused] confirms that hooks were suspended for
// the named session. [Resumed] confirms hooks were
// re-enabled. Both print the session ID so the
// user can verify which session was affected.
//
// # Wrap-Up
//
// [WrappedUp] confirms the end-of-session
// persistence ceremony completed. This is the
// final output before context files are committed.
package session

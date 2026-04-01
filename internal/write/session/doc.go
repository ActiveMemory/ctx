//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package session provides terminal output for session lifecycle
// commands (ctx pause, ctx resume, ctx wrap-up, ctx system session-event).
//
// [Event] confirms a session start or end event was recorded.
// [Paused] confirms hooks were suspended for the session.
// [Resumed] confirms hooks were re-enabled. [WrappedUp] confirms
// the end-of-session persistence ceremony completed.
package session

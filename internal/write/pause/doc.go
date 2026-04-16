//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pause provides terminal output for the context
// pause command (ctx pause).
//
// When a user pauses context hooks for the current
// session, the CLI confirms the action through this
// package. The pause command suspends all hook-based
// nudges and context injections until the session is
// explicitly resumed.
//
// # Output
//
// [Confirmed] prints a confirmation message that
// includes the session ID whose hooks were paused.
// It accepts a *cobra.Command for output routing
// and the session identifier string.
//
// A nil *cobra.Command is treated as a no-op so
// callers do not need nil guards.
package pause

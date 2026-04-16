//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package provenance resolves session and git identity
// for hook relay output. The check-reminders hook uses
// these helpers to inject provenance lines into the
// user-facing nudge box, giving the agent visibility
// into which session and git state produced the output.
//
// # Session ID
//
// [ShortSessionID] truncates a full session UUID to
// ShortIDLen characters for compact display. Returns
// "unknown" when the input is empty.
//
// # Provenance Emission
//
// [Emit] prints a provenance line to stdout containing
// the short session ID, current git branch, short commit
// hash, and a context-free suffix. When session token
// data is available, it appends a "Context: N% free"
// indicator so the agent sees how much context window
// remains at the start of each prompt.
//
// # Context Window Tracking
//
// [ContextFreePct] reads the token info for the given
// session and computes the percentage of the model's
// context window that is still free. Returns 0 when no
// token data is available (first prompt, new session,
// or read error), which the caller treats as "no suffix
// to render".
//
// # Utility
//
// [DefaultVal] returns the given value or "unknown" when
// empty. Used to normalize git branch and commit fields.
package provenance

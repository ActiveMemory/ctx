//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package wrap defines marker constants for the
// end-of-session wrap-up suppression mechanism.
//
// When a user completes the ctx wrap-up ceremony
// (persisting decisions, learnings, and tasks), ctx
// writes a short-lived marker file to the state
// directory. While the marker exists, the wrap-up
// nudge is suppressed so the user is not prompted
// again within the same working session.
//
// # Key Constants
//
//   - [Marker] ("ctx-wrapped-up") — the state file
//     name written after a successful wrap-up.
//   - [Content] ("wrapped-up") — the content stored
//     in the marker file.
//   - [ExpiryHours] (2) — hours the marker
//     suppresses further nudges. After expiry the
//     nudge engine considers the session stale and
//     may prompt again.
//
// # How It Works
//
// The wrap-up command writes a file named [Marker]
// containing [Content] to the state directory. The
// nudge engine checks the file's modification time;
// if it is younger than [ExpiryHours] the nudge is
// skipped.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package wrap

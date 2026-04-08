//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package normalize implements the "ctx remind normalize"
// subcommand for compacting reminder IDs.
//
// Reminders use stable auto-incrementing IDs that persist
// across dismissals. When reminders are dismissed, their
// IDs leave gaps (e.g., 1, 4, 7). Normalize closes those
// gaps by reassigning IDs sequentially starting from 1,
// preserving the current order.
//
// This is a deliberate user action, not automatic: it
// invalidates any IDs that were previously displayed in
// hook relays or session transcripts. The user should
// only run this when the gap cosmetics bother them,
// not as routine maintenance.
//
// The command reads reminders from the JSON store,
// reassigns IDs in-place, writes back, and prints a
// confirmation with the count of reminders renumbered.
package normalize

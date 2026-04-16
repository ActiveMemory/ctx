//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package normalize implements the "ctx remind
// normalize" subcommand for compacting reminder IDs.
//
// # Behavior
//
// Reminders use stable auto-incrementing IDs that
// persist across dismissals. When reminders are
// dismissed, their IDs leave gaps (e.g., 1, 4, 7).
// Normalize closes those gaps by reassigning IDs
// sequentially starting from 1, preserving the
// current order.
//
// This is a deliberate user action, not automatic:
// it invalidates any IDs that were previously
// displayed in hook relays or session transcripts.
// The user should only run this when the gap
// cosmetics bother them, not as routine
// maintenance.
//
// When no reminders exist, the command prints a
// notice and exits without writing.
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// On success, prints a confirmation showing the
// count of reminders renumbered. When the store
// is empty, prints a "no reminders" notice instead.
//
// # Delegation
//
// The command reads reminders via [store.Read],
// reassigns IDs in-place, and writes back via
// [store.Write]. Output is routed through the
// [remind] write package.
package normalize

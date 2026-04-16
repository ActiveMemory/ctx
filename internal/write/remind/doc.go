//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package remind provides terminal output for the
// session reminder commands (ctx remind add, list,
// dismiss).
//
// Reminders are short messages attached to a session
// with an optional date gate. The output functions
// cover the full CRUD lifecycle.
//
// # Adding
//
// [Added] prints a confirmation that includes the
// reminder ID, message text, and optional "after"
// date suffix when a date gate is set.
//
// # Listing
//
// [Item] renders a single reminder with its ID,
// message, and a "not yet due" annotation when the
// gate date is in the future relative to today.
// [None] handles the empty-list case.
//
// # Dismissing
//
// [Dismissed] confirms removal of a single reminder
// by ID and message. [DismissedAll] reports bulk
// dismissal with a count of removed items.
//
// # Maintenance
//
// [Normalized] confirms that reminder IDs were
// renumbered sequentially after gaps formed from
// deletions.
package remind

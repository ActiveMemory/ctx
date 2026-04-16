//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package add implements the "ctx remind add"
// subcommand for creating new reminders.
//
// # Behavior
//
// The command accepts exactly one positional argument:
// the reminder message text. It assigns the next
// available stable ID, timestamps the creation in
// UTC, and appends the reminder to the JSON store.
//
// An optional --after flag sets a date gate in
// YYYY-MM-DD format. When set, the reminder will
// only be surfaced by hooks after that date has
// arrived. The date is validated at creation time;
// invalid formats produce an immediate error.
//
// The Run function is exported so the parent remind
// command can reuse it as a default action.
//
// # Flags
//
//	--after, -a <date>   Date gate in YYYY-MM-DD
//	                      format. The reminder is
//	                      hidden until this date.
//
// # Output
//
// On success, prints a confirmation showing the
// assigned reminder ID, message text, and the after
// date (if set). On failure, returns an error for
// invalid dates or read/write problems.
//
// # Delegation
//
// Reminder persistence is handled by the core/store
// package via [store.Read] and [store.Write]. ID
// assignment uses [store.NextID]. Output is routed
// through [remind.Added].
package add

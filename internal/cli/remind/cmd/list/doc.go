//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package list implements the "ctx remind list"
// subcommand for displaying all pending reminders.
//
// # Behavior
//
// The command reads all reminders from the JSON
// store and prints each one with its stable ID,
// message text, and date annotations. Reminders
// with an "after" date gate are annotated to show
// whether the gate date has passed relative to
// today.
//
// When no reminders exist, the command prints a
// notice and exits. The Run function is exported
// so the parent remind command can reuse it as
// the default action when invoked without a
// subcommand.
//
// The command is aliased for shorter invocation.
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// Prints one line per reminder showing the ID,
// message, and date status. Date-gated reminders
// show whether they are active (gate date passed)
// or pending (gate date in the future). When the
// store is empty, prints a "no reminders" notice.
//
// # Delegation
//
// Reminder reading is handled by [store.Read].
// Per-item formatting and date comparison use
// [remind.Item] with today's date. Empty-state
// output goes through [remind.None].
package list

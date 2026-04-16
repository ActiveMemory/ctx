//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reminder defines file and template constants
// for the session reminder subsystem, which lets users
// set time-based or turn-based reminders that fire
// during an active Claude Code session.
//
// # How Reminders Work
//
// Users create reminders via "ctx remind". Each
// reminder is stored in File ("reminders.json") inside
// .context/. During the session, hooks check the
// reminder list and inject any due reminders into the
// agent's next system prompt.
//
// # Template Variables
//
// VarList ("ReminderList") is the key injected into
// hook message templates. It contains the formatted
// list of due reminders so the nudge message can
// display them to the agent.
//
// # Key Constants
//
//   - File ("reminders.json") -- the session-scoped
//     reminder storage file in .context/.
//   - VarList ("ReminderList") -- template variable
//     for the formatted reminder list injected into
//     hook messages.
//
// # Why Centralize
//
// The reminder file name and template variable key
// are referenced by both the CLI command handler and
// the hook executor. Defining them here prevents
// drift between the two sides.
package reminder

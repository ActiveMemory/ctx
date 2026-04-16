//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_reminder implements the
// **`ctx system check-reminder`** hidden hook, which
// surfaces pending reminders at session start so the
// agent can act on deferred tasks.
//
// # What It Does
//
// The hook loads all stored reminders from the reminder
// store, filters to those whose "after" date is today
// or earlier, and emits a nudge box listing each due
// reminder with its ID and message. A dismiss hint is
// appended so the agent knows how to clear reminders
// after acting on them.
//
// Before checking reminders, the hook unconditionally
// emits session provenance (session ID, branch, commit)
// so the agent always has orientation context even when
// no reminders are due.
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// Always: provenance information (session, branch,
// commit). On due reminders: a nudge box listing each
// reminder ID and message with dismiss instructions.
// On no due reminders: no additional output.
//
// # Throttling
//
// No daily throttle; reminders are checked on every
// session start. The provenance output is also
// unconditional and fires even when hooks are paused.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble], emits provenance
// through [core/provenance.Emit], loads reminders from
// [remind/core/store.Read], filters by date, and emits
// the nudge via [core/nudge.LoadAndEmit].
package check_reminder

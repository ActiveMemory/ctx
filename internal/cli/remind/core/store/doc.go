//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package store manages the reminder persistence layer,
// reading and writing the reminders.json file in the context
// directory.
//
// # Storage Format
//
// Reminders are stored as a JSON array of [Reminder] objects
// in .context/reminders.json. Each reminder carries an
// auto-incremented ID, a message string, an ISO 8601
// creation timestamp, and an optional date gate (After)
// that defers the reminder until a specific date.
//
// # CRUD Operations
//
//   - [Read] loads all reminders from the JSON file.
//     Returns nil (not an error) when the file does not
//     exist, treating a missing file as an empty store.
//   - [Write] persists the full reminder slice as
//     indented JSON using atomic file writes.
//   - [NextID] scans existing reminders and returns
//     max(ID) + 1 for the next insertion.
//   - [Path] resolves the absolute path to the store
//     file using [rc.ContextDir] as the base.
//
// # File Safety
//
// Reads use [io.SafeReadUserFile] for path validation.
// Writes use [io.SafeWriteFile] for atomic persistence.
// Both prevent path traversal and partial-write corruption.
package store

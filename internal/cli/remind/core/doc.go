//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides business logic for session-scoped
// reminders that persist across agent conversations.
//
// Reminders are stored as JSON in the context directory
// and can be created, listed, and dismissed through the
// CLI. Each reminder has an auto-incremented ID, a
// message, a creation timestamp, and an optional trigger
// date.
//
// # Subpackages
//
// The core package delegates to two subpackages:
//
// The store subpackage handles persistence. [store.Read]
// loads reminders from the JSON file, returning nil when
// the file is absent. [store.Write] serializes the full
// reminder slice back to disk. [store.NextID] scans
// existing reminders to compute the next sequential ID.
// [store.Path] returns the absolute path to the
// reminders JSON file.
//
// The dismiss subpackage handles removal. [dismiss.Many]
// removes one or more reminders by ID, validating all
// IDs before any deletion to avoid partial removal.
// [dismiss.All] clears every active reminder.
//
// # Data Flow
//
// The cmd/ layer parses flags and arguments, then calls
// store or dismiss functions. Output is delegated to
// the write/remind package for formatted messages.
package core

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package lock implements the "ctx journal lock" command.
//
// # Overview
//
// The lock command protects journal entries from being
// overwritten during export regeneration. Once an entry
// is locked, "ctx journal export --regenerate" skips it
// regardless of other flags.
//
// This is useful when an entry has been manually enriched
// with metadata, tags, or summary edits that should not
// be lost during a bulk re-export.
//
// # Flags
//
//	--all    Lock every journal entry in the journal
//	         directory. Without this flag, one or more
//	         filename patterns must be provided as
//	         positional arguments.
//
// # Arguments
//
// Positional arguments are glob patterns matched against
// journal filenames. At least one pattern is required
// unless --all is set.
//
// # Behavior
//
// [Cmd] builds the cobra.Command and registers the --all
// flag. [Run] delegates to the shared lock/unlock core
// in journal/core/lock with lock=true.
//
// The core logic loads .state.json from the journal
// directory, marks matched entries as locked, and
// persists the updated state file.
//
// # Output
//
// Each locked entry is confirmed on stdout with its
// filename. A summary line reports the total number
// of entries locked.
package lock

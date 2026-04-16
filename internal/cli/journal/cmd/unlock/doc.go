//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package unlock implements the "ctx journal unlock"
// command.
//
// # Overview
//
// The unlock command removes lock protection from
// journal entries, allowing "ctx journal export
// --regenerate" to overwrite them again. This reverses
// the effect of "ctx journal lock".
//
// # Flags
//
//	--all    Unlock every journal entry in the journal
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
// in journal/core/lock with lock=false.
//
// The core logic loads .state.json from the journal
// directory, clears the locked mark on matched entries,
// and persists the updated state file.
//
// # Output
//
// Each unlocked entry is confirmed on stdout with its
// filename. A summary line reports the total number
// of entries unlocked.
package unlock

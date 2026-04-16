//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cooldown prevents redundant context emissions
// within a single agent session.
//
// When the agent command outputs a context packet, it
// writes a tombstone file to mark the emission time. On
// subsequent invocations within the same session, the
// cooldown check skips output if the tombstone is still
// fresh, avoiding noise in short-lived tool loops.
//
// # Active
//
// [Active] checks whether the cooldown tombstone for a
// given session identifier exists and was modified within
// the specified duration. It returns false when the
// session string is empty, the duration is non-positive,
// or the tombstone file is missing or stale.
//
// # TouchTombstone
//
// [TouchTombstone] creates or updates the tombstone file
// for a session. It writes an empty file with restricted
// permissions (fs.PermSecret) to the state directory
// under .context/state/. Write failures are logged as
// warnings but do not propagate errors, so the agent
// command never fails due to cooldown bookkeeping.
//
// # TombstonePath
//
// [TombstonePath] returns the absolute filesystem path
// for a session's tombstone. The file is placed in the
// .context/state/ directory with a prefix defined by
// agent.TombstonePrefix. The state directory is created
// on demand with restricted permissions.
//
// # Data Flow
//
// The agent command's Run function checks Active before
// assembling a context packet. If Active returns true,
// Run exits early. Otherwise it builds the packet,
// emits it, and calls TouchTombstone to start the
// cooldown window.
package cooldown

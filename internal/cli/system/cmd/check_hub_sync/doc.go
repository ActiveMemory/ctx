//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_hub_sync implements the
// **`ctx system check-hub-sync`** hidden hook, which
// pulls new entries from a registered ctx Hub at
// session start.
//
// # What It Does
//
// When a hub connection is configured, the hook syncs
// new entries from the remote hub to the local
// .context/hub/ directory. If no hub is configured or
// there are no new entries, the hook is silent. When
// new entries are pulled, a nudge message is emitted
// to inform the agent of the incoming content.
//
// This is internal plumbing invoked by Claude Code
// hooks; not intended for manual use.
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// On new entries: a nudge message describing what was
// synced. On no hub, no new entries, or throttled:
// no output.
//
// # Throttling
//
// The hook is throttled to fire at most once per day
// using a marker file in the state directory.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble], checks for hub
// connectivity through [core/hubsync.Connected], and
// delegates the actual sync to [core/hubsync.Sync].
// The daily throttle marker is managed via
// [io.TouchFile].
package check_hub_sync

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync implements hub-to-local entry
// synchronisation for the ctx connection sync command.
//
// # Overview
//
// This package pulls new journal entries from a remote
// hub and writes them as markdown files into the local
// .context/hub/ directory. It tracks progress with a
// sequence-based sync state so only new entries are
// fetched on each invocation.
//
// # Behavior
//
// [Run] acquires a file lock, dials the hub over gRPC,
// fetches entries newer than the last-seen sequence, and
// writes them as markdown into .context/hub/.
//
// # Data Flow
//
// When [Run] is called it performs these steps:
//
//  1. Loads connection config (hub address, token,
//     subscribed types) via the config sub-package.
//  2. Acquires a file-based lock to prevent concurrent
//     syncs from colliding.
//  3. Reads the persisted sync state to obtain the
//     last-seen sequence number.
//  4. Dials the hub via gRPC, requesting all entries
//     after the last sequence for the subscribed types.
//  5. Renders received entries as markdown through the
//     render sub-package.
//  6. Updates the sync state with the highest sequence
//     number from the batch.
//  7. Reports the count of synced entries to the user.
//
// # Lock File
//
// A lock file at .context/hub/.sync.lock prevents two
// sync processes from running at the same time. The
// lock is released via a deferred cleanup function
// returned by loadState.
//
// # Internal Types
//
//   - state -- tracks the LastSequence field, persisted
//     as JSON in .context/hub/.sync_state.json.
package sync

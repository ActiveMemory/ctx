//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync implements the "ctx connection sync"
// subcommand that pulls new entries from a connected
// ctx Hub into the local .context/hub/ directory.
//
// # What It Does
//
// Loads the connection config and sync state, connects
// to the hub, pulls entries added since the last sync
// sequence number, renders them as markdown files in
// .context/hub/, and updates the sync state so the
// next run only fetches new entries.
//
// # Flags
//
// None. The command accepts no arguments. Connection
// settings are read from .context/.connect.enc and
// sync state is tracked in a lock-guarded state file.
//
// # Output
//
// Prints a summary line showing how many entries were
// synced (e.g. "Synced 3 entries."). When there are
// no new entries it prints "Synced 0 entries."
//
// # Delegation
//
// [Cmd] builds the cobra.Command and delegates to
// [coreSync.Run] which handles config loading, state
// locking, gRPC client setup, entry pulling, markdown
// rendering, and state persistence.
package sync

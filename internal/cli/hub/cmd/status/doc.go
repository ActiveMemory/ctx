//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package status implements the "ctx hub status"
// subcommand that displays cluster-level information
// about a running ctx Hub.
//
// # What It Does
//
// Loads the connection config, connects to the hub,
// calls the Status RPC, and prints a cluster summary
// including the node role, hub address, total stored
// entries, and number of registered projects.
//
// # Flags
//
// None. The command accepts no arguments. Connection
// settings are read from .context/.connect.enc.
//
// # Output
//
// A human-readable status block including:
//
//   - Node role (active or follower)
//   - Hub address (host:port)
//   - Total entries stored
//   - Number of registered projects
//
// # Delegation
//
// [Cmd] builds the cobra.Command and delegates
// directly to [coreStatus.Run] which handles config
// loading, gRPC client setup, the status call, role
// determination, and output formatting.
package status

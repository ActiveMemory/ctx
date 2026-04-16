//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package status implements the "ctx connection status"
// subcommand that displays the current hub connection
// state and entry statistics.
//
// # What It Does
//
// Loads the encrypted connection config, connects to
// the hub, calls the Status RPC, and prints a summary
// showing the hub address, total stored entries, and
// number of connected clients.
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
//   - Hub address (host:port)
//   - Total entries stored in the hub
//   - Number of currently connected clients
//
// # Delegation
//
// [Cmd] builds the cobra.Command and delegates
// directly to [coreStatus.Run] which handles config
// loading, gRPC client setup, the status call, and
// output formatting.
package status

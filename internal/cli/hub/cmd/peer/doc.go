//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package peer implements the "ctx hub peer" subcommand
// for managing peer nodes in a ctx Hub cluster.
//
// # What It Does
//
// Adds or removes a peer address from the hub's
// cluster membership. This is used when scaling a
// hub deployment across multiple nodes for Raft-based
// leader election and replication.
//
// # Arguments
//
// Requires exactly two positional arguments:
//
//   - args[0]: action -- "add" or "remove"
//   - args[1]: address -- peer gRPC address
//     (host:port)
//
// # Flags
//
// None.
//
// # Output
//
// Prints a confirmation line indicating the peer
// was added or removed, e.g. "Peer added: host:9090"
// or "Peer removed: host:9090".
//
// # Delegation
//
// [Cmd] builds the cobra.Command and delegates
// directly to [corePeer.Run] which validates the
// action string and writes the confirmation message.
package peer

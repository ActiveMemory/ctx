//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package publish implements entry publishing to the hub
// for the ctx connection publish command.
//
// # Run
//
// [Run] sends local entries to the ctx Hub via the
// Publish RPC. It accepts a slice of hub.PublishEntry
// values prepared by the cmd/ layer.
//
// The execution flow is:
//
//  1. Load the encrypted connection config via
//     connectCfg.Load to obtain the hub address and
//     bearer token.
//  2. Dial the hub with hub.NewClient, establishing a
//     gRPC connection.
//  3. Call client.Publish with the entries, sending
//     them in a single batch RPC.
//  4. Print a confirmation showing the number of
//     published entries via writeConnect.Published.
//
// The function returns an error if config loading,
// connection setup, or the publish RPC fails. The gRPC
// connection is closed via a deferred Close call.
//
// # Data Flow
//
// The cmd/ layer prepares PublishEntry values from
// command arguments and passes them to Run. Run handles
// config, networking, and output. Future versions may
// support reading entries from local context files with
// a --new flag.
package publish

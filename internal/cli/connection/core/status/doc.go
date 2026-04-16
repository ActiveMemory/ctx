//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package status implements hub status display for the
// ctx connection status command.
//
// # Run
//
// [Run] queries the ctx Hub for connection health and
// entry statistics, then prints a summary to the
// terminal.
//
// The execution flow is:
//
//  1. Load the encrypted connection config via
//     connectCfg.Load to obtain the hub address and
//     bearer token.
//  2. Dial the hub with hub.NewClient, establishing a
//     gRPC connection.
//  3. Call client.Status to retrieve server-side
//     statistics including total entry count and the
//     number of connected clients.
//  4. Print the hub address, total entries, and
//     connected client count via writeConnect.Status.
//
// The function returns an error if config loading,
// connection setup, or the status RPC fails. The gRPC
// connection is closed via a deferred Close call.
//
// # Data Flow
//
// The cmd/ layer calls Run as the cobra RunE function.
// Run handles config loading, networking, and output
// delegation. The write/connect package formats the
// final user-facing output.
package status

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package listen implements real-time hub entry streaming
// for the ctx connection listen command.
//
// # Run
//
// [Run] opens a persistent gRPC stream to the ctx Hub
// and writes each received entry to .context/hub/ as it
// arrives. The function blocks until the user presses
// Ctrl-C or an unrecoverable error occurs.
//
// The execution flow is:
//
//  1. Load the encrypted connection config via
//     connectCfg.Load to obtain the hub address and
//     bearer token.
//  2. Dial the hub with hub.NewClient, establishing a
//     gRPC connection.
//  3. Set up a signal handler for os.Interrupt using
//     signal.NotifyContext so Ctrl-C cancels the
//     stream context.
//  4. Print a "listening" status message via
//     writeConnect.Listening.
//  5. Call client.Listen with the configured entry type
//     filters and a callback that writes each received
//     EntryMsg to disk via render.WriteEntries, then
//     prints a confirmation via writeConnect.EntryReceived.
//  6. On context cancellation (Ctrl-C), return nil.
//     On any other error, propagate it to the cmd/ layer.
//
// # Data Flow
//
// The cmd/ layer calls Run as the cobra RunE function.
// Run loads config, dials the hub, and streams entries.
// The write/ layer handles all user-facing output. The
// render subpackage persists entries to .context/hub/.
package listen

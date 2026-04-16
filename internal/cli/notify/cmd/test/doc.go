//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package test implements the "ctx hook notify test"
// command.
//
// # Overview
//
// The test command sends a test notification to the
// configured webhook URL to verify that the notification
// pipeline is working end-to-end. It loads the saved
// webhook, sends a test payload, and reports the HTTP
// status code and success/failure result.
//
// # Flags
//
// This command accepts no flags.
//
// # Behavior
//
// [Cmd] builds a simple cobra.Command with no flags.
// [Run] delegates to coreTest.Send which loads the
// encrypted webhook URL, sends a test HTTP request,
// and returns the result. Then it dispatches to the
// appropriate output based on the result:
//
//   - If no webhook is configured, prints a "no
//     webhook" message.
//   - If the notification was filtered (e.g. by rate
//     limiting), prints a "filtered" notice.
//   - In all cases where a request was made, prints
//     the HTTP status code and whether it was
//     successful.
//
// # Output
//
// Prints the HTTP status code and a pass/fail
// indicator, along with the encryption method used.
// If no webhook is configured, prints a message
// directing the user to run "ctx hook notify setup".
package test

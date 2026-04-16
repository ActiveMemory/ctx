//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package setup implements the "ctx hook notify setup"
// command.
//
// # Overview
//
// The setup command configures webhook-based
// notifications for ctx hooks. It prompts the user to
// enter a webhook URL via stdin, validates the input,
// and saves the URL in encrypted form for use by
// subsequent hook executions.
//
// # Flags
//
// This command accepts no flags.
//
// # Behavior
//
// [Cmd] builds a simple cobra.Command with no flags.
// [Run] is exported for testability and accepts an
// *os.File for stdin injection. It performs these steps:
//
//  1. Prints a prompt asking for the webhook URL.
//  2. Reads one line from stdin.
//  3. Validates that the input is non-empty.
//  4. Saves the webhook URL via iNotify.SaveWebhook,
//     which encrypts and persists it.
//  5. Prints a confirmation with the masked URL and
//     the encryption method used.
//
// If stdin is empty or the URL is blank, the command
// returns an appropriate error.
//
// # Output
//
// Prints a setup prompt, then a confirmation showing
// the masked webhook URL (only the last few characters
// visible) and the encryption algorithm used.
package setup

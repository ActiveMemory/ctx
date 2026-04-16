//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package subscribe implements the "ctx connection subscribe"
// subcommand that configures which entry types this
// project receives from the hub.
//
// # What It Does
//
// Updates the subscription type list in the encrypted
// connection config (.context/.connect.enc). Subsequent
// listen and sync operations will only receive entries
// matching the subscribed types.
//
// # Arguments
//
// One or more positional arguments specifying entry
// types to subscribe to (e.g. "decision", "learning",
// "task", "convention").
//
// # Flags
//
// None.
//
// # Output
//
// Prints a confirmation line listing the subscribed
// types.
//
// # Delegation
//
// [Cmd] builds the cobra.Command and delegates
// directly to [coreSub.Run] which loads the existing
// config, updates the types list, saves the config
// back, and prints confirmation.
package subscribe

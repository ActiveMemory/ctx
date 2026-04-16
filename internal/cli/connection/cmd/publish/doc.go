//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package publish implements the "ctx connection publish"
// subcommand that sends context entries to a connected
// ctx Hub.
//
// # What It Does
//
// Takes a type and content as positional arguments,
// wraps them into a timestamped publish entry, and
// sends them to the hub via the Publish RPC. This
// is the manual publish path; the --share flag on
// "ctx add" uses the same core logic automatically.
//
// # Arguments
//
// Requires exactly two positional arguments:
//
//   - args[0]: entry type (e.g. "decision",
//     "learning")
//   - args[1]: entry content text
//
// # Flags
//
// None. Connection settings are read from the
// encrypted config at .context/.connect.enc.
//
// # Output
//
// Prints a confirmation line showing how many
// entries were published (always 1 for this
// command).
//
// # Delegation
//
// [Cmd] builds the cobra.Command, constructs a
// [hub.PublishEntry] with the current timestamp,
// and delegates to [corePub.Run] for config loading,
// gRPC client setup, and the publish call.
package publish

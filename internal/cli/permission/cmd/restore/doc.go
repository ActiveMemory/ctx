//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package restore implements the "ctx permission
// restore" subcommand for resetting Claude Code
// permissions to a known-good baseline.
//
// # Behavior
//
// The command reads the golden image (a previously
// saved snapshot of settings.local.json) and replaces
// the current settings.local.json with it. This
// reverts any permission drift that accumulated
// during interactive sessions.
//
// Before overwriting, the command computes a diff
// between the golden and current permission lists
// (both allow and deny), reporting which entries
// were restored and which were dropped. If the
// files already match, it prints a match notice
// and exits.
//
// When settings.local.json does not exist, the
// golden image is copied in directly. When the
// golden image itself is missing, the command
// returns an error instructing the user to run
// "ctx permission snapshot" first.
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// Prints a diff summary showing restored and dropped
// permission entries for both allow and deny lists,
// followed by a completion notice. Special cases
// print match or missing-file notices.
//
// # Delegation
//
// Permission diff logic is provided by
// [diff.StringSlices]. File I/O uses
// [io.SafeReadUserFile] and [io.SafeWriteFile].
// Output is routed through the [restore] write
// package.
package restore

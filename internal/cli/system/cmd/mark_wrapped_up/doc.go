//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mark_wrapped_up implements the hidden
// "ctx system mark-wrapped-up" cobra subcommand.
//
// This command writes a marker file that suppresses
// context-checkpoint nudges for a configured expiry
// period after a wrap-up ceremony completes. Without
// this marker, hooks such as persistence reminders
// and ceremony prompts would continue to fire even
// though the session has already been wrapped up.
//
// # Behavior
//
// When invoked the command:
//
//   - Computes the marker path inside the temporary
//     state directory using the configured wrap marker
//     filename.
//   - Writes a fixed-content marker file with secret
//     permissions (owner-only read/write).
//   - Prints a confirmation message to stdout.
//
// Downstream hooks check for the marker's existence
// and age before emitting nudges. If the marker is
// present and younger than the configured expiry,
// those hooks silently skip.
//
// # Flags
//
// None. The command takes no arguments.
//
// # Output
//
// Prints a single confirmation line indicating the
// wrap-up marker was created.
//
// # Delegation
//
// State directory resolution uses system/core/state.
// File writing delegates to internal/io.SafeWriteFile.
// Output formatting uses write/session.WrappedUp.
package mark_wrapped_up

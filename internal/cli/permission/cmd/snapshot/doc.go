//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package snapshot implements the "ctx permission
// snapshot" subcommand for saving the current Claude
// Code permissions as a golden image.
//
// # Behavior
//
// The command reads the current settings.local.json
// and saves a byte-for-byte copy as the golden image
// file. This golden image serves as the baseline for
// future "ctx permission restore" operations.
//
// If settings.local.json does not exist, the command
// returns an error prompting the user to configure
// permissions first. If a golden image already exists,
// it is overwritten and the output notes that the
// snapshot was updated rather than created.
//
// The snapshot is stored at the path defined by
// [claude.SettingsGolden] and written with standard
// file permissions.
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// On success, prints a confirmation line showing the
// golden image path and whether it was created or
// updated. On failure, returns an error for missing
// settings or write problems.
//
// # Delegation
//
// File I/O uses [io.SafeReadUserFile] and
// [io.SafeWriteFile]. Output is routed through the
// [restore.SnapshotDone] writer. Path constants come
// from the [claude] config package.
package snapshot

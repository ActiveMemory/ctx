//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package normalize implements the "ctx pad normalize"
// subcommand for compacting stable entry IDs.
//
// # Behavior
//
// Pad entries carry stable IDs that persist across
// additions and deletions. Over time, gaps accumulate
// (e.g., 1, 3, 7) as entries are removed. Normalize
// closes those gaps by reassigning IDs sequentially
// starting from 1, preserving the current file order.
//
// This is a deliberate user action, not automatic: it
// invalidates any IDs the user may have noted or that
// appear in prior session transcripts. The command
// should only be run when gap cosmetics matter, not
// as routine maintenance.
//
// When the pad is empty, the command prints an empty
// notice and exits without writing.
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// On success, prints a confirmation line showing the
// count of entries renumbered. When the pad has no
// entries, prints an empty-pad notice instead.
//
// # Delegation
//
// The underlying ID reassignment logic lives in
// [parse.Normalize]; this package handles the CLI
// wiring, file I/O via [store.WriteEntriesWithIDs],
// and user output via [writePad].
package normalize

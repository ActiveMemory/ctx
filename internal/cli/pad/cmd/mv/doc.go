//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mv implements the "ctx pad mv" subcommand
// for reordering scratchpad entries by position.
//
// # Behavior
//
// The command accepts exactly two positional arguments:
// a source position N and a destination position M,
// both 1-based. It extracts the entry at position N,
// removes it from its current slot, and inserts it at
// position M. All other entries shift accordingly.
//
// Both positions are validated against the current
// entry count before any mutation occurs. Invalid
// indices produce an error without modifying the pad.
//
// Note that positions are display-order indices, not
// stable IDs. The command operates on the physical
// ordering of entries in the pad file.
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// On success, prints a one-line confirmation showing
// the source and destination positions. On failure,
// returns an error for out-of-range positions or
// read/write problems.
//
// # Delegation
//
// Index validation is handled by [validate.Index].
// Entry persistence goes through [store.WriteEntries].
// User-facing output is routed through [pad.EntryMoved].
package mv

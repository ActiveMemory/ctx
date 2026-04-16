//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package rm implements the "ctx pad rm" subcommand
// for removing entries from the scratchpad by stable
// ID.
//
// # Behavior
//
// The command accepts one or more positional arguments
// specifying entry IDs to remove. Arguments can be
// individual IDs or ranges (e.g., "3-5" expands to
// 3, 4, 5). All IDs are resolved against the current
// pad before any deletion occurs, preventing shift-
// induced mismatches when removing multiple entries
// in a single invocation.
//
// If any specified ID does not exist in the pad, the
// command returns an error without modifying anything.
// On success, the remaining entries keep their
// original stable IDs (no renumbering occurs).
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// Each successfully removed entry prints a one-line
// confirmation showing the removed ID. On failure,
// returns an error identifying the missing entry ID.
//
// # Delegation
//
// Argument parsing and range expansion are handled by
// [parse.IDs]. ID-to-index resolution uses
// [parse.FindByID]. Persistence goes through
// [store.WriteEntriesWithIDs]. User-facing output
// is routed through [pad.EntryRemoved].
package rm

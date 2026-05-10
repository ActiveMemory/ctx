//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package validate enforces noun-specific body-flag contracts on
// add subcommands. It centralises two pieces of policy:
//
//   - Required flags: noun-level commands declare which body
//     flags must be present. RequireBodyFlags wires
//     [cobra.Command.MarkFlagRequired] for each so cobra rejects
//     missing flags before RunE runs.
//
//   - Placeholder rejection: a closed set of placeholder values
//     (TBD, see chat, n/a, etc., plus whitespace-only) is rejected
//     at PreRunE time with a clear error naming the flag and the
//     offending value. Substring matches are not treated as
//     placeholders so legitimate prose containing the word "TBD"
//     still passes.
//
// The package is internal to the add core; noun-level subcommands
// (decision, learning) call RequireBodyFlags after constructing
// their command via [build.Cmd].
package validate

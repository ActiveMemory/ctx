//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package validate enforces body-flag contracts on add
// subcommands at PreRunE time. It centralises one policy:
//
//   - A closed set of placeholder values (TBD, see chat,
//     n/a, etc., plus whitespace-only) is rejected with a
//     clear error naming the flag and the offending value.
//     Cobra defaults string flags to "", which the same
//     empty-value check rejects — so omitting a required
//     body flag and supplying a placeholder fail through
//     the same code path. Substring matches are not
//     treated as placeholders so legitimate prose
//     containing the word "TBD" still passes.
//
// The package does not call [cobra.Command.MarkFlagRequired]:
// PreRunE is the single enforcement point, and routing
// missing-flag and placeholder-flag through the same path
// avoids both the discarded-error wart on MarkFlagRequired
// and the divergent error messages it would produce.
//
// The package is internal to the add core; noun-level
// subcommands (decision, learning) call [RequireBodyFlags]
// after constructing their command via [build.Cmd].
package validate

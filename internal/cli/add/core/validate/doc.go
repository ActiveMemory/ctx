//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package validate is a pure-function policy layer for add
// subcommand body flags.
//
// It exports two functions:
//
//   - [BodyFlags] reads each named flag from a cobra command and
//     returns an error if any value is empty, whitespace-only,
//     or matches the closed placeholder set (TBD, see chat, n/a,
//     etc.). It does not mutate the command.
//
//   - [RejectPlaceholder] is the per-value primitive used by
//     [BodyFlags] and is exported for tests and ad-hoc reuse.
//
// Cobra defaults string flags to "", so the empty-value check
// catches missing flags through the same code path as
// placeholder rejection. Substring matches are not placeholders
// — legitimate prose containing the word "TBD" still passes.
//
// Noun-level add subcommands (decision, learning) invoke
// [BodyFlags] from their own PreRunE so the wiring is visible
// at the call site. The package does not register hooks on the
// caller's command.
package validate

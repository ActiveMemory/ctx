//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sourcecoverage defines the typed error
// constructors for the source-coverage ledger writer
// ([github.com/ActiveMemory/ctx/internal/write/kb/sourcecoverage]).
//
// # Domain
//
// Two sentinels plus three wrapping I/O constructors cover the
// writer's failure surface:
//
//   - [ErrIllegalTransition]: Advance refused a (from, to)
//     pair the state machine does not allow.
//   - [ErrUnknownSource]: Advance referenced a Source not yet
//     present at a non-entry-point State.
//   - [IllegalTransition], [UnknownSource], [ReadLedger],
//     [MkdirLedgerDir], [WriteLedger] wrap I/O and sentinel-
//     bind failures with operator-friendly context.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb/sourcecoverage]
//     supplies the message + format-string constants.
//   - [github.com/ActiveMemory/ctx/internal/write/kb/sourcecoverage]
//     is the primary caller.
package sourcecoverage

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package evidence defines the typed error constructors for
// the evidence-index writer
// ([github.com/ActiveMemory/ctx/internal/write/kb/evidence]).
//
// # Domain
//
// Two sentinels plus a fistful of wrapping constructors cover
// the writer's full failure surface:
//
//   - [ErrDuplicateID]: Append called with an explicit row.ID
//     already present in the file.
//   - [ErrInvalidBand]: Confidence outside the four canonical
//     bands.
//   - [DuplicateID], [InvalidBand], [ParseEVNumber],
//     [ReadIndex], [MkdirDir], [OpenIndex], [WriteRow] wrap
//     parse, sentinel-bind, and I/O failures with
//     operator-friendly context.
//
// # Wrapping strategy
//
// The constructor functions use `fmt.Errorf` with `%w` so
// callers can `errors.Is` against the sentinel where
// applicable and `errors.Unwrap` to recover the underlying
// cause for diagnostic output.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb/evidence]
//     supplies the message + format-string constants.
//   - [github.com/ActiveMemory/ctx/internal/write/kb/evidence]
//     is the primary caller.
package evidence

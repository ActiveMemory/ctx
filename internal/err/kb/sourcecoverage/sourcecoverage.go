//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcecoverage

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/entity"
)

const (
	// ErrIllegalTransition signals that Advance was called with
	// a (from, to) state pair the source-coverage state machine
	// rejects.
	ErrIllegalTransition = entity.Sentinel(
		text.DescKeyErrKbSourcecoverageIllegalTransitionMsg,
	)
	// ErrUnknownSource signals that Advance referenced a Source
	// not yet present in the ledger AND the new State is not
	// one of the initial states (`discovered`, `admitted`).
	ErrUnknownSource = entity.Sentinel(
		text.DescKeyErrKbSourcecoverageUnknownSourceMsg,
	)
)

// IllegalTransition wraps [ErrIllegalTransition] with the
// offending from-state, to-state, and source name.
//
// Parameters:
//   - from: current state.
//   - to: rejected next state.
//   - source: the source identifier that triggered the call.
//
// Returns:
//   - error: wraps [ErrIllegalTransition] for [errors.Is].
func IllegalTransition(from, to, source string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbSourcecoverageIllegalTransition),
		ErrIllegalTransition, from, to, source,
	)
}

// UnknownSource wraps [ErrUnknownSource] with the offending
// source name and entry state.
//
// Parameters:
//   - source: the source identifier that triggered the call.
//   - state: the entry state the caller supplied.
//
// Returns:
//   - error: wraps [ErrUnknownSource] for [errors.Is].
func UnknownSource(source, state string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbSourcecoverageUnknownSource),
		ErrUnknownSource, source, state,
	)
}

// ReadLedger wraps a file-read failure on the source-coverage
// ledger.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ReadLedger(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbSourcecoverageReadLedger), cause,
	)
}

// MkdirLedgerDir wraps a directory-create failure on the
// ledger's parent directory.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func MkdirLedgerDir(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbSourcecoverageMkdirLedgerDir), cause,
	)
}

// WriteLedger wraps a file-write failure on the ledger file.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func WriteLedger(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbSourcecoverageWriteLedger), cause,
	)
}

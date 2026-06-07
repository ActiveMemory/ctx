//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// CheckIgnore wraps a failure to run git check-ignore (a real exec
// failure, not the normal exit-1 "path is not ignored" answer).
//
// Parameters:
//   - path: the path being checked
//   - cause: the underlying exec error
//
// Returns:
//   - error: "dream: git check-ignore <path>: <cause>"
func CheckIgnore(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamCheckIgnore), path, cause,
	)
}

// WriteScope returns an error when a write target resolves outside
// the dream's allowed write scope (dreams/, ideas/, or specs/ only
// via an accepted promote).
//
// Parameters:
//   - path: the refused write target
//
// Returns:
//   - error: "ctx-dream guard: write outside dream scope refused: <path>"
func WriteScope(path string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamWriteScope), path,
	)
}

// Leak returns an error when a write target resolves to a git-tracked
// path, violating the don't-leak invariant.
//
// Parameters:
//   - path: the refused write target
//
// Returns:
//   - error: "ctx-dream guard: write to tracked path refused: <path>"
func Leak(path string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamLeak), path,
	)
}

// ResolveRoot wraps a failure to resolve the project root.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: resolve project root: <cause>"
func ResolveRoot(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamResolveRoot), cause,
	)
}

// RelPath wraps a failure to compute a relative path.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: compute relative path: <cause>"
func RelPath(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamRelPath), cause,
	)
}

// ReadState wraps a failure to read the state file.
//
// Parameters:
//   - path: the state file path
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: read state <path>: <cause>"
func ReadState(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamReadState), path, cause,
	)
}

// WriteState wraps a failure to write the state file.
//
// Parameters:
//   - path: the state file path
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: write state <path>: <cause>"
func WriteState(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamWriteState), path, cause,
	)
}

// MarshalState wraps a failure to marshal the state slice to JSON.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: marshal state: <cause>"
func MarshalState(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamMarshalState), cause,
	)
}

// UnmarshalState wraps a failure to unmarshal the state file JSON.
//
// Parameters:
//   - path: the state file path
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: unmarshal state <path>: <cause>"
func UnmarshalState(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamUnmarshalState), path, cause,
	)
}

// AppendLedger wraps a failure to append to the ledger file.
//
// Parameters:
//   - path: the ledger file path
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: append ledger <path>: <cause>"
func AppendLedger(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamAppendLedger), path, cause,
	)
}

// ReadLedger wraps a failure to read the ledger file.
//
// Parameters:
//   - path: the ledger file path
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: read ledger <path>: <cause>"
func ReadLedger(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamReadLedger), path, cause,
	)
}

// MarshalEntry wraps a failure to marshal a ledger entry to JSON.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: marshal ledger entry: <cause>"
func MarshalEntry(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamMarshalEntry), cause,
	)
}

// Mkdir wraps a failure to create a dream notebook directory.
//
// Parameters:
//   - path: the directory path
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: create notebook directory <path>: <cause>"
func Mkdir(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamMkdir), path, cause,
	)
}

// InvalidProposal returns an error when a proposal carries an unknown
// status, action, or confidence value.
//
// Parameters:
//   - id: the proposal ID
//   - reason: which field was invalid and its value
//
// Returns:
//   - error: "dream: invalid proposal <id>: <reason>"
func InvalidProposal(id, reason string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamInvalidProposal), id, reason,
	)
}

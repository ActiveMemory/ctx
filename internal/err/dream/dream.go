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

// BackupFailed wraps a failure to back up a source file before a
// destructive mutation. The mutation must abort when this fires.
//
// Parameters:
//   - path: the source file that could not be backed up
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: backup failed for <path>: <cause>"
func BackupFailed(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamBackupFailed), path, cause,
	)
}

// ExecutorNotFound returns a fail-loud error when the configured
// executor binary is not on PATH.
//
// Parameters:
//   - name: the executor binary name
//   - cause: the underlying lookup error
//
// Returns:
//   - error: "[dream] FAIL: executor <name> not on PATH: <cause>"
func ExecutorNotFound(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamExecutorNotFound), name, cause,
	)
}

// ExecutorRun returns a fail-loud error when the executor ran but
// exited non-zero.
//
// Parameters:
//   - name: the executor binary name
//   - cause: the underlying run error
//
// Returns:
//   - error: "[dream] FAIL: executor <name> failed: <cause>"
func ExecutorRun(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamExecutorRun), name, cause,
	)
}

// GuardRefused wraps a guard refusal reason as an error so a
// disposition applier can abort a refused write.
//
// Parameters:
//   - reason: the registry-sourced refusal reason from a GuardDecision
//
// Returns:
//   - error: the reason verbatim
func GuardRefused(reason string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamGuardRefused), reason,
	)
}

// LockAcquire wraps a failure to acquire the dream pass lock.
//
// Parameters:
//   - path: the lock file path
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: acquire lock <path>: <cause>"
func LockAcquire(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamLockAcquire), path, cause,
	)
}

// MoveSource wraps a failure to relocate a source file (archive).
//
// Parameters:
//   - path: the source file path
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: move source <path>: <cause>"
func MoveSource(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamMoveSource), path, cause,
	)
}

// ProposalNotFound returns an error when no proposal with the given
// id exists in the scanned run directory.
//
// Parameters:
//   - id: the requested proposal ID
//   - dir: the run directory searched
//
// Returns:
//   - error: "dream: proposal <id> not found in <dir>"
func ProposalNotFound(id, dir string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamProposalNotFound), id, dir,
	)
}

// ReadProposals wraps a failure to read a proposals file.
//
// Parameters:
//   - path: the proposals file path
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: read proposals <path>: <cause>"
func ReadProposals(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamReadProposals), path, cause,
	)
}

// ReadSource wraps a failure to read a source idea file.
//
// Parameters:
//   - path: the source file path
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: read source <path>: <cause>"
func ReadSource(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamReadSource), path, cause,
	)
}

// ScanIdeas wraps a failure to walk the ideas/ directory.
//
// Parameters:
//   - path: the ideas directory path
//   - cause: the underlying error
//
// Returns:
//   - error: "dream: scan ideas <path>: <cause>"
func ScanIdeas(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamScanIdeas), path, cause,
	)
}

// UnknownAction returns an error when a disposition names an action
// the applier does not recognize.
//
// Parameters:
//   - action: the unrecognized action
//   - id: the proposal ID
//
// Returns:
//   - error: "dream: unknown action <action> for proposal <id>"
func UnknownAction(action, id string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDreamUnknownAction), action, id,
	)
}

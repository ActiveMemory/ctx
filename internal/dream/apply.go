//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"path/filepath"
	"time"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Reject records a rejection in the ledger with no mutation. The
// rejected proposal is not re-surfaced unless its source changes.
//
// Parameters:
//   - dreamsDir: absolute path to the dreams/ notebook directory
//   - p: the proposal being rejected
//   - note: optional human note
//
// Returns:
//   - ApplyResult: Performed=true, Action=the proposal's action
//   - error: non-nil on a ledger append failure
func Reject(
	dreamsDir string, p Proposal, note string,
) (ApplyResult, error) {
	if appendErr := AppendLedger(dreamsDir, LedgerEntry{
		ProposalID: p.ID,
		Decision:   cfgDream.DecisionRejected,
		Action:     p.Action,
		At:         time.Now().UTC(),
		Note:       note,
	}); appendErr != nil {
		return ApplyResult{}, appendErr
	}
	return ApplyResult{Performed: true, Action: p.Action}, nil
}

// Accept applies the proposal's own action. Mechanical actions
// (archive, mark-blog, keep) are performed here and recorded as
// accepted; generative actions (promote, merge) are recorded as
// accepted intent and deferred to the agent.
//
// Parameters:
//   - projectRoot: absolute path to the project root
//   - dreamsDir: absolute path to the dreams/ notebook directory
//   - p: the proposal being accepted
//   - note: optional human note
//
// Returns:
//   - ApplyResult: how the action was dispatched
//   - error: a guard refusal, mutation failure, or ledger failure
func Accept(
	projectRoot, dreamsDir string, p Proposal, note string,
) (ApplyResult, error) {
	return dispatch(
		projectRoot, dreamsDir, p, p.Action,
		cfgDream.DecisionAccepted, note,
	)
}

// Amend applies action in place of the proposal's recommended action,
// recording the decision as amended. Mechanical actions are performed;
// generative actions are deferred to the agent.
//
// Parameters:
//   - projectRoot: absolute path to the project root
//   - dreamsDir: absolute path to the dreams/ notebook directory
//   - p: the proposal being amended
//   - action: the action to apply instead of p.Action
//   - note: optional human note
//
// Returns:
//   - ApplyResult: how the action was dispatched
//   - error: an unknown action, guard refusal, mutation failure, or
//     ledger failure
func Amend(
	projectRoot, dreamsDir string,
	p Proposal, action cfgDream.ProposalAction, note string,
) (ApplyResult, error) {
	if !actionKnown(action) {
		return ApplyResult{}, errDream.UnknownAction(action, p.ID)
	}
	return dispatch(
		projectRoot, dreamsDir, p, action,
		cfgDream.DecisionAmended, note,
	)
}

// Backup snapshots an existing source .md into the dreams/ notebook
// before a destructive mutation. The backup target is itself guarded.
// Backup-before-mutate: a destructive op (merge/overwrite, completed by
// the agent from the full source) must call this and abort the mutation
// when it fails.
//
// Parameters:
//   - projectRoot: absolute path to the project root
//   - dreamsDir: absolute path to the dreams/ notebook directory
//   - source: relative or absolute path to the source .md
//
// Returns:
//   - error: BackupFailed wrapping a guard refusal or copy failure
func Backup(projectRoot, dreamsDir, source string) error {
	abs := source
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(projectRoot, source)
	}
	dst := filepath.Join(
		dreamsDir, filepath.Base(source)+cfgDream.BackupSuffix,
	)
	if guardErr := guard(
		projectRoot, dst, cfgDream.ActionMerge,
	); guardErr != nil {
		return errDream.BackupFailed(source, guardErr)
	}
	data, readErr := ctxIo.SafeReadUserFile(abs)
	if readErr != nil {
		return errDream.BackupFailed(source, readErr)
	}
	if writeErr := ctxIo.SafeWriteFileAtomic(
		dst, data, cfgFs.PermSecret,
	); writeErr != nil {
		return errDream.BackupFailed(source, writeErr)
	}
	return nil
}

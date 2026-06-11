//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"path/filepath"
	"time"

	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgToken "github.com/ActiveMemory/ctx/internal/config/token"
	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// dispatch performs the action and records the decision. Mechanical
// actions mutate here and report Performed; generative actions record
// accepted intent and report Generative for the agent to complete from
// the full source. Every mutation passes both guards (and a backup for
// destructive ops) before any ledger entry is written.
//
// Parameters:
//   - projectRoot: absolute path to the project root
//   - dreamsDir: absolute path to the dreams/ notebook directory
//   - p: the proposal being dispositioned
//   - action: the action to apply (proposal's own, or an amendment)
//   - decision: the recorded review outcome
//   - note: optional human note
//
// Returns:
//   - ApplyResult: how the action was dispatched
//   - error: an unknown action, guard refusal, mutation, or ledger
//     failure
func dispatch(
	projectRoot, dreamsDir string,
	p Proposal,
	action cfgDream.ProposalAction,
	decision cfgDream.Decision,
	note string,
) (ApplyResult, error) {
	generative, mutateErr := mutate(projectRoot, dreamsDir, p, action)
	if mutateErr != nil {
		return ApplyResult{}, mutateErr
	}
	if appendErr := AppendLedger(dreamsDir, LedgerEntry{
		ProposalID: p.ID,
		Decision:   decision,
		Action:     action,
		At:         time.Now().UTC(),
		Note:       note,
	}); appendErr != nil {
		return ApplyResult{}, appendErr
	}
	return ApplyResult{
		Performed:  !generative,
		Generative: generative,
		Action:     action,
	}, nil
}

// mutate carries out the file-system effect of an action. It returns
// generative=true for promote/merge (no mutation here; the agent owns
// it) and performs archive/mark-blog/keep mechanically.
//
// Parameters:
//   - projectRoot: absolute path to the project root
//   - dreamsDir: absolute path to the dreams/ notebook directory
//   - p: the proposal being dispositioned
//   - action: the action to apply
//
// Returns:
//   - bool: true when the action is generative (deferred to the agent)
//   - error: an unknown action, guard refusal, or mutation failure
func mutate(
	projectRoot, dreamsDir string,
	p Proposal, action cfgDream.ProposalAction,
) (bool, error) {
	switch action {
	case cfgDream.ActionArchive:
		return false, archive(projectRoot, p)
	case cfgDream.ActionMarkBlog:
		return false, markBlog(projectRoot, p)
	case cfgDream.ActionKeep:
		return false, nil
	case cfgDream.ActionPromote, cfgDream.ActionMerge:
		return true, nil
	default:
		return false, errDream.UnknownAction(action, p.ID)
	}
}

// guard runs both structural guards (write-scope then don't-leak) for
// a write target under the given action, returning an error built from
// the first refusal's registry-sourced reason. Every dream file
// mutation flows through here before touching disk.
//
// Parameters:
//   - projectRoot: absolute path to the project root
//   - target: the write target (absolute or relative to projectRoot)
//   - action: the disposition driving the write
//
// Returns:
//   - error: a GuardRefused error on refusal; a real exec/resolve
//     error; or nil when both guards allow the write
func guard(
	projectRoot, target string, action cfgDream.ProposalAction,
) error {
	scope, scopeErr := WriteScope(projectRoot, target, action)
	if scopeErr != nil {
		return scopeErr
	}
	if !scope.Allowed {
		return errDream.GuardRefused(scope.Reason)
	}
	leak, leakErr := Leak(projectRoot, target, action)
	if leakErr != nil {
		return leakErr
	}
	if !leak.Allowed {
		return errDream.GuardRefused(leak.Reason)
	}
	return nil
}

// archive moves the first target idea file into ideas/done/, a
// reversible relocation that needs no backup. The destination is
// guarded before the move.
//
// Parameters:
//   - projectRoot: absolute path to the project root
//   - p: the proposal whose first target is archived
//
// Returns:
//   - error: a guard refusal, a missing target, or a move failure
func archive(projectRoot string, p Proposal) error {
	src, srcErr := firstTarget(p)
	if srcErr != nil {
		return srcErr
	}
	absSrc := src
	if !filepath.IsAbs(absSrc) {
		absSrc = filepath.Join(projectRoot, src)
	}
	doneDir := filepath.Join(projectRoot, cfgDir.Ideas, cfgDir.Done)
	dst := filepath.Join(doneDir, filepath.Base(src))
	if guardErr := guard(
		projectRoot, dst, cfgDream.ActionArchive,
	); guardErr != nil {
		return guardErr
	}
	if mkErr := ctxIo.SafeMkdirAll(
		doneDir, cfgFs.PermRestrictedDir,
	); mkErr != nil {
		return errDream.MoveSource(src, mkErr)
	}
	if mvErr := ctxIo.SafeRename(absSrc, dst); mvErr != nil {
		return errDream.MoveSource(src, mvErr)
	}
	return nil
}

// markBlog tags the first target idea file as blog material in place
// by appending a marker line. The in-place write passes both guards.
//
// Parameters:
//   - projectRoot: absolute path to the project root
//   - p: the proposal whose first target is tagged
//
// Returns:
//   - error: a guard refusal, a missing target, or a write failure
func markBlog(projectRoot string, p Proposal) error {
	src, srcErr := firstTarget(p)
	if srcErr != nil {
		return srcErr
	}
	absSrc := src
	if !filepath.IsAbs(absSrc) {
		absSrc = filepath.Join(projectRoot, src)
	}
	if guardErr := guard(
		projectRoot, absSrc, cfgDream.ActionMarkBlog,
	); guardErr != nil {
		return guardErr
	}
	marker := cfgDream.BlogMarker + cfgToken.NewlineLF
	if appendErr := ctxIo.AppendBytes(
		absSrc, []byte(marker), cfgFs.PermSecret,
	); appendErr != nil {
		return errDream.MoveSource(src, appendErr)
	}
	return nil
}

// firstTarget returns the proposal's first target path, or an error
// when the proposal carries no target.
//
// Parameters:
//   - p: the proposal
//
// Returns:
//   - string: the first target path
//   - error: ProposalNotFound when Targets is empty
func firstTarget(p Proposal) (string, error) {
	if len(p.Targets) == 0 {
		return "", errDream.ProposalNotFound(p.ID, cfgDir.Ideas)
	}
	return p.Targets[0], nil
}

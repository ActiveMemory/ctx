//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
	execGit "github.com/ActiveMemory/ctx/internal/exec/git"
)

// WriteScope decides whether a write target is within the dream's
// sanctioned write scope. A target is allowed iff it resolves under
// dreams/ or ideas/ relative to projectRoot, OR under specs/ when the
// action is ActionPromote — the one sanctioned boundary crossing
// (deliberate declassification into a tracked spec).
//
// Parameters:
//   - projectRoot: absolute path to the project root
//   - target: the write target (absolute or relative to projectRoot)
//   - action: the disposition driving the write (gates the specs/
//     crossing)
//
// Returns:
//   - GuardDecision: Allowed plus a registry-sourced refusal Reason
//   - error: non-nil only when the path cannot be resolved
func WriteScope(
	projectRoot, target string, action cfgDream.ProposalAction,
) (GuardDecision, error) {
	rel, relErr := relUnderRoot(projectRoot, target)
	if relErr != nil {
		return GuardDecision{}, relErr
	}
	if underDir(rel, cfgDir.Dreams) || underDir(rel, cfgDir.Ideas) {
		return GuardDecision{Allowed: true}, nil
	}
	if action == cfgDream.ActionPromote && underDir(rel, cfgDir.Specs) {
		return GuardDecision{Allowed: true}, nil
	}
	return GuardDecision{
		Allowed: false,
		Reason:  errDream.WriteScope(target).Error(),
	}, nil
}

// Leak decides whether a write target satisfies the don't-leak
// invariant: a target is allowed iff git reports it as ignored, EXCEPT
// the specs/ promote crossing, which is allowed though tracked (the
// human's deliberate declassification). The check runs git check-ignore
// from projectRoot; a real exec failure is returned as an error, while
// a clean "not ignored" answer becomes a structured refusal.
//
// Parameters:
//   - projectRoot: absolute path to the project root (git working tree)
//   - target: the write target (absolute or relative to projectRoot)
//   - action: the disposition driving the write (exempts the specs/
//     crossing)
//
// Returns:
//   - GuardDecision: Allowed plus a registry-sourced refusal Reason
//   - error: non-nil only on a real git/exec failure
func Leak(
	projectRoot, target string, action cfgDream.ProposalAction,
) (GuardDecision, error) {
	rel, relErr := relUnderRoot(projectRoot, target)
	if relErr != nil {
		return GuardDecision{}, relErr
	}
	if action == cfgDream.ActionPromote && underDir(rel, cfgDir.Specs) {
		return GuardDecision{Allowed: true}, nil
	}
	ignored, checkErr := execGit.CheckIgnore(projectRoot, rel)
	if checkErr != nil {
		return GuardDecision{}, errDream.CheckIgnore(target, checkErr)
	}
	if ignored {
		return GuardDecision{Allowed: true}, nil
	}
	return GuardDecision{
		Allowed: false,
		Reason:  errDream.Leak(target).Error(),
	}, nil
}

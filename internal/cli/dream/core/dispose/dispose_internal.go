//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dispose

import (
	dreamPaths "github.com/ActiveMemory/ctx/internal/cli/dream/core/paths"
	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
	engine "github.com/ActiveMemory/ctx/internal/dream"
	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
)

// load resolves the dream working locations and finds the proposal
// with id in the latest run directory.
//
// Parameters:
//   - id: the proposal ID to locate
//
// Returns:
//   - dreamPaths.Resolved: the resolved working locations
//   - engine.Proposal: the matching proposal
//   - error: a resolution failure, no-runs, or proposal-not-found
func load(id string) (dreamPaths.Resolved, engine.Proposal, error) {
	loc, locErr := dreamPaths.Resolve()
	if locErr != nil {
		return dreamPaths.Resolved{}, engine.Proposal{}, locErr
	}
	runDir, runErr := engine.LatestRunDir(loc.Dreams)
	if runErr != nil {
		return loc, engine.Proposal{}, runErr
	}
	if runDir == "" {
		return loc, engine.Proposal{},
			errDream.ProposalNotFound(id, cfgDir.Dreams)
	}
	p, findErr := engine.FindProposal(runDir, id)
	if findErr != nil {
		return loc, engine.Proposal{}, findErr
	}
	return loc, p, nil
}

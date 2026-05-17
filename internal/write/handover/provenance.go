//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handover

import (
	cfgHandover "github.com/ActiveMemory/ctx/internal/config/handover"
	errHandover "github.com/ActiveMemory/ctx/internal/err/handover"
	"github.com/ActiveMemory/ctx/internal/git_meta"
)

// resolveProvenance picks the SHA / branch pair for a new
// handover. When override is non-empty, it is used verbatim
// for SHA; branch comes from git or "detached" when git is
// unavailable.
//
// Parameters:
//   - projectRoot: absolute path to the project root.
//   - override: optional explicit commit SHA (from --commit).
//
// Returns:
//   - string: short SHA.
//   - string: branch name.
//   - error: non-nil when ResolveHead fails and override is
//     empty.
func resolveProvenance(projectRoot, override string) (string, string, error) {
	if override != "" {
		ref, headErr := gitmeta.ResolveHead(projectRoot)
		if headErr != nil {
			return override, cfgHandover.BranchDetached, nil
		}
		return override, ref.Branch, nil
	}
	ref, headErr := gitmeta.ResolveHead(projectRoot)
	if headErr != nil {
		return "", "", errHandover.ResolveHead(headErr)
	}
	return ref.SHA, ref.Branch, nil
}

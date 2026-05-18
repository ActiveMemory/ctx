//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package gitmeta

import (
	"os"
	"strings"

	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	cfgGitmeta "github.com/ActiveMemory/ctx/internal/config/gitmeta"
	errGitmeta "github.com/ActiveMemory/ctx/internal/err/gitmeta"
	execGit "github.com/ActiveMemory/ctx/internal/exec/git"
)

// ResolveHead reads HEAD into a [HeadRef]. Environment
// overrides for CI replay are checked first:
//
//   - CTX_TASK_COMMIT, when non-empty, is used verbatim as the
//     SHA. Branch is still resolved from git (or "detached").
//   - GITHUB_SHA, when GITHUB_ACTIONS="true" and GITHUB_SHA is
//     non-empty, is truncated to the short form. Branch is
//     resolved from git (or "detached").
//
// Otherwise, `git rev-parse --short HEAD` returns the
// abbreviated SHA, and `git rev-parse --abbrev-ref HEAD`
// (via FlagShowCurrent) returns the current branch.
//
// Parameters:
//   - projectRoot: absolute path to the project root; passed
//     to git via -C so resolution is independent of the caller
//     CWD.
//
// Returns:
//   - HeadRef: resolved short SHA + branch.
//   - error: non-nil on resolution failure. [RequireGitTree]
//     must succeed before calling this.
func ResolveHead(projectRoot string) (HeadRef, error) {
	if v := strings.TrimSpace(os.Getenv(cfgGitmeta.EnvCtxTaskCommit)); v != "" {
		return HeadRef{
			SHA:    v,
			Branch: resolveBranchOrDetached(projectRoot),
		}, nil
	}
	if os.Getenv(cfgGitmeta.EnvGithubActions) == cfgGitmeta.GithubActionsTrue {
		if v := strings.TrimSpace(os.Getenv(cfgGitmeta.EnvGithubSHA)); v != "" {
			return HeadRef{
				SHA:    shortSHA(v),
				Branch: resolveBranchOrDetached(projectRoot),
			}, nil
		}
	}

	out, runErr := execGit.Run(
		cfgGit.FlagChangeDir, projectRoot,
		cfgGit.RevParse, cfgGit.FlagShort, cfgGitmeta.RefHEAD,
	)
	if runErr != nil {
		return HeadRef{}, errGitmeta.ResolveHeadFailed(runErr)
	}
	sha := strings.TrimSpace(string(out))
	if sha == "" {
		return HeadRef{}, errGitmeta.ErrResolveHeadEmpty
	}
	return HeadRef{
		SHA:    sha,
		Branch: resolveBranchOrDetached(projectRoot),
	}, nil
}

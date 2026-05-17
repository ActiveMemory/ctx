//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package gitmeta

// Environment variables honored as provenance overrides for
// CI replay.
const (
	// EnvCtxTaskCommit overrides the resolved short SHA with
	// a user-supplied value (typically used by CI replay).
	EnvCtxTaskCommit = "CTX_TASK_COMMIT"

	// EnvGithubActions is GitHub Actions' canonical truthy
	// marker; checked alongside EnvGithubSHA so the override
	// fires only inside a real Actions run.
	EnvGithubActions = "GITHUB_ACTIONS"

	// EnvGithubSHA is GitHub Actions' commit-SHA injection.
	EnvGithubSHA = "GITHUB_SHA"
)

// GithubActionsTrue is the literal value that GITHUB_ACTIONS
// carries when running under GitHub Actions.
const GithubActionsTrue = "true"

// BranchDetached is the canonical branch-name placeholder when
// HEAD is detached (points at a commit, not a symbolic ref).
// Used by [github.com/ActiveMemory/ctx/internal/git_meta.ResolveHead]
// and recorded verbatim in closeout / handover frontmatter.
const BranchDetached = "detached"

// RefHEAD is the git ref name for the current commit.
const RefHEAD = "HEAD"

// ShortLen is the truncation length for short SHAs (git's
// default --short width).
const ShortLen = 7

// Sentinel error-message constants. These back `errors.New`
// values in `internal/err/git_meta/` and are matched via
// `errors.Is` at the call site. They cannot use desc.Text
// because the sentinels are package-level vars evaluated
// before the embedded YAML lookup is populated; wrapping
// format strings live in commands/text/errors.yaml.
const (
	// ErrMsgMissingGitTree is the sentinel for the
	// "<projectRoot>/.git is absent" condition.
	ErrMsgMissingGitTree = "git working tree required"
	// ErrMsgResolveHeadEmpty signals that `git rev-parse --short
	// HEAD` returned an empty string (typically: unborn HEAD).
	ErrMsgResolveHeadEmpty = "resolve git head: empty output"
)

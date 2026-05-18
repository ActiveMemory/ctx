//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for gitmeta error wrappers. The matching YAML
// entries live in commands/text/errors.yaml; constructors in
// internal/err/gitmeta/ resolve them via desc.Text at error
// construction time.
const (
	// DescKeyErrGitmetaMissingGitTree is the text key for the
	// missing-.git/ wrap used by direct API callers.
	DescKeyErrGitmetaMissingGitTree = "err.gitmeta.missing-git-tree"
	// DescKeyErrGitmetaMissingGitTreeForCmd is the text key for
	// the missing-.git/ wrap used by the root PersistentPreRunE.
	DescKeyErrGitmetaMissingGitTreeForCmd = "err.gitmeta.missing-git-tree-for-cmd"
	// DescKeyErrGitmetaStatGitDir is the text key for the
	// non-ENOENT stat failure on the `.git` entry.
	DescKeyErrGitmetaStatGitDir = "err.gitmeta.stat-git-dir"
	// DescKeyErrGitmetaResolveHead is the text key for the
	// `git rev-parse --short HEAD` invocation failure.
	DescKeyErrGitmetaResolveHead = "err.gitmeta.resolve-head"
	// DescKeyErrGitmetaMissingGitTreeMsg is the text key for the
	// missing-git-tree sentinel's own `.Error()` string (the
	// prefix interpolated via `%w` by the wrapper formats).
	DescKeyErrGitmetaMissingGitTreeMsg = "err.gitmeta.missing-git-tree-msg"
	// DescKeyErrGitmetaResolveHeadEmpty is the text key for the
	// empty-HEAD-output sentinel.
	DescKeyErrGitmetaResolveHeadEmpty = "err.gitmeta.resolve-head-empty"
)

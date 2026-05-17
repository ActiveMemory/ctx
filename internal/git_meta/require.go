//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package gitmeta

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	errGitmeta "github.com/ActiveMemory/ctx/internal/err/git_meta"
)

// RequireGitTree returns nil when `<projectRoot>/.git` exists
// as a directory (regular repo) or a regular file (worktree
// pointer per git convention).
//
// Parameters:
//   - projectRoot: absolute path to the project root (parent of
//     `.context/`, by ctx convention).
//
// Returns:
//   - error: nil on success;
//     [errGitmeta.MissingGitTree]-wrapping error when `.git`
//     is absent (matchable via `errors.Is` against
//     [errGitmeta.ErrMissingGitTree]); a wrapped stat error
//     for other failures.
func RequireGitTree(projectRoot string) error {
	p := filepath.Join(projectRoot, cfgGit.DotDir)
	if _, err := os.Stat(p); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return errGitmeta.MissingGitTree(projectRoot)
		}
		return errGitmeta.StatGitDir(p, err)
	}
	return nil
}

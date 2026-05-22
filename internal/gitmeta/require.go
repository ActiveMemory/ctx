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
	errGitmeta "github.com/ActiveMemory/ctx/internal/err/gitmeta"
)

// RequireGitTree returns nil when `$PWD/.git` exists as a directory
// (regular repo) or a regular file (worktree pointer per git
// convention). Under the cwd-anchored resolution model
// (spec: specs/cwd-anchored-context.md) the project root is always
// the current working directory, so this check needs no walk and
// no caller-supplied root.
//
// Returns:
//   - error: nil on success;
//     [errGitmeta.MissingGitTree]-wrapping error when `.git` is
//     absent (matchable via `errors.Is` against
//     [errGitmeta.ErrMissingGitTree]); a wrapped stat error for
//     other failures.
func RequireGitTree() error {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return cwdErr
	}
	p := filepath.Join(cwd, cfgGit.DotDir)
	if _, err := os.Stat(p); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return errGitmeta.MissingGitTree(cwd)
		}
		return errGitmeta.StatGitDir(p, err)
	}
	return nil
}

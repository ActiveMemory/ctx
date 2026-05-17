//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package gitmeta

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgGitmeta "github.com/ActiveMemory/ctx/internal/config/git_meta"
)

// ErrMissingGitTree signals that `<projectRoot>/.git` is
// absent. The PersistentPreRunE catches this via `errors.Is`
// and wraps it with the failing subcommand name through
// [MissingGitTreeForCmd]; direct API callers may wrap with
// [MissingGitTree].
var ErrMissingGitTree = errors.New(cfgGitmeta.ErrMsgMissingGitTree)

// ErrResolveHeadEmpty signals that `git rev-parse --short HEAD`
// returned an empty string. Typically: unborn HEAD (repository
// initialized but no commit yet).
var ErrResolveHeadEmpty = errors.New(cfgGitmeta.ErrMsgResolveHeadEmpty)

// MissingGitTree wraps [ErrMissingGitTree] with the
// project-root path that was scanned. Used by direct API
// callers (i.e. anyone not invoked via cobra).
//
// Parameters:
//   - projectRoot: absolute path of the directory scanned.
//
// Returns:
//   - error: wrapping [ErrMissingGitTree] for [errors.Is]
//     matches at the call site.
func MissingGitTree(projectRoot string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrGitmetaMissingGitTree),
		ErrMissingGitTree, projectRoot,
	)
}

// MissingGitTreeForCmd wraps [ErrMissingGitTree] with both the
// failing subcommand name and the project-root path.
//
// Parameters:
//   - cmdName: the failing subcommand (`init`, `kb`, etc.).
//   - projectRoot: absolute path of the directory scanned.
//
// Returns:
//   - error: wrapping [ErrMissingGitTree] for [errors.Is]
//     matches at the call site.
func MissingGitTreeForCmd(cmdName, projectRoot string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrGitmetaMissingGitTreeForCmd),
		ErrMissingGitTree, cmdName, projectRoot,
	)
}

// StatGitDir wraps a non-ENOENT stat failure on the `.git`
// entry (permission denied, I/O error, etc.).
//
// Parameters:
//   - path: full path to `<projectRoot>/.git`.
//   - cause: underlying `os.Stat` error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func StatGitDir(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrGitmetaStatGitDir), path, cause,
	)
}

// ResolveHeadFailed wraps a `git rev-parse --short HEAD`
// invocation failure.
//
// Parameters:
//   - cause: underlying exec / git error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func ResolveHeadFailed(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrGitmetaResolveHead), cause)
}

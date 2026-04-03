//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	readHook "github.com/ActiveMemory/ctx/internal/assets/read/hook"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
	errTrace "github.com/ActiveMemory/ctx/internal/err/trace"
	"github.com/ActiveMemory/ctx/internal/exec/git"
	"github.com/ActiveMemory/ctx/internal/io"
	writeTrace "github.com/ActiveMemory/ctx/internal/write/trace"
)

// Enable installs both the prepare-commit-msg and post-commit hooks.
//
// Parameters:
//   - cmd: Cobra command for output stream
//
// Returns:
//   - error: non-nil on installation failure
func Enable(cmd *cobra.Command) error {
	prepScript, prepReadErr := readHook.TraceScript(
		cfgTrace.ScriptPrepareCommitMsg)
	if prepReadErr != nil {
		return prepReadErr
	}
	prepPath, prepErr := FilePath(cfgGit.HookPrepareCommitMsg)
	if prepErr != nil {
		return prepErr
	}
	if installErr := Install(
		prepPath, prepScript, cfgGit.HookPrepareCommitMsg,
	); installErr != nil {
		return installErr
	}

	postScript, postReadErr := readHook.TraceScript(cfgTrace.ScriptPostCommit)
	if postReadErr != nil {
		return postReadErr
	}
	postPath, postErr := FilePath(cfgGit.HookPostCommit)
	if postErr != nil {
		return postErr
	}
	if installErr := Install(
		postPath, postScript, cfgGit.HookPostCommit,
	); installErr != nil {
		return installErr
	}

	writeTrace.HooksEnabled(cmd)
	return nil
}

// Disable removes both the prepare-commit-msg and post-commit hooks if they
// were installed by ctx.
//
// Parameters:
//   - cmd: Cobra command for output stream
//
// Returns:
//   - error: non-nil on removal failure
func Disable(cmd *cobra.Command) error {
	prepPath, prepErr := FilePath(cfgGit.HookPrepareCommitMsg)
	if prepErr != nil {
		return prepErr
	}
	Remove(prepPath)

	postPath, postErr := FilePath(cfgGit.HookPostCommit)
	if postErr != nil {
		return postErr
	}
	Remove(postPath)

	writeTrace.HooksDisabled(cmd)
	return nil
}

// Install writes the hook script to path, checking for existing non-ctx hooks.
//
// Parameters:
//   - path: absolute path to the hook file
//   - script: hook script content to write
//   - name: hook name for error messages
//
// Returns:
//   - error: non-nil if a non-ctx hook already exists or write fails
func Install(path, script, name string) error {
	if _, statErr := io.SafeStat(path); statErr == nil {
		existing, readErr := io.SafeReadUserFile(path)
		if readErr == nil && !strings.Contains(
			string(existing), cfgTrace.CtxTraceMarker,
		) {
			return errTrace.HookExists(name, path)
		}
	}
	if writeErr := io.SafeWriteFile(
		path, []byte(script), cfgFs.PermExec,
	); writeErr != nil {
		return errTrace.HookWrite(name, writeErr)
	}
	return nil
}

// Remove removes the hook at path if it was installed by ctx.
//
// Parameters:
//   - path: absolute path to the hook file
func Remove(path string) {
	existing, err := io.SafeReadUserFile(path)
	if err != nil {
		return
	}
	if strings.Contains(string(existing), cfgTrace.CtxTraceMarker) {
		_ = os.Remove(path)
	}
}

// FilePath returns the absolute path to a git hook by name.
//
// Parameters:
//   - hookName: name of the git hook (e.g. "prepare-commit-msg")
//
// Returns:
//   - string: absolute path to the hook file
//   - error: non-nil if git rev-parse fails
func FilePath(hookName string) (string, error) {
	out, err := git.Run(cfgGit.RevParse, cfgGit.FlagGitDir)
	if err != nil {
		return "", errTrace.GitDir(err)
	}
	gitDir := strings.TrimSpace(string(out))
	return filepath.Join(gitDir, cfgGit.HooksDir, hookName), nil
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package git

import (
	"errors"
	"os/exec"
	"strings"
	"time"

	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	errGit "github.com/ActiveMemory/ctx/internal/err/git"
)

// Run executes a git command with the given arguments and returns
// its combined stdout output. LookPath is checked on every call.
//
// Parameters:
//   - args: git subcommand and flags (e.g. "log", "--oneline")
//
// Returns:
//   - []byte: raw git stdout
//   - error: non-nil if git is not found or the command fails
func Run(args ...string) ([]byte, error) {
	if _, lookErr := exec.LookPath(cfgGit.Binary); lookErr != nil {
		return nil, errGit.NotFound()
	}
	//nolint:gosec // G204: args are validated by callers
	return exec.Command(cfgGit.Binary, args...).Output()
}

// CheckIgnore reports whether path is ignored by git, running
// `git check-ignore -q -- <path>` from within dir. git check-ignore
// signals its answer through the exit code: 0 means the path is
// ignored, 1 means it is not, and 128+ means a real failure. This
// helper maps exit 0/1 to a clean bool and surfaces only genuine
// failures (git missing, not a repo) as errors.
//
// Parameters:
//   - dir: directory to run the check from (the repo working tree)
//   - path: path to test for ignore status (absolute or relative)
//
// Returns:
//   - bool: true when git reports the path as ignored
//   - error: non-nil only on a real exec failure (not on exit 1)
func CheckIgnore(dir, path string) (bool, error) {
	if _, lookErr := exec.LookPath(cfgGit.Binary); lookErr != nil {
		return false, errGit.NotFound()
	}
	//nolint:gosec // G204: binary is fixed; dir/path validated by callers
	cmd := exec.Command(
		cfgGit.Binary, cfgGit.FlagChangeDir, dir,
		cfgGit.CheckIgnore, cfgGit.FlagQuiet,
		cfgGit.FlagPathSep, path,
	)
	runErr := cmd.Run()
	if runErr == nil {
		return true, nil
	}
	if exitErr, ok := errors.AsType[*exec.ExitError](
		runErr,
	); ok && exitErr.ExitCode() == cfgGit.CheckIgnoreNotIgnored {
		return false, nil
	}
	return false, runErr
}

// Root returns the repository root directory for the current
// working directory.
//
// Returns:
//   - string: absolute path to the repository root
//   - error: non-nil if git is not found or CWD is not in a repo
func Root() (string, error) {
	out, runErr := Run(cfgGit.RevParse, cfgGit.FlagShowToplevel)
	if runErr != nil {
		return "", errGit.NotInRepo(runErr)
	}
	return strings.TrimSpace(string(out)), nil
}

// RemoteURL returns the origin remote URL for a directory.
// Returns an empty string on any error (best-effort).
//
// Parameters:
//   - dir: directory path to query
//
// Returns:
//   - string: remote URL, or empty string on any error
func RemoteURL(dir string) string {
	if dir == "" {
		return ""
	}
	out, runErr := Run(
		cfgGit.FlagChangeDir, dir,
		cfgGit.Remote, cfgGit.RemoteGetURL, cfgGit.RemoteOrigin,
	)
	if runErr != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// LogSince runs git log with a --since filter derived from t.
//
// Parameters:
//   - t: reference time for --since
//   - extraArgs: additional literal git log flags
//
// Returns:
//   - []byte: raw git output
//   - error: non-nil if git is not found or the command fails
func LogSince(
	t time.Time, extraArgs ...string,
) ([]byte, error) {
	args := []string{
		cfgGit.Log, cfgGit.FlagSince, t.Format(time.RFC3339),
	}
	args = append(args, extraArgs...)
	return Run(args...)
}

// LastCommitMessage returns the full message of the most recent
// commit.
//
// Returns:
//   - []byte: raw commit message
//   - error: non-nil if git is not found or the command fails
func LastCommitMessage() ([]byte, error) {
	return Run(cfgGit.Log, cfgGit.FlagLast, cfgGit.FormatBody)
}

// ShortHead returns the abbreviated commit hash for HEAD.
//
// Returns:
//   - string: short commit hash (7-8 chars), or empty on error
func ShortHead() string {
	out, err := Run(cfgGit.RevParse, cfgGit.FlagShort, cfgGit.RefHead)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// CurrentBranch returns the current branch name.
//
// Returns:
//   - string: branch name, or empty if detached or on error
func CurrentBranch() string {
	out, err := Run(cfgGit.Branch, cfgGit.FlagShowCurrent)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// DiffTreeHead returns the list of files changed in HEAD.
//
// Returns:
//   - []byte: newline-separated file paths
//   - error: non-nil if git is not found or the command fails
func DiffTreeHead() ([]byte, error) {
	return Run(
		cfgGit.DiffTree, cfgGit.FlagNoCommitID,
		cfgGit.FlagNameOnly, cfgGit.FlagRecursive, cfgGit.RefHead,
	)
}

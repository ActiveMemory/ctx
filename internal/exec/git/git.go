//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package git

import (
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

// DiffTreeHead returns the list of files changed in HEAD.
//
// Returns:
//   - []byte: newline-separated file paths
//   - error: non-nil if git is not found or the command fails
func DiffTreeHead() ([]byte, error) {
	return Run(
		cfgGit.DiffTree, cfgGit.FlagNoCommitID,
		cfgGit.FlagNameOnly, cfgGit.FlagRecursive, "HEAD",
	)
}

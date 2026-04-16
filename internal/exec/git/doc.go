//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package git wraps git command execution behind
// typed functions.
//
// All exec.Command calls for git are centralized
// here. LookPath is checked on every call. Callers
// never import os/exec directly for git operations.
//
// # Running Commands
//
// Run executes a git command with the given arguments
// and returns raw stdout output.
//
//	out, err := git.Run("log", "--oneline")
//
// # Repository Queries
//
// Root returns the repository root directory for the
// current working directory. RemoteURL returns the
// origin remote URL for a given directory path (best
// effort, returns empty on error).
//
//	root, err := git.Root()
//	url := git.RemoteURL("/path/to/repo")
//
// # Log and Diff
//
// LogSince runs git log with a --since time filter.
// LastCommitMessage returns the full message of the
// most recent commit. DiffTreeHead lists files
// changed in HEAD.
//
//	out, err := git.LogSince(since, "--oneline")
//	msg, err := git.LastCommitMessage()
//	files, err := git.DiffTreeHead()
//
// # HEAD Queries
//
// ShortHead returns the abbreviated commit hash for
// HEAD. CurrentBranch returns the current branch
// name (empty if detached or on error).
//
//	hash := git.ShortHead()
//	branch := git.CurrentBranch()
package git

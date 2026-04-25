//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package git

// Binary is the git executable name.
const Binary = "git"

// Subcommand names passed as the first argument to git.
const (
	Branch   = "branch"
	Diff     = "diff"
	DiffTree = "diff-tree"
	Log      = "log"
	Remote   = "remote"
	RevParse = "rev-parse"
)

// Hook names used in .git/hooks/.
const (
	HookPrepareCommitMsg = "prepare-commit-msg"
	HookPostCommit       = "post-commit"
	HooksDir             = "hooks"
)

// Rev-parse flags.
const (
	FlagShort        = "--short"
	FlagShowCurrent  = "--show-current"
	FlagShowToplevel = "--show-toplevel"
	FlagGitDir       = "--git-dir"
)

// Common flags and format strings for git commands.
const (
	FlagCached         = "--cached"
	FlagChangeDir      = "-C"
	FlagLast           = "-1"
	FlagNoCommitID     = "--no-commit-id"
	FlagNameOnly       = "--name-only"
	FlagOneline        = "--oneline"
	FlagRecursive      = "-r"
	FlagSince          = "--since"
	FormatAuthor       = "--format=%aN"
	FormatBody         = "--format=%B"
	FormatEmpty        = "--format="
	FormatDateISO      = "--format=%ci"
	FormatHashDateSubj = "--format=%H %ci %s"
	FormatHashSubj     = "--format=%H %s"
	FormatSubject      = "--format=%s"
	FormatTrailerValue = "--format=%%(trailers:key=%s,valueonly)"
	// FlagPathSep is the separator between flags and paths.
	FlagPathSep = "--"
	// FlagLastN is the format string for limiting git log
	// output to the last N commits (e.g. "-5").
	FlagLastN = "-%d"
)

// Ref constants for addressing commits and branches.
const (
	// RefHead is the symbolic reference for the current commit.
	RefHead = "HEAD"
)

// Remote subcommands and arguments.
const (
	RemoteGetURL = "get-url"
	RemoteOrigin = "origin"
)

// PathSeparator is the separator git uses in file paths (always forward slash).
const PathSeparator = "/"

// Commit trailer keys for structured metadata in commit messages.
const (
	// TrailerSpec is the commit trailer for spec references.
	TrailerSpec = "Spec: specs/"
	// TrailerSignedOffBy is the commit trailer for sign-off.
	TrailerSignedOffBy = "Signed-off-by:"
)

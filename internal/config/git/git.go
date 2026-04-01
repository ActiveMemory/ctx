//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package git

// Binary is the git executable name.
const Binary = "git"

// Git subcommands.
const (
	DiffTree = "diff-tree"
	Log      = "log"
	Remote   = "remote"
	RevParse = "rev-parse"
)

// Git hook names.
const (
	HookPrepareCommitMsg = "prepare-commit-msg"
	HookPostCommit       = "post-commit"
	HooksDir             = "hooks"
)

// Git subcommands (additional).
const (
	Diff = "diff"
)

// Git rev-parse flags.
const (
	FlagShowToplevel = "--show-toplevel"
	FlagGitDir       = "--git-dir"
)

// Git flags.
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
)

// Git remote subcommands and arguments.
const (
	RemoteGetURL = "get-url"
	RemoteOrigin = "origin"
)

// PathSeparator is the separator git uses in file paths (always forward slash).
const PathSeparator = "/"

// Git commit trailers.
const (
	// TrailerSpec is the commit trailer for spec references.
	TrailerSpec = "Spec: specs/"
	// TrailerSignedOffBy is the commit trailer for sign-off.
	TrailerSignedOffBy = "Signed-off-by:"
)

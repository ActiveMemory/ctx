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

// Git rev-parse flags.
const (
	// FlagShowToplevel is a git flag.
	FlagShowToplevel = "--show-toplevel"
)

// Git flags.
const (
	FlagChangeDir  = "-C"
	FlagLast       = "-1"
	FlagNoCommitID = "--no-commit-id"
	FlagNameOnly   = "--name-only"
	FlagOneline    = "--oneline"
	FlagRecursive  = "-r"
	FlagSince      = "--since"
	FormatAuthor   = "--format=%aN"
	FormatBody     = "--format=%B"
	FormatEmpty    = "--format="
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

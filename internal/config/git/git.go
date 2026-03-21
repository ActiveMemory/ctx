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
	Log      = "log"
	RevParse = "rev-parse"
)

// Git rev-parse flags.
const (
	FlagShowToplevel = "--show-toplevel"
)

// Git log flags.
const (
	FlagNoCommitID = "--no-commit-id"
	FlagNameOnly   = "--name-only"
	FlagOneline    = "--oneline"
	FlagSince      = "--since"
	FormatAuthor   = "--format=%aN"
	FormatEmpty    = "--format="
)

// PathSeparator is the separator git uses in file paths (always forward slash).
const PathSeparator = "/"

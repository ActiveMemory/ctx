//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package git centralizes constants for invoking the
// git binary: the executable name, subcommands, flags,
// format strings, hook names, ref identifiers, remote
// arguments, and commit trailer keys.
//
// ctx relies on git for branch detection, commit
// metadata extraction, diff generation, and hook
// installation. Every git argument or format string
// used in exec.Command calls is defined here so that
// changes to git invocation are single-point edits
// with compile-time verification.
//
// # Binary and Metadata Directory
//
//   - Binary ("git") -- the executable name passed
//     to exec.LookPath and exec.Command
//   - DotDir (".git") -- the metadata directory (or
//     file, in worktrees) used to detect whether a
//     directory is a git repository
//
// # Subcommands
//
//   - Branch, Diff, DiffTree, Log, Remote, RevParse
//     -- first arguments to the git binary
//
// # Hook Names
//
//   - HookPrepareCommitMsg -- the prepare-commit-msg
//     hook for injecting commit trailers
//   - HookPostCommit -- the post-commit hook for
//     session event recording
//   - HooksDir ("hooks") -- the subdirectory under
//     .git/ where hooks live
//
// # Rev-Parse and Common Flags
//
// Rev-parse flags (FlagShort, FlagShowCurrent,
// FlagShowToplevel, FlagGitDir) and general flags
// (FlagCached, FlagChangeDir, FlagNameOnly,
// FlagOneline, FlagSince, etc.) are defined as named
// constants.
//
// # Format Strings
//
// Git --format templates for extracting commit data:
//
//   - FormatAuthor, FormatBody, FormatDateISO,
//     FormatHashDateSubj, FormatHashSubj,
//     FormatSubject, FormatTrailerValue
//   - FormatEmpty -- suppresses default output
//
// # Refs and Remotes
//
//   - RefHead ("HEAD") -- symbolic reference for the
//     current commit
//   - RemoteGetURL, RemoteOrigin -- arguments for
//     git remote commands
//   - PathSeparator ("/") -- git's path separator
//
// # Commit Trailers
//
//   - TrailerSpec ("Spec: specs/") -- links a commit
//     to its design spec
//   - TrailerSignedOffBy ("Signed-off-by:") -- the
//     standard sign-off trailer
//
// # Why Centralized
//
// Git arguments are scattered across branch detection,
// commit parsing, hook installation, journal import,
// and diff generation. Centralizing them prevents
// silent breakage from typos in flag strings and makes
// it easy to audit every git invocation in the
// codebase.
package git

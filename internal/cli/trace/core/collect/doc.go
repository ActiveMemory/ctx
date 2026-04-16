//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package collect records context references after a git
// commit. It is called from the post-commit hook and the
// "ctx trace collect" CLI subcommand.
//
// # Recording Flow
//
// [RecordCommit] performs a three-step process for each
// commit:
//
//  1. Read context refs from the commit trailer. The
//     trailer is the single source of truth, set by the
//     prepare-commit-msg hook during commit creation.
//  2. Write a history entry to .context/trace/ containing
//     the commit hash, refs, and commit message.
//  3. Truncate the pending refs file in .context/state/
//     so stale refs never leak into future commits.
//
// # No-Trailer Case
//
// When the commit has no context trailer (e.g., the user
// committed with --no-verify or an external tool), step 1
// returns an empty ref list. In this case, no history
// entry is written, but pending refs are still truncated.
// This prevents stale refs accumulated during the current
// commit window from attaching to the next commit.
//
// # Pending State Lifecycle
//
// Pending context refs are accumulated by hooks between
// commits. Each commit consumes (truncates) the pending
// state regardless of whether a trailer was present.
// This ensures a clean slate for the next commit window.
package collect

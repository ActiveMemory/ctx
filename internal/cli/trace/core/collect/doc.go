//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package collect implements the core logic for recording context refs
// after a commit.
//
// [RecordCommit] reads refs from the commit trailer (the single source
// of truth set by the prepare-commit-msg hook), writes a history entry,
// and truncates pending state so stale refs never leak into future
// commits.
//
// Key exports: [RecordCommit].
// Called by the collect CLI subcommand and the post-commit hook path.
package collect

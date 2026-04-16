//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package provenance provides terminal output for
// session and git identity lines emitted by hooks.
//
// Before any conditional hook logic runs, the hook
// system emits a provenance line that identifies the
// current session, git branch, and commit. This
// package formats and prints that line.
//
// # Output
//
// [Line] prints a single provenance line containing
// the short session ID, branch name, commit hash,
// and an optional context-free percentage suffix.
// The suffix is formatted by [ContextSuffix], which
// returns an empty string when the percentage is
// zero or out of range.
//
// Together they produce output like:
//
//	[abc123] main @ def456 | Context: 45% free
package provenance

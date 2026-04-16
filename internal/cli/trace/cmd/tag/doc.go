//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tag implements the "ctx trace tag" cobra
// subcommand.
//
// This command attaches a free-text context note to a
// specific commit by writing an override entry to the
// trace directory. Use it to retroactively annotate
// commits that were made without automatic tracing or
// to add additional context to an existing commit.
//
// # Usage
//
//	ctx trace tag <commit> --note <text>
//
// # Arguments
//
// Exactly one positional argument is required:
//
//   - commit: a commit ref or hash to tag (e.g.
//     "HEAD", "abc1234"). The ref is resolved to a
//     full SHA before recording.
//
// # Flags
//
//	--note   Required. The context note to attach
//	         to the commit. Must be non-empty.
//
// # Behavior
//
// The command:
//
//   - Validates that --note is non-empty.
//   - Resolves the commit ref to a full hash via
//     trace.ResolveCommitHash.
//   - Builds an OverrideEntry with the hash and
//     the quoted note as a ref.
//   - Writes the entry to the trace directory via
//     trace.WriteOverride.
//   - Prints a confirmation with the short hash and
//     the attached note.
//
// # Output
//
// A single confirmation line showing the short
// commit hash and the note that was attached.
//
// # Delegation
//
// Hash resolution and override writing use the trace
// package. Output formatting uses write/trace.
package tag

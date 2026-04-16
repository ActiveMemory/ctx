//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package show implements the "ctx trace" cobra
// subcommand (the default trace display handler).
//
// This command shows context refs attached to recent
// commits or a specific commit. It is the primary way
// to inspect what decisions, tasks, and learnings
// were recorded alongside code changes.
//
// # Usage
//
//	ctx trace [commit] [--last N] [--json]
//
// # Arguments
//
// An optional positional argument:
//
//   - commit: a commit hash or ref to inspect.
//     When omitted, defaults to showing the last
//     N commits (controlled by --last or a built-in
//     default).
//
// # Flags
//
//	--last, -l   Number of recent commits to show.
//	             When set to a positive value, takes
//	             precedence over positional arguments.
//	             Defaults to the configured default
//	             (typically 5).
//	--json       Format output as JSON instead of a
//	             human-readable table.
//
// # Behavior
//
// The command resolves the context and trace
// directories, then dispatches to one of two modes:
//
//   - Last-N mode: calls core/show.Last to display
//     context refs for the most recent N commits.
//   - Single-commit mode: calls core/show.Commit to
//     display context refs for a specific hash.
//
// In both modes, refs are collected from the trace
// history file and any override entries.
//
// # Output
//
// A table (or JSON array) of commits with their
// associated context refs, including short hash,
// subject line, and ref details.
//
// # Delegation
//
// Display logic is in trace/core/show. Directory
// resolution uses rc and config/dir.
package show

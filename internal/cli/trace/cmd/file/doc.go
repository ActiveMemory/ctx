//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package file implements the "ctx trace file" cobra
// subcommand.
//
// This command shows context refs attached to commits
// that touched a specific file. It answers the
// question "why was this file changed?" by displaying
// the decisions, tasks, and learnings recorded at
// each commit.
//
// # Usage
//
//	ctx trace file <path> [--last N]
//
// # Arguments
//
// Exactly one positional argument is required:
//
//   - path: the file to trace. An optional line-range
//     suffix is supported (e.g. "src/auth.go:42-60")
//     and is stripped before querying git log.
//
// # Flags
//
//	--last, -l   Maximum number of commits to show.
//	             Defaults to the configured default
//	             for file traces.
//
// # Behavior
//
// The command:
//
//   - Parses the path argument, stripping any
//     line-range suffix.
//   - Runs git log for the resolved file path to
//     retrieve commits that touched it.
//   - For each commit, collects context refs from
//     both the trace history file and any override
//     entries.
//   - Prints the results as a table with commit hash,
//     date, and associated context refs.
//
// # Output
//
// A table of commits with their context refs. Each
// row shows the short commit hash, date, and the
// refs that were active when the commit was made.
//
// # Delegation
//
// Path parsing and git-log tracing are handled by
// trace/core/file. Path resolution uses rc.ContextDir
// and config/dir.
package file

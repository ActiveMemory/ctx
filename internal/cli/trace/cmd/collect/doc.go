//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package collect implements the hidden
// "ctx trace collect" cobra subcommand.
//
// This command gathers context refs from all sources
// and outputs them as a git commit trailer, or
// records refs from a commit trailer into persistent
// history.
//
// # Usage
//
//	ctx trace collect
//	ctx trace collect --record <trailer>
//
// # Flags
//
//	--record   When set, records the given commit
//	           trailer string into trace history and
//	           truncates the pending state file.
//	           When omitted, collects pending refs
//	           and prints the trailer to stdout.
//
// # Behavior
//
// Without --record the command:
//
//   - Reads the context directory path.
//   - Calls trace.Collect to gather pending context
//     refs from all registered sources (completed
//     tasks, decisions, learnings, etc.).
//   - Formats the refs as a git trailer line via
//     trace.FormatTrailer.
//   - Prints the trailer to stdout for the
//     prepare-commit-msg hook to inject.
//
// With --record the command delegates to
// core/collect.RecordCommit which parses the trailer,
// appends entries to the trace history file, and
// clears the pending state.
//
// # Output
//
// Without --record: a single trailer line suitable
// for appending to a commit message.
// With --record: no output (silent persistence).
//
// # Delegation
//
// Ref collection and formatting use the trace
// package. Recording uses trace/core/collect.
// Output formatting uses write/trace.
package collect

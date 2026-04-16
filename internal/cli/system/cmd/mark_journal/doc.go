//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mark_journal implements the hidden
// "ctx system mark-journal" cobra subcommand.
//
// This command tracks processing stages for journal
// files. Each journal entry passes through stages
// such as exported, enriched, and normalized; this
// command records which stage a file has reached or
// queries the current stage value.
//
// # Usage
//
//	ctx system mark-journal <filename> <stage>
//	ctx system mark-journal --check <filename> <stage>
//
// # Arguments
//
// The command requires exactly two positional args:
//
//   - filename: the journal file to mark or check.
//   - stage: the processing stage name (one of the
//     valid stages defined in journal/state).
//
// # Flags
//
//	--check   Query the current stage value instead
//	          of setting it. Prints the stored value
//	          for the given stage without modifying
//	          the state file.
//
// # Behavior
//
// Without --check the command delegates to
// core/journal.MarkStage to persist the stage, then
// prints a confirmation message.
//
// With --check the command calls
// core/journal.CheckStage and prints the current
// value for the requested stage.
//
// # Output
//
// Prints a one-line confirmation when marking, or
// the current stage value when checking.
//
// # Delegation
//
// Stage persistence is handled by
// system/core/journal. Output formatting is handled
// by write/mark_journal.
package mark_journal

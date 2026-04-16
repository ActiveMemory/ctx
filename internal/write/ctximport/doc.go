//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ctximport provides terminal output for memory
// import operations (ctx import).
//
// # Scan Phase
//
// [NoEntries] prints a notice when the source file
// contains no importable entries. [ScanHeader] prints
// the source filename and discovered entry count.
//
// # Classification Phase
//
// [EntrySkipped] prints a block for entries classified
// as "skip" (not promotable). [EntryClassified] prints
// a block for entries that matched a target file with
// keywords during dry-run preview. [EntryAdded] prints
// a block for entries that were successfully promoted
// to a context file.
//
// # Error Handling
//
// [ErrPromote] prints a promotion error to stderr when
// an entry cannot be written to its target file.
//
// # Summary
//
// [Summary] prints the full import summary with totals
// broken down by type (conventions, decisions, learnings,
// tasks), plus counts of skipped and duplicate entries.
// The summary adjusts its wording for dry-run mode.
//
// # Message Categories
//
//   - Info: scan results, classifications, promotions
//   - Error: promotion failures (stderr)
//   - Summary: aggregate counts with type breakdown
package ctximport

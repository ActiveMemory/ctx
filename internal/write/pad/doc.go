//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pad provides terminal output for the encrypted
// scratchpad command (ctx pad).
//
// The scratchpad stores short text entries and binary
// blobs, each encrypted at rest. Every mutation and
// query has a dedicated output function so the command
// layer stays free of presentation logic.
//
// # Text Entries
//
// [EntryAdded], [EntryUpdated], and [EntryRemoved] emit
// confirmation lines when entries change. [EntryMoved]
// reports position changes. [EntryShow] prints a single
// entry and [EntryList] prints a formatted list item.
// [Normalized] confirms sequential ID renumbering.
//
// # Binary Blobs
//
// [BlobWritten] confirms a blob was written to disk with
// its byte count and path. [BlobShow] prints raw blob
// data to stdout for piping.
//
// # Import and Export
//
// [ImportDone] reports the number of entries imported.
// [ImportNone] handles the empty-import case.
// [ImportBlobAdded] confirms each blob file imported.
// [ImportBlobSummary] closes import with counts.
// [ErrImportBlobSkipped] and [ErrImportBlobTooLarge]
// report per-blob failures to stderr.
//
// [ExportPlan] previews each blob in dry-run mode.
// [ExportDone] confirms each exported blob.
// [ExportSummary] closes the operation with totals.
//
// # Merge
//
// [MergeAdded] and [MergeDupe] report per-entry merge
// results. [MergeBlobConflict] warns about label
// collisions. [MergeBinaryWarning] flags binary data
// in a source file. [MergeSummary] closes the merge
// with added and skipped counts, adjusting for
// dry-run mode.
//
// # Tags and State
//
// [TagsItem] prints a tag with its count. [TagsJSON]
// emits JSON-encoded tag data. [TagsNone] handles the
// empty case. [Empty] reports an empty scratchpad.
// [KeyCreated] confirms encryption key generation.
//
// # Conflict Resolution
//
// [ResolveSide] renders OURS/THEIRS conflict blocks
// with numbered entries for interactive resolution.
package pad

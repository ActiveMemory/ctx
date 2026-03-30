//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pad provides terminal output for the encrypted scratchpad
// command (ctx pad).
//
// The scratchpad supports text entries and binary blobs, each with
// its own output group:
//
//   - Text entries: [EntryAdded], [EntryUpdated], [EntryRemoved],
//     [EntryMoved], [EntryShow], [EntryList]
//   - Binary blobs: [BlobWritten], [BlobShow]
//   - Import/export: [ImportDone], [ImportNone], [ImportBlobAdded],
//     [ExportPlan], [ExportDone], [ExportSummary]
//   - Merge: [MergeAdded], [MergeDupe], [MergeBlobConflict],
//     [MergeSummary]
//   - State: [Empty], [KeyCreated]
//
// Example:
//
//	write.EntryAdded(cmd, index)
//	write.EntryShow(cmd, formatted)
package pad

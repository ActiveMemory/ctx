//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for scratchpad merge output.
const (
	DescKeyWritePadMergeAdded          = "write.pad-merge-added"
	DescKeyWritePadMergeBinaryWarning  = "write.pad-merge-binary-warning"
	DescKeyWritePadMergeBlobConflict   = "write.pad-merge-blob-conflict"
	DescKeyWritePadMergeDone1Entry     = "write.pad-merge-done-1-entry"
	DescKeyWritePadMergeDoneNEntries   = "write.pad-merge-done-n-entries"
	DescKeyWritePadMergeDryRun1Entry   = "write.pad-merge-dry-run-1-entry"
	DescKeyWritePadMergeDryRunNEntries = "write.pad-merge-dry-run-n-entries"
	DescKeyWritePadMergeDupe           = "write.pad-merge-dupe"
	DescKeyWritePadMergeNone           = "write.pad-merge-none"
	DescKeyWritePadMergeNoneNew        = "write.pad-merge-none-new"
	DescKeyWritePadMergeSkipped1       = "write.pad-merge-skipped-1"
	DescKeyWritePadMergeSkippedN       = "write.pad-merge-skipped-n"
)

// DescKeys for scratchpad blob import output.
const (
	DescKeyWritePadImportBlobAdded    = "write.pad-import-blob-added"
	DescKeyWritePadImportBlobNone     = "write.pad-import-blob-none"
	DescKeyWritePadImportBlobSkipped  = "write.pad-import-blob-skipped"
	DescKeyWritePadImportBlobSummary  = "write.pad-import-blob-summary"
	DescKeyWritePadImportBlobTooLarge = "write.pad-import-blob-too-large"
	DescKeyWritePadImportCloseWarning = "write.pad-import-close-warning"
	DescKeyWritePadImportDone         = "write.pad-import-done"
	DescKeyWritePadImportNone         = "write.pad-import-none"
)

// DescKeys for scratchpad entry mutation output.
const (
	DescKeyWritePadEntryAdded   = "write.pad-entry-added"
	DescKeyWritePadEntryMoved   = "write.pad-entry-moved"
	DescKeyWritePadEntryRemoved = "write.pad-entry-removed"
	DescKeyWritePadEntryUpdated = "write.pad-entry-updated"
)

// DescKeys for scratchpad export output.
const (
	DescKeyWritePadExportDone        = "write.pad-export-done"
	DescKeyWritePadExportNone        = "write.pad-export-none"
	DescKeyWritePadExportPlan        = "write.pad-export-plan"
	DescKeyWritePadExportSummary     = "write.pad-export-summary"
	DescKeyWritePadExportVerbDone    = "write.pad-export-verb-done"
	DescKeyWritePadExportVerbDryRun  = "write.pad-export-verb-dry-run"
	DescKeyWritePadExportWriteFailed = "write.pad-export-write-failed"
)

// DescKeys for scratchpad list and blob output.
const (
	DescKeyWritePadBlobWritten = "write.pad-blob-written"
	DescKeyWritePadEmpty       = "write.pad-empty"
	DescKeyWritePadListItem    = "write.pad-list-item"
)

// DescKeys for scratchpad conflict resolution.
const (
	DescKeyWritePadResolveEntry  = "write.pad-resolve-entry"
	DescKeyWritePadResolveHeader = "write.pad-resolve-header"
)

// DescKeys for scratchpad operations.
const (
	DescKeyWritePadKeyCreated = "write.pad-key-created"
)

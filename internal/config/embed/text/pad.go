//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for scratchpad merge output.
const (
	// DescKeyWritePadMergeAdded is the text key for write pad merge added
	// messages.
	DescKeyWritePadMergeAdded = "write.pad-merge-added"
	// DescKeyWritePadMergeBinaryWarning is the text key for write pad merge
	// binary warning messages.
	DescKeyWritePadMergeBinaryWarning = "write.pad-merge-binary-warning"
	// DescKeyWritePadMergeBlobConflict is the text key for write pad merge blob
	// conflict messages.
	DescKeyWritePadMergeBlobConflict = "write.pad-merge-blob-conflict"
	// DescKeyWritePadMergeDone1Entry is the text key for write pad merge
	// done1entry messages.
	DescKeyWritePadMergeDone1Entry = "write.pad-merge-done-1-entry"
	// DescKeyWritePadMergeDoneNEntries is the text key for write pad merge done n
	// entries messages.
	DescKeyWritePadMergeDoneNEntries = "write.pad-merge-done-n-entries"
	// DescKeyWritePadMergeDryRun1Entry is the text key for write pad merge dry
	// run1entry messages.
	DescKeyWritePadMergeDryRun1Entry = "write.pad-merge-dry-run-1-entry"
	// DescKeyWritePadMergeDryRunNEntries is the text key for write pad merge dry
	// run n entries messages.
	DescKeyWritePadMergeDryRunNEntries = "write.pad-merge-dry-run-n-entries"
	// DescKeyWritePadMergeDupe is the text key for write pad merge dupe messages.
	DescKeyWritePadMergeDupe = "write.pad-merge-dupe"
	// DescKeyWritePadMergeNone is the text key for write pad merge none messages.
	DescKeyWritePadMergeNone = "write.pad-merge-none"
	// DescKeyWritePadMergeNoneNew is the text key for write pad merge none new
	// messages.
	DescKeyWritePadMergeNoneNew = "write.pad-merge-none-new"
	// DescKeyWritePadMergeSkipped1 is the text key for write pad merge skipped1
	// messages.
	DescKeyWritePadMergeSkipped1 = "write.pad-merge-skipped-1"
	// DescKeyWritePadMergeSkippedN is the text key for write pad merge skipped n
	// messages.
	DescKeyWritePadMergeSkippedN = "write.pad-merge-skipped-n"
)

// DescKeys for scratchpad blob import output.
const (
	// DescKeyWritePadImportBlobAdded is the text key for write pad import blob
	// added messages.
	DescKeyWritePadImportBlobAdded = "write.pad-import-blob-added"
	// DescKeyWritePadImportBlobNone is the text key for write pad import blob
	// none messages.
	DescKeyWritePadImportBlobNone = "write.pad-import-blob-none"
	// DescKeyWritePadImportBlobSkipped is the text key for write pad import blob
	// skipped messages.
	DescKeyWritePadImportBlobSkipped = "write.pad-import-blob-skipped"
	// DescKeyWritePadImportBlobSummary is the text key for write pad import blob
	// summary messages.
	DescKeyWritePadImportBlobSummary = "write.pad-import-blob-summary"
	// DescKeyWritePadImportBlobTooLarge is the text key for write pad import blob
	// too large messages.
	DescKeyWritePadImportBlobTooLarge = "write.pad-import-blob-too-large"
	// DescKeyWritePadImportCloseWarning is the text key for write pad import
	// close warning messages.
	DescKeyWritePadImportCloseWarning = "write.pad-import-close-warning"
	// DescKeyWritePadImportDone is the text key for write pad import done
	// messages.
	DescKeyWritePadImportDone = "write.pad-import-done"
	// DescKeyWritePadImportNone is the text key for write pad import none
	// messages.
	DescKeyWritePadImportNone = "write.pad-import-none"
)

// DescKeys for scratchpad entry mutation output.
const (
	// DescKeyWritePadEntryAdded is the text key for write pad entry added
	// messages.
	DescKeyWritePadEntryAdded = "write.pad-entry-added"
	// DescKeyWritePadEntryMoved is the text key for write pad entry moved
	// messages.
	DescKeyWritePadEntryMoved = "write.pad-entry-moved"
	// DescKeyWritePadEntryRemoved is the text key for write pad entry removed
	// messages.
	DescKeyWritePadEntryRemoved = "write.pad-entry-removed"
	// DescKeyWritePadEntryUpdated is the text key for write pad entry updated
	// messages.
	DescKeyWritePadEntryUpdated = "write.pad-entry-updated"
)

// DescKeys for scratchpad export output.
const (
	// DescKeyWritePadExportDone is the text key for write pad export done
	// messages.
	DescKeyWritePadExportDone = "write.pad-export-done"
	// DescKeyWritePadExportNone is the text key for write pad export none
	// messages.
	DescKeyWritePadExportNone = "write.pad-export-none"
	// DescKeyWritePadExportPlan is the text key for write pad export plan
	// messages.
	DescKeyWritePadExportPlan = "write.pad-export-plan"
	// DescKeyWritePadExportSummary is the text key for write pad export summary
	// messages.
	DescKeyWritePadExportSummary = "write.pad-export-summary"
	// DescKeyWritePadExportVerbDone is the text key for write pad export verb
	// done messages.
	DescKeyWritePadExportVerbDone = "write.pad-export-verb-done"
	// DescKeyWritePadExportVerbDryRun is the text key for write pad export verb
	// dry run messages.
	DescKeyWritePadExportVerbDryRun = "write.pad-export-verb-dry-run"
	// DescKeyWritePadExportWriteFailed is the text key for write pad export write
	// failed messages.
	DescKeyWritePadExportWriteFailed = "write.pad-export-write-failed"
)

// DescKeys for scratchpad list and blob output.
const (
	// DescKeyWritePadBlobWritten is the text key for write pad blob written
	// messages.
	DescKeyWritePadBlobWritten = "write.pad-blob-written"
	// DescKeyWritePadEmpty is the text key for write pad empty messages.
	DescKeyWritePadEmpty = "write.pad-empty"
	// DescKeyWritePadListItem is the text key for write pad list item messages.
	DescKeyWritePadListItem = "write.pad-list-item"
)

// DescKeys for scratchpad conflict resolution.
const (
	// DescKeyWritePadResolveEntry is the text key for write pad resolve entry
	// messages.
	DescKeyWritePadResolveEntry = "write.pad-resolve-entry"
	// DescKeyWritePadResolveHeader is the text key for write pad resolve header
	// messages.
	DescKeyWritePadResolveHeader = "write.pad-resolve-header"
)

// DescKeys for scratchpad operations.
const (
	// DescKeyWritePadKeyCreated is the text key for write pad key created
	// messages.
	DescKeyWritePadKeyCreated = "write.pad-key-created"
)

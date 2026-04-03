//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for import operations.
const (
	DescKeyImportCountConvention = "import.count-convention"
	DescKeyImportCountDecision   = "import.count-decision"
	DescKeyImportCountLearning   = "import.count-learning"
	DescKeyImportCountTask       = "import.count-task"
)

// DescKeys for memory import write output.
const (
	DescKeyWriteImportAdded          = "write.import-added"
	DescKeyWriteImportBreakdown      = "write.import-breakdown"
	DescKeyWriteImportClassified     = "write.import-classified"
	DescKeyWriteImportClassifiedSkip = "write.import-classified-skip"
	DescKeyWriteImportDuplicates     = "write.import-duplicates"
	DescKeyWriteImportEntry          = "write.import-entry"
	DescKeyWriteImportFound          = "write.import-found"
	DescKeyWriteImportNoEntries      = "write.import-no-entries"
	DescKeyWriteImportScanning       = "write.import-scanning"
	DescKeyWriteImportSkipped        = "write.import-skipped"
	DescKeyWriteImportErrorPromote   = "write.import-error-promote"
	DescKeyWriteImportSummary        = "write.import-summary"
	DescKeyWriteImportSummaryDryRun  = "write.import-summary-dry-run"
)

// DescKeys for journal import write output.
const (
	DescKeyWriteJournalImportNothing      = "write.journal-import-nothing"
	DescKeyWriteJournalImportPartNew      = "write.journal-import-part-new"
	DescKeyWriteJournalImportPartRegen    = "write.journal-import-part-regen"
	DescKeyWriteJournalImportPartSkip     = "write.journal-import-part-skip"
	DescKeyWriteJournalImportPartSkipLock = "write.journal-import-part-skip-locked"
	DescKeyWriteJournalImportSummary      = "write.journal-import-summary"
	DescKeyWriteJournalImportVerb         = "write.journal-import-verb"
	DescKeyWriteJournalImportVerbDryRun   = "write.journal-import-verb-dry-run"
)

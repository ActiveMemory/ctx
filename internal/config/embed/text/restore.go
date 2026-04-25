//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for restore operations write output.
const (
	// DescKeyWriteRestoreAdded is the text key for write restore added messages.
	DescKeyWriteRestoreAdded = "write.restore-added"
	// DescKeyWriteRestoreDenyDroppedHeader is the text key for write restore deny
	// dropped header messages.
	DescKeyWriteRestoreDenyDroppedHeader = "write.restore-deny-dropped-header"
	// DescKeyWriteRestoreDenyRestoredHeader is the text key for write restore
	// deny restored header messages.
	DescKeyWriteRestoreDenyRestoredHeader = "write.restore-deny-restored-header"
	// DescKeyWriteRestoreDone is the text key for write restore done messages.
	DescKeyWriteRestoreDone = "write.restore-done"
	// DescKeyWriteRestoreDroppedHeader is the text key for write restore dropped
	// header messages.
	DescKeyWriteRestoreDroppedHeader = "write.restore-dropped-header"
	// DescKeyWriteRestoreMatch is the text key for write restore match messages.
	DescKeyWriteRestoreMatch = "write.restore-match"
	// DescKeyWriteRestoreNoLocal is the text key for write restore no local
	// messages.
	DescKeyWriteRestoreNoLocal = "write.restore-no-local"
	// DescKeyWriteRestorePermMatch is the text key for write restore perm match
	// messages.
	DescKeyWriteRestorePermMatch = "write.restore-perm-match"
	// DescKeyWriteRestoreRemoved is the text key for write restore removed
	// messages.
	DescKeyWriteRestoreRemoved = "write.restore-removed"
	// DescKeyWriteRestoreRestoredHeader is the text key for write restore
	// restored header messages.
	DescKeyWriteRestoreRestoredHeader = "write.restore-restored-header"
	// DescKeyWriteSnapshotSaved is the text key for the first-time
	// golden snapshot save confirmation.
	DescKeyWriteSnapshotSaved = "write.snapshot-saved"
	// DescKeyWriteSnapshotUpdated is the text key for the subsequent
	// golden snapshot update confirmation.
	DescKeyWriteSnapshotUpdated = "write.snapshot-updated"
)

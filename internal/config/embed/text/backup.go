//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for backup operations.
const (
	DescKeyBackupBoxTitle       = "backup.box-title"
	DescKeyBackupNoMarker       = "backup.no-marker"
	DescKeyBackupRelayMessage   = "backup.relay-message"
	DescKeyBackupRelayPrefix    = "backup.relay-prefix"
	DescKeyBackupRunHint        = "backup.run-hint"
	DescKeyBackupSMBNotMounted  = "backup.smb-not-mounted"
	DescKeyBackupSMBUnavailable = "backup.smb-unavailable"
	DescKeyBackupStale          = "backup.stale"
)

// DescKeys for backup result write output.
const (
	DescKeyWriteBackupResult  = "write.backup-result"
	DescKeyWriteBackupSMBDest = "write.backup-smb-dest"
)

// DescKeys for snapshot write output.
const (
	DescKeyWriteSnapshotSaved   = "write.snapshot-saved"
	DescKeyWriteSnapshotUpdated = "write.snapshot-updated"
)

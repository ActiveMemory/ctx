//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for the surviving err/backup constructors (see
// internal/err/backup). The former ctx-backup-specific keys were
// removed when the command was deprecated; these four remain
// because init's `.bak` writer, task archival, and bootstrap all
// still use the package.
const (
	// DescKeyErrBackupContextDirNotFound is the text key for the
	// "context directory not found" bootstrap error.
	DescKeyErrBackupContextDirNotFound = "err.backup.context-dir-not-found"
	// DescKeyErrBackupCreateArchiveDir is the text key for task
	// archive directory creation failures.
	DescKeyErrBackupCreateArchiveDir = "err.backup.create-archive-dir"
	// DescKeyErrBackupCreateBackup is the text key for the `.bak`
	// file creation failure (ctx init --reset).
	DescKeyErrBackupCreateBackup = "err.backup.create-backup"
	// DescKeyErrBackupWriteArchive is the text key for task archive
	// write failures.
	DescKeyErrBackupWriteArchive = "err.backup.write-archive"
)

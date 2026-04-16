//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package archive provides terminal output for task
// archival and snapshot operations (ctx task archive,
// ctx task snapshot).
//
// # Task Archival
//
// Functions cover the full lifecycle of an archive
// operation. [DryRun] previews what would be archived
// with counts and a content preview. [Success] reports
// the number of tasks archived and the output file path.
// [NoCompleted] handles the empty case when no completed
// tasks exist. [Skipping] explains why a specific parent
// task was excluded due to incomplete children, and
// [SkipIncomplete] summarizes the total skip count.
//
// # Task Snapshots
//
// [SnapshotSaved] confirms a snapshot was written and
// prints the output file path. [SnapshotContent] formats
// the snapshot body with a name header, creation
// timestamp, and separator, returning the assembled
// string for the caller to write to disk.
//
// # Message Categories
//
//   - Info: archive and snapshot confirmations
//   - Warning: skip notices for incomplete tasks
//
// # Usage
//
//	if dryRun {
//	    archive.DryRun(cmd, count, pending, preview, sep)
//	    return
//	}
//	archive.Success(cmd, count, path, pending)
package archive

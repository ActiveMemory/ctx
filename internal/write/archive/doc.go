//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package archive provides terminal output for task archival and
// snapshot operations.
//
// All functions take a *cobra.Command for output routing. The package
// handles two related workflows:
//
//   - Task archival: [DryRun] previews what would be archived,
//     [Success] reports completion, [NoCompleted] handles the
//     empty case, and [Skipping]/[SkipIncomplete] explain why
//     specific tasks were excluded.
//   - Task snapshots: [SnapshotSaved] confirms the write path
//     and [SnapshotContent] formats the snapshot body with a
//     timestamp header and separator.
//
// Example usage from a command's Run function:
//
//	if dryRun {
//	    write.DryRun(cmd, tasks, archivePath)
//	    return
//	}
//	write.Success(cmd, count, archivePath)
package archive

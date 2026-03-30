//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package scan queries the filesystem and git history for changes
// since a reference time.
//
// [FindContextChanges] returns .context/ files modified after the
// reference time. [SummarizeCodeChanges] extracts commit counts,
// latest message, affected directories, and authors from git log.
// [UniqueTopDirs] deduplicates directory paths from git output.
package scan

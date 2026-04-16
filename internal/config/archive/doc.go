//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package archive centralizes constants for task archival,
// project backups, and snapshot management.
//
// ctx archives completed tasks into dated markdown files
// and backs up entire project contexts to SMB shares as
// compressed tarballs. This package defines file-naming
// patterns, backup scopes, staleness thresholds, and
// template variables used by both subsystems.
//
// # Task Snapshots
//
// When a user archives tasks, ctx creates a timestamped
// markdown file in .context/archive/:
//
//   - [ScopeTasks] identifies the task archive scope.
//   - [SnapshotFilenameFormat] and [SnapshotTimeFormat]
//     produce filenames like tasks-snapshot-2026-04-15-0930.md.
//   - [DefaultSnapshotName] provides the fallback name.
//   - [TplFilename] and [DateSep] control the general
//     archive filename template and header formatting.
//
// # Backup System
//
// Full project backups are tar.gz archives written to an
// SMB share:
//
//   - [BackupScopeProject], [BackupScopeGlobal], and
//     [BackupScopeAll] control what gets backed up.
//   - [TplProjectArchive] and [TplGlobalArchive] name the
//     output files with timestamps.
//   - [BackupDefaultSubdir] sets the target subdirectory
//     on the SMB share (ctx-sessions).
//   - [BackupMaxAgeDays] triggers a nudge when the last
//     backup is older than 2 days.
//   - [BackupMarkerFile] and [BackupMarkerDir] track the
//     last successful backup timestamp.
//   - [BackupThrottleID] ensures the staleness nudge
//     fires at most once per day.
//
// # Task Parsing
//
// [SubTaskMinIndent] defines the minimum indentation (2
// spaces) for a line to be treated as a subtask rather
// than a top-level task during archive parsing.
//
// # Writer Identifiers
//
// [WriterGzip] and [WriterTar] label the compression and
// archival stages for structured error reporting.
//
// # Why Centralized
//
// Filename templates, scopes, and staleness thresholds
// are shared between the archive command, the backup
// skill, and the staleness hook. Centralizing them
// prevents drift and makes the naming scheme auditable.
package archive

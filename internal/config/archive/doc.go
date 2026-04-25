//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package archive centralizes constants for task archival and
// snapshot management.
//
// When a user archives tasks, ctx creates a timestamped markdown
// file in .context/archive/:
//
//   - [ScopeTasks] identifies the task archive scope.
//   - [SnapshotFilenameFormat] and [SnapshotTimeFormat]
//     produce filenames like tasks-snapshot-2026-04-15-0930.md.
//   - [DefaultSnapshotName] provides the fallback name.
//   - [TplFilename] and [DateSep] control the general
//     archive filename template and header formatting.
//
// [SubTaskMinIndent] defines the minimum indentation (2 spaces)
// for a line to be treated as a subtask rather than a top-level
// task during archive parsing.
//
// Filename templates, scopes, and date formatting are shared
// between the task archive command and the tidy/compact helpers.
// Centralizing them prevents drift and makes the naming scheme
// auditable.
package archive

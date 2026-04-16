//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task moves completed tasks from their original
// phase sections in TASKS.md into the Completed section at
// the bottom of the file.
//
// The compaction algorithm scans TASKS.md for checked items
// ("- [x]") that appear outside the Completed section. Each
// match includes its nested content: indented lines below
// the task line are treated as subtasks or details and move
// together with their parent. Tasks with at least one
// unchecked subtask are skipped to avoid orphaning
// in-progress work.
//
// # Archival
//
// When the archive flag is set, compacted tasks are also
// written to a dated file in .context/archive/ with a
// standardized heading. This preserves traceability for
// completed work without cluttering the active task list.
//
// # Write Safety
//
// File writes use [internal/io.SafeWriteFile] which writes
// atomically (temp + rename) to avoid partial writes on
// crash. The function returns the count of tasks moved so
// callers can report progress.
package task
